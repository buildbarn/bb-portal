package dbcleanupservice

import (
	"context"

	"github.com/buildbarn/bb-portal/ent/gen/ent/testtarget"
	"github.com/buildbarn/bb-portal/internal/database/sqlc"
	"github.com/buildbarn/bb-storage/pkg/util"
	"go.opentelemetry.io/otel/attribute"
)

// RemoveOrphanedTestTargets remove the test target marker from targest
// which no longer are referred to as tests by any invocation.
func (dc *DbCleanupService) RemoveOrphanedTestTargets(ctx context.Context) error {
	ctx, span := dc.tracer.Start(ctx, "DbCleanupService.RemoveOrphanedTestTargets")
	defer span.End()

	start, count, err := dc.nextSlice(ctx, testtarget.Table)
	if err != nil {
		return err
	}

	deleted, err := dc.db.Sqlc().DeleteOrphanedTestTargetsFromPages(ctx, sqlc.DeleteOrphanedTestTargetsFromPagesParams{
		FromPage: start,
		Pages:    count,
	})
	if err != nil {
		return util.StatusWrap(err, "Failed to remove orphaned test targets")
	}

	span.SetAttributes(attribute.Int64("deleted_test_targets", deleted))
	return nil
}
