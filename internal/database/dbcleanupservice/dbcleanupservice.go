package dbcleanupservice

import (
	"context"
	"log/slog"
	"math/rand/v2"
	"time"

	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/ent/gen/ent/build"
	"github.com/buildbarn/bb-portal/ent/gen/ent/incompletebuildlog"
	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-portal/internal/database/dbauthservice"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-storage/pkg/clock"
	"github.com/buildbarn/bb-storage/pkg/program"
	"github.com/buildbarn/bb-storage/pkg/util"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// DbCleanupService a service that performs periodic cleanup of the
// database to remove old data that is no longer needed. This includes:
//
//  1. Locking unfinished invocations that have not received any new event
//     metadata for a certain period of time (invocationMessageTimeout) and
//     setting their end time if it is not set.
//  2. Compacting logs by normalizing and compressing them.
//  3. Deleting incomplete logs that have been compacted.
//  4. Removing invocations that have completed and whose completion time
//     is older than invocationRetention.
//  5. Removing builds that do not have any associated invocations.
//  6. Removing old TargetKindMappings.
//  7. Removing unused targets.
type DbCleanupService struct {
	db                       database.Client
	clock                    clock.Clock
	cleanupInterval          time.Duration
	invocationMessageTimeout time.Duration
	invocationRetention      time.Duration
	tracer                   trace.Tracer
}

// NewDbCleanupService creates a new DbCleanupService.
func NewDbCleanupService(
	db database.Client,
	clock clock.Clock,
	cleanupConfiguration *bb_portal.BuildEventStreamService_DatabaseCleanupConfiguration,
	tracerProvider trace.TracerProvider,
) (*DbCleanupService, error) {
	cleanupInterval := cleanupConfiguration.CleanupInterval
	if err := cleanupInterval.CheckValid(); err != nil {
		return nil, util.StatusWrap(err, "Failed to parse cleanupInterval parameter time")
	}

	invocationMessageTimeout := cleanupConfiguration.InvocationMessageTimeout
	if err := invocationMessageTimeout.CheckValid(); err != nil {
		return nil, util.StatusWrap(err, "Failed to parse invocationMessageTimeout parameter time")
	}

	invocationRetention := cleanupConfiguration.InvocationRetention
	if err := invocationRetention.CheckValid(); err != nil {
		return nil, util.StatusWrap(err, "Failed to parse invocationRetention parameter time")
	}

	return &DbCleanupService{
		db:                       db,
		clock:                    clock,
		cleanupInterval:          cleanupInterval.AsDuration(),
		invocationMessageTimeout: invocationMessageTimeout.AsDuration(),
		invocationRetention:      invocationRetention.AsDuration(),
		tracer:                   tracerProvider.Tracer("github.com/buildbarn/bb-portal/internal/database/dbcleanupservice"),
	}, nil
}

// StartDbCleanupService starts a goroutine that performs periodic
// cleanup of the database.
func (dc *DbCleanupService) StartDbCleanupService(ctx context.Context, group program.Group) {
	group.Go(func(ctx context.Context, siblingsGroup, dependenciesGroup program.Group) error {
		ctx = dbauthservice.NewContextWithDbAuthServiceBypass(ctx)
		for {
			// Add 5% jitter to the cleanup interval
			timeToSleep := dc.cleanupInterval + time.Duration((rand.Float64()*0.1-0.05)*float64(dc.cleanupInterval))
			time.Sleep(timeToSleep)
			select {
			case <-ctx.Done():
				return nil
			default:
				if err := dc.LockInvocationsWithNoRecentEvents(ctx); err != nil {
					slog.Warn("Failed to lock unfinished invocations with no recent events", "err", err)
				}
				if err := dc.UpdateInvocationEndedAtFromEvents(ctx); err != nil {
					slog.Warn("Failed to update invocation ended_at from event metadata", "err", err)
				}
				if err := dc.CompactLogs(ctx); err != nil {
					slog.Warn("Failed to compact logs", "err", err)
				}
				if err := dc.DeleteIncompleteLogs(ctx); err != nil {
					slog.Warn("Failed to delete incomplete logs", "err", err)
				}
				if err := dc.RemoveOldInvocations(ctx); err != nil {
					slog.Warn("Failed to remove old invocations", "err", err)
				}
				if err := dc.RemoveBuildsWithoutInvocations(ctx); err != nil {
					slog.Warn("Failed to remove builds without invocations", "err", err)
				}
				if err := dc.RemoveTargetKindMappings(ctx); err != nil {
					slog.Warn("Failed to remove old TargetKindMappings", "err", err)
				}
				if err := dc.RemoveUnusedTargets(ctx); err != nil {
					slog.Warn("Failed to remove unused targets", "err", err)
				}
			}
		}
	})
}

// DeleteIncompleteLogs deletes logs which have had their incomplete
// build logs normalized.
func (dc *DbCleanupService) DeleteIncompleteLogs(ctx context.Context) error {
	ctx, span := dc.tracer.Start(ctx, "DbCleanupService.DeleteIncompleteLogs")
	defer span.End()

	deletedRows, err := dc.db.Ent().IncompleteBuildLog.Delete().
		Where(
			incompletebuildlog.HasBazelInvocationWith(
				bazelinvocation.HasBuildLogChunks(),
			),
		).
		Exec(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Could not delete incompleted build logs")
		return util.StatusWrap(err, "Could not delete incompleted build logs")
	}

	span.SetAttributes(attribute.KeyValue{Key: "incomplete_logs_deleted", Value: attribute.IntValue(deletedRows)})

	return nil
}

// RemoveOldInvocations removes invocations that have completed before a
// certain cutoff time.
func (dc *DbCleanupService) RemoveOldInvocations(ctx context.Context) error {
	ctx, span := dc.tracer.Start(ctx, "DbCleanupService.RemoveOldInvocations")
	defer span.End()

	cutoffTime := dc.clock.Now().UTC().Add(-dc.invocationRetention)
	deletedInvocation, err := dc.db.Ent().BazelInvocation.Delete().
		Where(
			bazelinvocation.BepCompleted(true),
			bazelinvocation.EndedAtLT(cutoffTime),
		).
		Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to remove old invocations")
	}

	span.SetAttributes(attribute.KeyValue{Key: "deleted_invocations", Value: attribute.IntValue(deletedInvocation)})

	return nil
}

// RemoveBuildsWithoutInvocations removes builds that do not have any
// associated invocations.
func (dc *DbCleanupService) RemoveBuildsWithoutInvocations(ctx context.Context) error {
	ctx, span := dc.tracer.Start(ctx, "DbCleanupService.RemoveBuildsWithoutInvocations")
	defer span.End()

	deletedBuilds, err := dc.db.Ent().Build.Delete().
		Where(
			build.Not(build.HasInvocations()),
		).
		Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to remove builds without invocations")
	}

	span.SetAttributes(attribute.KeyValue{Key: "deleted_builds", Value: attribute.IntValue(deletedBuilds)})

	return nil
}
