package bes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"

	build "google.golang.org/genproto/googleapis/devtools/build/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/internal/database/buildeventrecorder"
	"github.com/buildbarn/bb-portal/internal/database/dbauthservice"
	"github.com/buildbarn/bb-portal/pkg/authmetadataextraction"
	"github.com/buildbarn/bb-portal/pkg/events"
	"github.com/buildbarn/bb-portal/pkg/processing"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-storage/pkg/auth"
	auth_configuration "github.com/buildbarn/bb-storage/pkg/auth/configuration"
	bb_grpc "github.com/buildbarn/bb-storage/pkg/grpc"
	"github.com/buildbarn/bb-storage/pkg/program"
	"github.com/buildbarn/bb-storage/pkg/util"

	"go.opentelemetry.io/otel/trace"
)

// BuildEventServer implements the Build Event Service.
// It receives events and forwards them to a BuildEventChannel.
// TODO: Should this support forwarding events? Users might want to create their own
// tooling that reacts to build events, and it would be useful if this service could
// forward events to those.
type BuildEventServer struct {
	db                     *ent.Client
	instanceNameAuthorizer auth.Authorizer
	blobArchiver           processing.BlobMultiArchiver
	saveTargetDataLevel    *bb_portal.BuildEventStreamService_SaveTargetDataLevel
	tracerProvider         trace.TracerProvider
	extractors             *authmetadataextraction.AuthMetadataExtractors
	uuidGenerator          util.UUIDGenerator
}

// NewBuildEventServer creates a new BuildEventServer
func NewBuildEventServer(db *ent.Client, blobArchiver processing.BlobMultiArchiver, configuration *bb_portal.ApplicationConfiguration, dependenciesGroup program.Group, grpcClientFactory bb_grpc.ClientFactory, tracerProvider trace.TracerProvider, uuidGenerator util.UUIDGenerator) (*BuildEventServer, error) {
	if configuration.InstanceNameAuthorizer == nil {
		return nil, status.Error(codes.NotFound, "No InstanceNameAuthorizer configured")
	}
	instanceNameAuthorizer, err := auth_configuration.DefaultAuthorizerFactory.NewAuthorizerFromConfiguration(configuration.InstanceNameAuthorizer, dependenciesGroup, grpcClientFactory)
	if err != nil {
		return nil, util.StatusWrap(err, "Failed to create InstanceNameAuthorizer")
	}

	besConfiguration := configuration.BesServiceConfiguration
	if besConfiguration == nil {
		return nil, fmt.Errorf("No BesServiceConfiguration configured")
	}

	saveTargetDataLevel := besConfiguration.SaveTargetDataLevel
	if saveTargetDataLevel == nil || saveTargetDataLevel.Level == nil {
		return nil, fmt.Errorf("No saveTargetDataLevel configured")
	}

	extractors, err := authmetadataextraction.AuthMetadataExtractorsFromConfiguration(besConfiguration.AuthMetadataKeyConfiguration, dependenciesGroup)
	if err != nil {
		return nil, util.StatusWrap(err, "Failed to create AutheMetadataExtractors")
	}

	return &BuildEventServer{
		instanceNameAuthorizer: instanceNameAuthorizer,
		db:                     db,
		blobArchiver:           blobArchiver,
		saveTargetDataLevel:    saveTargetDataLevel,
		tracerProvider:         tracerProvider,
		extractors:             extractors,
		uuidGenerator:          uuidGenerator,
	}, nil
}

// PublishLifecycleEvent handles life cycle events.
func (s BuildEventServer) PublishLifecycleEvent(ctx context.Context, request *build.PublishLifecycleEventRequest) (*emptypb.Empty, error) {
	slog.InfoContext(ctx, "Received event", "event", protojson.Format(request.BuildEvent.GetEvent()))
	return &emptypb.Empty{}, nil
}

// PublishBuildToolEventStream handles a build tool event stream.
func (s BuildEventServer) PublishBuildToolEventStream(stream build.PublishBuildEvent_PublishBuildToolEventStreamServer) error {
	// We can safely bypass authorization checks here, as we check that the
	// user is allowed to upload to the instance name when creating the
	// BuildEventRecorder.
	ctx := dbauthservice.NewContextWithDbAuthServiceBypass(stream.Context())
	slog.InfoContext(ctx, "Stream started", "event", ctx)

	var buildEventRecorder *buildeventrecorder.BuildEventRecorder = nil

	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				slog.InfoContext(ctx, "Stream finished", "event", ctx)
				return nil
			}
			slog.ErrorContext(ctx, "Recv failed", "err", err)
			return err
		}

		if buildEventRecorder == nil {
			buildEventRecorder, err = buildeventrecorder.NewBuildEventRecorder(
				ctx,
				s.db,
				s.instanceNameAuthorizer,
				s.blobArchiver,
				s.saveTargetDataLevel,
				s.tracerProvider,
				req.GetProjectId(),
				req.GetOrderedBuildEvent().GetStreamId().GetInvocationId(),
				req.GetOrderedBuildEvent().GetStreamId().GetBuildId(),
				true, // isRealTime
				s.extractors,
				s.uuidGenerator,
			)
			if err != nil {
				return util.StatusWrap(err, "Failed to create BuildEventRecorder")
			}

			buildEventRecorder.StartLoggingConnectionMetadata(ctx)
		}

		err = s.handleBuildEvent(ctx, buildEventRecorder, req.OrderedBuildEvent.GetEvent(), req.OrderedBuildEvent.SequenceNumber)
		if err != nil {
			return util.StatusWrap(err, "Failed to handle build event")
		}

		err = ackBuildEventStreamMessage(ctx, stream, req.GetOrderedBuildEvent().GetStreamId(), req.GetOrderedBuildEvent().GetSequenceNumber())
		if err != nil {
			return util.StatusWrap(err, "Failed to ACK build event")
		}
	}
}

func (s BuildEventServer) handleBuildEvent(ctx context.Context, buildEventRecorder *buildeventrecorder.BuildEventRecorder, event *build.BuildEvent, sequenceNumber int64) error {
	if event.GetBazelEvent() == nil {
		return nil
	}
	var bazelEvent bes.BuildEvent
	err := event.GetBazelEvent().UnmarshalTo(&bazelEvent)
	if err != nil {
		return util.StatusWrap(err, "Failed to unmarshal BES event")
	}
	// TODO (isakstenstrom): Remove this and send the raw BES event instead. This can only be
	// done when we no longer need JSON serialization of events, like we do for
	// BazelInvocationProblems.
	buildEvent := events.NewBuildEvent(&bazelEvent, json.RawMessage(protojson.Format(&bazelEvent)))
	if err = buildEventRecorder.RecordEvent(ctx, &buildEvent, sequenceNumber); err != nil {
		return util.StatusWrap(err, "Failed to record build event")
	}
	return nil
}

func ackBuildEventStreamMessage(ctx context.Context, stream build.PublishBuildEvent_PublishBuildToolEventStreamServer, streamID *build.StreamId, sequenceNumber int64) error {
	err := stream.Send(&build.PublishBuildToolEventStreamResponse{
		StreamId:       streamID,
		SequenceNumber: sequenceNumber,
	})
	if err != nil {
		return util.StatusWrap(err, "Failed to send ACK for build event")
	}
	return nil
}
