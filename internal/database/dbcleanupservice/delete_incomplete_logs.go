package dbcleanupservice

import (
	"context"

	"github.com/buildbarn/bb-portal/ent/gen/ent/incompletebuildlog"
	"github.com/buildbarn/bb-portal/internal/database/sqlc"
	"github.com/buildbarn/bb-storage/pkg/util"
)

// DeleteIncompleteLogs deletes logs which have had their incomplete
// build logs normalized.
func (dc *DbCleanupService) DeleteIncompleteLogs(ctx context.Context) (int64, error) {
	start, count, err := dc.nextSlice(ctx, incompletebuildlog.Table)
	if err != nil {
		return 0, err
	}

	deleted, err := dc.batcher.Batch(ctx, func(ctx context.Context, limit int64) (int64, error) {
		return dc.db.Sqlc().DeleteIncompleteLogsFromPages(ctx, sqlc.DeleteIncompleteLogsFromPagesParams{
			FromPage:   start,
			Pages:      count,
			BatchLimit: limit,
		})
	})
	if err != nil {
		return 0, util.StatusWrap(err, "Could not delete incompleted build logs")
	}

	return deleted, nil
}
