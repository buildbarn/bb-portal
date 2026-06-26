package dbcleanupservice

import (
	"context"

	"github.com/buildbarn/bb-portal/ent/gen/ent/targetkindmapping"
	"github.com/buildbarn/bb-portal/internal/database/sqlc"
	"github.com/buildbarn/bb-storage/pkg/util"
)

// RemoveTargetKindMappings removes TargetKindMappings for BazelInvocations
// that are completed.
func (dc *DbCleanupService) RemoveTargetKindMappings(ctx context.Context) (int64, error) {
	start, count, err := dc.nextSlice(ctx, targetkindmapping.Table)
	if err != nil {
		return 0, err
	}

	deleted, err := dc.batcher.Batch(ctx, func(ctx context.Context, limit int64) (int64, error) {
		return dc.db.Sqlc().DeleteTargetKindMappingsFromPages(ctx, sqlc.DeleteTargetKindMappingsFromPagesParams{
			FromPage:   start,
			Pages:      count,
			BatchLimit: limit,
		})
	})
	if err != nil {
		return 0, util.StatusWrap(err, "Failed to remove old TargetKindMappings")
	}

	return deleted, nil
}
