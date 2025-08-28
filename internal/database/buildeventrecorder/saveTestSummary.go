package buildeventrecorder

import (
	"context"
	"fmt"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/ent/gen/ent/testcollection"
	"github.com/buildbarn/bb-portal/ent/gen/ent/testsummary"
	"github.com/buildbarn/bb-storage/pkg/util"
)

func (r *BuildEventRecorder) saveTestSummary(ctx context.Context, tx *ent.Tx, testSummary *bes.TestSummary, label string) error {
	if r.saveTargetDataLevel.GetNone() != nil {
		return nil
	}

	if label == "" {
		return fmt.Errorf("missing label for TestSummary BES message")
	}
	if testSummary == nil {
		return fmt.Errorf("missing TestSummary BES message")
	}

	update := tx.TestCollection.Update().
		Where(
			testcollection.HasBazelInvocationWith(bazelinvocation.ID(r.InvocationDbID)),
			testcollection.LabelEQ(label),
		).
		SetOverallStatus(testcollection.OverallStatus(bes.TestStatus_name[int32(testSummary.OverallStatus)])).
		SetDurationMs(testSummary.TotalRunDuration.AsDuration().Milliseconds())

	if r.saveTargetDataLevel.GetEnriched() != nil {
		testSummaryDb, err := tx.TestSummary.Create().
			SetOverallStatus(testsummary.OverallStatus(bes.TestStatus_name[int32(testSummary.OverallStatus)])).
			SetAttemptCount(testSummary.AttemptCount).
			SetRunCount(testSummary.RunCount).
			SetShardCount(testSummary.ShardCount).
			SetFirstStartTime(testSummary.FirstStartTime.AsTime().Unix()).
			SetLastStopTime(testSummary.FirstStartTime.AsTime().Unix()).
			SetTotalRunCount(testSummary.TotalRunCount).
			SetTotalNumCached(testSummary.TotalNumCached).
			SetTotalRunDuration(testSummary.TotalRunDuration.AsDuration().Milliseconds()).
			SetLabel(label).
			Save(ctx)
		if err != nil {
			return util.StatusWrap(err, "Failed to save test summary to database")
		}

		update.SetTestSummary(testSummaryDb)
	}

	updatedRows, err := update.Save(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to update test collection with test summary")
	}
	if updatedRows == 0 {
		return fmt.Errorf("no test collection found for label %q", label)
	}
	return nil
}
