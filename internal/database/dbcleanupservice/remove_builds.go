package dbcleanupservice

import (
	"context"

	"github.com/buildbarn/bb-portal/ent/gen/ent/build"
	"github.com/buildbarn/bb-storage/pkg/util"
)

// RemoveBuildsWithoutInvocations removes builds that do not have any
// associated invocations.
func (dc *DbCleanupService) RemoveBuildsWithoutInvocations(ctx context.Context) (int64, error) {
	deletedBuilds, err := dc.db.Ent().Build.Delete().
		Where(
			build.Not(build.HasInvocations()),
		).
		Exec(ctx)
	if err != nil {
		return 0, util.StatusWrap(err, "Failed to remove builds without invocations")
	}

	return int64(deletedBuilds), nil
}
