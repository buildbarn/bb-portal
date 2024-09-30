package processing

import (
	"context"

	"github.com/buildbarn/bb-portal/pkg/summary"
)

// SummarizeActor struct.
type SummarizeActor struct{}

// Summarize function.
func (SummarizeActor) Summarize(ctx context.Context, eventFileURL string) (*summary.Summary, error) {
	return summary.Summarize(ctx, eventFileURL)
}
