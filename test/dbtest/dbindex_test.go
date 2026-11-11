package dbtest

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	atlas "ariga.io/atlas/sql/schema"
	entschema "entgo.io/ent/dialect/sql/schema"

	"github.com/buildbarn/bb-portal/ent/gen/ent/migrate"
	"github.com/buildbarn/bb-portal/internal/database/dbauthservice"
	"github.com/buildbarn/bb-portal/internal/database/embedded"
	"github.com/buildbarn/bb-portal/test/testutils"
	"github.com/stretchr/testify/require"
)

var dbProvider *embedded.DatabaseProvider

// postgresMaxIdentifierLength is the maximum length in bytes of an
// identifier (NAMEDATALEN - 1). Postgres silently truncates longer ones.
const postgresMaxIdentifierLength = 63

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

	t.Run("IndexNamesFitIdentifierLimit", func(t *testing.T) {
		// A truncated index name does not match the name Ent expects,
		// so Ent recreates the index on every migration. Ent shortens
		// long names of explicit indexes itself, but not those of the
		// implicit indexes of unique columns. These are not listed in
		// migrate.Tables, so the names are taken from the schema Ent
		// derives.
		db := testutils.SetupTestDB(t, dbProvider)
		var indexCount int
		var tooLong []string
		checkIndexNames := entschema.WithDiffHook(func(next entschema.Differ) entschema.Differ {
			return entschema.DiffFunc(func(current, desired *atlas.Schema) ([]atlas.Change, error) {
				for _, table := range desired.Tables {
					for _, index := range table.Indexes {
						indexCount++
						if len(index.Name) > postgresMaxIdentifierLength {
							tooLong = append(tooLong, fmt.Sprintf("%s (%d bytes, table %s)", index.Name, len(index.Name), table.Name))
						}
					}
				}
				return next.Diff(current, desired)
			})
		})
		require.NoError(t, db.Ent().Schema.WriteTo(ctx, io.Discard, checkIndexNames))
		require.NotZero(t, indexCount, "the diff hook did not see any indexes")
		require.Empty(t, tooLong)
	})

	t.Run("SchemaCreateIsIdempotent", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		var pendingChanges bytes.Buffer
		require.NoError(t, db.Ent().Schema.WriteTo(ctx, &pendingChanges, migrate.WithDropIndex(true)))
		require.Empty(t, pendingChanges.String())
	})
}
