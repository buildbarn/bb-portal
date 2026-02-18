package dbcleanupservice_test

import (
	"context"
	"testing"
	"time"

	"github.com/buildbarn/bb-portal/internal/database/dbauthservice"
	"github.com/buildbarn/bb-portal/internal/mock"
	"github.com/buildbarn/bb-portal/test/testutils"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/mock/gomock"
)

func TestRemoveOldInvocations(t *testing.T) {
	ctrl, ctx := gomock.WithContext(context.Background(), t)
	ctx = dbauthservice.NewContextWithDbAuthServiceBypass(ctx)
	clock := mock.NewMockClock(ctrl)
	cleanupTime := time.Unix(1600000200, 0)
	traceProvider := noop.NewTracerProvider()

	t.Run("NoInvocations", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		client := db.Ent()

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.RemoveOldInvocations(ctx)
		require.NoError(t, err)

		count, err := client.BazelInvocation.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})

	t.Run("InvocationNotCompleted", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		client := db.Ent()
		instanceName := testutils.CreateInstanceName(ctx, t, client, "testInstance")
		_, err := testutils.StartCreateInvocation(client, instanceName).
			SetCreatedTimestamp(cleanupTime.Add(-60 * time.Minute)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.RemoveOldInvocations(ctx)
		require.NoError(t, err)

		count, err := client.BazelInvocation.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})

	t.Run("InvocationCompletedButNotOld", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		client := db.Ent()
		instanceName := testutils.CreateInstanceName(ctx, t, client, "testInstance")
		_, err := testutils.StartCreateInvocation(client, instanceName).
			SetBepCompleted(true).
			SetCreatedTimestamp(cleanupTime.Add(-15 * time.Minute)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.RemoveOldInvocations(ctx)
		require.NoError(t, err)

		count, err := client.BazelInvocation.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})

	t.Run("InvocationCompletedAndOld", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		client := db.Ent()
		instanceName := testutils.CreateInstanceName(ctx, t, client, "testInstance")
		_, err := testutils.StartCreateInvocation(client, instanceName).
			SetBepCompleted(true).
			SetCreatedTimestamp(cleanupTime.Add(-60 * time.Minute)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)

		err = cleanup.RemoveOldInvocations(ctx)
		require.NoError(t, err)

		count, err := client.BazelInvocation.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})

	t.Run("MultipleInvocationsMixed", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		client := db.Ent()
		instanceName := testutils.CreateInstanceName(ctx, t, client, "testInstance")
		// Old and completed
		_, err := testutils.StartCreateInvocation(client, instanceName).
			SetBepCompleted(true).
			SetCreatedTimestamp(cleanupTime.Add(-60 * time.Minute)).
			Save(ctx)
		require.NoError(t, err)
		// Not completed
		_, err = testutils.StartCreateInvocation(client, instanceName).
			SetCreatedTimestamp(cleanupTime.Add(-60 * time.Minute)).
			Save(ctx)
		require.NoError(t, err)
		// Completed but not old
		_, err = testutils.StartCreateInvocation(client, instanceName).
			SetBepCompleted(true).
			SetCreatedTimestamp(cleanupTime.Add(-15 * time.Minute)).
			Save(ctx)
		require.NoError(t, err)
		// Not completed and not old
		_, err = testutils.StartCreateInvocation(client, instanceName).
			SetCreatedTimestamp(cleanupTime.Add(-15 * time.Minute)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.RemoveOldInvocations(ctx)
		require.NoError(t, err)

		count, err := client.BazelInvocation.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 3, count)
	})
}
