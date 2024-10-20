package bes

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/buildbarn/bb-portal/pkg/events"
	"github.com/buildbarn/bb-portal/pkg/processing"
	"github.com/buildbarn/bb-portal/pkg/summary"
	"github.com/buildbarn/bb-portal/third_party/bazel/gen/bes"
	"google.golang.org/genproto/googleapis/devtools/build/v1"
)

// BuildEventChannel handles a single BuildEvent stream
type BuildEventChannel interface {
	// HandleBuildEvent processes a single BuildEvent
	// This method should be called for each received event.
	HandleBuildEvent(event *build.BuildEvent) error

	// Finalize does post-processing of a stream of BuildEvents.
	// This method should be called after receiving the EOF event.
	Finalize() error
}

type buildEventChannel struct {
	ctx        context.Context
	streamID   *build.StreamId
	summarizer *summary.Summarizer
	workflow   *processing.Workflow
}

// HandleBuildEvent implements BuildEventChannel.HandleBuildEvent.
func (c *buildEventChannel) HandleBuildEvent(event *build.BuildEvent) error {
	if event.GetBazelEvent() == nil {
		return nil
	}

	var bazelEvent bes.BuildEvent
	err := event.GetBazelEvent().UnmarshalTo(&bazelEvent)
	if err != nil {
		slog.ErrorContext(c.ctx, "UnmarshalTo failed", "err", err)
		return err
	}
	buildEvent := events.NewBuildEvent(&bazelEvent, json.RawMessage(protojson.Format(&bazelEvent)))
	if err = c.summarizer.ProcessEvent(&buildEvent); err != nil {
		slog.ErrorContext(c.ctx, "ProcessEvent failed", "err", err)
		return fmt.Errorf("could not process event (%s): , %w", buildEvent, err)
	}
	return nil
}

// Finalize implements BuildEventChannel.Finalize.
func (c *buildEventChannel) Finalize() error {
	summaryReport, err := c.summarizer.FinishProcessing()
	if err != nil {
		slog.ErrorContext(c.ctx, "FinishProcessing failed", "err", err)
		return err
	}

	// Hack for eventFile being required
	summaryReport.EventFileURL = fmt.Sprintf(
		"grpc://localhost:8082/google.devtools.build.v1/PublishLifecycleEvent?streamID=%s",
		c.streamID.String(),
	)

	slog.InfoContext(c.ctx, "Saving invocation", "id", c.streamID.String())
	startTime := time.Now()
	invocation, err := c.workflow.SaveSummary(c.ctx, summaryReport)
	if err != nil {
		slog.ErrorContext(c.ctx, "SaveSummary failed", "err", err)
		return err
	}
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	slog.InfoContext(c.ctx, fmt.Sprintf("Saved invocation in %v", elapsedTime.String()), "id", invocation.InvocationID)
	return nil
}

// noOpBuildEventChannel is an implementation of BuildEventChannel which does no processing of events.
// It is used when receiving a stream of events that we wish to ack without processing.
type noOpBuildEventChannel struct{}

// HandleBuildEvent implements BuildEventChannel.HandleBuildEvent.
func (c *noOpBuildEventChannel) HandleBuildEvent(event *build.BuildEvent) error {
	return nil
}

// Finalize implements BuildEventChannel.Finalize.
func (c *noOpBuildEventChannel) Finalize() error {
	return nil
}
