package processing

import (
	"context"
	"log/slog"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/pkg/archive"
)

// Workflow struct.
type Workflow struct {
	SummarizeActor
	SaveActor
}

// New Worflow constructor
func New(db *ent.Client, blobArchiver archive.BlobMultiArchiver) *Workflow {
	return &Workflow{
		SummarizeActor: SummarizeActor{},
		SaveActor: SaveActor{
			db:           db,
			blobArchiver: blobArchiver,
		},
	}
}

// ProcessFile function.
func (w Workflow) ProcessFile(ctx context.Context, file string) (*ent.BazelInvocation, error) {
	summary, err := w.Summarize(ctx, file)
	if err != nil {
		return nil, err
	}

	if !summary.BEPCompleted {
		slog.Info("File does not have a final event; will reprocess on next write")
		return nil, nil
	}
	return w.SaveSummary(ctx, summary)
}
