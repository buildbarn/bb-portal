package processing_test

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"

	"github.com/buildbarn/bb-portal/ent/gen/ent/enttest"
	"github.com/buildbarn/bb-portal/pkg/processing"
	prometheusmetrics "github.com/buildbarn/bb-portal/pkg/prometheus_metrics"
)

var (
	snapshotForAPITest   = flag.Bool("snapshot-for-api-tests", false, "Update the DB snapshot used for API tests")
	testDBDriver         = flag.String("test-db-driver", "sqlite3", "Driver to use: sqlite3 or postgres")
	testDBDataSourceName = flag.String("test-db-url", "file:ent?mode=memory&_fk=1", "DB connection string. Default is in-memory. Can set to sqlite3 file or postgres server")
)

const inputFixtureBaseDir = "../summary/testdata/"

func TestWorkflow_ProcessFile(t *testing.T) {
	// Ensure Prometheus metrics are registered
	prometheusmetrics.RegisterMetrics()

	if *snapshotForAPITest {
		*testDBDriver = "sqlite3"
		dbFile := "../../internal/graphql/testdata/snapshot.db"
		t.Log("Ignoring --test-db-driver and --test-db-url since --snapshot-for-api-tests was enabled")
		err := os.Remove(dbFile)
		if !errors.Is(err, os.ErrNotExist) {
			require.NoError(t, err)
		}
		t.Logf("Updating DB snapshot in %s", dbFile)
		*testDBDataSourceName = fmt.Sprintf("file:%s?_fk=1", dbFile)
	}

	db := enttest.Open(t, *testDBDriver, *testDBDataSourceName)
	defer func() {
		require.NoError(t, db.Close())
	}()

	worker := processing.New(db, processing.BlobMultiArchiver{})
	ctx := context.Background()

	dirEntries, err := os.ReadDir(inputFixtureBaseDir)
	require.NoError(t, err)

	for _, dirEntry := range dirEntries {
		if filepath.Ext(dirEntry.Name()) != ".ndjson" {
			continue
		}

		name := dirEntry.Name()
		t.Run(name, func(t *testing.T) {
			file := filepath.Join(inputFixtureBaseDir, name)
			invocation, err := worker.ProcessFile(ctx, file)
			require.NoError(t, err)
			require.NotEmpty(t, invocation.InvocationID)
		})
	}
}
