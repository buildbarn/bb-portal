package dbcleanupservice

import (
	"context"
	"database/sql"
	"log/slog"
	"math/rand/v2"
	"strings"
	"time"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/ent/gen/ent/build"
	"github.com/buildbarn/bb-portal/ent/gen/ent/connectionmetadata"
	"github.com/buildbarn/bb-portal/ent/gen/ent/eventmetadata"
	"github.com/buildbarn/bb-portal/ent/gen/ent/incompletebuildlog"
	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-portal/internal/database/common"
	"github.com/buildbarn/bb-portal/internal/database/dbauthservice"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-storage/pkg/clock"
	"github.com/buildbarn/bb-storage/pkg/program"
	"github.com/buildbarn/bb-storage/pkg/util"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// DbCleanupService a service that performs periodic cleanup of the
// database to remove old data that is no longer needed. This includes:
//
//  1. Locking unfinished invocations that have not received any new event
//     metadata for a certain period of time (invocationMessageTimeout).
//  2. Removing event metadata for invocations that have completed and
//     whose completion time is older than invocationMessageTimeout.
//  3. Removing invocations that have completed and whose completion time
//     is older than invocationRetention.
//  4. Removing builds that do not have any associated invocations.
type DbCleanupService struct {
	db                          database.Client
	clock                       clock.Clock
	cleanupInterval             time.Duration
	invocationConnectionTimeout time.Duration
	invocationMessageTimeout    time.Duration
	invocationRetention         time.Duration
	tracer                      trace.Tracer
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

	invocationConnectionTimeout := cleanupConfiguration.InvocationConnectionTimeout
	if err := invocationConnectionTimeout.CheckValid(); err != nil {
		return nil, util.StatusWrap(err, "Failed to parse invocationConnectionTimeout parameter time")
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
		db:                          db,
		clock:                       clock,
		cleanupInterval:             cleanupInterval.AsDuration(),
		invocationConnectionTimeout: invocationConnectionTimeout.AsDuration(),
		invocationMessageTimeout:    invocationMessageTimeout.AsDuration(),
		invocationRetention:         invocationRetention.AsDuration(),
		tracer:                      tracerProvider.Tracer("github.com/buildbarn/bb-portal/internal/database/dbcleanupservice"),
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
				slog.Info("Starting database cleanup")
				if err := dc.LockInvocationsWithNoRecentConnections(ctx); err != nil {
					slog.Warn("Failed to lock unfinished invocations with no recent connections", "err", err)
				}
				if err := dc.LockInvocationsWithNoRecentEvents(ctx); err != nil {
					slog.Warn("Failed to lock unfinished invocations with no recent events", "err", err)
				}
				if err := dc.RemoveOldInvocationConnections(ctx); err != nil {
					slog.Warn("Failed to remove old invocation connections", "err", err)
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
				slog.Info("Finished database cleanup")
			}
		}
	})
}

// LockInvocationsWithNoRecentConnections locks invocations where the gRPC
// stream has been interrupted, and no new connection has been made within a
// certain period of time.
func (dc *DbCleanupService) LockInvocationsWithNoRecentConnections(ctx context.Context) error {
	slog.Info("Locking unfinished invocations with no recent connections")

	ctx, span := dc.tracer.Start(ctx, "DbCleanupService.LockInvocationsWithNoRecentEvents")
	defer span.End()

	cutoffTime := dc.clock.Now().UTC().Add(-dc.invocationConnectionTimeout)

	invocationsUpdated, err := dc.db.Ent().BazelInvocation.Update().
		Where(
			bazelinvocation.BepCompleted(false),
			bazelinvocation.HasConnectionMetadataWith(
				connectionmetadata.ConnectionLastOpenAtLT(cutoffTime),
			),
		).
		SetBepCompleted(true).
		SetEndedAt(cutoffTime).
		Save(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to lock invocations")
	}

	slog.Info("Locked unfinished invocations", "count", invocationsUpdated)
	return nil
}

