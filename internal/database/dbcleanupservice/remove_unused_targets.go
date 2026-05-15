package dbcleanupservice

import (
	"context"

	"github.com/buildbarn/bb-portal/ent/gen/ent/target"
	"github.com/buildbarn/bb-portal/internal/database/sqlc"
	"github.com/buildbarn/bb-storage/pkg/util"
)

// RemoveUnusedTargets removes Targets that have no InvocationTargets or
// TargetKindMappings.
func (dc *DbCleanupService) RemoveUnusedTargets(ctx context.Context) (int64, error) {
	start, count, err := dc.nextSlice(ctx, target.Table)
	if err != nil {
		return 0, err
	}

	deleted, err := dc.db.Sqlc().DeleteUnusedTargetsFromPages(ctx, sqlc.DeleteUnusedTargetsFromPagesParams{
		FromPage: start,
		Pages:    count,
	})
	if err != nil {
		return 0, util.StatusWrap(err, "Failed to remove unused Targets")
	}

	return deleted, nil
}
