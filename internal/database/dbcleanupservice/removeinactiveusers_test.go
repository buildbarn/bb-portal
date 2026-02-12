package dbcleanupservice_test

import (
	"context"
	"testing"
	"time"

	"github.com/buildbarn/bb-portal/internal/database/dbauthservice"
	"github.com/buildbarn/bb-portal/internal/mock"
	"github.com/buildbarn/bb-portal/test/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/mock/gomock"
)

func TestRemoveInactiveUser(t *testing.T) {
	ctrl, ctx := gomock.WithContext(context.Background(), t)
	ctx = dbauthservice.NewContextWithDbAuthServiceBypass(ctx)
	clock := mock.NewMockClock(ctrl)
	traceProvider := noop.NewTracerProvider()

	t.Run("NoUsers", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		client := db.Ent()

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		err = cleanup.RemoveInactiveUsers(ctx)
		require.NoError(t, err)

		count, err := client.AuthenticatedUser.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})

	t.Run("InactiveUser", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		client := db.Ent()

		_, err := client.AuthenticatedUser.
			Create().
			SetExternalID(uuid.NewString()).
			SetUserUUID(uuid.New()).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		err = cleanup.RemoveInactiveUsers(ctx)
		require.NoError(t, err)

		count, err := client.AuthenticatedUser.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})

	t.Run("ActiveUser", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		client := db.Ent()

		user, err := client.AuthenticatedUser.
			Create().
			SetExternalID(uuid.NewString()).
			SetUserUUID(uuid.New()).
			Save(ctx)
		require.NoError(t, err)
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		_, err = client.BazelInvocation.
			Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetStartedAt(time.Now()).
			SetAuthenticatedUser(user).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		err = cleanup.RemoveInactiveUsers(ctx)
		require.NoError(t, err)

		count, err := client.AuthenticatedUser.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})

	t.Run("InactiveAndActiveUser", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		client := db.Ent()

		// Inactive user
		_, err := client.AuthenticatedUser.
			Create().
			SetExternalID(uuid.NewString()).
			SetUserUUID(uuid.New()).
			Save(ctx)
		require.NoError(t, err)

		activeUser, err := client.AuthenticatedUser.
			Create().
			SetExternalID(uuid.NewString()).
			SetUserUUID(uuid.New()).
			Save(ctx)
		require.NoError(t, err)
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		_, err = client.BazelInvocation.
			Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetStartedAt(time.Now()).
			SetAuthenticatedUser(activeUser).
			Save(ctx)
		require.NoError(t, err)

		count, err := client.AuthenticatedUser.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 2, count)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		err = cleanup.RemoveInactiveUsers(ctx)
		require.NoError(t, err)

		count, err = client.AuthenticatedUser.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})
}