// LockInvocationsWithNoRecentEvents locks invocations that have not received any new
// events in a certain period of time.
func (dc *DbCleanupService) LockInvocationsWithNoRecentEvents(ctx context.Context) error {
	slog.Info("Locking unfinished invocations that has no recent events")

	ctx, span := dc.tracer.Start(ctx, "DbCleanupService.LockInvocationsWithNoRecentEvents")
	defer span.End()

	cutoffTime := dc.clock.Now().UTC().Add(-dc.invocationMessageTimeout)

	tx, err := dc.db.BeginTx(ctx, &entsql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return util.StatusWrap(err, "Failed to create transaction")
	}

	var invocationsToLock []struct {
		InvocationDbID int64  `sql:"invocation_db_id"`
		MaxTime        string `sql:"max_time"`
	}

	err = tx.Ent().EventMetadata.
		Query().
		Modify(func(sel *entsql.Selector) {
			sel.Select(
				entsql.As(sel.C(eventmetadata.BazelInvocationColumn), "invocation_db_id"),
				entsql.As(entsql.Max(sel.C(eventmetadata.FieldEventReceivedAt)), "max_time"),
			)
			sel.GroupBy(sel.C(eventmetadata.BazelInvocationColumn))
			sel.Having(entsql.LT(entsql.Max(sel.C(eventmetadata.FieldEventReceivedAt)), cutoffTime))
		}).
		Scan(ctx, &invocationsToLock)
	if err != nil {
		return common.RollbackAndWrapError(tx, util.StatusWrap(err, "Failed to query invocations to lock"))
	}

	if len(invocationsToLock) == 0 {
		slog.Info("No invocations to lock")
		return common.RollbackAndWrapError(tx, nil)
	}

	for _, r := range invocationsToLock {
		// Sqlite does not fully support the RFC 3339 format, so we need to replace the space
		// between the date and time with a 'T'.
		endedAt, err := time.Parse(time.RFC3339Nano, strings.Replace(r.MaxTime, " ", "T", 1))
		if err != nil {
			return common.RollbackAndWrapError(tx, util.StatusWrapf(err, "Failed to parse time %s", r.MaxTime))
		}
		err = tx.Ent().BazelInvocation.
			Update().
			Where(
				bazelinvocation.IDEQ(r.InvocationDbID),
				bazelinvocation.EndedAtIsNil(),
			).
			SetEndedAt(endedAt).
			Exec(ctx)
		if err != nil {
			return common.RollbackAndWrapError(tx, util.StatusWrapf(err, "Failed to set ended_at for invocation %s", r.InvocationDbID))
		}
	}

	invocationIDs := make([]int64, 0, len(invocationsToLock))
	for _, r := range invocationsToLock {
		invocationIDs = append(invocationIDs, r.InvocationDbID)
	}

	invocationsUpdated, err := tx.Ent().BazelInvocation.
		Update().
		Where(bazelinvocation.IDIn(invocationIDs...)).
		SetBepCompleted(true).
		Save(ctx)
	if err != nil {
		return common.RollbackAndWrapError(tx, util.StatusWrap(err, "Failed to lock invocations"))
	}
	err = tx.Commit()
	if err != nil {
		return util.StatusWrap(err, "Failed to commit transaction")
	}
	slog.Info("Locked unfinished invocations", "count", invocationsUpdated)
	return nil
}

// RemoveOldInvocationConnections removes InvocationConnections for invocations
// that have completed.
func (dc *DbCleanupService) RemoveOldInvocationConnections(ctx context.Context) error {
	slog.Info("Locking old invocation connections")

	ctx, span := dc.tracer.Start(ctx, "DbCleanupService.RemoveOldInvocationConnections")
	defer span.End()

	deletedRows, err := dc.db.Ent().ConnectionMetadata.Delete().
		Where(
			connectionmetadata.HasBazelInvocationWith(
				bazelinvocation.BepCompleted(true),
			),
		).Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to lock invocations")
	}

	slog.Info("Removed old invocations connections", "count", deletedRows)
	return nil
}

// DeleteIncompleteLogs deletes logs which have had their incomplete
// build logs normalized.
func (dc *DbCleanupService) DeleteIncompleteLogs(ctx context.Context) error {
	ctx, span := dc.tracer.Start(ctx, "DbCleanupService.DeleteIncompleteLogs")
	defer span.End()

	_, err := dc.db.Ent().IncompleteBuildLog.Delete().
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

	return nil
}

// RemoveOldInvocations removes invocations that have completed before a
// certain cutoff time.
func (dc *DbCleanupService) RemoveOldInvocations(ctx context.Context) error {
	slog.Info("Removing old invocations")

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
	slog.Info("Removed old invocations", "count", deletedInvocation)
	return nil
}

// RemoveBuildsWithoutInvocations removes builds that do not have any
// associated invocations.
func (dc *DbCleanupService) RemoveBuildsWithoutInvocations(ctx context.Context) error {
	slog.Info("Removing builds without invocations")

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
	slog.Info("Removed builds without invocations", "count", deletedBuilds)
	return nil
}
