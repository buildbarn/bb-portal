package buildeventrecorder

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/ent/gen/ent/connectionmetadata"
	apicommon "github.com/buildbarn/bb-portal/internal/api/common"
	"github.com/buildbarn/bb-portal/internal/database/common"
	"github.com/buildbarn/bb-portal/pkg/processing"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-portal/pkg/summary/detectors"
	"github.com/buildbarn/bb-storage/pkg/auth"
	"github.com/buildbarn/bb-storage/pkg/util"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// BuildEventRecorder contains information about a Bazel invocation that is
// required when processing individual build events.
type BuildEventRecorder struct {
	db                  *ent.Client
	problemDetector     detectors.ProblemDetector
	blobArchiver        processing.BlobMultiArchiver
	saveTargetDataLevel *bb_portal.BuildEventStreamService_SaveTargetDataLevel
	tracer              trace.Tracer

	InstanceName           string
	InstanceNameDbID       int
	InvocationID           string
	InvocationDbID         int
	CorrelatedInvocationID string
	IsRealTime             bool
}

// NewBuildEventRecorder creates a new BuildEventRecorder
func NewBuildEventRecorder(
	ctx context.Context,
	db *ent.Client,
	instanceNameAuthorizer auth.Authorizer,
	blobArchiver processing.BlobMultiArchiver,
	saveTargetDataLevel *bb_portal.BuildEventStreamService_SaveTargetDataLevel,
	tracerProvider trace.TracerProvider,
	instanceName string,
	invocationID string,
	correlatedInvocationID string,
	isRealTime bool,
) (*BuildEventRecorder, error) {
	if invocationID == "" {
		return nil, status.Error(codes.InvalidArgument, "Invocation ID is required")
	}
	if !apicommon.IsInstanceNameAllowed(ctx, instanceNameAuthorizer, instanceName) {
		return nil, status.Error(codes.PermissionDenied, fmt.Sprintf("Instance name %q is not allowed", instanceName))
	}

	tracer := tracerProvider.Tracer("github.com/buildbarn/bb-portal/internal/database/buildeventrecorder")

	instanceNameDbID, invocationDbID, err := findOrCreateInvocation(ctx, db, invocationID, instanceName, tracer)
	if err != nil {
		return nil, util.StatusWrap(err, "Failed to find or create bazel invocation")
	}

	return &BuildEventRecorder{
		db:                  db,
		problemDetector:     detectors.NewProblemDetector(),
		blobArchiver:        blobArchiver,
		saveTargetDataLevel: saveTargetDataLevel,
		tracer:              tracer,

		InstanceName:           instanceName,
		InstanceNameDbID:       instanceNameDbID,
		InvocationID:           invocationID,
		InvocationDbID:         invocationDbID,
		CorrelatedInvocationID: correlatedInvocationID,
		IsRealTime:             isRealTime,
	}, nil
}

func findOrCreateInvocation(
	ctx context.Context,
	db *ent.Client,
	invocationID string,
	instanceName string,
	tracer trace.Tracer,
) (int, int, error) {
	ctx, span := tracer.Start(ctx,
		fmt.Sprintf("BuildEventRecorder.findOrCreateInvocation"),
		trace.WithAttributes(
			attribute.String("invocation.id", invocationID),
			attribute.String("invocation.instance_name", instanceName),
		),
	)
	defer span.End()

	invocationIDUUID, err := uuid.Parse(invocationID)
	if err != nil {
		return 0, 0, util.StatusWrap(err, "Failed to parse invocation ID")
	}

	tx, err := db.BeginTx(ctx, &entsql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return 0, 0, util.StatusWrap(err, "Failed to create transaction")
	}

	instanceNameDbID, err := tx.InstanceName.Create().
		SetName(instanceName).
		OnConflictColumns("name").
		Ignore().
		ID(ctx)
	if err != nil {
		return 0, 0, common.RollbackAndWrapError(tx, util.StatusWrap(err, "Failed to create or find instance name"))
	}

	invocationDbID, err := tx.BazelInvocation.Create().
		SetInvocationID(invocationIDUUID).
		SetInstanceNameID(instanceNameDbID).
		OnConflictColumns(bazelinvocation.FieldInvocationID).
		Ignore().
		ID(ctx)
	if err != nil {
		return 0, 0, common.RollbackAndWrapError(tx, util.StatusWrap(err, "Failed to create or find bazel invocation"))
	}

	invocationDb, err := tx.BazelInvocation.Query().
		Where(bazelinvocation.ID(invocationDbID)).
		Only(ctx)
	if err != nil {
		return 0, 0, common.RollbackAndWrapError(tx, util.StatusWrap(err, "Failed to query bazel invocation by ID"))
	}

	instanceNameDb, err := invocationDb.InstanceName(ctx)
	if err != nil {
		return 0, 0, common.RollbackAndWrapError(tx, util.StatusWrap(err, "Failed to get existing invocation instance name"))
	}

	if instanceNameDb.ID != instanceNameDbID {
		return 0, 0, common.RollbackAndWrapError(tx, status.Errorf(codes.FailedPrecondition, "Invocation with ID %q already exists with different instance name. Possible UUID collision.", invocationID))
	}
	if invocationDb.BepCompleted {
		return 0, 0, common.RollbackAndWrapError(tx, status.Errorf(codes.FailedPrecondition, "Invocation with ID %q already exists and is locked for writing", invocationID))
	}

	err = tx.Commit()
	if err != nil {
		return 0, 0, util.StatusWrap(err, "Failed to commit transaction")
	}

	return instanceNameDbID, invocationDb.ID, nil
}

// StartLoggingConnectionMetadata starts a goroutine that periodically
// logs connection metadata for the invocation associated with this
// BuildEventRecorder. It continues to do so until the provided context
// is canceled.
func (r *BuildEventRecorder) StartLoggingConnectionMetadata(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				err := r.db.ConnectionMetadata.Create().
					SetBazelInvocationID(r.InvocationDbID).
					SetConnectionLastOpenAt(time.Now()).
					OnConflictColumns(connectionmetadata.BazelInvocationColumn).
					UpdateNewValues().
					Exec(ctx)
				if err != nil {
					slog.Warn("Failed to log connection metadata", "err", err)
				}
				time.Sleep(1 * time.Second)
			}
		}
	}()
}
