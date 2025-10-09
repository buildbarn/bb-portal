package dbcleanupservice_test

import (
	"context"
	"testing"
	"time"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/internal/database/dbcleanupservice"
	"github.com/buildbarn/bb-portal/internal/mock"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-storage/pkg/clock"
	_ "github.com/buildbarn/bb-storage/pkg/clock"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/durationpb"
)

func setupTestDB(t *testing.T) *ent.Client {
	client, err := ent.Open("sqlite3", "file:dbCleanupTestDatabase?mode=memory&cache=shared&_fk=1")
	require.NoError(t, err)
	t.Cleanup(func() { client.Close() })
	err = client.Schema.Create(context.Background())
	require.NoError(t, err)
	return client
}

func getNewDbCleanupService(client *ent.Client, clock clock.Clock, traceProvider trace.TracerProvider) (*dbcleanupservice.DbCleanupService, error) {
	cleanupConfiguration := &bb_portal.BuildEventStreamService_DatabaseCleanupConfiguration{
		CleanupInterval:             durationpb.New(1 * time.Minute),
		InvocationConnectionTimeout: durationpb.New(30 * time.Second),
		InvocationMessageTimeout:    durationpb.New(30 * time.Second),
		InvocationRetention:         durationpb.New(30 * time.Minute),
	}
	return dbcleanupservice.NewDbCleanupService(client, clock, cleanupConfiguration, traceProvider)
}

