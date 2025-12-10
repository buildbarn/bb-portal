package buildeventrecorder

import (
	"context"
	"fmt"
	"time"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/ent/gen/ent/testcollection"
	"github.com/buildbarn/bb-portal/ent/gen/ent/testresultbes"
	"github.com/buildbarn/bb-storage/pkg/util"
)

func (r *BuildEventRecorder) saveTimingChildren(ctx context.Context, tx *ent.Tx, children []*bes.TestResult_ExecutionInfo_TimingBreakdown) ([]*ent.TimingChild, error) {
	if children == nil || len(children) == 0 {
		return nil, nil
	}

	tc, err := tx.TimingChild.MapCreateBulk(children, func(create *ent.TimingChildCreate, i int) {
		child := children[i]
		create.
			SetName(child.Name).
			SetTime(child.Time.AsDuration().String())
	}).Save(ctx)
	if err != nil {
		return nil, util.StatusWrap(err, "failed to save timing children to database")
	}
	return tc, nil
}

func (r *BuildEventRecorder) saveTimingBreakdown(ctx context.Context, tx *ent.Tx, timingBreakdown *bes.TestResult_ExecutionInfo_TimingBreakdown) (*ent.TimingBreakdown, error) {
	if timingBreakdown == nil {
		return nil, nil
	}

	timingChildren, err := r.saveTimingChildren(ctx, tx, timingBreakdown.Child)
	if err != nil {
		return nil, util.StatusWrap(err, "failed to save timing breakdown children to database")
	}
	tb, err := tx.TimingBreakdown.Create().
		SetName(timingBreakdown.Name).
		SetTime(timingBreakdown.Time.AsDuration().String()).
		AddChild(timingChildren...).
		Save(ctx)
	if err != nil {
		return nil, util.StatusWrap(err, "failed to save timing breakdown to database")
	}
	return tb, nil
}

func (r *BuildEventRecorder) saveResourceUsage(ctx context.Context, tx *ent.Tx, resourceUsages []*bes.TestResult_ExecutionInfo_ResourceUsage) ([]*ent.ResourceUsage, error) {
	if resourceUsages == nil || len(resourceUsages) == 0 {
		return nil, nil
	}

	ru, err := tx.ResourceUsage.MapCreateBulk(resourceUsages, func(create *ent.ResourceUsageCreate, i int) {
		resouceUsage := resourceUsages[i]
		create.
			SetName(resouceUsage.Name).
			SetValue(string(resouceUsage.Value))
	}).Save(ctx)
	if err != nil {
		return nil, util.StatusWrap(err, "failed to save resource usage to database")
	}
	return ru, nil
}

func (r *BuildEventRecorder) saveExecutionInfo(ctx context.Context, tx *ent.Tx, executionInfo *bes.TestResult_ExecutionInfo) (*ent.ExectionInfo, error) {
	if executionInfo == nil {
		return nil, nil
	}

	timingBreakdown, err := r.saveTimingBreakdown(ctx, tx, executionInfo.TimingBreakdown)
	if err != nil {
		return nil, util.StatusWrap(err, "failed to save timing breakdown")
	}
	resourceUsages, err := r.saveResourceUsage(ctx, tx, executionInfo.ResourceUsage)
	if err != nil {
		return nil, util.StatusWrap(err, "failed to save resource usage")
	}
	create := tx.ExectionInfo.Create().
		SetStrategy(executionInfo.Strategy).
		SetCachedRemotely(executionInfo.CachedRemotely).
		SetExitCode(executionInfo.ExitCode).
		SetHostname(executionInfo.Hostname).
		AddResourceUsage(resourceUsages...)
	if timingBreakdown != nil {
		create.SetTimingBreakdown(timingBreakdown)
	}
	ei, err := create.Save(ctx)
	if err != nil {
		return nil, util.StatusWrap(err, "failed to save execution info to database")
	}
	return ei, nil
}

