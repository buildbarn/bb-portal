package dbcleanupservice

import (
	"context"

	"github.com/buildbarn/bb-portal/ent/gen/ent/invocationartifactgraph"
	"github.com/buildbarn/bb-portal/internal/database/sqlc"
	"github.com/buildbarn/bb-storage/pkg/util"
)

// RemoveOldArtifacts deletes invocation_artifact_graphs rows whose
// owning invocation completed more than artifact_retention ago. Aggregate
// ArtifactMetrics on the invocation row stay intact. Paced via
// nextSlice() so a single tick processes one slice of the table; the
// cleanup loop walks the full table roughly once per hour.
func (dc *DbCleanupService) RemoveOldArtifacts(ctx context.Context) (int64, error) {
	cutoff := dc.clock.Now().UTC().Add(-dc.artifactRetention)
	start, count, err := dc.nextSlice(ctx, invocationartifactgraph.Table)
	if err != nil {
		return 0, err
	}

	deleted, err := dc.db.Sqlc().DeleteOldInvocationArtifactGraphsFromPages(ctx, sqlc.DeleteOldInvocationArtifactGraphsFromPagesParams{
		FromPage:   start,
		Pages:      count,
		CutoffTime: cutoff,
	})
	if err != nil {
		return 0, util.StatusWrap(err, "Failed to remove old invocation artifact graphs")
	}

	return deleted, nil
}
