package cal

import (
	"context"
	"fmt"
	"log"

	"entgo.io/ent/dialect/sql"
	remoteexecution "github.com/bazelbuild/remote-apis/build/bazel/remote/execution/v2"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/digest"
	apicommon "github.com/buildbarn/bb-portal/internal/api/common"
	"github.com/buildbarn/bb-portal/internal/database"
	prometheusmetrics "github.com/buildbarn/bb-portal/pkg/prometheus_metrics"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-remote-execution/pkg/proto/completedactionlogger"
	"github.com/buildbarn/bb-storage/pkg/auth"
	"github.com/buildbarn/bb-storage/pkg/blobstore"
	"github.com/buildbarn/bb-storage/pkg/blobstore/buffer"
	pkgDigest "github.com/buildbarn/bb-storage/pkg/digest"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CompletedActionLoggerServer implements the Completed Action Logger
// service. It receives completed actions and stores them in the database
type CompletedActionLoggerServer struct {
	db            database.Client
	logAuthorizer auth.Authorizer
	cas           blobstore.BlobAccess
}

// NewCompletedActionLoggerServer creates a new CompletedActionLoggerServer
func NewCompletedActionLoggerServer(
	db database.Client,
	configuration *bb_portal.CompletedActionLoggerService,
	logAuthorizer auth.Authorizer,
	cas blobstore.BlobAccess,
) (*CompletedActionLoggerServer, error) {
	return &CompletedActionLoggerServer{
		db:            db,
		logAuthorizer: logAuthorizer,
		cas:           cas,
	}, nil
}

func createActionDigest(ctx context.Context, s *CompletedActionLoggerServer, req *completedactionlogger.CompletedAction) (*int64, error) {
	e := s.db.Ent()
	actionDigestID, err := e.Digest.Create().
		SetRev2InstanceName(req.GetInstanceName()).
		SetHash([]byte(req.GetHistoricalExecuteResponse().GetActionDigest().GetHash())).
		SetSizeBytes(req.GetHistoricalExecuteResponse().GetActionDigest().GetSizeBytes()).
		SetDigestFunction(int16(req.GetDigestFunction())).
		OnConflictColumns(
			digest.FieldRev2InstanceName,
			digest.FieldHash,
			digest.FieldSizeBytes,
			digest.FieldDigestFunction,
		).
		UpdateNewValues().
		ID(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to create action digest: %w", err)
	}

	return &actionDigestID, nil
}

func addHistoricalExecuteResponseToCAS(ctx context.Context, req *completedactionlogger.CompletedAction, s *CompletedActionLoggerServer) (*pkgDigest.Digest, error) {
	// Serialize execute response
	historicalExecuteResponse := req.GetHistoricalExecuteResponse()

	data, err := proto.Marshal(historicalExecuteResponse)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal historical execute response: %w", err)
	}
	digestFunction := pkgDigest.MustNewFunction(
		req.GetInstanceName(),
		req.GetDigestFunction(),
	)

	// Create a digest for the serialized ExecuteResponse.
	generator := digestFunction.NewGenerator(int64(len(data)))

	if _, err := generator.Write(data); err != nil {
		return nil, fmt.Errorf("Failed to generate execute response digest: %w", err)
	}

	executeResponseDigest := generator.Sum()

	if err := s.cas.Put(
		ctx,
		executeResponseDigest,
		buffer.NewValidatedBufferFromByteSlice(data),
	); err != nil {
		return nil, fmt.Errorf("Failed to put execute response in CAS: %w", err)
	}

	return &executeResponseDigest, nil
}

func createHistoricalExecuteResponseDigest(ctx context.Context, req *completedactionlogger.CompletedAction, s *CompletedActionLoggerServer) (*int64, error) {
	executeResponseDigest, err := addHistoricalExecuteResponseToCAS(ctx, req, s)
	if err != nil {
		return nil, fmt.Errorf("Failed to put execute response in CAS: %w", err)
	}

	e := s.db.Ent()
	historicalExecuteResponseDigestID, err := e.Digest.Create().
		SetRev2InstanceName(req.GetInstanceName()).
		SetHash([]byte(executeResponseDigest.GetHashBytes())).
		SetSizeBytes(executeResponseDigest.GetSizeBytes()).
		SetDigestFunction(int16(executeResponseDigest.GetDigestFunction().GetEnumValue())).
		OnConflictColumns(
			digest.FieldRev2InstanceName,
			digest.FieldHash,
			digest.FieldSizeBytes,
			digest.FieldDigestFunction,
		).
		UpdateNewValues().
		ID(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to upsert historical execute response digest: %w", err)
	}

	return &historicalExecuteResponseDigestID, nil
}

// LogCompletedActions receives completed actions and stores them in the database
func (s *CompletedActionLoggerServer) LogCompletedActions(stream completedactionlogger.CompletedActionLogger_LogCompletedActionsServer) error {
	ctx := stream.Context()
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		if !apicommon.IsInstanceNameAllowed(ctx, s.logAuthorizer, req.GetInstanceName()) {
			log.Printf("Instance name %q is not allowed", req.GetInstanceName())
		} else {
			auxMetadata := req.GetHistoricalExecuteResponse().GetExecuteResponse().GetResult().GetExecutionMetadata().GetAuxiliaryMetadata()
			for _, anyMsg := range auxMetadata {
				switch {
				case anyMsg.MessageIs(&remoteexecution.RequestMetadata{}):
					var msg remoteexecution.RequestMetadata
					if err := anyMsg.UnmarshalTo(&msg); err != nil {
						log.Printf("Failed to unmarshal RequestMetadata: %v", err)
						return err
					}

					actionDigestID, err := createActionDigest(ctx, s, req)
					if err != nil {
						return fmt.Errorf("Failed to create action digest: %w", err)
					}

					historicalExecuteResponseDigestID, err := createHistoricalExecuteResponseDigest(ctx, req, s)
					if err != nil {
						return fmt.Errorf("Failed to create historical execute response digest: %w", err)
					}

					// Add the completed action
					e := s.db.Ent()
					_, err = e.CompletedAction.Create().
						SetInvocationID(msg.GetToolInvocationId()).
						SetTargetID(msg.GetTargetId()).
						SetActionDigestID(*actionDigestID).
						SetHistoricalExecuteResponseDigestID(*historicalExecuteResponseDigestID).
						SetActionID(msg.GetActionId()).
						SetToolName(msg.GetToolDetails().GetToolName()).
						SetMnemonic(msg.GetActionMnemonic()).
						OnConflict(
							sql.ConflictColumns(
								"completed_action_historical_execute_response_digest",
							),
						).
						Ignore().
						ID(ctx)
					if err != nil {
						if ent.IsConstraintError(err) {
							log.Println("On Conflict Triggered: CompletedAction already exists. Ignoring insert.")
							return nil
						}
						return fmt.Errorf("Failed to insert completed action: %w", err)
					}
					prometheusmetrics.CompletedActions.Inc()
					if err := stream.Send(&emptypb.Empty{}); err != nil {
						return fmt.Errorf("Failed to send response: %w", err)
					}
				}
			}
		}
	}
}