func (r *BuildEventRecorder) saveTestResultBES(ctx context.Context, tx *ent.Tx, testResult *bes.TestResult, label string, testCollectionDbID int) error {
	if label == "" {
		return fmt.Errorf("missing label for Test Result BES message")
	}
	if testResult == nil {
		return fmt.Errorf("missing test result for label %q", label)
	}

	executionInfo, err := r.saveExecutionInfo(ctx, tx, testResult.ExecutionInfo)
	if err != nil {
		return util.StatusWrap(err, "failed to save execution info for test result")
	}

	err = tx.TestResultBES.Create().
		SetLabel(label).
		SetTestStatus(testresultbes.TestStatus(bes.TestStatus_name[int32(testResult.Status)])).
		SetStatusDetails(testResult.StatusDetails).
		SetWarning(testResult.GetWarning()).
		SetCachedLocally(testResult.GetCachedLocally()).
		SetTestAttemptDuration(testResult.TestAttemptDuration.AsDuration().Milliseconds()).
		SetTestAttemptStart(testResult.TestAttemptStart.AsTime().String()).
		SetExecutionInfo(executionInfo).
		SetTestCollectionID(testCollectionDbID).
		Exec(ctx)
		// TODO: implement test action output AddTestActionOutput()
	if err != nil {
		return util.StatusWrapf(err, "failed to save test result for label %q", label)
	}
	return nil
}

func (r *BuildEventRecorder) createTestCollection(ctx context.Context, tx *ent.Tx, testResult *bes.TestResult, label string) (int, error) {
	cachedLocally := testResult.GetCachedLocally()
	cachedRemotely := testResult.GetExecutionInfo().GetCachedRemotely()
	strategy := testResult.GetExecutionInfo().GetStrategy()

	testCollectionDb, err := tx.TestCollection.Create().
		SetBazelInvocationID(r.InvocationDbID).
		SetLabel(label).
		SetStrategy(strategy).
		SetCachedLocally(cachedLocally).
		SetCachedRemotely(cachedRemotely).
		SetFirstSeen(time.Now().UTC()).
		Save(ctx)
	if err != nil {
		return 0, util.StatusWrap(err, "failed to save test collection to database")
	}

	return testCollectionDb.ID, nil
}

func (r *BuildEventRecorder) updateTestCollection(ctx context.Context, tx *ent.Tx, testResult *bes.TestResult, label string, testCollectionDb *ent.TestCollection) (int, error) {
	cachedLocally := testResult.GetCachedLocally()
	cachedRemotely := testResult.GetExecutionInfo().GetCachedRemotely()
	strategy := testResult.GetExecutionInfo().GetStrategy()

	update := tx.TestCollection.UpdateOneID(testCollectionDb.ID)
	if !cachedLocally {
		update.SetCachedLocally(false)
	}
	if !cachedRemotely {
		update.SetCachedRemotely(false)
	}
	if testCollectionDb.Strategy == "INITIALIZED" {
		update.SetStrategy(strategy)
	} else if testCollectionDb.Strategy != strategy {
		update.SetStrategy("indeterminate")
	}
	testCollectionDb, err := update.Save(ctx)
	if err != nil {
		return 0, util.StatusWrap(err, "failed to update test collection with test result")
	}

	return testCollectionDb.ID, nil
}

func (r *BuildEventRecorder) saveTestResult(ctx context.Context, tx *ent.Tx, testResult *bes.TestResult, label string) error {
	if r.saveTestDataLevel.GetNone() != nil {
		return nil
	}

	if label == "" {
		return fmt.Errorf("missing label for Test Result BES message")
	}
	if testResult == nil {
		return fmt.Errorf("missing test result for label %q", label)
	}

	var testCollectionDbID int

	testCollectionDb, err := tx.TestCollection.Query().Where(
		testcollection.HasBazelInvocationWith(bazelinvocation.ID(r.InvocationDbID)),
		testcollection.LabelEQ(label),
	).Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return util.StatusWrap(err, "failed to query test collection from database")
	}
	if ent.IsNotFound(err) {
		testCollectionDbID, err = r.createTestCollection(ctx, tx, testResult, label)
		if err != nil {
			return util.StatusWrap(err, "failed to create test collection in database")
		}
	} else {
		testCollectionDbID, err = r.updateTestCollection(ctx, tx, testResult, label, testCollectionDb)
		if err != nil {
			return util.StatusWrap(err, "failed to update test collection in database")
		}
	}

	if r.saveTestDataLevel.GetEnriched() != nil {
		err = r.saveTestResultBES(ctx, tx, testResult, label, testCollectionDbID)
		if err != nil {
			return util.StatusWrap(err, "failed to save test result BES message to database")
		}
	}

	return nil
}
