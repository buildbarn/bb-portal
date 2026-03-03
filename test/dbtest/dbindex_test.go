package dbtest

import (
	"context"
	"fmt"
	"os"
	"testing"

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
	os.Exit(m.Run())
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
}
