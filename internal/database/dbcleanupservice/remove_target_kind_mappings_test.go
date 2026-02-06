package dbcleanupservice_test

import (
	"context"
	"testing"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/internal/database/dbauthservice"
	"github.com/buildbarn/bb-portal/internal/mock"
	"github.com/buildbarn/bb-portal/test/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/mock/gomock"
)

func createTestTarget(t *testing.T, ctx context.Context, client *ent.Client, targetLabel string, instanceNameDbID int64) *ent.Target {
	target, err := client.Target.Create().
		SetInstanceNameID(instanceNameDbID).
		SetLabel(targetLabel).
		SetAspect("").
		SetTargetKind("testkind").
		Save(ctx)
	require.NoError(t, err)
	return target
}

func TestRemoveTargetKindMappings(t *testing.T) {
	ctrl, ctx := gomock.WithContext(context.Background(), t)
	ctx = dbauthservice.NewContextWithDbAuthServiceBypass(ctx)
	clock := mock.NewMockClock(ctrl)
	traceProvider := noop.NewTracerProvider()

	t.Run("NoInvocations", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		client := db.Ent()

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		err = cleanup.RemoveTargetKindMappings(ctx)
		require.NoError(t, err)

		count, err := client.TargetKindMapping.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})

	t.Run("UnfinishedInvocation", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		client := db.Ent()

		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		target := createTestTarget(t, ctx, client, "targetName", instanceNameDbID)

		invocation, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetBepCompleted(false).
			Save(ctx)
		require.NoError(t, err)

		_, err = client.TargetKindMapping.Create().
			SetBazelInvocation(invocation).
			SetTarget(target).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		err = cleanup.RemoveTargetKindMappings(ctx)
		require.NoError(t, err)

		count, err := client.TargetKindMapping.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})

	t.Run("NoTargetKindMappings", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		client := db.Ent()

		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		_, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetBepCompleted(true).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		err = cleanup.RemoveTargetKindMappings(ctx)
		require.NoError(t, err)

		count, err := client.TargetKindMapping.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})

	t.Run("SingleTargetKindMapping", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		client := db.Ent()

		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		target := createTestTarget(t, ctx, client, "targetName", instanceNameDbID)

		invocation, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetBepCompleted(true).
			Save(ctx)
		require.NoError(t, err)

		_, err = client.TargetKindMapping.Create().
			SetBazelInvocation(invocation).
			SetTarget(target).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		err = cleanup.RemoveTargetKindMappings(ctx)
		require.NoError(t, err)

		count, err := client.TargetKindMapping.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})

	t.Run("MultipleTargetKindMappings", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		client := db.Ent()

		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		target1 := createTestTarget(t, ctx, client, "targetName1", instanceNameDbID)
		target2 := createTestTarget(t, ctx, client, "targetName2", instanceNameDbID)

		invocation1, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetBepCompleted(true).
			Save(ctx)
		require.NoError(t, err)

		invocation2, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetBepCompleted(true).
			Save(ctx)
		require.NoError(t, err)

		_, err = client.TargetKindMapping.Create().
			SetBazelInvocation(invocation1).
			SetTarget(target1).
			Save(ctx)
		require.NoError(t, err)

		_, err = client.TargetKindMapping.Create().
			SetBazelInvocation(invocation2).
			SetTarget(target1).
			Save(ctx)
		require.NoError(t, err)

		_, err = client.TargetKindMapping.Create().
			SetBazelInvocation(invocation2).
			SetTarget(target2).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		err = cleanup.RemoveTargetKindMappings(ctx)
		require.NoError(t, err)

		count, err := client.TargetKindMapping.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})
}
