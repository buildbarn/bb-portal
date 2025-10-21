package dbcleanupservice

import (
	"context"
	"log/slog"

	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/ent/gen/ent/targetkindmapping"
	"github.com/buildbarn/bb-storage/pkg/util"
)

// RemoveTargetKindMappings removes TargetKindMappings for BazelInvocations
// that are completed.
func (dc *DbCleanupService) RemoveTargetKindMappings(ctx context.Context) error {
	ctx, span := dc.tracer.Start(ctx, "DbCleanupService.RemoveTargetKindMappings")
	defer span.End()

	deletedRows, err := dc.db.TargetKindMapping.Delete().
		Where(
			targetkindmapping.HasBazelInvocationWith(
				bazelinvocation.BepCompletedEQ(true),
			),
		).
		Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to remove old TargetKindMappings")
	}

	slog.Info("Removed old TargetKindMappings", "count", deletedRows)
	return nil
}
