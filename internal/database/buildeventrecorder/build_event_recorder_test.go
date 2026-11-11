package buildeventrecorder_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-portal/internal/database/buildeventrecorder"
	"github.com/buildbarn/bb-portal/internal/database/dbauthservice"
	"github.com/buildbarn/bb-portal/internal/database/embedded"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-portal/test/testutils"
	"github.com/buildbarn/bb-storage/pkg/auth"
	"github.com/buildbarn/bb-storage/pkg/digest"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	_ "github.com/buildbarn/bb-portal/ent/gen/ent/runtime"
)

var dbProvider *embedded.DatabaseProvider

func TestMain(m *testing.M) {
	var err error
	dbProvider, err = embedded.NewDatabaseProvider(os.Stderr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not start embedded DB: %v\n", err)
		os.Exit(1)
	}
	defer dbProvider.Cleanup()
	m.Run()
}

func TestFindOrCreateInstanceName(t *testing.T) {
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

type transactionFailureClient struct {
	database.Client
	failuresRemaining int
	beginAttempts     int
}

func (c *transactionFailureClient) BeginTx(ctx context.Context, opts *sql.TxOptions) (database.Tx, error) {
	c.beginAttempts++
	if c.failuresRemaining > 0 {
		c.failuresRemaining--
		return nil, status.Error(codes.Unavailable, "transient transaction failure")
	}
	return c.Client.BeginTx(ctx, opts)
}

func TestSaveBatchRetriesTransientFailureWithConfiguredTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(dbauthservice.NewContextWithDbAuthServiceBypass(context.Background()), 5*time.Second)
	defer cancel()
	db := &transactionFailureClient{Client: testutils.SetupTestDB(t, dbProvider)}
	recorder, err := buildeventrecorder.NewBuildEventRecorder(
		ctx, db, auth.NewStaticAuthorizer(func(digest.InstanceName) bool { return true }),
		&bb_portal.BuildEventStreamService_SaveDataLevel{
			Level: &bb_portal.BuildEventStreamService_SaveDataLevel_Basic{Basic: &emptypb.Empty{}},
		},
		noop.NewTracerProvider(), "test", uuid.NewString(),
		&buildeventrecorder.DataExtractors{}, "", time.Second,
	)
	require.NoError(t, err)

	// Fail the first attempt to save the event, after construction succeeds.
	// No events are committed before the failure, so the retry budget cannot
	// be reset by forward progress.
	db.failuresRemaining = 1
	db.beginAttempts = 0
	batch := []buildeventrecorder.BuildEventWithInfo{{
		SequenceNumber: 1,
		Event: &bes.BuildEvent{
			Id:      &bes.BuildEventId{Id: &bes.BuildEventId_Fetch{Fetch: &bes.BuildEventId_FetchId{}}},
			Payload: &bes.BuildEvent_Fetch{Fetch: &bes.Fetch{}},
		},
	}}
	require.NoError(t, recorder.SaveBatch(ctx, batch))
	require.Equal(t, 2, db.beginAttempts)
	require.EqualValues(t, 1, db.Ent().BazelInvocation.Query().OnlyX(ctx).NumFetches)

	// Retrying the same batch must retain the usual duplicate-event behavior.
	require.NoError(t, recorder.SaveBatch(ctx, batch))
	require.EqualValues(t, 1, db.Ent().BazelInvocation.Query().OnlyX(ctx).NumFetches)
}
