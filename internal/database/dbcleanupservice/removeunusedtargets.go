package dbcleanupservice

import (
	"context"

	"github.com/buildbarn/bb-portal/ent/gen/ent/target"
	"github.com/buildbarn/bb-storage/pkg/util"
	"go.opentelemetry.io/otel/attribute"
)

// RemoveUnusedTargets removes Targets that have no InvocationTargets or
// TargetKindMappings.
func (dc *DbCleanupService) RemoveUnusedTargets(ctx context.Context) error {
	ctx, span := dc.tracer.Start(ctx, "DbCleanupService.RemoveUnusedTargets")
	defer span.End()

	deletedRows, err := dc.db.Ent().Target.Delete().
		Where(
			target.Not(target.HasInvocationTargets()),
			target.Not(target.HasTargetKindMappings()),
		).
		Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to remove unused Targets")
	}

	span.SetAttributes(attribute.KeyValue{Key: "unused_targets_removed", Value: attribute.IntValue(deletedRows)})

	return nil
}
