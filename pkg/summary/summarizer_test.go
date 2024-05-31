package summary_test

import (
	"context"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/buildbarn/bb-portal/pkg/summary"
	"github.com/buildbarn/bb-portal/pkg/testkit"
)

var update = flag.Bool("update-golden", false, "update golden (.out.json) files")

func TestSummarize(t *testing.T) {
	require.NoError(t, testDirectory(t, "testdata"))
}

func TestBuildEventProcessor_Original_Testdata(t *testing.T) {
	err := testDirectory(t, "testdata/original")
	if os.IsNotExist(err) {
		t.Skip("Skipping tests using original samples and snapshots.")
		return
	}

	require.NoError(t, err)
}

func testDirectory(t *testing.T, baseDir string) error {
	dirEntries, err := os.ReadDir(baseDir)
	if err != nil {
		return err
	}

	snapshotPath := filepath.Join(baseDir, "snapshots")

	for _, dirEntry := range dirEntries {
		if filepath.Ext(dirEntry.Name()) != ".ndjson" {
			continue
		}

		name := dirEntry.Name()

		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			file := filepath.Join(baseDir, name)
			gotSummary, err := summary.Summarize(ctx, file)

			require.NoError(t, err)
			testkit.CheckAgainstGoldenFile(t, map[string]interface{}{"summary": gotSummary}, snapshotPath, name, update, nil)
		})
	}

	return nil
}
