package dbcleanupservice

import (
	"context"

	"github.com/buildbarn/bb-portal/ent/gen/ent/completedaction"
	"github.com/buildbarn/bb-portal/internal/database/sqlc"
	prometheusmetrics "github.com/buildbarn/bb-portal/pkg/prometheus_metrics"
	"github.com/buildbarn/bb-storage/pkg/util"
)

// RemoveOldCompletedActions removes actions that have completed before
// a certain cutoff time
func (dc *DbCleanupService) RemoveOldCompletedActions(ctx context.Context) (int64, error) {
	cutoffTime := dc.clock.Now().UTC().Add(-dc.completedActionRetention)
	start, count, err := dc.nextSlice(ctx, completedaction.Table)
	if err != nil {
		return 0, err
	}

	deleted, err := dc.batcher.Batch(ctx, func(ctx context.Context, limit int64) (int64, error) {
		return dc.db.Sqlc().DeleteOldCompletedActions(ctx, sqlc.DeleteOldCompletedActionsParams{
			FromPage:   start,
			Pages:      count,
			CutoffTime: cutoffTime,
			BatchLimit: limit,
		})
	})
	if err != nil {
		return 0, util.StatusWrap(err, "Failed to remove old actions")
	}
	prometheusmetrics.SyncCompletedActions(ctx, dc.db.Ent())
	return deleted, nil
}
