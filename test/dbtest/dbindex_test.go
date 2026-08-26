package dbtest

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-portal/internal/database/dbauthservice"
	"github.com/buildbarn/bb-portal/internal/database/embedded"
	"github.com/buildbarn/bb-portal/test/testutils"
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

func TestDatabaseProperties(t *testing.T) {
	ctx := context.Background()
	ctx = dbauthservice.NewContextWithDbAuthServiceBypass(ctx)

	t.Run("AllForeignKeysHaveIndexes", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		rows, err := db.Sqlc().SelectForeignKeysWithoutIndexes(ctx)
		require.NoError(t, err)
		require.Empty(t, rows)
	})

	t.Run("NoRedundantIndexes", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		rows, err := db.Sqlc().SelectRedundantIndexes(ctx)
		require.NoError(t, err)
		require.Empty(t, rows)
	})
}

func TestRenameLegacyActionsTable(t *testing.T) {
	ctx := context.Background()
	connection, err := dbProvider.CreateDatabase()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, connection.Close()) })

	_, err = connection.ExecContext(ctx, `CREATE TABLE actions (id bigint PRIMARY KEY)`)
	require.NoError(t, err)
	require.NoError(t, database.RenameLegacyActionsTable(ctx, connection))
	require.NoError(t, database.RenameLegacyActionsTable(ctx, connection))

	var legacyTable, renamedTable *string
	require.NoError(t, connection.QueryRowContext(
		ctx,
		`SELECT to_regclass('actions')::text, to_regclass('action_executions')::text`,
	).Scan(&legacyTable, &renamedTable))
	require.Nil(t, legacyTable)
	require.NotNil(t, renamedTable)
	require.Equal(t, "action_executions", *renamedTable)
}
