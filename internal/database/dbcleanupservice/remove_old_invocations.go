package dbcleanupservice

import (
	"context"

	"github.com/buildbarn/bb-portal/internal/database/sqlc"
	"github.com/buildbarn/bb-storage/pkg/util"
	"go.opentelemetry.io/otel/attribute"
)

var deleteOldInvocationsBatchSize int64 = 1000

// RemoveOldInvocations removes invocations that have completed before a
// certain cutoff time.
func (dc *DbCleanupService) RemoveOldInvocations(ctx context.Context) error {
	ctx, span := dc.tracer.Start(ctx, "DbCleanupService.RemoveOldInvocations")
	defer span.End()

	cutoffTime := dc.clock.Now().UTC().Add(-dc.invocationRetention)
	var totalDeleted int64
	for {
		deleted, err := dc.db.Sqlc().DeleteOldInvocations(ctx, sqlc.DeleteOldInvocationsParams{
			BatchSize:  deleteOldInvocationsBatchSize,
			CutoffTime: cutoffTime,
		})
		if err != nil {
			return util.StatusWrap(err, "Failed to remove old invocations")
		}
		totalDeleted += deleted
		if deleted < deleteOldInvocationsBatchSize {
			break
		}
	}

	span.SetAttributes(attribute.Int64("deleted_invocations", totalDeleted))
	return nil
}
