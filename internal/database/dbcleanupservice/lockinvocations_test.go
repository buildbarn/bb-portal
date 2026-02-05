package dbcleanupservice_test

import (
	"context"
	"testing"
	"time"

	// Needed to avoid cyclic dependencies in ent (https://entgo.io/docs/privacy#privacy-policy-registration)

	_ "github.com/buildbarn/bb-portal/ent/gen/ent/runtime"
	"github.com/buildbarn/bb-portal/internal/database/dbauthservice"
	"github.com/buildbarn/bb-portal/internal/mock"
	"github.com/buildbarn/bb-portal/test/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/mock/gomock"
)

func TestLockInvocationsWithNoRecentEvents(t *testing.T) {
	ctrl, ctx := gomock.WithContext(context.Background(), t)
	ctx = dbauthservice.NewContextWithDbAuthServiceBypass(ctx)
	clock := mock.NewMockClock(ctrl)
	cleanupTime := time.Unix(1600000060, 0)
	traceProvider := noop.NewTracerProvider()

	t.Run("NoInvocations-NoEventMetadata", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentEvents(ctx)
		require.NoError(t, err)
	})

	t.Run("UnfinishedInvocation-NoEventMetadata", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		startInv, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
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
		db := testutils.SetupTestDB(t, dbProvider)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		startInv, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetBepCompleted(true).
			SetEndedAt(cleanupTime.Add(-60 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentEvents(ctx)
		require.NoError(t, err)

		resInv, err := client.BazelInvocation.Get(ctx, startInv.ID)
		require.NoError(t, err)
		require.Equal(t, startInv.EndedAt.UTC(), resInv.EndedAt.UTC())
		require.True(t, resInv.BepCompleted)
	})

	t.Run("UnfinishedInvocation-RecentEventMetadata", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		invocationDb, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.EventMetadata.Create().
			SetBazelInvocationID(invocationDb.ID).
			SetHandled(roaringBytes(1)).
			SetVersion(0).
			SetEventReceivedAt(cleanupTime.Add(-20 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentEvents(ctx)
		require.NoError(t, err)

		inv, err := client.BazelInvocation.Get(ctx, invocationDb.ID)
		require.NoError(t, err)
		require.False(t, inv.BepCompleted)
	})

	t.Run("FinishedInvocation-RecentEventMetadata", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		invocationDb, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetBepCompleted(true).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.EventMetadata.Create().
			SetBazelInvocationID(invocationDb.ID).
			SetHandled(roaringBytes(1)).
			SetVersion(0).
			SetEventReceivedAt(cleanupTime.Add(-20 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentEvents(ctx)
		require.NoError(t, err)

		inv, err := client.BazelInvocation.Get(ctx, invocationDb.ID)
		require.NoError(t, err)
		require.True(t, inv.BepCompleted)
	})

	t.Run("UnfinishedInvocation-OldEventMetadata", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		invocationDb, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.EventMetadata.Create().
			SetBazelInvocationID(invocationDb.ID).
			SetHandled(roaringBytes(1)).
			SetVersion(0).
			SetEventReceivedAt(cleanupTime.Add(-70 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentEvents(ctx)
		require.NoError(t, err)

		inv, err := client.BazelInvocation.Get(ctx, invocationDb.ID)
		require.NoError(t, err)
		require.True(t, inv.BepCompleted)
	})

	t.Run("FinishedInvocation-OldEventMetadata", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		invocationDb, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetBepCompleted(true).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.EventMetadata.Create().
			SetBazelInvocationID(invocationDb.ID).
			SetHandled(roaringBytes(1)).
			SetVersion(0).
			SetEventReceivedAt(cleanupTime.Add(-70 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentEvents(ctx)
		require.NoError(t, err)

		inv, err := client.BazelInvocation.Get(ctx, invocationDb.ID)
		require.NoError(t, err)
		require.True(t, inv.BepCompleted)
	})

	t.Run("MultipleMixedInvocations", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")

		// old event metadata -> should be locked
		invOld, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.EventMetadata.Create().
			SetBazelInvocationID(invOld.ID).
			SetHandled(roaringBytes(1)).
			SetVersion(0).
			SetEventReceivedAt(cleanupTime.Add(-70 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		// recent event metadata -> should remain unlocked
		invRecent, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.EventMetadata.Create().
			SetBazelInvocationID(invRecent.ID).
			SetHandled(roaringBytes(1)).
			SetVersion(0).
			SetEventReceivedAt(cleanupTime.Add(-20 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)

		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentEvents(ctx)
		require.NoError(t, err)

		gotOld, err := client.BazelInvocation.Get(ctx, invOld.ID)
		require.NoError(t, err)
		require.True(t, gotOld.BepCompleted)

		gotRecent, err := client.BazelInvocation.Get(ctx, invRecent.ID)
		require.NoError(t, err)
		require.False(t, gotRecent.BepCompleted)
	})
}

func TestUpdateInvocationEndedAtFromEvents(t *testing.T) {
	ctrl, ctx := gomock.WithContext(context.Background(), t)
	ctx = dbauthservice.NewContextWithDbAuthServiceBypass(ctx)
	clock := mock.NewMockClock(ctrl)
	traceProvider := noop.NewTracerProvider()

	t.Run("NoInvocations", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		err = cleanup.UpdateInvocationEndedAtFromEvents(ctx)
		require.NoError(t, err)
	})

	t.Run("NoEndedAt-NoEventMetadata", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		startInv, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetBepCompleted(true).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		err = cleanup.UpdateInvocationEndedAtFromEvents(ctx)
		require.NoError(t, err)

		resInv, err := client.BazelInvocation.Get(ctx, startInv.ID)
		require.NoError(t, err)
		require.Nil(t, resInv.EndedAt)
	})

	t.Run("WithEndedAt-NoEventMetadata", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		invEndedAt := time.Unix(1600000000, 0)
		startInv, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetBepCompleted(true).
			SetEndedAt(invEndedAt).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		err = cleanup.UpdateInvocationEndedAtFromEvents(ctx)
		require.NoError(t, err)

		inv, err := client.BazelInvocation.Get(ctx, startInv.ID)
		require.NoError(t, err)
		// Should be unchanged
		require.Equal(t, invEndedAt.UTC(), inv.EndedAt.UTC())
	})

	t.Run("NoEndedAt-WithEventMetadata", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		eventEndedAt := time.Unix(1600000100, 0)
		invocationDb, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetBepCompleted(true).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.EventMetadata.Create().
			SetBazelInvocationID(invocationDb.ID).
			SetHandled(roaringBytes(1)).
			SetVersion(0).
			SetEventReceivedAt(eventEndedAt).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		err = cleanup.UpdateInvocationEndedAtFromEvents(ctx)
		require.NoError(t, err)

		inv, err := client.BazelInvocation.Get(ctx, invocationDb.ID)
		require.NoError(t, err)
		// Should be updated
		require.Equal(t, eventEndedAt.UTC(), inv.EndedAt.UTC())
	})

	t.Run("WithEndedAt-WithEventMetadata", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		invEndedAt := time.Unix(1600000000, 0)
		eventEndedAt := time.Unix(1600000100, 0)
		invocationDb, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetBepCompleted(true).
			SetEndedAt(invEndedAt).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.EventMetadata.Create().
			SetBazelInvocationID(invocationDb.ID).
			SetHandled(roaringBytes(1)).
			SetVersion(0).
			SetEventReceivedAt(eventEndedAt).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		err = cleanup.UpdateInvocationEndedAtFromEvents(ctx)
		require.NoError(t, err)

		inv, err := client.BazelInvocation.Get(ctx, invocationDb.ID)
		require.NoError(t, err)
		// Should be unchanged
		require.Equal(t, invEndedAt.UTC(), inv.EndedAt.UTC())
	})

	t.Run("MultipleMixedInvocations", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")

		// No EndedAt, With EventMetadata -> should be updated
		invToUpdate, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetBepCompleted(true).
			Save(ctx)
		require.NoError(t, err)
		eventEndedAt := time.Unix(1600000200, 0)
		_, err = client.EventMetadata.Create().
			SetBazelInvocationID(invToUpdate.ID).
			SetHandled(roaringBytes(1)).
			SetVersion(0).
			SetEventReceivedAt(eventEndedAt).
			Save(ctx)
		require.NoError(t, err)

		// With EndedAt, With EventMetadata -> should remain unchanged
		invUnchanged, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetBepCompleted(true).
			SetEndedAt(time.Unix(1600000000, 0)).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.EventMetadata.Create().
			SetBazelInvocationID(invUnchanged.ID).
			SetHandled(roaringBytes(1)).
			SetVersion(0).
			SetEventReceivedAt(time.Unix(1600000300, 0)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		err = cleanup.UpdateInvocationEndedAtFromEvents(ctx)
		require.NoError(t, err)

		gotUpdated, err := client.BazelInvocation.Get(ctx, invToUpdate.ID)
		require.NoError(t, err)
		require.Equal(t, eventEndedAt.UTC(), gotUpdated.EndedAt.UTC())

		gotUnchanged, err := client.BazelInvocation.Get(ctx, invUnchanged.ID)
		require.NoError(t, err)
		require.Equal(t, time.Unix(1600000000, 0).UTC(), gotUnchanged.EndedAt.UTC())
	})

	t.Run("DontUpdateUnfinishedInvocations", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")

		invToNotUpdate, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			Save(ctx)
		require.NoError(t, err)
		eventEndedAt := time.Unix(1600000200, 0)
		_, err = client.EventMetadata.Create().
			SetBazelInvocationID(invToNotUpdate.ID).
			SetHandled(roaringBytes(1)).
			SetVersion(0).
			SetEventReceivedAt(eventEndedAt).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		err = cleanup.UpdateInvocationEndedAtFromEvents(ctx)
		require.NoError(t, err)

		gotNotUpdated, err := client.BazelInvocation.Get(ctx, invToNotUpdate.ID)
		require.NoError(t, err)
		require.Nil(t, gotNotUpdated.EndedAt)
	})
}
