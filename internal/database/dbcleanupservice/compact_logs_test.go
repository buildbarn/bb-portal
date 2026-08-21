package dbcleanupservice_test

import (
	"context"
	"testing"

	// Needed to avoid cyclic dependencies in ent (https://entgo.io/docs/privacy#privacy-policy-registration)
	_ "github.com/buildbarn/bb-portal/ent/gen/ent/runtime"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/ent/gen/ent/buildlogchunk"
	"github.com/buildbarn/bb-portal/ent/gen/ent/incompletebuildlog"
	"github.com/buildbarn/bb-portal/internal/database/dbauthservice"
	"github.com/buildbarn/bb-portal/internal/mock"
	"github.com/buildbarn/bb-portal/test/testutils"
	"github.com/klauspost/compress/zstd"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/mock/gomock"
)

func populateIncompleteBuildLog(t *testing.T, ctx context.Context, client *ent.Client, invocationDbID int64) {
	logSnippets := []string{
		"\u001b[32mComputing main repo mapping:\u001b[0m \n\r\u001b[1A\u001b[K\u001b[32mLoading:\u001b[0m \n\r\u001b[1A\u001b[K\u001b[32mLoading:\u001b[0m 0 packages loaded\n",
		"\r\u001b[1A\u001b[K\u001b[35mWARNING: \u001b[0mBuild options --dynamic_mode, --extra_execution_platforms, and --extra_toolchains have changed, discarding analysis cache (this can be expensive, see https://bazel.build/advanced/performance/iteration-speed).\n\u001b[32mAnalyzing:\u001b[0m target //:hello (0 packages loaded)\n",
		"\r\u001b[1A\u001b[K\u001b[32mAnalyzing:\u001b[0m target //:hello (0 packages loaded, 0 targets configured)\n\r\u001b[1A\u001b[K\u001b[32mAnalyzing:\u001b[0m target //:hello (0 packages loaded, 0 targets configured)\n\n",
		"\r\u001b[1A\u001b[K\r\u001b[1A\u001b[K\u001b[32mINFO: \u001b[0mAnalyzed target //:hello (0 packages loaded, 2 targets configured).\n\n",
		"\r\u001b[1A\u001b[K\u001b[32mINFO: \u001b[0mFound 1 target...\n\u001b[32m[2 / 2]\u001b[0m no actions running\n",
		"\r\u001b[1A\u001b[KTarget //:hello up-to-date:\n\u001b[32m[2 / 2]\u001b[0m no actions running\n\r\u001b[1A\u001b[K  bazel-bin/hello.sh\n\u001b[32m[2 / 2]\u001b[0m no actions running\n\r\u001b[1A\u001b[K\u001b[32mINFO: \u001b[0mElapsed time: 0.137s, Critical Path: 0.02s\n\u001b[32m[2 / 2]\u001b[0m no actions running\n\r\u001b[1A\u001b[K\u001b[32mINFO: \u001b[0m2 processes: 1 internal, 1 linux-sandbox.\n\u001b[32m[2 / 2]\u001b[0m no actions running\n\r\u001b[1A\u001b[K\u001b[32mINFO: \u001b[0mBuild completed successfully, 2 total actions\n\u001b[32mINFO:\u001b[0m \n\r\u001b[1A\u001b[K\u001b[32mINFO:\u001b[0m \n",
	}
	for i, snippet := range logSnippets {
		_, err := client.IncompleteBuildLog.Create().
			SetBazelInvocationID(invocationDbID).
			// LogSnippetID is 1-indexed.
			SetSnippetID(int32(i + 1)).
			SetLogSnippet([]byte(snippet)).
			Save(ctx)
		require.NoError(t, err)
	}
}

