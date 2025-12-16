package dbcleanupservice_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

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
		CleanupInterval:             durationpb.New(1 * time.Minute),
		InvocationConnectionTimeout: durationpb.New(30 * time.Second),
		InvocationMessageTimeout:    durationpb.New(30 * time.Second),
		InvocationRetention:         durationpb.New(30 * time.Minute),
	}
	return dbcleanupservice.NewDbCleanupService(db, clock, cleanupConfiguration, traceProvider)
}

func createInstanceName(t *testing.T, ctx context.Context, client *ent.Client, name string) int {
	in, err := client.InstanceName.Create().
		SetName(name).
		Save(ctx)
	require.NoError(t, err)
	return in.ID
}

func populateIncompleteBuildLog(t *testing.T, ctx context.Context, client *ent.Client, invocationDbID int) {
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
			SetLogSnippet(snippet).
			Save(ctx)
		require.NoError(t, err)
	}
}

func TestLockInvocationsWithNoRecentConnections(t *testing.T) {
	ctrl, ctx := gomock.WithContext(context.Background(), t)
	ctx = dbauthservice.NewContextWithDbAuthServiceBypass(ctx)
	clock := mock.NewMockClock(ctrl)
	cleanupTime := time.Unix(1600000300, 0)
	traceProvider := noop.NewTracerProvider()

	t.Run("NoInvocations", func(t *testing.T) {
		db := setupTestDB(t)
		client := db.Ent()

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)

		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentConnections(ctx)
		require.NoError(t, err)

		count, err := client.BazelInvocation.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})

	t.Run("InvocationsWithNoConnectionMetadata", func(t *testing.T) {
		db := setupTestDB(t)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		inv1, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetBepCompleted(false).
			Save(ctx)
		require.NoError(t, err)
		inv2, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetBepCompleted(true).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)

		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentConnections(ctx)
		require.NoError(t, err)

		count, err := client.BazelInvocation.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 2, count)

		got1, err := client.BazelInvocation.Get(ctx, inv1.ID)
		require.NoError(t, err)
		require.False(t, got1.BepCompleted)

		got2, err := client.BazelInvocation.Get(ctx, inv2.ID)
		require.NoError(t, err)
		require.True(t, got2.BepCompleted)
	})

	t.Run("UnfinishedInvocationNoRecentConnections", func(t *testing.T) {
		db := setupTestDB(t)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		inv, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			Save(ctx)
		require.NoError(t, err)

		_, err = client.ConnectionMetadata.Create().
			SetBazelInvocationID(inv.ID).
			SetConnectionLastOpenAt(cleanupTime.Add(-60 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)

		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentConnections(ctx)
		require.NoError(t, err)

		got, err := client.BazelInvocation.Get(ctx, inv.ID)
		require.NoError(t, err)
		require.True(t, got.BepCompleted)
	})

	t.Run("UnfinishedInvocationRecentConnection", func(t *testing.T) {
		db := setupTestDB(t)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		inv, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			Save(ctx)
		require.NoError(t, err)

		_, err = client.ConnectionMetadata.Create().
			SetBazelInvocationID(inv.ID).
			SetConnectionLastOpenAt(cleanupTime.Add(-10 * time.Second)). // recent
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)

		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentConnections(ctx)
		require.NoError(t, err)

		got, err := client.BazelInvocation.Get(ctx, inv.ID)
		require.NoError(t, err)
		require.False(t, got.BepCompleted)
	})

	t.Run("FinishedInvocationWithOldConnection", func(t *testing.T) {
		db := setupTestDB(t)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		inv, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetBepCompleted(true).
			Save(ctx)
		require.NoError(t, err)

		_, err = client.ConnectionMetadata.Create().
			SetBazelInvocationID(inv.ID).
			SetConnectionLastOpenAt(cleanupTime.Add(-60 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)

		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentConnections(ctx)
		require.NoError(t, err)

		got, err := client.BazelInvocation.Get(ctx, inv.ID)
		require.NoError(t, err)
		require.True(t, got.BepCompleted)
	})

	t.Run("MultipleMixedInvocations", func(t *testing.T) {
		db := setupTestDB(t)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		// old connection -> should be locked
		invOld, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.ConnectionMetadata.Create().
			SetBazelInvocationID(invOld.ID).
			SetConnectionLastOpenAt(cleanupTime.Add(-60 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		// recent connection -> should remain unlocked
		invRecent, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.ConnectionMetadata.Create().
			SetBazelInvocationID(invRecent.ID).
			SetConnectionLastOpenAt(cleanupTime.Add(-10 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)

		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentConnections(ctx)
		require.NoError(t, err)

		gotOld, err := client.BazelInvocation.Get(ctx, invOld.ID)
		require.NoError(t, err)
		require.True(t, gotOld.BepCompleted)

		gotRecent, err := client.BazelInvocation.Get(ctx, invRecent.ID)
		require.NoError(t, err)
		require.False(t, gotRecent.BepCompleted)
	})
}

func TestLockInvocationsWithNoRecentEvents(t *testing.T) {
	ctrl, ctx := gomock.WithContext(context.Background(), t)
	ctx = dbauthservice.NewContextWithDbAuthServiceBypass(ctx)
	clock := mock.NewMockClock(ctrl)
	cleanupTime := time.Unix(1600000060, 0)
	traceProvider := noop.NewTracerProvider()

	t.Run("NoInvocations-NoEventMetadata", func(t *testing.T) {
		db := setupTestDB(t)
		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentEvents(ctx)
		require.NoError(t, err)
	})

	t.Run("UnfinishedInvocation-NoEventMetadata", func(t *testing.T) {
		db := setupTestDB(t)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		startInv, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentEvents(ctx)
		require.NoError(t, err)

		resInv, err := client.BazelInvocation.Get(ctx, startInv.ID)
		require.NoError(t, err)
		require.Nil(t, resInv.EndedAt)
		require.False(t, resInv.BepCompleted)
	})

	t.Run("FinishedInvocation-NoEventMetadata", func(t *testing.T) {
		db := setupTestDB(t)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		startInv, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetBepCompleted(true).
			SetEndedAt(cleanupTime.Add(-60 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentEvents(ctx)
		require.NoError(t, err)

		resInv, err := client.BazelInvocation.Get(ctx, startInv.ID)
		require.NoError(t, err)
		require.Equal(t, startInv.EndedAt.UTC(), resInv.EndedAt.UTC())
		require.True(t, resInv.BepCompleted)
	})

	t.Run("UnfinishedInvocationWithoutEndedAt", func(t *testing.T) {
		db := setupTestDB(t)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		invocationDb, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			Save(ctx)
		require.NoError(t, err)
		em, err := client.EventMetadata.Create().
			SetBazelInvocationID(invocationDb.ID).
			SetSequenceNumber(1).
			SetEventHash([]byte("hash")).
			SetEventReceivedAt(cleanupTime.Add(-60 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentEvents(ctx)
		require.NoError(t, err)

		inv, err := client.BazelInvocation.Get(ctx, invocationDb.ID)
		require.NoError(t, err)
		require.Equal(t, em.EventReceivedAt.UTC(), inv.EndedAt.UTC())
		require.True(t, inv.BepCompleted)
	})

	t.Run("UnfinishedInvocationWithEndedAt", func(t *testing.T) {
		db := setupTestDB(t)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		invocationDb, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetEndedAt(cleanupTime.Add(-60 * time.Second)).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.EventMetadata.Create().
			SetBazelInvocationID(invocationDb.ID).
			SetSequenceNumber(1).
			SetEventHash([]byte("hash")).
			SetEventReceivedAt(cleanupTime.Add(-50 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentEvents(ctx)
		require.NoError(t, err)

		inv, err := client.BazelInvocation.Get(ctx, invocationDb.ID)
		require.NoError(t, err)
		require.Equal(t, cleanupTime.Add(-60*time.Second).UTC(), inv.EndedAt.UTC())
		require.True(t, inv.BepCompleted)
	})

	t.Run("FinishedInvocationWithoutEndedAt", func(t *testing.T) {
		db := setupTestDB(t)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		invocationDb, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetBepCompleted(true).
			Save(ctx)
		require.NoError(t, err)
		em, err := client.EventMetadata.Create().
			SetBazelInvocationID(invocationDb.ID).
			SetSequenceNumber(1).
			SetEventHash([]byte("hash")).
			SetEventReceivedAt(cleanupTime.Add(-50 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentEvents(ctx)
		require.NoError(t, err)

		inv, err := client.BazelInvocation.Get(ctx, invocationDb.ID)
		require.NoError(t, err)
		require.Equal(t, em.EventReceivedAt.UTC(), inv.EndedAt.UTC())
		require.True(t, inv.BepCompleted)
	})

	t.Run("FinishedInvocationWithEndedAt", func(t *testing.T) {
		db := setupTestDB(t)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		invocationDb, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetEndedAt(cleanupTime.Add(-60 * time.Second)).
			SetBepCompleted(true).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.EventMetadata.Create().
			SetBazelInvocationID(invocationDb.ID).
			SetSequenceNumber(1).
			SetEventHash([]byte("hash")).
			SetEventReceivedAt(cleanupTime.Add(-50 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentEvents(ctx)
		require.NoError(t, err)

		inv, err := client.BazelInvocation.Get(ctx, invocationDb.ID)
		require.NoError(t, err)
		require.Equal(t, cleanupTime.Add(-60*time.Second).UTC(), inv.EndedAt.UTC())
		require.True(t, inv.BepCompleted)
	})

	t.Run("UnfinishedInvocation-MultipleEventMetadata", func(t *testing.T) {
		db := setupTestDB(t)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		invocationDb, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.EventMetadata.Create().
			SetBazelInvocationID(invocationDb.ID).
			SetSequenceNumber(1).
			SetEventHash([]byte("hash1")).
			SetEventReceivedAt(cleanupTime.Add(-59 * time.Second)).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.EventMetadata.Create().
			SetBazelInvocationID(invocationDb.ID).
			SetSequenceNumber(2).
			SetEventHash([]byte("hash2")).
			SetEventReceivedAt(cleanupTime.Add(-58 * time.Second)).
			Save(ctx)
		require.NoError(t, err)
		em3, err := client.EventMetadata.Create().
			SetBazelInvocationID(invocationDb.ID).
			SetSequenceNumber(3).
			SetEventHash([]byte("hash3")).
			SetEventReceivedAt(cleanupTime.Add(-57 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentEvents(ctx)
		require.NoError(t, err)

		inv, err := client.BazelInvocation.Get(ctx, invocationDb.ID)
		require.NoError(t, err)
		require.Equal(t, em3.EventReceivedAt.UTC(), inv.EndedAt.UTC())
		require.True(t, inv.BepCompleted)
	})

	t.Run("MultipleUnfinishedInvocations", func(t *testing.T) {
		db := setupTestDB(t)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		invocationDb1, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			Save(ctx)
		require.NoError(t, err)
		em1, err := client.EventMetadata.Create().
			SetBazelInvocationID(invocationDb1.ID).
			SetSequenceNumber(1).
			SetEventHash([]byte("hash1")).
			SetEventReceivedAt(cleanupTime.Add(-60 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		invocationDb2, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			Save(ctx)
		require.NoError(t, err)
		em2, err := client.EventMetadata.Create().
			SetBazelInvocationID(invocationDb2.ID).
			SetSequenceNumber(1).
			SetEventHash([]byte("hash2")).
			SetEventReceivedAt(cleanupTime.Add(-50 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentEvents(ctx)
		require.NoError(t, err)

		inv1, err := client.BazelInvocation.Get(ctx, invocationDb1.ID)
		require.NoError(t, err)
		require.Equal(t, em1.EventReceivedAt.UTC(), inv1.EndedAt.UTC())
		require.True(t, inv1.BepCompleted)

		inv2, err := client.BazelInvocation.Get(ctx, invocationDb2.ID)
		require.NoError(t, err)
		require.Equal(t, em2.EventReceivedAt.UTC(), inv2.EndedAt.UTC())
		require.True(t, inv2.BepCompleted)
	})

	t.Run("UnfinishedInvocation-RecentEventMetadata", func(t *testing.T) {
		db := setupTestDB(t)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		invocationDb, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.EventMetadata.Create().
			SetBazelInvocationID(invocationDb.ID).
			SetSequenceNumber(1).
			SetEventHash([]byte("hash")).
			SetEventReceivedAt(cleanupTime.Add(-20 * time.Second)).
			Save(ctx)
		_, err = client.EventMetadata.Create().
			SetBazelInvocationID(invocationDb.ID).
			SetSequenceNumber(2).
			SetEventHash([]byte("hash")).
			SetEventReceivedAt(cleanupTime.Add(-10 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.LockInvocationsWithNoRecentEvents(ctx)
		require.NoError(t, err)

		inv, err := client.BazelInvocation.Get(ctx, invocationDb.ID)
		require.NoError(t, err)
		require.Nil(t, inv.EndedAt)
		require.False(t, inv.BepCompleted)
	})
}

func TestRemoveOldInvocationConnections(t *testing.T) {
	ctrl, ctx := gomock.WithContext(context.Background(), t)
	ctx = dbauthservice.NewContextWithDbAuthServiceBypass(ctx)
	clock := mock.NewMockClock(ctrl)
	traceProvider := noop.NewTracerProvider()

	t.Run("NoConnections", func(t *testing.T) {
		db := setupTestDB(t)
		client := db.Ent()

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)

		err = cleanup.RemoveOldInvocationConnections(ctx)
		require.NoError(t, err)

		count, err := client.ConnectionMetadata.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})

	t.Run("ConnectionForCompletedInvocation", func(t *testing.T) {
		db := setupTestDB(t)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		inv, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetBepCompleted(true).
			Save(ctx)
		require.NoError(t, err)

		_, err = client.ConnectionMetadata.Create().
			SetBazelInvocationID(inv.ID).
			SetConnectionLastOpenAt(time.Now().UTC()).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)

		err = cleanup.RemoveOldInvocationConnections(ctx)
		require.NoError(t, err)

		count, err := client.ConnectionMetadata.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, count)

		// Invocation should still exist.
		invCount, err := client.BazelInvocation.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, invCount)
	})

	t.Run("ConnectionForUnfinishedInvocation", func(t *testing.T) {
		db := setupTestDB(t)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		inv, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetBepCompleted(false).
			Save(ctx)
		require.NoError(t, err)

		_, err = client.ConnectionMetadata.Create().
			SetBazelInvocationID(inv.ID).
			SetConnectionLastOpenAt(time.Now().UTC()).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)

		err = cleanup.RemoveOldInvocationConnections(ctx)
		require.NoError(t, err)

		count, err := client.ConnectionMetadata.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})

	t.Run("MultipleMixedConnections", func(t *testing.T) {
		db := setupTestDB(t)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		// Completed invocation
		invDone, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetBepCompleted(true).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.ConnectionMetadata.Create().
			SetBazelInvocationID(invDone.ID).
			SetConnectionLastOpenAt(time.Now().UTC()).
			Save(ctx)
		require.NoError(t, err)

		// Unfinished invocation
		invNotDone, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetBepCompleted(false).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.ConnectionMetadata.Create().
			SetBazelInvocationID(invNotDone.ID).
			SetConnectionLastOpenAt(time.Now().UTC()).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)

		err = cleanup.RemoveOldInvocationConnections(ctx)
		require.NoError(t, err)

		// Only the connection for the unfinished invocation should remain.
		conns, err := client.ConnectionMetadata.Query().All(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, len(conns))
		invDbId, err := conns[0].QueryBazelInvocation().OnlyID(ctx)
		require.NoError(t, err)
		require.Equal(t, int(invNotDone.ID), invDbId)
	})
}

func TestRemoveOldEventMetadata(t *testing.T) {
	ctrl, ctx := gomock.WithContext(context.Background(), t)
	ctx = dbauthservice.NewContextWithDbAuthServiceBypass(ctx)
	clock := mock.NewMockClock(ctrl)
	cleanupTime := time.Unix(1600000120, 0)
	traceProvider := noop.NewTracerProvider()

	t.Run("NoEventMetadata", func(t *testing.T) {
		db := setupTestDB(t)
		client := db.Ent()
		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.RemoveOldEventMetadata(ctx)
		require.NoError(t, err)
		count, err := client.EventMetadata.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})

	t.Run("EventMetadataNotOld", func(t *testing.T) {
		db := setupTestDB(t)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		inv, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetBepCompleted(true).
			SetEndedAt(cleanupTime.Add(-15 * time.Second)).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.EventMetadata.Create().
			SetBazelInvocationID(inv.ID).
			SetSequenceNumber(1).
			SetEventHash([]byte("hash")).
			SetEventReceivedAt(cleanupTime.Add(-15 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.RemoveOldEventMetadata(ctx)
		require.NoError(t, err)
		count, err := client.EventMetadata.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})

	t.Run("EventMetadataOld", func(t *testing.T) {
		db := setupTestDB(t)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		inv, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetBepCompleted(true).
			SetEndedAt(cleanupTime.Add(-60 * time.Second)).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.EventMetadata.Create().
			SetBazelInvocationID(inv.ID).
			SetSequenceNumber(1).
			SetEventHash([]byte("hash")).
			SetEventReceivedAt(cleanupTime.Add(-60 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.RemoveOldEventMetadata(ctx)
		require.NoError(t, err)
		count, err := client.EventMetadata.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})

	t.Run("EventMetadataOldButInvocationNotCompleted", func(t *testing.T) {
		db := setupTestDB(t)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		inv, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetEndedAt(cleanupTime.Add(-60 * time.Second)).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.EventMetadata.Create().
			SetBazelInvocationID(inv.ID).
			SetSequenceNumber(1).
			SetEventHash([]byte("hash")).
			SetEventReceivedAt(cleanupTime.Add(-60 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.RemoveOldEventMetadata(ctx)
		require.NoError(t, err)
		count, err := client.EventMetadata.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})

	t.Run("EventMetadataOldButInvocationEndedAtAfterCutoff", func(t *testing.T) {
		db := setupTestDB(t)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		inv, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetBepCompleted(true).
			SetEndedAt(cleanupTime.Add(-15 * time.Second)).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.EventMetadata.Create().
			SetBazelInvocationID(inv.ID).
			SetSequenceNumber(1).
			SetEventHash([]byte("hash")).
			SetEventReceivedAt(cleanupTime.Add(-60 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.RemoveOldEventMetadata(ctx)
		require.NoError(t, err)
		count, err := client.EventMetadata.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})

	t.Run("MultipleEventMetadata", func(t *testing.T) {
		db := setupTestDB(t)
		client := db.Ent()
		instanceNameDbID := createInstanceName(t, ctx, client, "testInstance")
		inv, err := client.BazelInvocation.Create().
			SetInvocationID(uuid.New()).
			SetInstanceNameID(instanceNameDbID).
			SetBepCompleted(true).
			SetEndedAt(cleanupTime.Add(-60 * time.Second)).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.EventMetadata.Create().
			SetBazelInvocationID(inv.ID).
			SetSequenceNumber(1).
			SetEventHash([]byte("hash1")).
			SetEventReceivedAt(cleanupTime.Add(-60 * time.Second)).
			Save(ctx)
		require.NoError(t, err)
		_, err = client.EventMetadata.Create().
			SetBazelInvocationID(inv.ID).
			SetSequenceNumber(2).
			SetEventHash([]byte("hash2")).
			SetEventReceivedAt(cleanupTime.Add(-15 * time.Second)).
			Save(ctx)
		require.NoError(t, err)

		cleanup, err := getNewDbCleanupService(db, clock, traceProvider)
		require.NoError(t, err)
		clock.EXPECT().Now().Return(cleanupTime)
		err = cleanup.RemoveOldEventMetadata(ctx)
		require.NoError(t, err)
		count, err := client.EventMetadata.Query().Count(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})
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
