package dbcleanupservice

import (
	"context"

	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/internal/database/sqlc"
	prometheusmetrics "github.com/buildbarn/bb-portal/pkg/prometheus_metrics"
	"github.com/buildbarn/bb-storage/pkg/util"
)

// RemoveOldInvocations removes invocations that have completed before a
// certain cutoff time.
func (dc *DbCleanupService) RemoveOldInvocations(ctx context.Context) (int64, error) {
	cutoffTime := dc.clock.Now().UTC().Add(-dc.invocationRetention)
	start, count, err := dc.nextSlice(ctx, bazelinvocation.Table)
	if err != nil {
		return 0, err
	}

	deleted, err := dc.batcher.Batch(ctx, func(ctx context.Context, limit int64) (int64, error) {
		return dc.db.Sqlc().DeleteOldInvocationsFromPages(ctx, sqlc.DeleteOldInvocationsFromPagesParams{
			FromPage:   start,
			Pages:      count,
			CutoffTime: cutoffTime,
			BatchLimit: limit,
		})
	})
	if err != nil {
		return 0, util.StatusWrap(err, "Failed to remove old invocations")
	}

	prometheusmetrics.SyncInvocations(ctx, dc.db.Ent())
	return deleted, nil
}
