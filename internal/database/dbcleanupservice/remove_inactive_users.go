package dbcleanupservice

import (
	"context"

	"github.com/buildbarn/bb-portal/ent/gen/ent/authenticateduser"
	prometheusmetrics "github.com/buildbarn/bb-portal/pkg/prometheus_metrics"
	"github.com/buildbarn/bb-storage/pkg/util"
)

// RemoveInactiveUsers removes all users with no associated invocations.
func (dc *DbCleanupService) RemoveInactiveUsers(ctx context.Context) (int64, error) {
	deletedUsers, err := dc.db.Ent().AuthenticatedUser.Delete().
		Where(authenticateduser.Not(authenticateduser.HasBazelInvocations())).
		Exec(ctx)
	if err != nil {
		return 0, util.StatusWrap(err, "Failed to remove inactive users")
	}

	prometheusmetrics.SyncAuthenticatedUsersCount(ctx, dc.db.Ent())

	return int64(deletedUsers), nil
}
