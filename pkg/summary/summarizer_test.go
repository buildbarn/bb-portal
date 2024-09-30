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

// update flag
var update = flag.Bool("update-golden", false, "update golden (.out.json) files")

// TestSummarize
func TestSummarize(t *testing.T) {
	require.NoError(t, testDirectory(t, "testdata"))
}

// TestBuildEventProcessor_Original_Testdata
func TestBuildEventProcessor_Original_Testdata(t *testing.T) {
	err := testDirectory(t, "testdata/original")
	if os.IsNotExist(err) {
		t.Skip("Skipping tests using original samples and snapshots.")
		return
	}

	require.NoError(t, err)
}

// testDirectory
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

			// the timestamps and duration related fields can and will differ between runs, so we need to zero these out
			for label, target := range gotSummary.Targets {
				target.DurationInMs = 0
				target.Configuration.StartTimeInMs = 0
				target.Completion.EndTimeInMs = 0
				gotSummary.Targets[label] = target
			}

			for label, test := range gotSummary.Tests {
				test.DurationMs = 0
				test.TestSummary.FirstStartTime = 0
				test.TestSummary.LastStopTime = 0
				for i := range test.TestResults {
					test.TestResults[i].TestAttemptDuration = 0
					test.TestResults[i].ExecutionInfo.TimingBreakdown.Time = "0"
					for z := range test.TestResults[i].ExecutionInfo.TimingBreakdown.Child {
						test.TestResults[i].ExecutionInfo.TimingBreakdown.Child[z].Time = "0"
					}
				}
				gotSummary.Tests[label] = test
			}

			require.NoError(t, err)

			// i don't know that this actually has any effect, but setting it just in case
			opts := testkit.CompareOptions{
				DateTimeAgnostic: true,
			}
			testkit.CheckAgainstGoldenFile(t, map[string]interface{}{"summary": gotSummary}, snapshotPath, name, update, &opts)
		})
	}

	return nil
}
