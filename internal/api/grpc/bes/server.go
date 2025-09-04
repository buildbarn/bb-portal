package bes

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sort"

	build "google.golang.org/genproto/googleapis/devtools/build/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/internal/api/common"
	"github.com/buildbarn/bb-portal/pkg/processing"
	"github.com/buildbarn/bb-storage/pkg/auth"
	auth_configuration "github.com/buildbarn/bb-storage/pkg/auth/configuration"
	bb_grpc "github.com/buildbarn/bb-storage/pkg/grpc"
	auth_pb "github.com/buildbarn/bb-storage/pkg/proto/configuration/auth"
	"github.com/buildbarn/bb-storage/pkg/util"
)

// BuildEventServer implements the Build Event Service.
// It receives events and forwards them to a BuildEventChannel.
// TODO: Should this support forwarding events? Users might want to create their own
// tooling that reacts to build events, and it would be useful if this service could
// forward events to those.
type BuildEventServer struct {
	handler                *BuildEventHandler
	instanceNameAuthorizer auth.Authorizer
}

// NewBuildEventServer creates a new BuildEventServer
func NewBuildEventServer(db *ent.Client, blobArchiver processing.BlobMultiArchiver, authorizerConfiguration *auth_pb.AuthorizerConfiguration, grpcClientFactory bb_grpc.ClientFactory) (build.PublishBuildEventServer, error) {
	if authorizerConfiguration == nil {
		return nil, status.Error(codes.NotFound, "No InstanceNameAuthorizer configured")
	}
	instanceNameAuthorizer, err := auth_configuration.DefaultAuthorizerFactory.NewAuthorizerFromConfiguration(authorizerConfiguration, grpcClientFactory)
	if err != nil {
		return nil, util.StatusWrap(err, "Failed to create InstanceNameAuthorizer")
	}

	return &BuildEventServer{
		handler:                NewBuildEventHandler(processing.New(db, blobArchiver)),
		instanceNameAuthorizer: instanceNameAuthorizer,
	}, nil
}

// PublishLifecycleEvent handles life cycle events.
func (s BuildEventServer) PublishLifecycleEvent(ctx context.Context, request *build.PublishLifecycleEventRequest) (*emptypb.Empty, error) {
	slog.InfoContext(ctx, "Received event", "event", protojson.Format(request.BuildEvent.GetEvent()))
	return &emptypb.Empty{}, nil
}

// PublishBuildToolEventStream handles a build tool event stream.
func (s BuildEventServer) PublishBuildToolEventStream(stream build.PublishBuildEvent_PublishBuildToolEventStreamServer) error {
	slog.InfoContext(stream.Context(), "Stream started", "event", stream.Context())

	// List of SequenceIds we've received.
	// We'll want to ack these once all events are received, as we don't support resumption.
	seqNrs := make([]int64, 0)

	ack := func(streamID *build.StreamId, sequenceNumber int64, isClosing bool) {
		if err := stream.Send(&build.PublishBuildToolEventStreamResponse{
			StreamId:       streamID,
			SequenceNumber: sequenceNumber,
		}); err != nil {

			// with the option --bes_upload_mode=fully_async or nowait_for_upload_complete
			// its not an error when the send fails. the bes gracefully terminated the close
			// i.e. sent an EOF. for long running builds that take a while to save to the db (> 1s)
			// the context is processed in the background, so by the time we are acknowledging these
			// requests, the client connection may have already timed out and these errors can be
			// safely ignored
			grpcErr := status.Convert(err)
			if isClosing &&
				grpcErr.Code() == codes.Unavailable &&
				grpcErr.Message() == "transport is closing" {
				return
			}

			slog.ErrorContext(
				stream.Context(),
				"Send failed",
				"err",
				err,
				"streamid",
				streamID,
				"sequenceNumber",
				sequenceNumber,
			)
		}
	}

	var streamID *build.StreamId
	reqCh := make(chan *build.PublishBuildToolEventStreamRequest)
	errCh := make(chan error)
	var eventCh BuildEventChannel

	go func() {
		for {
			req, err := stream.Recv()
			if err != nil {
				errCh <- err
				return
			}
			reqCh <- req
		}
	}()

	for {
		select {
		case err := <-errCh:
			if err == io.EOF {
				slog.InfoContext(stream.Context(), "Stream finished", "event", stream.Context())

				if eventCh == nil {
					slog.WarnContext(stream.Context(), "No event channel found for stream event", "event", stream.Context())
					return nil
				}

				// Validate that all events were received
				sort.Slice(seqNrs, func(i, j int) bool { return seqNrs[i] < seqNrs[j] })

				// TODO: Find out if initial sequence number can be != 1
				expected := int64(1)
				for _, seqNr := range seqNrs {
					if seqNr != expected {
						return status.Error(codes.Unknown, fmt.Sprintf("received unexpected sequence number %d, expected %d", seqNr, expected))
					}
					expected++
				}

				err := eventCh.Finalize()
				if err != nil {
					return err
				}

				// Ack all events
				for _, seqNr := range seqNrs {
					ack(streamID, seqNr, true)
				}

				return nil
			}

			slog.ErrorContext(stream.Context(), "Recv failed", "err", err)
			return err

		case req := <-reqCh:
			if !common.IsInstanceNameAllowed(stream.Context(), s.instanceNameAuthorizer, req.ProjectId) {
				return status.Error(codes.PermissionDenied, fmt.Sprintf("Instance name %q is not allowed", req.ProjectId))
			}

			// First event
			if streamID == nil {
				streamID = req.OrderedBuildEvent.GetStreamId()
				eventCh = s.handler.CreateEventChannel(stream.Context(), req.OrderedBuildEvent)
			}

			seqNrs = append(seqNrs, req.OrderedBuildEvent.GetSequenceNumber())

			if err := eventCh.HandleBuildEvent(req.OrderedBuildEvent.Event, req.ProjectId); err != nil {
				slog.ErrorContext(stream.Context(), "HandleBuildEvent failed", "err", err)
				return err
			}
		}
	}
}
