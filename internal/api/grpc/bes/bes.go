package bes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"time"

	build "google.golang.org/genproto/googleapis/devtools/build/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/pkg/archive"
	"github.com/buildbarn/bb-portal/pkg/events"
	"github.com/buildbarn/bb-portal/pkg/processing"
	"github.com/buildbarn/bb-portal/pkg/summary"
	"github.com/buildbarn/bb-portal/third_party/bazel/gen/bes"
)

// BES A type for the Build Event Service.
type BES struct {
	db           *ent.Client
	blobArchiver archive.BlobMultiArchiver
}

// New BES initializer function.
func New(db *ent.Client, blobArchiver archive.BlobMultiArchiver) build.PublishBuildEventServer {
	return &BES{
		db:           db,
		blobArchiver: blobArchiver,
	}
}

// PublishLifecycleEvent Publush a life cycle event.
func (b BES) PublishLifecycleEvent(ctx context.Context, request *build.PublishLifecycleEventRequest) (*emptypb.Empty, error) {
	slog.InfoContext(ctx, "Received event", "event", protojson.Format(request.BuildEvent.GetEvent()))
	return &emptypb.Empty{}, nil
}

// PublishBuildToolEventStream Public a build tool event stream.
func (b BES) PublishBuildToolEventStream(stream build.PublishBuildEvent_PublishBuildToolEventStreamServer) error {
	slog.InfoContext(stream.Context(), "Stream started", "event", stream.Context())

	summarizer := summary.NewSummarizer()

	ack := func(req *build.PublishBuildToolEventStreamRequest) {
		if err := stream.Send(&build.PublishBuildToolEventStreamResponse{
			StreamId:       req.OrderedBuildEvent.StreamId,
			SequenceNumber: req.OrderedBuildEvent.SequenceNumber,
		}); err != nil {
			slog.ErrorContext(stream.Context(), "Send failed", "err", err)
		}
	}

	var streamID *build.StreamId
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			slog.InfoContext(stream.Context(), "Stream finished", "event", stream.Context())
			break
		}
		if err != nil {
			slog.ErrorContext(stream.Context(), "Recv failed", "err", err)
			return err
		}

		if streamID == nil {
			streamID = req.GetOrderedBuildEvent().GetStreamId()
		}

		err = processBazelEvent(stream.Context(), req.OrderedBuildEvent.Event, summarizer)
		if err != nil {
			return err
		}

		ack(req)
	}

	summaryReport, err := summarizer.FinishProcessing()
	if err != nil {
		slog.ErrorContext(stream.Context(), "FinishProcessing failed", "err", err)
		return err
	}

	// Hack for eventFile being required
	summaryReport.EventFileURL = fmt.Sprintf(
		"grpc://localhost:8082/google.devtools.build.v1/PublishLifecycleEvent?streamID=%s",
		streamID.String(),
	)

	workflow := processing.New(b.db, b.blobArchiver)
	slog.InfoContext(stream.Context(), "Saving invocation", "id", streamID.String())
	startTime := time.Now()
	invocation, err := workflow.SaveSummary(stream.Context(), summaryReport)
	if err != nil {
		slog.ErrorContext(stream.Context(), "SaveSummary failed", "err", err)
		return err
	}
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	slog.InfoContext(stream.Context(), fmt.Sprintf("Saved invocation in %v", elapsedTime.String()), "id", invocation.InvocationID)
	return nil
}

// Process a bazel Event.
func processBazelEvent(ctx context.Context, event *build.BuildEvent, summarizer *summary.Summarizer) error {
	if event.GetBazelEvent() == nil {
		return nil
	}

	var bazelEvent bes.BuildEvent
	err := event.GetBazelEvent().UnmarshalTo(&bazelEvent)
	if err != nil {
		slog.ErrorContext(ctx, "UnmarshalTo failed", "err", err)
		return err
	}
	buildEvent := events.NewBuildEvent(&bazelEvent, json.RawMessage(protojson.Format(&bazelEvent)))
	if err = summarizer.ProcessEvent(&buildEvent); err != nil {
		slog.ErrorContext(ctx, "ProcessEvent failed", "err", err)
		return fmt.Errorf("could not process event (%s): , %w", buildEvent, err)
	}
	return nil
}
