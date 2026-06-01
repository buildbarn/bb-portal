package dbcleanupservice

import (
	"context"

	"github.com/buildbarn/bb-portal/ent/gen/ent/incompleteartifactgraph"
	"github.com/buildbarn/bb-portal/internal/database/sqlc"
	"github.com/buildbarn/bb-storage/pkg/util"
)

// DeleteIncompleteArtifactGraphs deletes staged artifact-graph events for
// invocations whose graph has already been compacted into a blob. Paced
// via nextSlice() like DeleteIncompleteLogs.
func (dc *DbCleanupService) DeleteIncompleteArtifactGraphs(ctx context.Context) (int64, error) {
	start, count, err := dc.nextSlice(ctx, incompleteartifactgraph.Table)
	if err != nil {
		return 0, err
	}

	deleted, err := dc.db.Sqlc().DeleteIncompleteArtifactGraphsFromPages(ctx, sqlc.DeleteIncompleteArtifactGraphsFromPagesParams{
		FromPage: start,
		Pages:    count,
	})
	if err != nil {
		return 0, util.StatusWrap(err, "Could not delete staged artifact graph events")
	}

	return deleted, nil
}
