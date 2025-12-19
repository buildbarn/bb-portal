package dbcleanupservice_test

import (
	"context"
	"testing"

	"github.com/buildbarn/bb-portal/ent/gen/ent/invocationtarget"
	"github.com/buildbarn/bb-portal/internal/database/dbauthservice"
	"github.com/buildbarn/bb-portal/internal/mock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/mock/gomock"
)

func TestRemoveUnusedTargets(t *testing.T) {
	ctrl, ctx := gomock.WithContext(context.Background(), t)
	ctx = dbauthservice.NewContextWithDbAuthServiceBypass(ctx)
	clock := mock.NewMockClock(ctrl)
	traceProvider := noop.NewTracerProvider()

	t.Run("NoTargets", func(t *testing.T) {
		db := setupTestDB(t)
		client := db.Ent()

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		err = cleanup.RemoveUnusedTargets(ctx)
		require.NoError(t, err)

		count, err := client.Target.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})

	t.Run("UnusedTarget", func(t *testing.T) {
		db := setupTestDB(t)
		client := db.Ent()

		instanceName, err := client.InstanceName.Create().
			SetName("instance").
			Save(ctx)
		require.NoError(t, err)

		_, err = client.Target.Create().
			SetInstanceName(instanceName).
			SetLabel("testLabel").
			SetAspect("testAspect").
			SetTargetKind("testTargetKind").
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		err = cleanup.RemoveUnusedTargets(ctx)
		require.NoError(t, err)

		count, err := client.Target.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})

	t.Run("TargetWithInvocationTarget", func(t *testing.T) {
		db := setupTestDB(t)
		client := db.Ent()

		instanceName, err := client.InstanceName.Create().
			SetName("instance").
			Save(ctx)
		require.NoError(t, err)

		target, err := client.Target.Create().
			SetInstanceName(instanceName).
			SetLabel("testLabel").
			SetAspect("testAspect").
			SetTargetKind("testTargetKind").
			Save(ctx)
		require.NoError(t, err)

		invocation, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceName(instanceName).
			Save(ctx)
		require.NoError(t, err)

		_, err = client.InvocationTarget.Create().
			SetBazelInvocation(invocation).
			SetTarget(target).
			SetAbortReason(invocationtarget.AbortReasonNONE).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		err = cleanup.RemoveUnusedTargets(ctx)
		require.NoError(t, err)

		count, err := client.Target.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})

	t.Run("TargetWithTargetKindMapping", func(t *testing.T) {
		db := setupTestDB(t)
		client := db.Ent()

		instanceName, err := client.InstanceName.Create().
			SetName("instance").
			Save(ctx)
		require.NoError(t, err)

		target, err := client.Target.Create().
			SetInstanceName(instanceName).
			SetLabel("testLabel").
			SetAspect("testAspect").
			SetTargetKind("testTargetKind").
			Save(ctx)
		require.NoError(t, err)

		invocation, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceName(instanceName).
			Save(ctx)
		require.NoError(t, err)

		_, err = client.TargetKindMapping.Create().
			SetBazelInvocation(invocation).
			SetTarget(target).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		err = cleanup.RemoveUnusedTargets(ctx)
		require.NoError(t, err)

		count, err := client.Target.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})

	t.Run("TargetWithInvocationTargetAndTargetKindMapping", func(t *testing.T) {
		db := setupTestDB(t)
		client := db.Ent()

		instanceName, err := client.InstanceName.Create().
			SetName("instance").
			Save(ctx)
		require.NoError(t, err)

		target, err := client.Target.Create().
			SetInstanceName(instanceName).
			SetLabel("testLabel").
			SetAspect("testAspect").
			SetTargetKind("testTargetKind").
			Save(ctx)
		require.NoError(t, err)

		invocation, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceName(instanceName).
			Save(ctx)
		require.NoError(t, err)

		_, err = client.InvocationTarget.Create().
			SetBazelInvocation(invocation).
			SetTarget(target).
			SetAbortReason(invocationtarget.AbortReasonNONE).
			Save(ctx)
		require.NoError(t, err)

		_, err = client.TargetKindMapping.Create().
			SetBazelInvocation(invocation).
			SetTarget(target).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		err = cleanup.RemoveUnusedTargets(ctx)
		require.NoError(t, err)

		count, err := client.Target.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})
}
