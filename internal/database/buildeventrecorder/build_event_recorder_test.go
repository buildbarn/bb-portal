package buildeventrecorder_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/buildbarn/bb-portal/internal/database/buildeventrecorder"
	"github.com/buildbarn/bb-portal/internal/database/dbauthservice"
	"github.com/buildbarn/bb-portal/internal/database/embedded"
	"github.com/buildbarn/bb-portal/test/testutils"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	_ "github.com/buildbarn/bb-portal/ent/gen/ent/runtime"
)

var dbProvider *embedded.DatabaseProvider

func TestMain(m *testing.M) {
	var err error
	dbProvider, err = embedded.NewDatabaseProvider(os.Stderr)
	if err != nil {
		// Embedded postgres needs Bazel runfiles to find its binaries; under
		// plain `go test`, this fails. Continue so tests that don't need the
		// DB still run; DB-using tests will skip themselves below.
		fmt.Fprintf(os.Stderr, "WARN: embedded DB unavailable; DB-using tests will skip: %v\n", err)
		dbProvider = nil
	} else {
		defer dbProvider.Cleanup()
	}
	os.Exit(m.Run())
}

// skipIfNoDB skips the test if the embedded postgres provider failed to
// initialize (typically when running `go test` outside Bazel runfiles).
func skipIfNoDB(t *testing.T) {
	t.Helper()
	if dbProvider == nil {
		t.Skip("embedded postgres unavailable; needs Bazel runfiles")
	}
}

func TestFindOrCreateInstanceName(t *testing.T) {
	skipIfNoDB(t)
	ctx := context.Background()

	t.Run("CreatesAndReturnsIdempotently", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)

		firstID, err := buildeventrecorder.FindOrCreateInstanceName(ctx, db, "testInstance")
		require.NoError(t, err)
		secondID, err := buildeventrecorder.FindOrCreateInstanceName(ctx, db, "testInstance")
		require.NoError(t, err)
		require.Equal(t, firstID, secondID)
	})
}

func TestFindOrCreateInvocation(t *testing.T) {
	skipIfNoDB(t)
	ctx := dbauthservice.NewContextWithDbAuthServiceBypass(context.Background())

	invocationsGaugeVec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "test_bazel_invocations_total",
			Help: "Test gauge for bazel invocations.",
		},
		[]string{"user_type"},
	)

	t.Run("ReconnectExistingUnlocked", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		instanceNameDbID, err := buildeventrecorder.FindOrCreateInstanceName(ctx, db, "testInstance")
		require.NoError(t, err)
		invocationId := uuid.New()

		firstID, err := buildeventrecorder.FindOrCreateInvocation(
			ctx, db, invocationId, instanceNameDbID, nil, invocationsGaugeVec,
		)
		require.NoError(t, err)
		secondID, err := buildeventrecorder.FindOrCreateInvocation(
			ctx, db, invocationId, instanceNameDbID, nil, invocationsGaugeVec,
		)
		require.NoError(t, err, "Failed to reconnect to unlocked invocation")
		require.Equal(t, firstID, secondID)
	})

	t.Run("FailToReconnectLockedInvocation", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		client := db.Ent()
		instanceNameDbID, err := buildeventrecorder.FindOrCreateInstanceName(ctx, db, "testInstance")
		require.NoError(t, err)
		invocationID := uuid.New()

		id, err := buildeventrecorder.FindOrCreateInvocation(
			ctx, db, invocationID, instanceNameDbID, nil, invocationsGaugeVec,
		)
		require.NoError(t, err)
		err = client.BazelInvocation.UpdateOneID(id).SetBepCompleted(true).Exec(ctx)
		require.NoError(t, err)
		_, err = buildeventrecorder.FindOrCreateInvocation(
			ctx, db, invocationID, instanceNameDbID, nil, invocationsGaugeVec,
		)
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok, "Expected a gRPC status error")
		require.Equal(t, codes.FailedPrecondition, st.Code())
		require.Contains(t, err.Error(), "locked for writing")
	})
}
