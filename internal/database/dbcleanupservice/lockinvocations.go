package dbcleanupservice

import (
	"context"

	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/ent/gen/ent/eventmetadata"
	"github.com/buildbarn/bb-storage/pkg/util"
	"go.opentelemetry.io/otel/attribute"
)

// LockInvocationsWithNoRecentEvents locks invocations that have not received any new
// events in a certain period of time.
func (dc *DbCleanupService) LockInvocationsWithNoRecentEvents(ctx context.Context) error {
	ctx, span := dc.tracer.Start(ctx, "DbCleanupService.LockInvocationsWithNoRecentEvents")
	defer span.End()

	cutoffTime := dc.clock.Now().UTC().Add(-dc.invocationMessageTimeout)

	updatedRows, err := dc.db.Ent().BazelInvocation.Update().
		Where(
			bazelinvocation.BepCompleted(false),
			bazelinvocation.HasEventMetadataWith(
				eventmetadata.EventReceivedAtLTE(cutoffTime),
			),
		).
		SetBepCompleted(true).
		Save(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to lock invocations")
	}

	span.SetAttributes(attribute.KeyValue{Key: "invocations_locked", Value: attribute.IntValue(updatedRows)})

	return nil
}

// UpdateInvocationEndedAtFromEvents updates the ended_at timestamp of locked
// invocations based on the latest event metadata.
func (dc *DbCleanupService) UpdateInvocationEndedAtFromEvents(ctx context.Context) error {
	ctx, span := dc.tracer.Start(ctx, "DbCleanupService.UpdateInvocationEndedAtFromEvents")
	defer span.End()

	updatedRows, err := dc.db.Sqlc().UpdateCompletedInvocationWithEndTimeFromEventMetadata(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to update invocation ended_at from event metadata")
	}

	span.SetAttributes(attribute.KeyValue{Key: "invocations_updated", Value: attribute.Int64Value(updatedRows)})

	return nil
}
