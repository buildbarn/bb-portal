package bes

import (
	"context"
	"fmt"
	"io"
	"math"
	"time"

	build "google.golang.org/genproto/googleapis/devtools/build/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-portal/internal/database/buildeventrecorder"
	"github.com/buildbarn/bb-portal/internal/database/dbauthservice"
	"github.com/buildbarn/bb-portal/pkg/authmetadataextraction"
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
type BuildEventServer struct {
	db                     database.Client
	instanceNameAuthorizer auth.Authorizer
	saveDataLevel          *bb_portal.BuildEventStreamService_SaveDataLevel
	tracerProvider         trace.TracerProvider
	extractors             *authmetadataextraction.AuthMetadataExtractors
	uuidGenerator          util.UUIDGenerator
}

// NewBuildEventServer creates a new BuildEventServer
func NewBuildEventServer(db database.Client, configuration *bb_portal.ApplicationConfiguration, dependenciesGroup program.Group, grpcClientFactory bb_grpc.ClientFactory, tracerProvider trace.TracerProvider, uuidGenerator util.UUIDGenerator) (*BuildEventServer, error) {
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

	saveDataLevel := besConfiguration.SaveDataLevel
	if saveDataLevel == nil || saveDataLevel.Level == nil {
		return nil, fmt.Errorf("No saveDataLevel configured")
	}

	extractors, err := authmetadataextraction.AuthMetadataExtractorsFromConfiguration(besConfiguration.AuthMetadataKeyConfiguration, dependenciesGroup)
	if err != nil {
		return nil, util.StatusWrap(err, "Failed to create AutheMetadataExtractors")
	}

	return &BuildEventServer{
		instanceNameAuthorizer: instanceNameAuthorizer,
		db:                     db,
		saveDataLevel:          saveDataLevel,
		tracerProvider:         tracerProvider,
		extractors:             extractors,
		uuidGenerator:          uuidGenerator,
	}, nil
}

// PublishLifecycleEvent handles life cycle events.
func (s *BuildEventServer) PublishLifecycleEvent(ctx context.Context, request *build.PublishLifecycleEventRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func requestToBuildEventWithInfo(req *build.PublishBuildToolEventStreamRequest) (*buildeventrecorder.BuildEventWithInfo, error) {
	var bazelEvent bes.BuildEvent
	obe := req.GetOrderedBuildEvent()
	be := obe.GetEvent()
	sid := obe.GetStreamId()
	sequenceNumber64 := obe.GetSequenceNumber()
	if obe == nil || be == nil || sid == nil {
		return nil, status.Error(codes.InvalidArgument, "Missing expected inputs.")
	}
	if sequenceNumber64 > math.MaxUint32 || sequenceNumber64 <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Sequence number out of range: %d", sequenceNumber64)
	}
	sequenceNumber := uint32(sequenceNumber64)

	if be.GetBazelEvent() == nil {
		return &buildeventrecorder.BuildEventWithInfo{
			SequenceNumber: sequenceNumber,
		}, nil
	}

	if err := be.GetBazelEvent().UnmarshalTo(&bazelEvent); err != nil {
		return nil, status.Error(codes.InvalidArgument, "Could not unmarshall bazel event")
	}

	// Add the event to the batch.
	return &buildeventrecorder.BuildEventWithInfo{
		Event:          &bazelEvent,
		SequenceNumber: sequenceNumber,
		AddedAt:        time.Now(),
	}, nil
}

// PublishBuildToolEventStream handles a build tool event stream.
func (s *BuildEventServer) PublishBuildToolEventStream(stream build.PublishBuildEvent_PublishBuildToolEventStreamServer) error {
	// We can safely bypass authorization checks here, as we check that the
	// user is allowed to upload to the instance name when creating the
	// BuildEventRecorder.
	ctx := dbauthservice.NewContextWithDbAuthServiceBypass(stream.Context())

	// Synchronously block for the first request to simplify
	// initialization logic.
	req, err := stream.Recv()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return util.StatusWrap(err, "Failed to recieve the initial event from the stream")
	}

	data, err := requestToBuildEventWithInfo(req)
	if err != nil {
		return util.StatusWrap(err, "Failed to convert the initial request to internal build event")
	}
	batcher := buildeventrecorder.NewBatcher[buildeventrecorder.BuildEventWithInfo]()
	if data != nil {
		// If the first message is a build event add it to the batch.
		batcher.Add(*data)
	}
	streamID := req.GetOrderedBuildEvent().GetStreamId()

	// initialize with the first message
	buildEventRecorder, err := buildeventrecorder.NewBuildEventRecorder(
		ctx,
		s.db,
		s.instanceNameAuthorizer,
		s.saveDataLevel,
		s.tracerProvider,
		req.GetProjectId(),
		streamID.GetInvocationId(),
		streamID.GetBuildId(),
		true, // isRealTime
		s.extractors,
		s.uuidGenerator,
	)
	if err != nil {
		return util.StatusWrap(err, "Could not initialize build event recorder")
	}
	buildEventRecorder.StartLoggingConnectionMetadata(ctx)

	errChan := make(chan error, 1)
	go func() {
		defer close(errChan)
		for {
			// This call blocks, but will be cancelled when we return
			// the parent function.
			req, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				errChan <- err
				return
			}
			data, err := requestToBuildEventWithInfo(req)
			if err != nil {
				errChan <- err
				return
			}
			if data != nil {
				batcher.Add(*data)
			}
		}
	}()

	processBatch := func(batch []buildeventrecorder.BuildEventWithInfo) error {
		if err := buildEventRecorder.SaveBatch(ctx, batch); err != nil {
			return util.StatusWrap(err, "Failed to commit build event batch")
		}
		for _, meta := range batch {
			err := stream.Send(&build.PublishBuildToolEventStreamResponse{
				StreamId:       streamID,
				SequenceNumber: int64(meta.SequenceNumber),
			})
			if err != nil {
				return util.StatusWrap(err, "Failed to send ACK for build event")
			}
		}
		return nil
	}

	prevBatch := make([]buildeventrecorder.BuildEventWithInfo, 0)
	for {
		// Block until receive an error or we've received a batch of
		// requests.
		select {
		case err, more := <-errChan:
			if err != nil {
				return util.StatusWrap(err, "Failed to receive build event")
			}
			// Input stream is closed without error. Do last batch.
			if !more {
				nextBatch := batcher.Swap(prevBatch)
				return processBatch(nextBatch)
			}
		case <-batcher.Ready():
			nextBatch := batcher.Swap(prevBatch)
			if err = processBatch(nextBatch); err != nil {
				return util.StatusWrap(err, "Failed to process build event batch")
			}
			prevBatch = nextBatch[:0]
		}
	}
}
