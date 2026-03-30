package dbcleanupservice_test

import (
	"context"
	"testing"

	"github.com/buildbarn/bb-portal/ent/gen/ent/testtarget"
	"github.com/buildbarn/bb-portal/internal/database/dbauthservice"
	"github.com/buildbarn/bb-portal/internal/mock"
	"github.com/buildbarn/bb-portal/test/testutils"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/mock/gomock"
)

func TestDbCleanupService_RemoveOrphanedTestTargets(t *testing.T) {
	ctrl, ctx := gomock.WithContext(context.Background(), t)
	ctx = dbauthservice.NewContextWithDbAuthServiceBypass(ctx)
	db := testutils.SetupTestDB(t, dbProvider)
	mockClock := mock.NewMockClock(ctrl)

	cleanupService, err := getNewDbCleanupService(db, mockClock, noop.NewTracerProvider())
	require.NoError(t, err)

	instance := testutils.CreateInstanceName(ctx, t, db.Ent(), "orphan_test_instance")
	invocation, err := testutils.StartCreateInvocation(db.Ent(), instance).Save(ctx)
	require.NoError(t, err)

	validTarget, err := db.Ent().Target.Create().SetInstanceName(instance).SetLabel("//app:valid").SetTargetKind("cc_test rule").SetAspect("").Save(ctx)
	require.NoError(t, err)

	validTestTarget, err := db.Ent().TestTarget.Create().SetTarget(validTarget).Save(ctx)
	require.NoError(t, err)

	validInvTarget, err := db.Ent().InvocationTarget.Create().SetBazelInvocation(invocation).SetTarget(validTarget).SetAbortReason("NONE").Save(ctx)
	require.NoError(t, err)

	_, err = db.Ent().TestSummary.Create().SetInvocationTarget(validInvTarget).SetOverallStatus("PASSED").Save(ctx)
	require.NoError(t, err)

	orphanedTarget, err := db.Ent().Target.Create().SetInstanceName(instance).SetLabel("//app:orphaned").SetTargetKind("cc_test rule").SetAspect("").Save(ctx)
	require.NoError(t, err)

	orphanedTestTarget, err := db.Ent().TestTarget.Create().SetTarget(orphanedTarget).Save(ctx)
	require.NoError(t, err)

	count, err := db.Ent().TestTarget.Query().Count(ctx)
	require.NoError(t, err)
	require.Equal(t, 2, count, "Both TestTargets should exist before cleanup")

	err = cleanupService.RemoveOrphanedTestTargets(ctx)
	require.NoError(t, err)

	remainingCount, err := db.Ent().TestTarget.Query().Count(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, remainingCount, "Exactly one TestTarget should remain")

	validExists, err := db.Ent().TestTarget.Query().Where(testtarget.IDEQ(validTestTarget.ID)).Exist(ctx)
	require.NoError(t, err)
	require.True(t, validExists, "The valid TestTarget must survive")

	orphanedExists, err := db.Ent().TestTarget.Query().Where(testtarget.IDEQ(orphanedTestTarget.ID)).Exist(ctx)
	require.NoError(t, err)
	require.False(t, orphanedExists, "The orphaned TestTarget must be deleted")
}
