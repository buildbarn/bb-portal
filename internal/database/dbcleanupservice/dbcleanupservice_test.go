package dbcleanupservice_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/RoaringBitmap/roaring"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/klauspost/compress/zstd"

	// Needed to avoid cyclic dependencies in ent (https://entgo.io/docs/privacy#privacy-policy-registration)
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/ent/gen/ent/buildlogchunk"
	"github.com/buildbarn/bb-portal/ent/gen/ent/incompletebuildlog"
	_ "github.com/buildbarn/bb-portal/ent/gen/ent/runtime"
	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-portal/internal/database/dbauthservice"
	"github.com/buildbarn/bb-portal/internal/database/dbcleanupservice"
	"github.com/buildbarn/bb-portal/internal/database/embedded"
	"github.com/buildbarn/bb-portal/internal/mock"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-storage/pkg/clock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/durationpb"
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

func getNewDbCleanupService(db database.Client, clock clock.Clock, traceProvider trace.TracerProvider) (*dbcleanupservice.DbCleanupService, error) {
	cleanupConfiguration := &bb_portal.BuildEventStreamService_DatabaseCleanupConfiguration{
		CleanupInterval:          durationpb.New(1 * time.Minute),
		InvocationMessageTimeout: durationpb.New(30 * time.Second),
		InvocationRetention:      durationpb.New(30 * time.Minute),
	}
	return dbcleanupservice.NewDbCleanupService(db, clock, cleanupConfiguration, traceProvider)
}

func createInstanceName(t *testing.T, ctx context.Context, client *ent.Client, name string) int64 {
	in, err := client.InstanceName.Create().
		SetName(name).
		Save(ctx)
	require.NoError(t, err)
	return in.ID
}

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

func roaringBytes(params ...uint32) []byte {
	bitmap := roaring.NewBitmap()
	bitmap.AddMany(params)
	bytes, _ := bitmap.ToBytes()
	return bytes
}

func TestCompactLogs(t *testing.T) {
	ctrl, ctx := gomock.WithContext(context.Background(), t)
	ctx = dbauthservice.NewContextWithDbAuthServiceBypass(ctx)
	clock := mock.NewMockClock(ctrl)
	traceProvider := noop.NewTracerProvider()

	t.Run("FinishedInvocationWithoutIncompleteLog", func(t *testing.T) {
		db := setupTestDB(t)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		inv, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetBepCompleted(true).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		err = cleanup.CompactLogs(ctx)
		require.NoError(t, err)
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
		db := setupTestDB(t)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		inv, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetBepCompleted(true).
			Save(ctx)
		require.NoError(t, err)

		populateIncompleteBuildLog(t, ctx, client, inv.ID)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		// Delete attempt before compaction should not delete logs.
		err = cleanup.DeleteIncompleteLogs(ctx)
		require.NoError(t, err)
		requireIncompleteLogCount(t, client, 6)
		// Compaction should not delete logs
		err = cleanup.CompactLogs(ctx)
		require.NoError(t, err)
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
		err = cleanup.DeleteIncompleteLogs(ctx)
		require.NoError(t, err)
		requireIncompleteLogCount(t, client, 0)
	})

	t.Run("UnfinishedInvocationWithIncompleteLog", func(t *testing.T) {
		db := setupTestDB(t)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		inv, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetBepCompleted(false).
			Save(ctx)
		require.NoError(t, err)

		populateIncompleteBuildLog(t, ctx, client, inv.ID)
		requireIncompleteLogCount(t, client, 6)
		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		err = cleanup.CompactLogs(ctx)
		require.NoError(t, err)
		requireIncompleteLogCount(t, client, 6)
		err = cleanup.DeleteIncompleteLogs(ctx)
		require.NoError(t, err)
		requireIncompleteLogCount(t, client, 6)
	})
}

