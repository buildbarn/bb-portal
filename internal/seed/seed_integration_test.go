package seed

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	seedsqlc "github.com/buildbarn/bb-portal/internal/seed/sqlc"
	// Needed to avoid cyclic dependencies in ent.
	_ "github.com/buildbarn/bb-portal/ent/gen/ent/runtime"
	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-portal/internal/database/dbauthservice"
	"github.com/buildbarn/bb-portal/internal/database/embedded"
	"github.com/buildbarn/bb-portal/internal/database/sqlc"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
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

func setupSeedDB(t *testing.T) (*seedsqlc.Queries, database.Client) {
	conn, err := dbProvider.CreateDatabase()
	require.NoError(t, err)
	t.Cleanup(func() { conn.Close() })

	db, err := database.New("postgres", conn)
	require.NoError(t, err)
	require.NoError(t, db.Ent().Schema.Create(context.Background()))

	return seedsqlc.New(conn), db
}

func testSeedParameters() Parameters {
	return Parameters{
		instances:              2,
		users:                  3,
		invocationsPerInstance: 10,
		targetsPerInvocation:   7,
		from:                   time.Now().Add(-time.Hour),
		timeSpan:               time.Hour,
		buildKey:               "build_id",
	}
}

func TestSeed(t *testing.T) {
	ctx := dbauthservice.NewContextWithDbAuthServiceBypass(context.Background())
	seedDB, db := setupSeedDB(t)
	params := testSeedParameters()

	require.NoError(t, PerformSeed(ctx, seedDB, params))

	// The total number of invocations should match the configuration.
	invocationTargetCount, err := db.Ent().InvocationTarget.Query().Count(ctx)
	require.NoError(t, err)
	require.Equal(t, int(params.instances)*int(params.invocationsPerInstance)*int(params.targetsPerInvocation), invocationTargetCount)

	userCount, err := db.Ent().AuthenticatedUser.Query().Count(ctx)
	require.NoError(t, err)
	require.Equal(t, int(params.users), userCount)

	instanceNameCount, err := db.Ent().InstanceName.Query().Count(ctx)
	require.NoError(t, err)
	require.Equal(t, int(params.instances), instanceNameCount)

	// Every test invocation target should have exactly one summary and one
	// result.
	testSummaryCount, err := db.Ent().TestSummary.Query().Count(ctx)
	require.NoError(t, err)
	testResultCount, err := db.Ent().TestResult.Query().Count(ctx)
	require.NoError(t, err)

	testTargetCount, err := db.Ent().TestTarget.Query().Count(ctx)
	require.NoError(t, err)
	require.Positive(t, testTargetCount)
	require.Equal(t, testTargetCount*int(params.invocationsPerInstance), testSummaryCount)
	require.Equal(t, testSummaryCount, testResultCount)

	// Source controls and invocation tags are only seeded when the
	// metadata extractor produces them, which is not the case here.
	sourceControlCount, err := db.Ent().SourceControl.Query().Count(ctx)
	require.NoError(t, err)
	require.Empty(t, sourceControlCount)
	invocationTagCount, err := db.Ent().InvocationTag.Query().Count(ctx)
	require.NoError(t, err)
	require.Empty(t, invocationTagCount)

	// Builds are only seeded when the build key is part of the build tag
	// keys, which requires an extractor, which is not present here.
	buildCount, err := db.Ent().Build.Query().Count(ctx)
	require.NoError(t, err)
	require.Empty(t, buildCount)
}

// TestSeedSequenceAdvancement verifies that after seeding, the identity
// sequences of all seeded tables have been advanced past the explicitly
// assigned seed identifiers. Without this, the real application would
// collide with seed identifiers when inserting new rows.
func TestSeedSequenceAdvancement(t *testing.T) {
	ctx := dbauthservice.NewContextWithDbAuthServiceBypass(context.Background())
	seedDB, db := setupSeedDB(t)
	params := testSeedParameters()

	require.NoError(t, PerformSeed(ctx, seedDB, params))

	// Inserting new rows via the application querier should not collide
	// with the explicitly assigned seed identifiers.
	authenticatedUser, err := db.Sqlc().CreateAuthenticatedUser(ctx, sqlc.CreateAuthenticatedUserParams{
		UserUuid:   uuid.NewSHA1(uuid.NameSpaceURL, []byte("new_user")),
		ExternalID: "new_user",
	})
	require.NoError(t, err)
	require.Greater(t, authenticatedUser.ID, int64(params.users))

	instanceNameID, err := db.Sqlc().CreateInstanceName(ctx, "new_instance")
	require.NoError(t, err)
	require.Greater(t, instanceNameID, int64(params.instances))
}

// TestSeedUsesClosedConnection verifies that a database error during the
// destructive reset is propagated as an error instead of being silently
// ignored, as was the case in the ent based implementation.
func TestSeedUsesClosedConnection(t *testing.T) {
	ctx := dbauthservice.NewContextWithDbAuthServiceBypass(context.Background())
	conn, err := dbProvider.CreateDatabase()
	require.NoError(t, err)

	db, err := database.New("postgres", conn)
	require.NoError(t, err)
	require.NoError(t, db.Ent().Schema.Create(ctx))
	seedDB := seedsqlc.New(conn)

	require.NoError(t, PerformSeed(ctx, seedDB, testSeedParameters()))

	// Closing the connection should cause the next truncate to fail.
	require.NoError(t, conn.Close())
	require.Error(t, seedDB.SeedTruncateEverythingCascade(ctx))
}
