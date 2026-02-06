package dbcleanupservice

import (
	"context"

	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/internal/database/sqlc"
	"github.com/buildbarn/bb-storage/pkg/util"
	"go.opentelemetry.io/otel/attribute"
)

// RemoveOldInvocations removes invocations that have completed before a
// certain cutoff time.
func (dc *DbCleanupService) RemoveOldInvocations(ctx context.Context) error {
	ctx, span := dc.tracer.Start(ctx, "DbCleanupService.RemoveOldInvocations")
	defer span.End()

	cutoffTime := dc.clock.Now().UTC().Add(-dc.invocationRetention)
	start, count, err := dc.nextSlice(ctx, bazelinvocation.Table)
	if err != nil {
		return err
	}

	deleted, err := dc.db.Sqlc().DeleteOldInvocationsFromPages(ctx, sqlc.DeleteOldInvocationsFromPagesParams{
		FromPage:   start,
		Pages:      count,
		CutoffTime: cutoffTime,
	})
	if err != nil {
		return util.StatusWrap(err, "Failed to remove old invocations")
	}

	span.SetAttributes(attribute.Int64("deleted_invocations", deleted))
	return nil
}
