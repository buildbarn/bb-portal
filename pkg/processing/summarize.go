package processing

import (
	"context"

	"github.com/buildbarn/bb-portal/pkg/summary"
)

type SummarizeActor struct{}

func (SummarizeActor) Summarize(ctx context.Context, eventFileURL string) (*summary.Summary, error) {
	return summary.Summarize(ctx, eventFileURL)
}
