package authschema_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-portal/internal/database/dbauthservice"
	"github.com/buildbarn/bb-portal/internal/database/embedded"
	"github.com/buildbarn/bb-portal/internal/mock"
	"github.com/buildbarn/bb-storage/pkg/auth"
	"github.com/buildbarn/bb-storage/pkg/jmespath"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	// Needed to avoid cyclic dependencies in ent (https://entgo.io/docs/privacy#privacy-policy-registration)
	"github.com/buildbarn/bb-portal/ent/gen/ent/invocationtarget"
	_ "github.com/buildbarn/bb-portal/ent/gen/ent/runtime"
)

var dbProvider *embedded.DatabaseProvider

func TestMain(m *testing.M) {
	var err error
	tmpDir, err := os.MkdirTemp(os.TempDir(), "embedded_db_test")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create temp dir: %v\n", err)
		os.Exit(1)
	}

	dbProvider, err = embedded.NewDatabaseProvider(tmpDir, os.Stderr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not start embedded DB: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		dbProvider.Cleanup()
		os.RemoveAll(tmpDir)
	}()

	code := m.Run()
	os.Exit(code)
}

func setupTestDB(t testing.TB) database.Client {
	conn, err := dbProvider.CreateDatabase()
	require.NoError(t, err)
	db, err := database.New("postgres", conn)
	require.NoError(t, err)
	t.Cleanup(func() { conn.Close() })
	err = db.Ent().Schema.Create(context.Background())
	require.NoError(t, err)
	return db
}

