package dbcleanupservice

import (
	"context"

	"github.com/buildbarn/bb-portal/ent/gen/ent/testtarget"
	"github.com/buildbarn/bb-portal/internal/database/sqlc"
	"github.com/buildbarn/bb-storage/pkg/util"
)

// RemoveOrphanedTestTargets remove the test target marker from targest
// which no longer are referred to as tests by any invocation.
func (dc *DbCleanupService) RemoveOrphanedTestTargets(ctx context.Context) (int64, error) {
	start, count, err := dc.nextSlice(ctx, testtarget.Table)
	if err != nil {
		return 0, err
	}

	deleted, err := dc.batcher.Batch(ctx, func(ctx context.Context, limit int64) (int64, error) {
		return dc.db.Sqlc().DeleteOrphanedTestTargetsFromPages(ctx, sqlc.DeleteOrphanedTestTargetsFromPagesParams{
			FromPage:   start,
			Pages:      count,
			BatchLimit: limit,
		})
	})
	if err != nil {
		return 0, util.StatusWrap(err, "Failed to remove orphaned test targets")
	}

	return deleted, nil
}
