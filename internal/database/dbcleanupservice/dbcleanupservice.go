package dbcleanupservice

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"time"

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
	code "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DbCleanupService a service that performs periodic cleanup of the database to
// remove old data that is no longer needed, and to compact data in a format
// more suitable for storage.
type DbCleanupService struct {
	db                       database.Client
	batcher                  Batcher
	counter                  int64
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
	batcher Batcher,
	cleanupConfiguration *bb_portal.Database_CleanupConfiguration,
	tracerProvider trace.TracerProvider,
) (*DbCleanupService, error) {
	if db == nil {
		return nil, status.Error(code.NotFound, "No Database configured")
	}

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
		counter:                  rand.Int64N(65536),
		clock:                    clock,
		batcher:                  batcher,
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
			dc.counter++
			startTime := dc.clock.Now()

			dc.performCleanup(ctx)

			// Add 5% jitter to the target interval
			jitter := time.Duration((rand.Float64()*0.1 - 0.05) * float64(dc.cleanupInterval))
			targetInterval := dc.cleanupInterval + jitter

			elapsed := dc.clock.Now().Sub(startTime)
			timeToSleep := max(targetInterval-elapsed, 0)

			select {
			case <-ctx.Done():
				return nil
			case <-time.After(timeToSleep):
			}
		}
	})
}

func (dc *DbCleanupService) performCleanup(ctx context.Context) {
	dc.executeTask(ctx, "LockInvocationsWithNoRecentEvents", "invocations_locked", dc.LockInvocationsWithNoRecentEvents)
	dc.executeTask(ctx, "UpdateInvocationEndedAtFromEvents", "updated_invocations", dc.UpdateInvocationEndedAtFromEvents)
	dc.executeTask(ctx, "CompactLogs", "compacted_logs", dc.CompactLogs)
	dc.executeTask(ctx, "RemoveIncompleteLogs", "deleted_logs", dc.RemoveIncompleteLogs)
	dc.executeTask(ctx, "RemoveOldInvocations", "deleted_invocations", dc.RemoveOldInvocations)
	dc.executeTask(ctx, "RemoveInactiveUsers", "deleted_users", dc.RemoveInactiveUsers)
	dc.executeTask(ctx, "RemoveBuildsWithoutInvocations", "deleted_builds", dc.RemoveBuildsWithoutInvocations)
	dc.executeTask(ctx, "RemoveTargetKindMappings", "removed_target_kind_mappings", dc.RemoveTargetKindMappings)
	dc.executeTask(ctx, "RemoveUnusedTargets", "removed_unused_targets", dc.RemoveUnusedTargets)
	dc.executeTask(ctx, "RemoveOrphanedTestTargets", "removed_test_targets", dc.RemoveOrphanedTestTargets)
	dc.executeTask(ctx, "RemoveUnusedFiles", "removed_files", dc.RemoveUnusedFiles)
	dc.executeTask(ctx, "RemoveUnusedFilePaths", "removed_file_paths", dc.RemoveUnusedFilePaths)
	dc.executeTask(ctx, "RemoveUnusedDigests", "removed_digests", dc.RemoveUnusedDigests)
}

// A helper function which records metrics and tracing attributes
// for the cleanup tasks
func (dc *DbCleanupService) executeTask(ctx context.Context, taskName, attributeKey string, task func(context.Context) (int64, error)) {
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
		slog.Warn(fmt.Sprintf("DbCleanupService operation %s failed", taskName), "err", err)
	}
}