func TestLockInvocationsWithNoRecentConnections(t *testing.T) {
	ctrl, ctx := gomock.WithContext(context.Background(), t)
	clock := mock.NewMockClock(ctrl)
	cleanupTime := time.Unix(1600000300, 0)
	traceProvider := noop.NewTracerProvider()

	t.Run("NoInvocations", func(t *testing.T) {
		client := setupTestDB(t)

		cleanup, err := getNewDbCleanupService(client, clock, traceProvider)
		require.NoError(t, err)

		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentConnections(ctx)
		require.NoError(t, err)

		count, err := client.BazelInvocation.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})

	t.Run("InvocationsWithNoConnectionMetadata", func(t *testing.T) {
		client := setupTestDB(t)
		inv1, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetBepCompleted(false).
			Save(ctx)
		require.NoError(t, err)
		inv2, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetBepCompleted(true).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(client, clock, traceProvider)
		require.NoError(t, err)

		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentConnections(ctx)
		require.NoError(t, err)

		count, err := client.BazelInvocation.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 2, count)

		got1, err := client.BazelInvocation.Get(ctx, inv1.ID)
		require.NoError(t, err)
		require.False(t, got1.BepCompleted)

		got2, err := client.BazelInvocation.Get(ctx, inv2.ID)
		require.NoError(t, err)
		require.True(t, got2.BepCompleted)
	})

	t.Run("UnfinishedInvocationNoRecentConnections", func(t *testing.T) {
		client := setupTestDB(t)
		inv, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			Save(ctx)
		require.NoError(t, err)

		_, err = client.ConnectionMetadata.Create().
			SetBazelInvocationID(inv.ID).
			SetConnectionLastOpenAt(cleanupTime.Add(-60 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(client, clock, traceProvider)
		require.NoError(t, err)

		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentConnections(ctx)
		require.NoError(t, err)

		got, err := client.BazelInvocation.Get(ctx, inv.ID)
		require.NoError(t, err)
		require.True(t, got.BepCompleted)
	})

	t.Run("UnfinishedInvocationRecentConnection", func(t *testing.T) {
		client := setupTestDB(t)
		inv, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			Save(ctx)
		require.NoError(t, err)

		_, err = client.ConnectionMetadata.Create().
			SetBazelInvocationID(inv.ID).
			SetConnectionLastOpenAt(cleanupTime.Add(-10 * time.Second)). // recent
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(client, clock, traceProvider)
		require.NoError(t, err)

		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentConnections(ctx)
		require.NoError(t, err)

		got, err := client.BazelInvocation.Get(ctx, inv.ID)
		require.NoError(t, err)
		require.False(t, got.BepCompleted)
	})

	t.Run("FinishedInvocationWithOldConnection", func(t *testing.T) {
		client := setupTestDB(t)
		inv, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetBepCompleted(true).
			Save(ctx)
		require.NoError(t, err)

		_, err = client.ConnectionMetadata.Create().
			SetBazelInvocationID(inv.ID).
			SetConnectionLastOpenAt(cleanupTime.Add(-60 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(client, clock, traceProvider)
		require.NoError(t, err)

		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentConnections(ctx)
		require.NoError(t, err)

		got, err := client.BazelInvocation.Get(ctx, inv.ID)
		require.NoError(t, err)
		require.True(t, got.BepCompleted)
	})

	t.Run("MultipleMixedInvocations", func(t *testing.T) {
		client := setupTestDB(t)
		// old connection -> should be locked
		invOld, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.ConnectionMetadata.Create().
			SetBazelInvocationID(invOld.ID).
			SetConnectionLastOpenAt(cleanupTime.Add(-60 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		// recent connection -> should remain unlocked
		invRecent, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.ConnectionMetadata.Create().
			SetBazelInvocationID(invRecent.ID).
			SetConnectionLastOpenAt(cleanupTime.Add(-10 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(client, clock, traceProvider)
		require.NoError(t, err)

		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentConnections(ctx)
		require.NoError(t, err)

		gotOld, err := client.BazelInvocation.Get(ctx, invOld.ID)
		require.NoError(t, err)
		require.True(t, gotOld.BepCompleted)

		gotRecent, err := client.BazelInvocation.Get(ctx, invRecent.ID)
		require.NoError(t, err)
		require.False(t, gotRecent.BepCompleted)
	})
}

func TestLockInvocationsWithNoRecentEvents(t *testing.T) {
	ctrl, ctx := gomock.WithContext(context.Background(), t)
	clock := mock.NewMockClock(ctrl)
	cleanupTime := time.Unix(1600000060, 0)
	traceProvider := noop.NewTracerProvider()

	t.Run("NoInvocations-NoEventMetadata", func(t *testing.T) {
		client := setupTestDB(t)
		cleanup, err := getNewDbCleanupService(client, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentEvents(ctx)
		require.NoError(t, err)
	})

	t.Run("UnfinishedInvocation-NoEventMetadata", func(t *testing.T) {
		client := setupTestDB(t)
		startInv, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(client, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentEvents(ctx)
		require.NoError(t, err)

		resInv, err := client.BazelInvocation.Get(ctx, startInv.ID)
		require.NoError(t, err)
		require.Nil(t, resInv.EndedAt)
		require.False(t, resInv.BepCompleted)
	})

	t.Run("FinishedInvocation-NoEventMetadata", func(t *testing.T) {
		client := setupTestDB(t)
		startInv, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetBepCompleted(true).
			SetEndedAt(cleanupTime.Add(-60 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(client, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentEvents(ctx)
		require.NoError(t, err)

		resInv, err := client.BazelInvocation.Get(ctx, startInv.ID)
		require.NoError(t, err)
		require.Equal(t, startInv.EndedAt.UTC(), resInv.EndedAt.UTC())
		require.True(t, resInv.BepCompleted)
	})

	t.Run("UnfinishedInvocationWithoutEndedAt", func(t *testing.T) {
		client := setupTestDB(t)
		invocationDb, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			Save(ctx)
		require.NoError(t, err)
		em, err := client.EventMetadata.Create().
			SetBazelInvocationID(invocationDb.ID).
			SetSequenceNumber(1).
			SetEventHash("hash").
			SetEventReceivedAt(cleanupTime.Add(-60 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(client, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentEvents(ctx)
		require.NoError(t, err)

		inv, err := client.BazelInvocation.Get(ctx, invocationDb.ID)
		require.NoError(t, err)
		require.Equal(t, em.EventReceivedAt.UTC(), inv.EndedAt.UTC())
		require.True(t, inv.BepCompleted)
	})

	t.Run("UnfinishedInvocationWithEndedAt", func(t *testing.T) {
		client := setupTestDB(t)
		invocationDb, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetEndedAt(cleanupTime.Add(-60 * time.Second)).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.EventMetadata.Create().
			SetBazelInvocationID(invocationDb.ID).
			SetSequenceNumber(1).
			SetEventHash("hash").
			SetEventReceivedAt(cleanupTime.Add(-50 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(client, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentEvents(ctx)
		require.NoError(t, err)

		inv, err := client.BazelInvocation.Get(ctx, invocationDb.ID)
		require.NoError(t, err)
		require.Equal(t, cleanupTime.Add(-60*time.Second).UTC(), inv.EndedAt.UTC())
		require.True(t, inv.BepCompleted)
	})

	t.Run("FinishedInvocationWithoutEndedAt", func(t *testing.T) {
		client := setupTestDB(t)
		invocationDb, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetBepCompleted(true).
			Save(ctx)
		require.NoError(t, err)
		em, err := client.EventMetadata.Create().
			SetBazelInvocationID(invocationDb.ID).
			SetSequenceNumber(1).
			SetEventHash("hash").
			SetEventReceivedAt(cleanupTime.Add(-50 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(client, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentEvents(ctx)
		require.NoError(t, err)

		inv, err := client.BazelInvocation.Get(ctx, invocationDb.ID)
		require.NoError(t, err)
		require.Equal(t, em.EventReceivedAt.UTC(), inv.EndedAt.UTC())
		require.True(t, inv.BepCompleted)
	})

	t.Run("FinishedInvocationWithEndedAt", func(t *testing.T) {
		client := setupTestDB(t)
		invocationDb, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetEndedAt(cleanupTime.Add(-60 * time.Second)).
			SetBepCompleted(true).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.EventMetadata.Create().
			SetBazelInvocationID(invocationDb.ID).
			SetSequenceNumber(1).
			SetEventHash("hash").
			SetEventReceivedAt(cleanupTime.Add(-50 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(client, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentEvents(ctx)
		require.NoError(t, err)

		inv, err := client.BazelInvocation.Get(ctx, invocationDb.ID)
		require.NoError(t, err)
		require.Equal(t, cleanupTime.Add(-60*time.Second).UTC(), inv.EndedAt.UTC())
		require.True(t, inv.BepCompleted)
	})

	t.Run("UnfinishedInvocation-MultipleEventMetadata", func(t *testing.T) {
		client := setupTestDB(t)
		invocationDb, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.EventMetadata.Create().
			SetBazelInvocationID(invocationDb.ID).
			SetSequenceNumber(1).
			SetEventHash("hash1").
			SetEventReceivedAt(cleanupTime.Add(-59 * time.Second)).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.EventMetadata.Create().
			SetBazelInvocationID(invocationDb.ID).
			SetSequenceNumber(2).
			SetEventHash("hash2").
			SetEventReceivedAt(cleanupTime.Add(-58 * time.Second)).
			Save(ctx)
		require.NoError(t, err)
		em3, err := client.EventMetadata.Create().
			SetBazelInvocationID(invocationDb.ID).
			SetSequenceNumber(3).
			SetEventHash("hash3").
			SetEventReceivedAt(cleanupTime.Add(-57 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(client, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentEvents(ctx)
		require.NoError(t, err)

		inv, err := client.BazelInvocation.Get(ctx, invocationDb.ID)
		require.NoError(t, err)
		require.Equal(t, em3.EventReceivedAt.UTC(), inv.EndedAt.UTC())
		require.True(t, inv.BepCompleted)
	})

	t.Run("MultipleUnfinishedInvocations", func(t *testing.T) {
		client := setupTestDB(t)
		invocationDb1, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			Save(ctx)
		require.NoError(t, err)
		em1, err := client.EventMetadata.Create().
			SetBazelInvocationID(invocationDb1.ID).
			SetSequenceNumber(1).
			SetEventHash("hash1").
			SetEventReceivedAt(cleanupTime.Add(-60 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		invocationDb2, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			Save(ctx)
		require.NoError(t, err)
		em2, err := client.EventMetadata.Create().
			SetBazelInvocationID(invocationDb2.ID).
			SetSequenceNumber(1).
			SetEventHash("hash2").
			SetEventReceivedAt(cleanupTime.Add(-50 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(client, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentEvents(ctx)
		require.NoError(t, err)

		inv1, err := client.BazelInvocation.Get(ctx, invocationDb1.ID)
		require.NoError(t, err)
		require.Equal(t, em1.EventReceivedAt.UTC(), inv1.EndedAt.UTC())
		require.True(t, inv1.BepCompleted)

		inv2, err := client.BazelInvocation.Get(ctx, invocationDb2.ID)
		require.NoError(t, err)
		require.Equal(t, em2.EventReceivedAt.UTC(), inv2.EndedAt.UTC())
		require.True(t, inv2.BepCompleted)
	})

	t.Run("UnfinishedInvocation-RecentEventMetadata", func(t *testing.T) {
		client := setupTestDB(t)
		invocationDb, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.EventMetadata.Create().
			SetBazelInvocationID(invocationDb.ID).
			SetSequenceNumber(1).
			SetEventHash("hash").
			SetEventReceivedAt(cleanupTime.Add(-20 * time.Second)).
			Save(ctx)
		_, err = client.EventMetadata.Create().
			SetBazelInvocationID(invocationDb.ID).
			SetSequenceNumber(2).
			SetEventHash("hash").
			SetEventReceivedAt(cleanupTime.Add(-10 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(client, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentEvents(ctx)
		require.NoError(t, err)

		inv, err := client.BazelInvocation.Get(ctx, invocationDb.ID)
		require.NoError(t, err)
		require.Nil(t, inv.EndedAt)
		require.False(t, inv.BepCompleted)
	})
}

func TestRemoveOldInvocationConnections(t *testing.T) {
	ctrl, ctx := gomock.WithContext(context.Background(), t)
	clock := mock.NewMockClock(ctrl)
	traceProvider := noop.NewTracerProvider()

	t.Run("NoConnections", func(t *testing.T) {
		client := setupTestDB(t)

		cleanup, err := getNewDbCleanupService(client, clock, traceProvider)
		require.NoError(t, err)

		err = cleanup.RemoveOldInvocationConnections(ctx)
		require.NoError(t, err)

		count, err := client.ConnectionMetadata.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})

	t.Run("ConnectionForCompletedInvocation", func(t *testing.T) {
		client := setupTestDB(t)
		inv, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetBepCompleted(true).
			Save(ctx)
		require.NoError(t, err)

		_, err = client.ConnectionMetadata.Create().
			SetBazelInvocationID(inv.ID).
			SetConnectionLastOpenAt(time.Now()).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(client, clock, traceProvider)
		require.NoError(t, err)

		err = cleanup.RemoveOldInvocationConnections(ctx)
		require.NoError(t, err)

		count, err := client.ConnectionMetadata.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, count)

		// Invocation should still exist.
		invCount, err := client.BazelInvocation.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, invCount)
	})

	t.Run("ConnectionForUnfinishedInvocation", func(t *testing.T) {
		client := setupTestDB(t)
		inv, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetBepCompleted(false).
			Save(ctx)
		require.NoError(t, err)

		_, err = client.ConnectionMetadata.Create().
			SetBazelInvocationID(inv.ID).
			SetConnectionLastOpenAt(time.Now()).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(client, clock, traceProvider)
		require.NoError(t, err)

		err = cleanup.RemoveOldInvocationConnections(ctx)
		require.NoError(t, err)

		count, err := client.ConnectionMetadata.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})

	t.Run("MultipleMixedConnections", func(t *testing.T) {
		client := setupTestDB(t)
		// Completed invocation
		invDone, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetBepCompleted(true).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.ConnectionMetadata.Create().
			SetBazelInvocationID(invDone.ID).
			SetConnectionLastOpenAt(time.Now()).
			Save(ctx)
		require.NoError(t, err)

		// Unfinished invocation
		invNotDone, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetBepCompleted(false).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.ConnectionMetadata.Create().
			SetBazelInvocationID(invNotDone.ID).
			SetConnectionLastOpenAt(time.Now()).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(client, clock, traceProvider)
		require.NoError(t, err)

		err = cleanup.RemoveOldInvocationConnections(ctx)
		require.NoError(t, err)

		// Only the connection for the unfinished invocation should remain.
		conns, err := client.ConnectionMetadata.Query().All(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, len(conns))
		invDbId, err := conns[0].QueryBazelInvocation().OnlyID(ctx)
		require.NoError(t, err)
		require.Equal(t, int(invNotDone.ID), invDbId)
	})
}

func TestRemoveOldEventMetadata(t *testing.T) {
	ctrl, ctx := gomock.WithContext(context.Background(), t)
	clock := mock.NewMockClock(ctrl)
	cleanupTime := time.Unix(1600000120, 0)
	traceProvider := noop.NewTracerProvider()

	t.Run("NoEventMetadata", func(t *testing.T) {
		client := setupTestDB(t)
		cleanup, err := getNewDbCleanupService(client, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.RemoveOldEventMetadata(ctx)
		require.NoError(t, err)
		count, err := client.EventMetadata.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})

	t.Run("EventMetadataNotOld", func(t *testing.T) {
		client := setupTestDB(t)
		inv, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetBepCompleted(true).
			SetEndedAt(cleanupTime.Add(-15 * time.Second)).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.EventMetadata.Create().
			SetBazelInvocationID(inv.ID).
			SetSequenceNumber(1).
			SetEventHash("hash").
			SetEventReceivedAt(cleanupTime.Add(-15 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(client, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.RemoveOldEventMetadata(ctx)
		require.NoError(t, err)
		count, err := client.EventMetadata.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})

	t.Run("EventMetadataOld", func(t *testing.T) {
		client := setupTestDB(t)
		inv, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetBepCompleted(true).
			SetEndedAt(cleanupTime.Add(-60 * time.Second)).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.EventMetadata.Create().
			SetBazelInvocationID(inv.ID).
			SetSequenceNumber(1).
			SetEventHash("hash").
			SetEventReceivedAt(cleanupTime.Add(-60 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(client, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.RemoveOldEventMetadata(ctx)
		require.NoError(t, err)
		count, err := client.EventMetadata.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})

	t.Run("EventMetadataOldButInvocationNotCompleted", func(t *testing.T) {
		client := setupTestDB(t)
		inv, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetEndedAt(cleanupTime.Add(-60 * time.Second)).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.EventMetadata.Create().
			SetBazelInvocationID(inv.ID).
			SetSequenceNumber(1).
			SetEventHash("hash").
			SetEventReceivedAt(cleanupTime.Add(-60 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(client, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.RemoveOldEventMetadata(ctx)
		require.NoError(t, err)
		count, err := client.EventMetadata.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})

	t.Run("EventMetadataOldButInvocationEndedAtAfterCutoff", func(t *testing.T) {
		client := setupTestDB(t)
		inv, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetBepCompleted(true).
			SetEndedAt(cleanupTime.Add(-15 * time.Second)).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.EventMetadata.Create().
			SetBazelInvocationID(inv.ID).
			SetSequenceNumber(1).
			SetEventHash("hash").
			SetEventReceivedAt(cleanupTime.Add(-60 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(client, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.RemoveOldEventMetadata(ctx)
		require.NoError(t, err)
		count, err := client.EventMetadata.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})

	t.Run("MultipleEventMetadata", func(t *testing.T) {
		client := setupTestDB(t)
		inv, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetBepCompleted(true).
			SetEndedAt(cleanupTime.Add(-60 * time.Second)).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.EventMetadata.Create().
			SetBazelInvocationID(inv.ID).
			SetSequenceNumber(1).
			SetEventHash("hash1").
			SetEventReceivedAt(cleanupTime.Add(-60 * time.Second)).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.EventMetadata.Create().
			SetBazelInvocationID(inv.ID).
			SetSequenceNumber(2).
			SetEventHash("hash2").
			SetEventReceivedAt(cleanupTime.Add(-15 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(client, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.RemoveOldEventMetadata(ctx)
		require.NoError(t, err)
		count, err := client.EventMetadata.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})
}

func TestRemoveOldInvocations(t *testing.T) {
	ctrl, ctx := gomock.WithContext(context.Background(), t)
	clock := mock.NewMockClock(ctrl)
	cleanupTime := time.Unix(1600000200, 0)
	traceProvider := noop.NewTracerProvider()

	t.Run("NoInvocations", func(t *testing.T) {
		client := setupTestDB(t)

		cleanup, err := getNewDbCleanupService(client, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.RemoveOldInvocations(ctx)
		require.NoError(t, err)

		count, err := client.BazelInvocation.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})

	t.Run("InvocationNotCompleted", func(t *testing.T) {
		client := setupTestDB(t)
		_, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetEndedAt(cleanupTime.Add(-60 * time.Minute)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(client, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.RemoveOldInvocations(ctx)
		require.NoError(t, err)

		count, err := client.BazelInvocation.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})

	t.Run("InvocationCompletedButNotOld", func(t *testing.T) {
		client := setupTestDB(t)
		_, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetBepCompleted(true).
			SetEndedAt(cleanupTime.Add(-15 * time.Minute)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(client, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.RemoveOldInvocations(ctx)
		require.NoError(t, err)

		count, err := client.BazelInvocation.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})

	t.Run("InvocationCompletedAndOld", func(t *testing.T) {
		client := setupTestDB(t)
		_, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetBepCompleted(true).
			SetEndedAt(cleanupTime.Add(-60 * time.Minute)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(client, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.RemoveOldInvocations(ctx)
		require.NoError(t, err)

		count, err := client.BazelInvocation.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})

	t.Run("MultipleInvocationsMixed", func(t *testing.T) {
		client := setupTestDB(t)
		// Old and completed
		_, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetBepCompleted(true).
			SetEndedAt(cleanupTime.Add(-60 * time.Minute)).
			Save(ctx)
		require.NoError(t, err)
		// Not completed
		_, err = client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetEndedAt(cleanupTime.Add(-60 * time.Minute)).
			Save(ctx)
		require.NoError(t, err)
		// Completed but not old
		_, err = client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetBepCompleted(true).
			SetEndedAt(cleanupTime.Add(-15 * time.Minute)).
			Save(ctx)
		require.NoError(t, err)
		// Not completed and not old
		_, err = client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetEndedAt(cleanupTime.Add(-15 * time.Minute)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(client, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.RemoveOldInvocations(ctx)
		require.NoError(t, err)

		count, err := client.BazelInvocation.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 3, count)
	})
}

func TestRemoveBuildsWithoutInvocations(t *testing.T) {
	ctrl, ctx := gomock.WithContext(context.Background(), t)
	clock := mock.NewMockClock(ctrl)
	traceProvider := noop.NewTracerProvider()

	t.Run("NoBuilds", func(t *testing.T) {
		client := setupTestDB(t)

		cleanup, err := getNewDbCleanupService(client, clock, traceProvider)
		require.NoError(t, err)
		err = cleanup.RemoveBuildsWithoutInvocations(ctx)
		require.NoError(t, err)

		count, err := client.Build.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})

	t.Run("BuildWithInvocation", func(t *testing.T) {
		client := setupTestDB(t)
		buildObj, err := client.Build.Create().SetBuildURL("1").SetBuildUUID(uuid.New()).SetInstanceName("").SetTimestamp(time.Now()).Save(ctx)
		require.NoError(t, err)
		_, err = client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetBuild(buildObj).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(client, clock, traceProvider)
		require.NoError(t, err)
		err = cleanup.RemoveBuildsWithoutInvocations(ctx)
		require.NoError(t, err)

		count, err := client.Build.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})

	t.Run("BuildWithoutInvocation", func(t *testing.T) {
		client := setupTestDB(t)
		_, err := client.Build.Create().SetBuildURL("1").SetBuildUUID(uuid.New()).SetInstanceName("").SetTimestamp(time.Now()).Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(client, clock, traceProvider)
		require.NoError(t, err)
		err = cleanup.RemoveBuildsWithoutInvocations(ctx)
		require.NoError(t, err)

		count, err := client.Build.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})

	t.Run("MultipleBuildsMixed", func(t *testing.T) {
		client := setupTestDB(t)
		// Build with invocation
		buildWithInv, err := client.Build.Create().SetBuildURL("1").SetBuildUUID(uuid.New()).SetInstanceName("").SetTimestamp(time.Now()).Save(ctx)
		require.NoError(t, err)
		_, err = client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetBuild(buildWithInv).
			Save(ctx)
		require.NoError(t, err)
		// Build without invocation
		_, err = client.Build.Create().SetBuildURL("2").SetBuildUUID(uuid.New()).SetInstanceName("").SetTimestamp(time.Now()).Save(ctx)
		require.NoError(t, err)
		// Another build without invocation
		_, err = client.Build.Create().SetBuildURL("3").SetBuildUUID(uuid.New()).SetInstanceName("").SetTimestamp(time.Now()).Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(client, clock, traceProvider)
		require.NoError(t, err)
		err = cleanup.RemoveBuildsWithoutInvocations(ctx)
		require.NoError(t, err)

		count, err := client.Build.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})
}
