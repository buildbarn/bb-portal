package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/buildbarn/bb-portal/ent/gen/ent/migrate"
	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-portal/internal/database/dbauthservice"
	"github.com/buildbarn/bb-portal/internal/database/embedded"
	prometheusmetrics "github.com/buildbarn/bb-portal/pkg/prometheus_metrics"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	bb_zstd "github.com/buildbarn/bb-storage/pkg/zstd"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

var dbProvider *embedded.DatabaseProvider

func TestMain(m *testing.M) {
	var err error
	dbProvider, err = embedded.NewDatabaseProvider(os.Stderr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not start embedded DB: %v\n", err)
		os.Exit(1)
	}
	code := m.Run()
	dbProvider.Cleanup()
	os.Exit(code)
}

// parseDatabaseConfiguration parses the JSON form of a Database message,
// with extraFields (each followed by a comma) added to a valid base.
// Configurations are parsed from JSON rather than built as Go literals so
// that an omitted field and an empty message are genuinely different
// inputs, as they are for operators writing jsonnet.
func parseDatabaseConfiguration(t *testing.T, extraFields string) *bb_portal.Database {
	var configuration bb_portal.Database
	require.NoError(t, protojson.Unmarshal([]byte(`{
		`+extraFields+`
		"cleanupConfiguration": {
			"cleanupInterval": "60s",
			"invocationMessageTimeout": "3600s",
			"invocationRetention": "604800s"
		}
	}`), &configuration))
	return &configuration
}

func newTestConnection(t *testing.T) *sql.DB {
	connection, err := dbProvider.CreateDatabase()
	require.NoError(t, err)
	t.Cleanup(func() { connection.Close() })
	return connection
}

func countPublicTables(t *testing.T, connection *sql.DB) int {
	var count int
	require.NoError(t, connection.QueryRow(
		"SELECT count(*) FROM information_schema.tables WHERE table_schema = 'public'",
	).Scan(&count))
	return count
}

func TestInitializeDatabaseSchemaManagement(t *testing.T) {
	for _, tc := range []struct {
		name        string
		extraFields string
		wantApplied bool
	}{
		{name: "Omitted", extraFields: "", wantApplied: true},
		{name: "EmptyMessage", extraFields: `"schemaManagement": {},`, wantApplied: true},
		{name: "Apply", extraFields: `"schemaManagement": {"apply": {}},`, wantApplied: true},
		{name: "ApplyAndExit", extraFields: `"schemaManagement": {"applyAndExit": {}},`, wantApplied: true},
		{name: "Skip", extraFields: `"schemaManagement": {"skip": {}},`, wantApplied: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			connection := newTestConnection(t)
			configuration := parseDatabaseConfiguration(t, tc.extraFields)

			dbClient, dbCleanupService, err := initializeDatabase(
				context.Background(),
				configuration,
				"postgres",
				connection,
				noop.NewTracerProvider(),
				bb_zstd.NewUnboundedPool(nil, nil),
			)
			require.NoError(t, err)

			// The remaining initialization must happen regardless of
			// whether the schema was applied.
			require.NotNil(t, dbClient)
			require.NotNil(t, dbCleanupService)

			if tc.wantApplied {
				require.NotZero(t, countPublicTables(t, connection), "Schema.Create should have created tables")
			} else {
				require.Zero(t, countPublicTables(t, connection), "Schema.Create should not have been called")
			}
		})
	}
}

// TestInitializeDatabaseSkipUsesExistingSchema simulates
// an externally managed schema: the schema is installed before startup,
// and initialization with schema management skipped must still sync metrics
// from it.
func TestInitializeDatabaseSkipUsesExistingSchema(t *testing.T) {
	ctx := dbauthservice.NewContextWithDbAuthServiceBypass(context.Background())
	connection := newTestConnection(t)

	externalMigrator, err := database.New("postgres", connection)
	require.NoError(t, err)
	require.NoError(t, externalMigrator.Ent().Schema.Create(ctx, migrate.WithDropIndex(true)))
	_, err = externalMigrator.Ent().AuthenticatedUser.Create().
		SetUserUUID(uuid.New()).
		SetExternalID("external-id").
		Save(ctx)
	require.NoError(t, err)

	// Use a sentinel so the assertion below cannot pass on a value
	// left behind by another test.
	prometheusmetrics.AuthenticatedUsersCount.Set(-1)

	dbClient, dbCleanupService, err := initializeDatabase(
		context.Background(),
		parseDatabaseConfiguration(t, `"schemaManagement": {"skip": {}},`),
		"postgres",
		connection,
		noop.NewTracerProvider(),
		bb_zstd.NewUnboundedPool(nil, nil),
	)
	require.NoError(t, err)
	require.NotNil(t, dbClient)
	require.NotNil(t, dbCleanupService)
	require.Equal(t, 1.0, testutil.ToFloat64(prometheusmetrics.AuthenticatedUsersCount))
}

func installSchema(t *testing.T, connection *sql.DB) {
	client, err := database.New("postgres", connection)
	require.NoError(t, err)
	require.NoError(t, client.Ent().Schema.Create(context.Background(), migrate.WithDropIndex(true)))
}

func execSQL(t *testing.T, connection *sql.DB, statements ...string) {
	for _, statement := range statements {
		_, err := connection.Exec(statement)
		require.NoError(t, err, statement)
	}
}

