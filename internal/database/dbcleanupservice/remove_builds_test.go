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

func TestRemoveBuildsWithoutInvocations(t *testing.T) {
	ctrl, ctx := gomock.WithContext(context.Background(), t)
	ctx = dbauthservice.NewContextWithDbAuthServiceBypass(ctx)
	clock := mock.NewMockClock(ctrl)
	traceProvider := noop.NewTracerProvider()

	t.Run("NoBuilds", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		client := db.Ent()

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		deleted, err := cleanup.RemoveBuildsWithoutInvocations(ctx)
		require.NoError(t, err)
		require.EqualValues(t, 0, deleted)

		count, err := client.Build.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})

	t.Run("BuildWithInvocation", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		client := db.Ent()
		instanceName := testutils.CreateInstanceName(ctx, t, client, "testInstance")

		buildObj, err := client.Build.Create().SetBuildUUID(uuid.New()).SetInstanceName(instanceName).SetTimestamp(time.Now().UTC()).Save(ctx)
		require.NoError(t, err)
		_, err = testutils.StartCreateInvocation(client, instanceName).
			SetBuild(buildObj).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		deleted, err := cleanup.RemoveBuildsWithoutInvocations(ctx)
		require.NoError(t, err)
		require.EqualValues(t, 0, deleted)

		count, err := client.Build.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})

	t.Run("BuildWithoutInvocation", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		client := db.Ent()
		instanceName := testutils.CreateInstanceName(ctx, t, client, "testInstance")

		_, err := client.Build.Create().SetBuildUUID(uuid.New()).SetInstanceName(instanceName).SetTimestamp(time.Now().UTC()).Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		deleted, err := cleanup.RemoveBuildsWithoutInvocations(ctx)
		require.NoError(t, err)
		require.EqualValues(t, 1, deleted)

		count, err := client.Build.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})

	t.Run("MultipleBuildsMixed", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		client := db.Ent()
		instanceName := testutils.CreateInstanceName(ctx, t, client, "testInstance")

		// Build with invocation
		buildWithInv, err := client.Build.Create().SetBuildUUID(uuid.New()).SetInstanceName(instanceName).SetTimestamp(time.Now().UTC()).Save(ctx)
		require.NoError(t, err)
		_, err = testutils.StartCreateInvocation(client, instanceName).
			SetBuild(buildWithInv).
			Save(ctx)
		require.NoError(t, err)
		// Build without invocation
		_, err = client.Build.Create().SetBuildUUID(uuid.New()).SetInstanceName(instanceName).SetTimestamp(time.Now().UTC()).Save(ctx)
		require.NoError(t, err)
		// Another build without invocation
		_, err = client.Build.Create().SetBuildUUID(uuid.New()).SetInstanceName(instanceName).SetTimestamp(time.Now().UTC()).Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		deleted, err := cleanup.RemoveBuildsWithoutInvocations(ctx)
		require.NoError(t, err)
		require.EqualValues(t, 2, deleted)

		count, err := client.Build.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})
}
