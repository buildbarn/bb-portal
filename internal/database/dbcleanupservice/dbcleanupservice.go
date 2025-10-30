package dbcleanupservice

import (
	"context"
	"database/sql"
	"log/slog"
	"math/rand/v2"
	"strings"
	"time"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/ent/gen/ent/build"
	"github.com/buildbarn/bb-portal/ent/gen/ent/connectionmetadata"
	"github.com/buildbarn/bb-portal/ent/gen/ent/eventmetadata"
	"github.com/buildbarn/bb-portal/internal/database/common"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-storage/pkg/clock"
	"github.com/buildbarn/bb-storage/pkg/program"
	"github.com/buildbarn/bb-storage/pkg/util"
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
	db                          *ent.Client
	clock                       clock.Clock
	cleanupInterval             time.Duration
	invocationConnectionTimeout time.Duration
	invocationMessageTimeout    time.Duration
	invocationRetention         time.Duration
}

// NewDbCleanupService creates a new DbCleanupService.
func NewDbCleanupService(db *ent.Client, clock clock.Clock, cleanupConfiguration *bb_portal.BuildEventStreamService_DatabaseCleanupConfiguration) (*DbCleanupService, error) {
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
	}, nil
}

// StartDbCleanupService starts a goroutine that performs periodic
// cleanup of the database.
func (dc *DbCleanupService) StartDbCleanupService(ctx context.Context, group program.Group) {
	group.Go(func(ctx context.Context, siblingsGroup, dependenciesGroup program.Group) error {
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
				if err := dc.RemoveOldEventMetadata(ctx); err != nil {
					slog.Warn("Failed to remove old event metadata", "err", err)
				}
				if err := dc.RemoveOldInvocations(ctx); err != nil {
					slog.Warn("Failed to remove old invocations", "err", err)
				}
				if err := dc.RemoveBuildsWithoutInvocations(ctx); err != nil {
					slog.Warn("Failed to remove builds without invocations", "err", err)
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

	cutoffTime := dc.clock.Now().Add(-dc.invocationConnectionTimeout)

	// TODO: Set the end time of the invocation based on the last
	// connection time.
	invocationsUpdated, err := dc.db.BazelInvocation.Update().
		Where(
			bazelinvocation.BepCompleted(false),
			bazelinvocation.HasConnectionMetadataWith(
				connectionmetadata.ConnectionLastOpenAtLT(cutoffTime),
			),
		).
		SetBepCompleted(true).
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

	cutoffTime := dc.clock.Now().Add(-dc.invocationMessageTimeout)

	tx, err := dc.db.BeginTx(ctx, &entsql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return util.StatusWrap(err, "Failed to create transaction")
	}

	var invocationsToLock []struct {
		InvocationDbID int    `sql:"invocation_db_id"`
		MaxTime        string `sql:"max_time"`
	}

	err = tx.EventMetadata.
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
		err = tx.BazelInvocation.
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

	invocationIDs := make([]int, 0, len(invocationsToLock))
	for _, r := range invocationsToLock {
		invocationIDs = append(invocationIDs, r.InvocationDbID)
	}

	invocationsUpdated, err := tx.BazelInvocation.
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

	deletedRows, err := dc.db.ConnectionMetadata.Delete().
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

// RemoveOldEventMetadata removes event metadata for invocations that have
// completed before a certain cutoff time.
func (dc *DbCleanupService) RemoveOldEventMetadata(ctx context.Context) error {
	slog.Info("Removing old event metadata")
	cutoffTime := dc.clock.Now().Add(-dc.invocationMessageTimeout)
	// Remove all event metadata that is for invocations that have
	// completed before the cutoff time.
	deletedEM, err := dc.db.EventMetadata.Delete().
		Where(
			eventmetadata.HasBazelInvocationWith(
				bazelinvocation.BepCompleted(true),
				bazelinvocation.EndedAtLT(cutoffTime),
			),
		).
		Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to remove old event metadata")
	}
	slog.Info("Removed old event metadata", "count", deletedEM)
	return nil
}

// RemoveOldInvocations removes invocations that have completed before a
// certain cutoff time.
func (dc *DbCleanupService) RemoveOldInvocations(ctx context.Context) error {
	slog.Info("Removing old invocations")
	cutoffTime := dc.clock.Now().Add(-dc.invocationRetention)
	deletedInvocation, err := dc.db.BazelInvocation.Delete().
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
	deletedBuilds, err := dc.db.Build.Delete().
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