func TestRemoveOldInvocations(t *testing.T) {
	ctrl, ctx := gomock.WithContext(context.Background(), t)
	ctx = dbauthservice.NewContextWithDbAuthServiceBypass(ctx)
	clock := mock.NewMockClock(ctrl)
	cleanupTime := time.Unix(1600000200, 0)
	traceProvider := noop.NewTracerProvider()

	t.Run("NoInvocations", func(t *testing.T) {
		db := setupTestDB(t)
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
		db := setupTestDB(t)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		_, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetEndedAt(cleanupTime.Add(-60 * time.Minute)).
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
		db := setupTestDB(t)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		_, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetBepCompleted(true).
			SetEndedAt(cleanupTime.Add(-15 * time.Minute)).
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
		db := setupTestDB(t)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		_, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetBepCompleted(true).
			SetEndedAt(cleanupTime.Add(-60 * time.Minute)).
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
		db := setupTestDB(t)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		// Old and completed
		_, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetBepCompleted(true).
			SetEndedAt(cleanupTime.Add(-60 * time.Minute)).
			Save(ctx)
		require.NoError(t, err)
		// Not completed
		_, err = client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetEndedAt(cleanupTime.Add(-60 * time.Minute)).
			Save(ctx)
		require.NoError(t, err)
		// Completed but not old
		_, err = client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetBepCompleted(true).
			SetEndedAt(cleanupTime.Add(-15 * time.Minute)).
			Save(ctx)
		require.NoError(t, err)
		// Not completed and not old
		_, err = client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetEndedAt(cleanupTime.Add(-15 * time.Minute)).
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

func TestRemoveBuildsWithoutInvocations(t *testing.T) {
	ctrl, ctx := gomock.WithContext(context.Background(), t)
	ctx = dbauthservice.NewContextWithDbAuthServiceBypass(ctx)
	clock := mock.NewMockClock(ctrl)
	traceProvider := noop.NewTracerProvider()

	t.Run("NoBuilds", func(t *testing.T) {
		db := setupTestDB(t)
		client := db.Ent()

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		err = cleanup.RemoveBuildsWithoutInvocations(ctx)
		require.NoError(t, err)

		count, err := client.Build.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})

	t.Run("BuildWithInvocation", func(t *testing.T) {
		db := setupTestDB(t)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")

		buildObj, err := client.Build.Create().SetBuildURL("1").SetBuildUUID(uuid.New()).SetInstanceNameID(instanceNameDbID).SetTimestamp(time.Now().UTC()).Save(ctx)
		require.NoError(t, err)
		_, err = client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetBuild(buildObj).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		err = cleanup.RemoveBuildsWithoutInvocations(ctx)
		require.NoError(t, err)

		count, err := client.Build.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})

	t.Run("BuildWithoutInvocation", func(t *testing.T) {
		db := setupTestDB(t)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")

		_, err := client.Build.Create().SetBuildURL("1").SetBuildUUID(uuid.New()).SetInstanceNameID(instanceNameDbID).SetTimestamp(time.Now().UTC()).Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		err = cleanup.RemoveBuildsWithoutInvocations(ctx)
		require.NoError(t, err)

		count, err := client.Build.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})

	t.Run("MultipleBuildsMixed", func(t *testing.T) {
		db := setupTestDB(t)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")

		// Build with invocation
		buildWithInv, err := client.Build.Create().SetBuildURL("1").SetBuildUUID(uuid.New()).SetInstanceNameID(instanceNameDbID).SetTimestamp(time.Now().UTC()).Save(ctx)
		require.NoError(t, err)
		_, err = client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetBuild(buildWithInv).
			Save(ctx)
		require.NoError(t, err)
		// Build without invocation
		_, err = client.Build.Create().SetBuildURL("2").SetBuildUUID(uuid.New()).SetInstanceNameID(instanceNameDbID).SetTimestamp(time.Now().UTC()).Save(ctx)
		require.NoError(t, err)
		// Another build without invocation
		_, err = client.Build.Create().SetBuildURL("3").SetBuildUUID(uuid.New()).SetInstanceNameID(instanceNameDbID).SetTimestamp(time.Now().UTC()).Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		err = cleanup.RemoveBuildsWithoutInvocations(ctx)
		require.NoError(t, err)

		count, err := client.Build.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})
}