func TestCompactLogs(t *testing.T) {
	ctrl, ctx := gomock.WithContext(context.Background(), t)
	ctx = dbauthservice.NewContextWithDbAuthServiceBypass(ctx)
	clock := mock.NewMockClock(ctrl)
	traceProvider := noop.NewTracerProvider()

	t.Run("FinishedInvocationWithoutIncompleteLog", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		client := db.Ent()
		instanceName := testutils.CreateInstanceName(ctx, t, client, "testInstance")
		inv, err := testutils.StartCreateInvocation(client, instanceName).
			SetBepCompleted(true).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		compacted, err := cleanup.CompactLogs(ctx)
		require.NoError(t, err)
		require.EqualValues(t, 0, compacted)
		count, err := client.BuildLogChunk.Query().Where(
			buildlogchunk.HasBazelInvocationWith(
				bazelinvocation.ID(inv.ID),
			),
		).Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, count)

		count, err = client.IncompleteBuildLog.Query().Where(
			incompletebuildlog.HasBazelInvocationWith(
				bazelinvocation.ID(inv.ID),
			),
		).Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})

	requireIncompleteLogCount := func(t *testing.T, client *ent.Client, expected int) {
		count, err := client.IncompleteBuildLog.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, expected, count)
	}

	t.Run("FinishedInvocationWithIncompleteLog", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		client := db.Ent()
		instanceName := testutils.CreateInstanceName(ctx, t, client, "testInstance")
		inv, err := testutils.StartCreateInvocation(client, instanceName).
			SetBepCompleted(true).
			Save(ctx)
		require.NoError(t, err)

		populateIncompleteBuildLog(t, ctx, client, inv.ID)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		// Delete attempt before compaction should not delete logs.
		deleted, err := cleanup.RemoveIncompleteLogs(ctx)
		require.NoError(t, err)
		require.EqualValues(t, 0, deleted)
		requireIncompleteLogCount(t, client, 6)
		// Compaction should not delete logs
		compacted, err := cleanup.CompactLogs(ctx)
		require.NoError(t, err)
		require.EqualValues(t, 1, compacted)
		requireIncompleteLogCount(t, client, 6)
		chunk, err := client.BuildLogChunk.Query().Where(
			buildlogchunk.HasBazelInvocationWith(
				bazelinvocation.IDEQ(inv.ID),
			),
		).Only(ctx)
		require.NoError(t, err)
		decoder, err := zstd.NewReader(nil)
		require.NoError(t, err)
		data, err := decoder.DecodeAll(chunk.Data, nil)
		require.NoError(t, err)
		require.Equal(t, "\x1b[35mWARNING: \x1b[0mBuild options --dynamic_mode, --extra_execution_platforms, and --extra_toolchains have changed, discarding analysis cache (this can be expensive, see https://bazel.build/advanced/performance/iteration-speed).\n\x1b[32mINFO: \x1b[0mAnalyzed target //:hello (0 packages loaded, 2 targets configured).\n\x1b[32mINFO: \x1b[0mFound 1 target...\nTarget //:hello up-to-date:\n  bazel-bin/hello.sh\n\x1b[32mINFO: \x1b[0mElapsed time: 0.137s, Critical Path: 0.02s\n\x1b[32mINFO: \x1b[0m2 processes: 1 internal, 1 linux-sandbox.\n\x1b[32mINFO: \x1b[0mBuild completed successfully, 2 total actions\n\x1b[32mINFO:\x1b[0m \n", string(data))
		// Now logs should be deleted
		deleted, err = cleanup.RemoveIncompleteLogs(ctx)
		require.NoError(t, err)
		require.EqualValues(t, 6, deleted)
		requireIncompleteLogCount(t, client, 0)
	})

	t.Run("UnfinishedInvocationWithIncompleteLog", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		client := db.Ent()
		instanceName := testutils.CreateInstanceName(ctx, t, client, "testInstance")
		inv, err := testutils.StartCreateInvocation(client, instanceName).
			SetBepCompleted(false).
			Save(ctx)
		require.NoError(t, err)

		populateIncompleteBuildLog(t, ctx, client, inv.ID)
		requireIncompleteLogCount(t, client, 6)
		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		compacted, err := cleanup.CompactLogs(ctx)
		require.NoError(t, err)
		require.EqualValues(t, 0, compacted)
		requireIncompleteLogCount(t, client, 6)
		deleted, err := cleanup.RemoveIncompleteLogs(ctx)
		require.NoError(t, err)
		require.EqualValues(t, 0, deleted)
		requireIncompleteLogCount(t, client, 6)
	})
}