func TestInitializeDatabaseVerify(t *testing.T) {
	// Each layout is relative to the schema bb_portal expects.
	current := func(t *testing.T, connection *sql.DB) { installSchema(t, connection) }
	behind := func(t *testing.T, connection *sql.DB) {
		installSchema(t, connection)
		execSQL(t, connection, `ALTER TABLE "authenticated_users" DROP COLUMN "display_name"`)
	}
	ahead := func(t *testing.T, connection *sql.DB) {
		installSchema(t, connection)
		execSQL(
			t, connection,
			`ALTER TABLE "authenticated_users" ADD COLUMN "added_by_newer_release" text NULL`,
			`CREATE INDEX "added_by_newer_release_idx" ON "authenticated_users" ("display_name")`,
		)
	}
	empty := func(t *testing.T, connection *sql.DB) {}

	for _, tc := range []struct {
		name    string
		layout  func(t *testing.T, connection *sql.DB)
		wantErr bool
	}{
		{name: "Current", layout: current},
		{name: "Behind", layout: behind, wantErr: true},
		{name: "Empty", layout: empty, wantErr: true},
		// A newer release may add objects before older replicas are
		// replaced, which must not prevent them from starting.
		{name: "Ahead", layout: ahead},
	} {
		t.Run(tc.name, func(t *testing.T) {
			connection := newTestConnection(t)
			tc.layout(t, connection)

			dbClient, dbCleanupService, err := initializeDatabase(
				context.Background(),
				parseDatabaseConfiguration(t, `"schemaManagement": {"verify": {}},`),
				"postgres",
				connection,
				noop.NewTracerProvider(),
				bb_zstd.NewUnboundedPool(nil, nil),
			)
			if tc.wantErr {
				require.Equal(t, codes.FailedPrecondition, status.Code(err), "error: %v", err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, dbClient)
			require.NotNil(t, dbCleanupService)
		})
	}
}

func TestInitializeDatabaseVerifyReportsMissingChanges(t *testing.T) {
	connection := newTestConnection(t)
	installSchema(t, connection)
	execSQL(t, connection, `ALTER TABLE "authenticated_users" DROP COLUMN "display_name"`)

	_, _, err := initializeDatabase(
		context.Background(),
		parseDatabaseConfiguration(t, `"schemaManagement": {"verify": {}},`),
		"postgres",
		connection,
		noop.NewTracerProvider(),
		bb_zstd.NewUnboundedPool(nil, nil),
	)
	require.ErrorContains(t, err, `ALTER TABLE "authenticated_users" ADD COLUMN "display_name"`)
}

func TestInitializeDatabaseVerifyComparisonError(t *testing.T) {
	// A cancelled context makes the comparison itself fail. Nothing else
	// on the verify path depends on the context.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	connection := newTestConnection(t)
	installSchema(t, connection)

	_, _, err := initializeDatabase(
		ctx,
		parseDatabaseConfiguration(t, `"schemaManagement": {"verify": {}},`),
		"postgres",
		connection,
		noop.NewTracerProvider(),
		bb_zstd.NewUnboundedPool(nil, nil),
	)
	require.ErrorContains(t, err, "Failed to compare the installed schema with the desired schema")
}

// TestInitializeDatabaseDMLOnlyRole runs startup as a role that may
// read and write rows, but not change the schema.
func TestInitializeDatabaseDMLOnlyRole(t *testing.T) {
	owner := newTestConnection(t)
	installSchema(t, owner)

	var databaseName string
	var port int
	require.NoError(t, owner.QueryRow("SELECT current_database(), inet_server_port()").Scan(&databaseName, &port))
	role := "dml_" + uuid.NewString()[:8]
	execSQL(
		t, owner,
		fmt.Sprintf(`CREATE ROLE %q LOGIN PASSWORD 'password'`, role),
		`REVOKE CREATE ON SCHEMA public FROM PUBLIC`,
		fmt.Sprintf(`GRANT USAGE ON SCHEMA public TO %q`, role),
		fmt.Sprintf(`GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO %q`, role),
	)
	initialize := func(t *testing.T, extraFields string) error {
		connection, err := sql.Open("pgx", fmt.Sprintf("postgres://%s:password@localhost:%d/%s?sslmode=disable", role, port, databaseName))
		require.NoError(t, err)
		t.Cleanup(func() { connection.Close() })
		_, _, err = initializeDatabase(
			context.Background(),
			parseDatabaseConfiguration(t, extraFields),
			"postgres",
			connection,
			noop.NewTracerProvider(),
			bb_zstd.NewUnboundedPool(nil, nil),
		)
		return err
	}

	// The subtests share the database, and run in order.
	t.Run("VerifyAcceptsCurrentSchema", func(t *testing.T) {
		require.NoError(t, initialize(t, `"schemaManagement": {"verify": {}},`))
	})
	execSQL(t, owner, `ALTER TABLE "authenticated_users" DROP COLUMN "display_name"`)
	t.Run("VerifyDetectsDrift", func(t *testing.T) {
		err := initialize(t, `"schemaManagement": {"verify": {}},`)
		require.Equal(t, codes.FailedPrecondition, status.Code(err), "error: %v", err)
	})
	t.Run("ApplyFailsOnDrift", func(t *testing.T) {
		require.ErrorContains(t, initialize(t, `"schemaManagement": {"apply": {}},`), "Could not automatically migrate to desired schema")
	})
}
