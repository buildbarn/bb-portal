package bes

import (
	"context"
	"io"
	"log/slog"

	build "google.golang.org/genproto/googleapis/devtools/build/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/pkg/processing"
)

// BuildEventServer implements the Build Event Service.
// It receives events and forwards them to a BuildEventChannel.
// TODO: Should this support forwarding events? Users might want to create their own
// tooling that reacts to build events, and it would be useful if this service could
// forward events to those.
type BuildEventServer struct {
	handler *BuildEventHandler
}

// NewBuildEventServer creates a new BuildEventServer
func NewBuildEventServer(db *ent.Client, blobArchiver processing.BlobMultiArchiver) build.PublishBuildEventServer {
	return &BuildEventServer{
		handler: NewBuildEventHandler(processing.New(db, blobArchiver)),
	}
}

// PublishLifecycleEvent handles life cycle events.
func (s BuildEventServer) PublishLifecycleEvent(ctx context.Context, request *build.PublishLifecycleEventRequest) (*emptypb.Empty, error) {
	slog.InfoContext(ctx, "Received event", "event", protojson.Format(request.BuildEvent.GetEvent()))
	return &emptypb.Empty{}, nil
}

// PublishBuildToolEventStream handles a build tool event stream.
func (s BuildEventServer) PublishBuildToolEventStream(stream build.PublishBuildEvent_PublishBuildToolEventStreamServer) error {
	slog.InfoContext(stream.Context(), "Stream started", "event", stream.Context())

	ack := func(req *build.PublishBuildToolEventStreamRequest) {
		if err := stream.Send(&build.PublishBuildToolEventStreamResponse{
			StreamId:       req.OrderedBuildEvent.StreamId,
			SequenceNumber: req.OrderedBuildEvent.SequenceNumber,
		}); err != nil {
			slog.ErrorContext(stream.Context(), "Send failed", "err", err)
		}
	}

	var streamID *build.StreamId
	reqCh := make(chan *build.PublishBuildToolEventStreamRequest)
	errCh := make(chan error)
	var eventCh *BuildEventChannel

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
					return nil
				}

				return eventCh.Finalize()
			}

			slog.ErrorContext(stream.Context(), "Recv failed", "err", err)
			return err

		case req := <-reqCh:
			// First request
			if streamID == nil {
				streamID = req.OrderedBuildEvent.GetStreamId()
				eventCh = s.handler.CreateEventChannel(stream.Context(), streamID)
			}

			if err := eventCh.HandleBuildEvent(req.OrderedBuildEvent.Event); err != nil {
				slog.ErrorContext(stream.Context(), "HandleBuildEvent failed", "err", err)
				return err
			}

			ack(req)
		}
	}
}
