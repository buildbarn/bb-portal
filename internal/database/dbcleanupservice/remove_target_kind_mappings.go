package dbcleanupservice

import (
	"context"

	"github.com/buildbarn/bb-portal/ent/gen/ent/targetkindmapping"
	"github.com/buildbarn/bb-portal/internal/database/sqlc"
	"github.com/buildbarn/bb-storage/pkg/util"
	"go.opentelemetry.io/otel/attribute"
)

// RemoveTargetKindMappings removes TargetKindMappings for BazelInvocations
// that are completed.
func (dc *DbCleanupService) RemoveTargetKindMappings(ctx context.Context) error {
	ctx, span := dc.tracer.Start(ctx, "DbCleanupService.RemoveTargetKindMappings")
	defer span.End()

	start, count, err := dc.nextSlice(ctx, targetkindmapping.Table)
	if err != nil {
		return err
	}

	deleted, err := dc.db.Sqlc().DeleteTargetKindMappingsFromPages(ctx, sqlc.DeleteTargetKindMappingsFromPagesParams{
		FromPage: start,
		Pages:    count,
	})
	if err != nil {
		return util.StatusWrap(err, "Failed to remove old TargetKindMappings")
	}

	span.SetAttributes(attribute.Int64("target_kind_mappings_removed", deleted))

	return nil
}
