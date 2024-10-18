package bes

import (
	"context"

	"github.com/buildbarn/bb-portal/pkg/processing"
	"github.com/buildbarn/bb-portal/pkg/summary"
	"google.golang.org/genproto/googleapis/devtools/build/v1"
)

// BuildEventHandler orchestrates the handling of incoming Build Event streams.
// For each incoming stream, and BuildEventChannel is created, which handles that stream.
// BuildEventHandler is responsible for managing the things that are common to these event streams.
type BuildEventHandler struct {
	workflow *processing.Workflow
}

// NewBuildEventHandler constructs a new BuildEventHandler
// TODO: Ensure we allow processing to complete before shutdown
// TODO: Cancel previous processing for an invocation if the client retries
// TODO: Write metrics
func NewBuildEventHandler(workflow *processing.Workflow) *BuildEventHandler {
	return &BuildEventHandler{
		workflow: workflow,
	}
}

// CreateEventChannel creates a new BuildEventChannel
func (h *BuildEventHandler) CreateEventChannel(ctx context.Context, streamID *build.StreamId) *BuildEventChannel {
	summarizer := summary.NewSummarizer()

	return &BuildEventChannel{
		ctx:        ctx,
		streamID:   streamID,
		summarizer: summarizer,
		workflow:   h.workflow,
	}
}
