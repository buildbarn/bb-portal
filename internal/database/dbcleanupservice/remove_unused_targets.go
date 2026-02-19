package dbcleanupservice

import (
	"context"

	"github.com/buildbarn/bb-portal/ent/gen/ent/target"
	"github.com/buildbarn/bb-portal/internal/database/sqlc"
	"github.com/buildbarn/bb-storage/pkg/util"
	"go.opentelemetry.io/otel/attribute"
)

// RemoveUnusedTargets removes Targets that have no InvocationTargets or
// TargetKindMappings.
func (dc *DbCleanupService) RemoveUnusedTargets(ctx context.Context) error {
	ctx, span := dc.tracer.Start(ctx, "DbCleanupService.RemoveUnusedTargets")
	defer span.End()

	start, count, err := dc.nextSlice(ctx, target.Table)
	if err != nil {
		return err
	}

	deleted, err := dc.db.Sqlc().DeleteUnusedTargetsFromPages(ctx, sqlc.DeleteUnusedTargetsFromPagesParams{
		FromPage: start,
		Pages:    count,
	})
	if err != nil {
		return util.StatusWrap(err, "Failed to remove unused Targets")
	}

	span.SetAttributes(attribute.Int64("unused_targets_removed", deleted))
	return nil
}