func TestPrivacy(t *testing.T) {
	ctrl, ctx := gomock.WithContext(context.Background(), t)
	clock := mock.NewMockClock(ctrl)
	db := setupTestDB(t).Ent()

	authorizer := auth.NewJMESPathExpressionAuthorizer(
		jmespath.MustCompile("instanceName == 'allowed1' || instanceName == 'allowed2'"),
	)

	dbAuthService := dbauthservice.NewDbAuthService(db, clock, authorizer, time.Second*5)
	ctx = dbauthservice.NewContextWithDbAuthService(ctx, dbAuthService)

	t.Run("EmptyDatabase", func(t *testing.T) {
		clock.EXPECT().Now().Return(time.Unix(10000000, 0)).Times(5)

		invocations, err := db.BazelInvocation.Query().IDs(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, len(invocations))

		users, err := db.AuthenticatedUser.Query().IDs(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, len(users))

		builds, err := db.Build.Query().IDs(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, len(builds))

		targets, err := db.Target.Query().IDs(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, len(targets))

		testSummaries, err := db.TestSummary.Query().IDs(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, len(testSummaries))
	})

	clock.EXPECT().Now().Return(time.Unix(20000000, 0)).Times(5)

	deniedInstance, err := db.InstanceName.Create().SetName("denied").Save(ctx)
	require.NoError(t, err)
	deniedInvocation, err := db.BazelInvocation.Create().SetInstanceName(deniedInstance).SetInvocationID(uuid.New()).Save(ctx)
	require.NoError(t, err)
	_, err = db.AuthenticatedUser.Create().SetUserUUID(uuid.New()).SetExternalID("denied_user").AddBazelInvocations(deniedInvocation).Save(ctx)
	require.NoError(t, err)
	_, err = db.Build.Create().SetInstanceName(deniedInstance).SetBuildUUID(uuid.New()).SetBuildURL("denied_build_url").SetTimestamp(time.Now()).Save(ctx)
	require.NoError(t, err)
	deniedTarget, err := db.Target.Create().SetInstanceName(deniedInstance).SetLabel("denied").SetAspect("aspect").SetTargetKind("targetKind").Save(ctx)
	require.NoError(t, err)
	deniedInvocationTarget, err := db.InvocationTarget.Create().SetBazelInvocation(deniedInvocation).SetTarget(deniedTarget).SetAbortReason(invocationtarget.AbortReasonNONE).Save(ctx)
	require.NoError(t, err)

	_, err = db.TestSummary.Create().SetInvocationTarget(deniedInvocationTarget).Save(ctx)
	require.NoError(t, err)

	t.Run("PopulatedDatabaseWithDeniedInstance", func(t *testing.T) {
		clock.EXPECT().Now().Return(time.Unix(30000000, 0)).Times(5)

		invocations, err := db.BazelInvocation.Query().IDs(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, len(invocations))

		users, err := db.AuthenticatedUser.Query().IDs(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, len(users))

		builds, err := db.Build.Query().IDs(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, len(builds))

		targets, err := db.Target.Query().IDs(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, len(targets))

		testSummaries, err := db.TestSummary.Query().IDs(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, len(testSummaries))
	})

	clock.EXPECT().Now().Return(time.Unix(40000000, 0)).Times(5)

	allowed1Instance, err := db.InstanceName.Create().SetName("allowed1").Save(ctx)
	require.NoError(t, err)
	allowed1Invocation, err := db.BazelInvocation.Create().SetInstanceName(allowed1Instance).SetInvocationID(uuid.New()).Save(ctx)
	require.NoError(t, err)
	allowed1User, err := db.AuthenticatedUser.Create().SetUserUUID(uuid.New()).SetExternalID("allowed1_user").AddBazelInvocations(allowed1Invocation).Save(ctx)
	require.NoError(t, err)
	allowed1Build, err := db.Build.Create().SetInstanceName(allowed1Instance).SetBuildUUID(uuid.New()).SetBuildURL("allowed1_build_url").SetTimestamp(time.Now()).Save(ctx)
	require.NoError(t, err)
	allowed1Target, err := db.Target.Create().SetInstanceName(allowed1Instance).SetLabel("allowed1").SetAspect("aspect").SetTargetKind("targetKind").Save(ctx)
	require.NoError(t, err)
	allowed1InvocationTarget, err := db.InvocationTarget.Create().SetBazelInvocation(allowed1Invocation).SetTarget(allowed1Target).SetAbortReason(invocationtarget.AbortReasonNONE).Save(ctx)
	require.NoError(t, err)
	allowed1TestSummary, err := db.TestSummary.Create().SetInvocationTarget(allowed1InvocationTarget).Save(ctx)
	require.NoError(t, err)

	t.Run("PopulatedDatabase", func(t *testing.T) {
		clock.EXPECT().Now().Return(time.Unix(50000000, 0)).Times(5)

		invocations, err := db.BazelInvocation.Query().IDs(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, len(invocations))
		require.Contains(t, invocations, allowed1Invocation.ID)

		users, err := db.AuthenticatedUser.Query().IDs(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, len(users))
		require.Contains(t, users, allowed1User.ID)

		builds, err := db.Build.Query().IDs(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, len(builds))
		require.Contains(t, builds, allowed1Build.ID)

		targets, err := db.Target.Query().IDs(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, len(targets))
		require.Contains(t, targets, allowed1Target.ID)

		testSummaries, err := db.TestSummary.Query().IDs(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, len(testSummaries))
		require.Contains(t, testSummaries, allowed1TestSummary.ID)
	})

	clock.EXPECT().Now().Return(time.Unix(50000000, 0)).Times(5)

	allowed2Instance, err := db.InstanceName.Create().SetName("allowed2").Save(ctx)
	require.NoError(t, err)
	allowed2Invocation, err := db.BazelInvocation.Create().SetInstanceName(allowed2Instance).SetInvocationID(uuid.New()).Save(ctx)
	require.NoError(t, err)
	allowed2User, err := db.AuthenticatedUser.Create().SetUserUUID(uuid.New()).SetExternalID("allowed2_user").AddBazelInvocations(allowed2Invocation).Save(ctx)
	require.NoError(t, err)
	allowed2Build, err := db.Build.Create().SetInstanceName(allowed2Instance).SetBuildUUID(uuid.New()).SetBuildURL("allowed2_build_url").SetTimestamp(time.Now()).Save(ctx)
	require.NoError(t, err)
	allowed2Target, err := db.Target.Create().SetInstanceName(allowed2Instance).SetLabel("allowed2").SetAspect("aspect").SetTargetKind("targetKind").Save(ctx)
	require.NoError(t, err)
	allowed2InvocationTarget, err := db.InvocationTarget.Create().SetBazelInvocation(allowed2Invocation).SetTarget(allowed2Target).SetAbortReason(invocationtarget.AbortReasonNONE).Save(ctx)
	require.NoError(t, err)
	allowed2TestSummary, err := db.TestSummary.Create().SetInvocationTarget(allowed2InvocationTarget).Save(ctx)
	require.NoError(t, err)

	t.Run("PopulatedDatabase2", func(t *testing.T) {
		clock.EXPECT().Now().Return(time.Unix(60000000, 0)).Times(5)

		invocations, err := db.BazelInvocation.Query().IDs(ctx)
		require.NoError(t, err)
		require.Equal(t, 2, len(invocations))
		require.Contains(t, invocations, allowed1Invocation.ID)
		require.Contains(t, invocations, allowed2Invocation.ID)

		users, err := db.AuthenticatedUser.Query().IDs(ctx)
		require.NoError(t, err)
		require.Equal(t, 2, len(users))
		require.Contains(t, users, allowed1User.ID)
		require.Contains(t, users, allowed2User.ID)

		builds, err := db.Build.Query().IDs(ctx)
		require.NoError(t, err)
		require.Equal(t, 2, len(builds))
		require.Contains(t, builds, allowed1Build.ID)
		require.Contains(t, builds, allowed2Build.ID)

		targets, err := db.Target.Query().IDs(ctx)
		require.NoError(t, err)
		require.Equal(t, 2, len(targets))
		require.Contains(t, targets, allowed1Target.ID)
		require.Contains(t, targets, allowed2Target.ID)

		testSummaries, err := db.TestSummary.Query().IDs(ctx)
		require.NoError(t, err)
		require.Equal(t, 2, len(testSummaries))
		require.Contains(t, testSummaries, allowed1TestSummary.ID)
		require.Contains(t, testSummaries, allowed2TestSummary.ID)
	})
}
