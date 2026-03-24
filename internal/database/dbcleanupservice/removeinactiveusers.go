package dbcleanupservice

import (
	"context"

	"github.com/buildbarn/bb-portal/ent/gen/ent/authenticateduser"
	prometheusmetrics "github.com/buildbarn/bb-portal/pkg/prometheus_metrics"
	"github.com/buildbarn/bb-storage/pkg/util"
	"go.opentelemetry.io/otel/attribute"
)

// RemoveInactiveUsers removes all users with no associated invocations.
func (dc *DbCleanupService) RemoveInactiveUsers(ctx context.Context) error {
	ctx, span := dc.tracer.Start(ctx, "DbCleanupService.RemoveInactiveUsers")
	defer span.End()

	deletedUsers, err := dc.db.Ent().AuthenticatedUser.Delete().
		Where(authenticateduser.Not(authenticateduser.HasBazelInvocations())).
		Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to remove inactive users")
	}

	span.SetAttributes(attribute.KeyValue{Key: "deleted_users", Value: attribute.IntValue(deletedUsers)})
	prometheusmetrics.SyncAuthenticatedUsersCount(ctx, dc.db.Ent())

	return nil
}
