package dbcleanupservice

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"time"

	"github.com/buildbarn/bb-portal/ent/gen/ent/build"
	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-portal/internal/database/dbauthservice"
	prometheusmetrics "github.com/buildbarn/bb-portal/pkg/prometheus_metrics"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-storage/pkg/clock"
	"github.com/buildbarn/bb-storage/pkg/program"
	"github.com/buildbarn/bb-storage/pkg/util"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
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
	counter                  int64
	clock                    clock.Clock
	cleanupInterval          time.Duration
	invocationMessageTimeout time.Duration
	invocationRetention      time.Duration
	artifactRetention        time.Duration
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

	artifactRetentionPb := cleanupConfiguration.ArtifactRetention
	artifactRetention := invocationRetention.AsDuration() // default: same as invocation_retention
	if artifactRetentionPb != nil {
		if err := artifactRetentionPb.CheckValid(); err != nil {
			return nil, util.StatusWrap(err, "Failed to parse artifactRetention parameter time")
		}
		artifactRetention = artifactRetentionPb.AsDuration()
		if artifactRetention > invocationRetention.AsDuration() {
			return nil, grpcstatus.Errorf(grpccodes.InvalidArgument,
				"artifact_retention (%s) must be <= invocation_retention (%s)",
				artifactRetention, invocationRetention.AsDuration())
		}
	}

	return &DbCleanupService{
		db:                       db,
		counter:                  rand.Int64N(65536),
		clock:                    clock,
		cleanupInterval:          cleanupInterval.AsDuration(),
		invocationMessageTimeout: invocationMessageTimeout.AsDuration(),
		invocationRetention:      invocationRetention.AsDuration(),
		artifactRetention:        artifactRetention,
		tracer:                   tracerProvider.Tracer("github.com/buildbarn/bb-portal/internal/database/dbcleanupservice"),
	}, nil
}

// StartDbCleanupService starts a goroutine that performs periodic
// cleanup of the database.
func (dc *DbCleanupService) StartDbCleanupService(ctx context.Context, group program.Group) {
	group.Go(func(ctx context.Context, siblingsGroup, dependenciesGroup program.Group) error {
		ctx = dbauthservice.NewContextWithDbAuthServiceBypass(ctx)
		for {
			dc.counter++
			// Add 5% jitter to the cleanup interval
			timeToSleep := dc.cleanupInterval + time.Duration((rand.Float64()*0.1-0.05)*float64(dc.cleanupInterval))
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(timeToSleep):
				if err := dc.executeTask(ctx, "LockInvocationsWithNoRecentEvents", "invocations_locked", dc.LockInvocationsWithNoRecentEvents); err != nil {
					slog.Warn("Failed to lock unfinished invocations with no recent events", "err", err)
				}
				if err := dc.executeTask(ctx, "UpdateInvocationEndedAtFromEvents", "updated_invocations", dc.UpdateInvocationEndedAtFromEvents); err != nil {
					slog.Warn("Failed to update invocation ended_at from event metadata", "err", err)
				}
				if err := dc.executeTask(ctx, "CompactLogs", "compacted_logs", dc.CompactLogs); err != nil {
					slog.Warn("Failed to compact logs", "err", err)
				}
				if err := dc.executeTask(ctx, "DeleteIncompleteLogs", "deleted_logs", dc.DeleteIncompleteLogs); err != nil {
					slog.Warn("Failed to delete incomplete logs", "err", err)
				}
				if err := dc.executeTask(ctx, "RemoveOldArtifacts", "deleted_artifacts", dc.RemoveOldArtifacts); err != nil {
					slog.Warn("Failed to remove old artifacts", "err", err)
				}
				if err := dc.executeTask(ctx, "RemoveOldInvocations", "deleted_invocations", dc.RemoveOldInvocations); err != nil {
					slog.Warn("Failed to remove old invocations", "err", err)
				}
				if err := dc.executeTask(ctx, "RemoveInactiveUsers", "deleted_users", dc.RemoveInactiveUsers); err != nil {
					slog.Warn("Failed to remove users without invocations")
				}
				if err := dc.executeTask(ctx, "RemoveBuildsWithoutInvocations", "deleted_builds", dc.RemoveBuildsWithoutInvocations); err != nil {
					slog.Warn("Failed to remove builds without invocations", "err", err)
				}
				if err := dc.executeTask(ctx, "RemoveTargetKindMappings", "removed_target_kind_mappings", dc.RemoveTargetKindMappings); err != nil {
					slog.Warn("Failed to remove old TargetKindMappings", "err", err)
				}
				if err := dc.executeTask(ctx, "RemoveUnusedTargets", "removed_unused_targets", dc.RemoveUnusedTargets); err != nil {
					slog.Warn("Failed to remove unused targets", "err", err)
				}
				if err := dc.executeTask(ctx, "RemoveOrphanedTestTargets", "removed_test_targets", dc.RemoveOrphanedTestTargets); err != nil {
					slog.Warn("Failed to remove orphaned test targets", "err", err)
				}
			}
		}
	})
}

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

// A helper function which records metrics and tracing attributes
// for the cleanup tasks
func (dc *DbCleanupService) executeTask(ctx context.Context, taskName, attributeKey string, task func(context.Context) (int64, error)) error {
	ctx, span := dc.tracer.Start(ctx, fmt.Sprintf("DbCleanupService.%s", taskName))
	defer span.End()
	start := dc.clock.Now()

	volume, err := task(ctx)
	prometheusmetrics.CleanupDurations.WithLabelValues(taskName).Add(dc.clock.Now().Sub(start).Seconds())
	prometheusmetrics.CleanupVolumes.WithLabelValues(taskName).Add(float64(volume))
	span.SetAttributes(attribute.Int64(attributeKey, volume))

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, fmt.Sprintf("An error occured during cleanup service %s", taskName))
	}

	return err
}
