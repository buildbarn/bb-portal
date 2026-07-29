package dbcleanupservice

import (
	"context"

	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/internal/database/sqlc"
	"github.com/buildbarn/bb-storage/pkg/util"
)

// LockInvocationsWithNoRecentEvents locks invocations that have not received any new
// events in a certain period of time.
func (dc *DbCleanupService) LockInvocationsWithNoRecentEvents(ctx context.Context) (int64, error) {
	cutoffTime := dc.clock.Now().UTC().Add(-dc.invocationMessageTimeout)
	start, count, err := dc.nextSlice(ctx, bazelinvocation.Table)
	if err != nil || count == 0 {
		return 0, err
	}

	updatedRows, err := dc.batcher.Batch(ctx, func(ctx context.Context, limit int64) (int64, error) {
		return dc.db.Sqlc().LockStaleInvocationsFromPages(
			ctx,
			sqlc.LockStaleInvocationsFromPagesParams{
				FromPage:   start,
				Pages:      count,
				CutoffTime: cutoffTime,
				BatchLimit: limit,
			},
		)
	})
	if err != nil {
		return 0, util.StatusWrap(err, "Failed to lock invocations")
	}

	return updatedRows, nil
}

// UpdateInvocationEndedAtFromEvents updates the ended_at timestamp of locked
// invocations based on the latest event metadata.
func (dc *DbCleanupService) UpdateInvocationEndedAtFromEvents(ctx context.Context) (int64, error) {
	updatedRows, err := dc.db.Sqlc().UpdateCompletedInvocationWithEndTimeFromEventMetadata(ctx)
	if err != nil {
		return 0, util.StatusWrap(err, "Failed to update invocation ended_at from event metadata")
	}

	return updatedRows, nil
}
