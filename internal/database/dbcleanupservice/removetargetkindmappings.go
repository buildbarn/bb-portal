package dbcleanupservice

import (
	"context"

	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/ent/gen/ent/targetkindmapping"
	"github.com/buildbarn/bb-storage/pkg/util"
	"go.opentelemetry.io/otel/attribute"
)

// RemoveTargetKindMappings removes TargetKindMappings for BazelInvocations
// that are completed.
func (dc *DbCleanupService) RemoveTargetKindMappings(ctx context.Context) error {
	ctx, span := dc.tracer.Start(ctx, "DbCleanupService.RemoveTargetKindMappings")
	defer span.End()

	deletedRows, err := dc.db.Ent().TargetKindMapping.Delete().
		Where(
			targetkindmapping.HasBazelInvocationWith(
				bazelinvocation.BepCompletedEQ(true),
			),
		).
		Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to remove old TargetKindMappings")
	}

	span.SetAttributes(attribute.KeyValue{Key: "target_kind_mappings_removed", Value: attribute.IntValue(deletedRows)})

	return nil
}
