package processing

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/blob"
	"github.com/buildbarn/bb-portal/ent/gen/ent/build"
	"github.com/buildbarn/bb-portal/ent/gen/ent/missdetail"
	"github.com/buildbarn/bb-portal/ent/gen/ent/targetcomplete"
	"github.com/buildbarn/bb-portal/ent/gen/ent/targetconfigured"
	"github.com/buildbarn/bb-portal/ent/gen/ent/targetpair"
	"github.com/buildbarn/bb-portal/ent/gen/ent/testcollection"
	"github.com/buildbarn/bb-portal/ent/gen/ent/testresultbes"
	"github.com/buildbarn/bb-portal/ent/gen/ent/testsummary"
	prometheusmetrics "github.com/buildbarn/bb-portal/pkg/prometheus_metrics"
	"github.com/buildbarn/bb-portal/pkg/summary"
	"github.com/buildbarn/bb-portal/pkg/summary/detectors"
)

// SaveActor The save actor struct with the db client and a blob archiver.
type SaveActor struct {
	db           *ent.Client
	blobArchiver BlobMultiArchiver
}

// Rollback a transaction on error
func rollback(tx *ent.Tx, err error) error {
	if rerr := tx.Rollback(); rerr != nil {
		err = fmt.Errorf("%w: %v", err, rerr)
	}
	return err
}

// SaveSummary saves an invocation summary to the database as a transaction or rollsback on failure.
func (act SaveActor) SaveSummary(ctx context.Context, summary *summary.Summary) (*ent.BazelInvocation, error) {
	tx, err := act.db.Tx(ctx)
	if err != nil {
		return nil, err
	}
	act.db = tx.Client()
	result, err := act.saveSummary(ctx, summary)
	if err != nil {
		slog.ErrorContext(ctx, "Error saving the invocation to the database.  Rolling back the transaction", "err", err)
		return nil, rollback(tx, err)
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// saveSummary saves an invocation summary to the database.
func (act SaveActor) saveSummary(ctx context.Context, summary *summary.Summary) (*ent.BazelInvocation, error) {
	// errors := []error{}
	if summary.InvocationID == "" {
		slog.ErrorContext(ctx, "No Invocation ID Found on summary", "ctx.Err()", ctx.Err())
		return nil, fmt.Errorf("no Invocation ID Found on summary")
	}
	eventFile, err := act.saveEventFile(ctx, summary)
	if err != nil {
		slog.ErrorContext(ctx, "failed to save event file", "id", summary.InvocationID, "err", err)
	}
	buildRecord, err := act.findOrCreateBuild(ctx, summary)
	if err != nil {
		slog.ErrorContext(ctx, "failed to find or create build", "summary.InvocationId", summary.InvocationID, "err", err)
	}
	metrics, err := act.saveMetrics(ctx, summary.Metrics)
	if err != nil {
		slog.ErrorContext(ctx, "failed to save metrics", "id", summary.InvocationID, "err", err)
	}
	targets, err := act.saveTargets(ctx, summary)
	if err != nil {
		slog.ErrorContext(ctx, "failed to save targets", "id", summary.InvocationID, "err", err)
	}
	tests, err := act.saveTests(ctx, summary)
	if err != nil {
		slog.ErrorContext(ctx, "failed to save tests", "id", summary.InvocationID, "err", err)
		tests = nil
	}
	sourcecontrol, err := act.saveSourceControl(ctx, summary)
	if err != nil {
		slog.ErrorContext(ctx, "failed to save source control information", "id", summary.InvocationID, "err", err)
	}
	bazelInvocation, err := act.saveBazelInvocation(ctx, summary, eventFile, buildRecord, metrics, tests, targets, sourcecontrol)
	if err != nil {
		slog.ErrorContext(ctx, "failed to save bazel invocation", "id", summary.InvocationID, "err", err)
		return nil, fmt.Errorf("could not save BazelInvocation: %w", err)
	}
	var detectedBlobs []detectors.BlobURI
	err = act.db.BazelInvocationProblem.MapCreateBulk(summary.Problems, func(create *ent.BazelInvocationProblemCreate, i int) {
		problem := summary.Problems[i]
		detectedBlobs = append(detectedBlobs, problem.DetectedBlobs...)
		create.
			SetProblemType(string(problem.ProblemType)).
			SetLabel(problem.Label).
			SetBepEvents(problem.BEPEvents).
			SetBazelInvocation(bazelInvocation)
	}).Exec(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to save bazel invocation problems ", "id", summary.InvocationID, "err", err)
		return nil, fmt.Errorf("could not save BazelInvocationProblems: %w", err)
	}
	missingBlobs, err := act.determineMissingBlobs(ctx, detectedBlobs)
	if err != nil {
		slog.ErrorContext(ctx, "failed to determine missing blobs", "id", summary.InvocationID, "err", err)
		return nil, err
	}
	err = act.db.Blob.MapCreateBulk(missingBlobs, func(create *ent.BlobCreate, i int) {
		b := missingBlobs[i]
		create.
			SetURI(string(b)).
			SetInstanceName(summary.InstanceName)
	}).Exec(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to save blobs", "id", summary.InvocationID, "err", err)
		return nil, fmt.Errorf("could not save Blobs: %w", err)
	}
	var archivedBlobs []ent.Blob
	archivedBlobs, err = act.blobArchiver.ArchiveBlobs(ctx, missingBlobs)
	if err != nil {
		slog.ErrorContext(ctx, "failed to archive", "id", summary.InvocationID, "err", err)
		return nil, fmt.Errorf("failed to archive blobs: %w", err)
	}
	for _, archivedBlob := range archivedBlobs {
		act.updateBlobRecord(ctx, archivedBlob)
	}
	return bazelInvocation, nil
}

func (act SaveActor) determineMissingBlobs(ctx context.Context, detectedBlobs []detectors.BlobURI) ([]detectors.BlobURI, error) {
	detectedBlobURIs := make([]string, 0, len(detectedBlobs))
	blobMap := make(map[string]struct{}, len(detectedBlobs))
	for _, detectedBlob := range detectedBlobs {
		detectedBlobURIs = append(detectedBlobURIs, string(detectedBlob))
	}
	foundInDB, err := act.db.Blob.Query().Where(blob.URIIn(detectedBlobURIs...)).All(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to query blobs", "err", err)
		return nil, fmt.Errorf("could not query Blobs: %w", err)
	}

	for _, foundBlob := range foundInDB {
		blobMap[foundBlob.URI] = struct{}{}
	}
	missingBlobs := make([]detectors.BlobURI, 0, len(detectedBlobs)-len(foundInDB))
	for _, detectedBlob := range detectedBlobs {
		if _, ok := blobMap[string(detectedBlob)]; ok {
			continue
		}
		missingBlobs = append(missingBlobs, detectedBlob)
	}
	return missingBlobs, nil
}

func (act SaveActor) saveBazelInvocation(
	ctx context.Context,
	summary *summary.Summary,
	eventFile *ent.EventFile,
	buildRecord *ent.Build,
	metrics *ent.Metrics,
	tests []*ent.TestCollection,
	targets []*ent.TargetPair,
	sourcecontrol *ent.SourceControl,
) (*ent.BazelInvocation, error) {
	if summary == nil {
		return nil, fmt.Errorf("no summary object provided")
	}
	uniqueID, err := uuid.Parse(summary.InvocationID)
	if err != nil {
		return nil, err
	}

	prometheusmetrics.Invocations.WithLabelValues(summary.Hostname, summary.UserLDAP, summary.StepLabel).Inc()

	create := act.db.BazelInvocation.Create().
		SetInvocationID(uniqueID).
		SetInstanceName(summary.InstanceName).
		SetProfileName(summary.ProfileName).
		SetStartedAt(summary.StartedAt).
		SetBuildLogs(summary.BuildLogs.String()).
		SetNillableEndedAt(summary.EndedAt).
		SetChangeNumber(summary.ChangeNumber).
		SetPatchsetNumber(summary.PatchsetNumber).
		SetSummary(*summary.InvocationSummary).
		SetBepCompleted(summary.BEPCompleted).
		SetStepLabel(summary.StepLabel).
		SetUserEmail(summary.UserEmail).
		SetUserLdap(summary.UserLDAP).
		SetCPU(summary.CPU).
		SetConfigurationMnemonic(summary.ConfigrationMnemonic).
		SetPlatformName(summary.PlatformName).
		SetNumFetches(summary.NumFetches).
		SetHostname(summary.Hostname).
		SetRelatedFiles(summary.RelatedFiles)

	if eventFile != nil {
		create = create.SetEventFile(eventFile)
	}
	if metrics != nil {
		create = create.SetMetrics(metrics)
	}
	if tests != nil {
		create = create.AddTestCollection(tests...)
	}
	if targets != nil {
		create = create.AddTargets(targets...)
	}
	if sourcecontrol != nil {
		create = create.SetSourceControl(sourcecontrol)
	}
	if buildRecord != nil {
		create = create.SetBuild(buildRecord)
	}

	return create.
		Save(ctx)
}

func (act SaveActor) saveEventFile(ctx context.Context, summary *summary.Summary) (*ent.EventFile, error) {
	eventFile, err := act.db.EventFile.Create().
		SetURL(summary.EventFileURL).
		SetModTime(time.Now()).              // TODO: Save modTime in summary?
		SetProtocol("BEP").                  // Legacy: used to detect other protocols, e.g. for codechecks.
		SetMimeType("application/x-ndjson"). // NOTE: Only ndjson supported right now, but we should be able to add binary support.
		SetStatus("SUCCESS").                // TODO: Keep workflow of DETECTED->IMPORTING->...?
		Save(ctx)
	return eventFile, err
}

// do we even really need these in the database?
func (act SaveActor) saveTargetConfiguration(ctx context.Context, targetConfiguration summary.TargetConfigured) (*ent.TargetConfigured, error) {
	return act.db.TargetConfigured.Create().
		SetTag(targetConfiguration.Tag).
		SetStartTimeInMs(targetConfiguration.StartTimeInMs).
		SetTargetKind(targetConfiguration.TargetKind).
		SetTestSize(targetconfigured.TestSize(targetConfiguration.TestSize.String())).
		Save(ctx)
}

func (act SaveActor) saveTestFiles(ctx context.Context, files []summary.TestFile) ([]*ent.TestFile, error) {
	return act.db.TestFile.MapCreateBulk(files, func(create *ent.TestFileCreate, i int) {
		file := files[i]
		create.SetDigest(file.Digest).
			SetFile(file.File).
			SetName(file.Name).
			SetLength(file.Length).
			SetPrefix(file.Prefix)
	}).Save(ctx)
}

func (act SaveActor) saveOutputGroup(ctx context.Context, ouputGroup summary.OutputGroup) (*ent.OutputGroup, error) {
	inlineFiles, err := act.saveTestFiles(ctx, ouputGroup.InlineFiles)
	if err != nil {
		slog.ErrorContext(ctx, "failed to save output group", "err", err)
		return nil, err
	}

	return act.db.OutputGroup.Create().
		SetName(ouputGroup.Name).
		SetIncomplete(ouputGroup.Incomplete).
		AddInlineFiles(inlineFiles...).
		// TODO: implement named set of files logic to recursively add files to this collection
		Save(ctx)
}

func (act SaveActor) saveTargetCompletion(ctx context.Context, targetCompletion summary.TargetComplete) (*ent.TargetComplete, error) {
	outputGroup, err := act.saveOutputGroup(ctx, targetCompletion.OutputGroup)
	if err != nil {
		return nil, err
	}
	importantOutput, err := act.saveTestFiles(ctx, targetCompletion.ImportantOutput)
	if err != nil {
		return nil, err
	}
	directoryOutpu, err := act.saveTestFiles(ctx, targetCompletion.DirectoryOutput)
	if err != nil {
		return nil, err
	}

	return act.db.TargetComplete.Create().
		SetSuccess(targetCompletion.Success).
		SetTargetKind(targetCompletion.TargetKind).
		SetTestSize(targetcomplete.TestSize(targetCompletion.TestSize.String())).
		SetTag(targetCompletion.Tag).
		SetEndTimeInMs(targetCompletion.EndTimeInMs).
		SetTestTimeout(targetCompletion.TestTimeout).
		SetTestTimeoutSeconds(targetCompletion.TestTimeoutSeconds).
		SetOutputGroup(outputGroup).
		AddImportantOutput(importantOutput...).
		AddDirectoryOutput(directoryOutpu...).
		Save(ctx)
}

func (act SaveActor) saveTargetPair(ctx context.Context, targetPair summary.TargetPair, label string, enrich, skipPrometheus, skipSave bool, threshold int64) (*ent.TargetPair, error) {
	if targetPair.DurationInMs > threshold && !skipPrometheus {
		prometheusmetrics.TargetDurations.WithLabelValues(label).Set(float64(targetPair.DurationInMs))
	}

	if skipSave {
		return nil, nil
	}

	create := act.db.TargetPair.Create().
		SetLabel(label).
		SetDurationInMs(targetPair.DurationInMs).
		SetSuccess(targetPair.Success).
		SetTargetKind(targetPair.TargetKind).
		SetTestSize(targetpair.TestSize(targetPair.TestSize.String()))

	if !targetPair.Success {

		reason := targetpair.AbortReason(targetPair.AbortReason.String())
		create = create.SetAbortReason(reason)
	}

	if enrich {
		configuration := targetPair.Configuration
		completion := targetPair.Completion

		targetConfiguration, err := act.saveTargetConfiguration(ctx, configuration)
		if err != nil {
			return nil, err
		}

		targetCompletion, err := act.saveTargetCompletion(ctx, completion)
		if err != nil {
			return nil, err
		}
		create.SetCompletion(targetCompletion).SetConfiguration(targetConfiguration)
	}

	return create.Save(ctx)
}

// TODO: is there a more effiient way to do bulk updates instead of sequentially adding everything to the database one object at a time?
// ironically, MapBulkCreate doesn't work for the map(string)TargetPair.  Its expecting an int index, not a label.
func (act SaveActor) saveTargets(ctx context.Context, summary *summary.Summary) ([]*ent.TargetPair, error) {
	// nothing to do, so return empty
	if summary.SkipPrometheusTargets && summary.SkipTargetData {
		return []*ent.TargetPair{}, nil
	}
	var result []*ent.TargetPair = make([]*ent.TargetPair, len(summary.Targets))
	i := 0
	for label, pair := range summary.Targets {
		targetPair, err := act.saveTargetPair(ctx, pair, label, summary.EnrichTargetData, summary.SkipPrometheusTargets, summary.SkipTargetData, summary.PrometheusTargetDurationSkipThreshold)
		if err != nil {
			return nil, err
		}
		result[i] = targetPair
		i++
	}
	if summary.SkipTargetData {
		return []*ent.TargetPair{}, nil
	}
	return result, nil
}

func (act SaveActor) saveTestSummary(ctx context.Context, testSummary summary.TestSummary, label string) (*ent.TestSummary, error) {
	return act.db.TestSummary.Create().
		SetOverallStatus(testsummary.OverallStatus(testSummary.Status.String())).
		SetAttemptCount(testSummary.AttemptCount).
		SetRunCount(testSummary.RunCount).
		SetShardCount(testSummary.ShardCount).
		SetFirstStartTime(testSummary.FirstStartTime).
		SetLastStopTime(testSummary.LastStopTime).
		SetTotalRunCount(testSummary.TotalRunCount).
		SetTotalNumCached(testSummary.TotalNumCached).
		SetTotalRunDuration(testSummary.TotalRunDuration).
		SetLabel(label).
		AddPassed().
		AddFailed().
		Save(ctx)
}

func (act SaveActor) saveTimingChildren(ctx context.Context, children []summary.TimingChild) ([]*ent.TimingChild, error) {
	return act.db.TimingChild.MapCreateBulk(children, func(create *ent.TimingChildCreate, i int) {
		child := children[i]
		create.
			SetName(child.Name).
			SetTime(child.Time)
	}).Save(ctx)
}

func (act SaveActor) saveTimingBreakdown(ctx context.Context, timingBreakdown summary.TimingBreakdown) (*ent.TimingBreakdown, error) {
	timingChildren, err := act.saveTimingChildren(ctx, timingBreakdown.Child)
	if err != nil {
		slog.ErrorContext(ctx, "failed to save timing breakdown", "err", err)
		return nil, err
	}
	return act.db.TimingBreakdown.Create().
		SetName(timingBreakdown.Name).
		SetTime(timingBreakdown.Time).
		AddChild(timingChildren...).
		Save(ctx)
}

func (act SaveActor) saveResourceUsage(ctx context.Context, resourceUsages []summary.ResourceUsage) ([]*ent.ResourceUsage, error) {
	return act.db.ResourceUsage.MapCreateBulk(resourceUsages, func(create *ent.ResourceUsageCreate, i int) {
		resouceUsage := resourceUsages[i]
		create.
			SetName(resouceUsage.Name).
			SetValue(resouceUsage.Value)
	}).Save(ctx)
}

func (act SaveActor) saveExecutionInfo(ctx context.Context, executionInfo summary.ExecutionInfo) (*ent.ExectionInfo, error) {
	timingBreakdown, err := act.saveTimingBreakdown(ctx, executionInfo.TimingBreakdown)
	if err != nil {
		return nil, err
	}
	resourceUsages, err := act.saveResourceUsage(ctx, executionInfo.ResourceUsage)
	if err != nil {
		return nil, err
	}
	return act.db.ExectionInfo.Create().
		SetStrategy(executionInfo.Strategy).
		SetCachedRemotely(executionInfo.CachedRemotely).
		SetExitCode(executionInfo.ExitCode).
		SetHostname(executionInfo.Hostname).
		SetTimingBreakdown(timingBreakdown).
		AddResourceUsage(resourceUsages...).
		Save(ctx)
}

func (act SaveActor) saveTestResults(ctx context.Context, testResults []summary.TestResult) ([]*ent.TestResultBES, error) {
	return act.db.TestResultBES.MapCreateBulk(testResults, func(create *ent.TestResultBESCreate, i int) {
		testResult := testResults[i]
		executionInfo, err := act.saveExecutionInfo(ctx, testResult.ExecutionInfo)
		if err != nil {
			slog.ErrorContext(ctx, "failed to save executioin info", "err", err)
			slog.Error("problem saving execution info object to database", "err", err)
			return
		}
		create.
			SetTestStatus(testresultbes.TestStatus(testResult.Status.String())).
			SetStatusDetails(testResult.StatusDetails).
			SetLabel(testResult.Label).
			SetWarning(testResult.Warning).
			SetCachedLocally(testResult.CachedLocally).
			SetTestAttemptDuration(testResult.TestAttemptDuration).
			SetTestAttemptStart(testResult.TestAttemptStart).
			SetExecutionInfo(executionInfo)
		// TODO: implement test action output AddTestActionOutput()
	}).Save(ctx)
}

func (act SaveActor) saveTestCollection(ctx context.Context, testCollection summary.TestsCollection, label string, encrich, skipDbSave, skipPrometheus bool) (*ent.TestCollection, error) {
	strStatus := testCollection.OverallStatus.String()
	strCacheHit := "miss"
	if testCollection.CachedLocally {
		strCacheHit = "local"
	}
	if testCollection.CachedRemotely {
		strCacheHit = "remote"
	}
	// if this is a cache hit, we don't care about the duration. Only track durations for non-cache hits
	if strCacheHit == "miss" && !skipPrometheus {
		prometheusmetrics.TestDurations.WithLabelValues(label, strStatus, testCollection.Strategy).Set(float64(testCollection.DurationMs))
	}

	if skipDbSave {
		return nil, nil
	}

	create := act.db.TestCollection.Create().
		SetLabel(label).
		SetOverallStatus(testcollection.OverallStatus(testCollection.OverallStatus.String())).
		SetStrategy(testCollection.Strategy).
		SetCachedLocally(testCollection.CachedLocally).
		SetCachedRemotely(testCollection.CachedRemotely).
		SetDurationMs(testCollection.DurationMs).
		SetFirstSeen((testCollection.FirstSeen))

	if encrich {

		testSummary, err := act.saveTestSummary(ctx, testCollection.TestSummary, label)
		if err != nil {
			return nil, err
		}
		testResults, err := act.saveTestResults(ctx, testCollection.TestResults)
		if err != nil {
			return nil, err
		}
		create.SetTestSummary(testSummary).AddTestResults(testResults...)
	}

	return create.Save(ctx)
}

func (act SaveActor) saveTests(ctx context.Context, summary *summary.Summary) ([]*ent.TestCollection, error) {
	if summary.SkipTargetData && summary.SkipPrometheusTargets {
		return []*ent.TestCollection{}, nil
	}
	var result []*ent.TestCollection = make([]*ent.TestCollection, len(summary.Tests))
	i := 0
	for label, collection := range summary.Tests {
		testCollection, err := act.saveTestCollection(ctx, collection, label, summary.EnrichTargetData, summary.SkipTargetData, summary.SkipPrometheusTargets)
		if err != nil {
			return nil, err
		}
		result[i] = testCollection
		i++
	}

	if summary.SkipTargetData {
		return []*ent.TestCollection{}, nil
	}
	return result, nil
}

func (act SaveActor) saveSourceControl(ctx context.Context, summary *summary.Summary) (*ent.SourceControl, error) {
	return act.db.SourceControl.Create().
		SetActor(summary.SourceControlData.Actor).
		SetRepoURL(summary.SourceControlData.RepositoryURL).
		SetCommitSha(summary.SourceControlData.CommitSHA).
		SetBranch(summary.SourceControlData.Branch).
		SetRefs(summary.SourceControlData.Refs).
		SetRunID(summary.SourceControlData.RunID).
		SetAction(summary.SourceControlData.Action).
		SetWorkflow(summary.SourceControlData.Workflow).
		SetWorkspace(summary.SourceControlData.Workspace).
		SetEventName(summary.SourceControlData.EventName).
		SetJob(summary.SourceControlData.Job).
		SetRunnerName(summary.SourceControlData.RunnerName).
		SetRunnerArch(summary.SourceControlData.RunnerArch).
		SetRunnerOs(summary.SourceControlData.RunnerOs).
		Save(ctx)
}

func (act SaveActor) saveMissDetails(ctx context.Context, missDetails []summary.MissDetail) ([]*ent.MissDetail, error) {
	return act.db.MissDetail.MapCreateBulk(missDetails, func(create *ent.MissDetailCreate, i int) {
		missDetal := missDetails[i]
		create.
			SetCount(missDetal.Count).
			SetReason(missdetail.Reason(missDetal.Reason.String()))
	}).Save(ctx)
}

func (act SaveActor) saveActionCacheStatistics(ctx context.Context, actionCacheStastics summary.ActionCacheStatistics) (*ent.ActionCacheStatistics, error) {
	missDetails, err := act.saveMissDetails(ctx, actionCacheStastics.MissDetails)
	if err != nil {
		return nil, err
	}

	// Calculate the cache hit percentage
	total := actionCacheStastics.Hits + actionCacheStastics.Misses
	var cacheHitPercentage float64
	if total > 0 {
		cacheHitPercentage = (float64(actionCacheStastics.Hits) / float64(total))
		prometheusmetrics.CacheHitRate.WithLabelValues("action_cache").Observe(cacheHitPercentage)
	}

	return act.db.ActionCacheStatistics.Create().
		SetSizeInBytes(actionCacheStastics.SizeInBytes).
		SetSaveTimeInMs(actionCacheStastics.SaveTimeInMs).
		SetHits(actionCacheStastics.Hits).
		SetMisses(actionCacheStastics.Misses).
		AddMissDetails(missDetails...).
		Save(ctx)
}

func (act SaveActor) saveRunnerCounts(ctx context.Context, runnerCounts []summary.RunnerCount) ([]*ent.RunnerCount, error) {
	return act.db.RunnerCount.MapCreateBulk(runnerCounts, func(create *ent.RunnerCountCreate, i int) {
		runnerCount := runnerCounts[i]
		create.
			// TODO is there a better type for unsigned int?
			SetActionsExecuted(int64(runnerCount.Count)).
			SetName(runnerCount.Name).
			SetExecKind(runnerCount.ExecKind)
	}).Save(ctx)
}

func (act SaveActor) saveActionDatas(ctx context.Context, actionDatas []summary.ActionData) ([]*ent.ActionData, error) {
	return act.db.ActionData.MapCreateBulk(actionDatas, func(create *ent.ActionDataCreate, i int) {
		actionData := actionDatas[i]
		create.
			SetActionsExecuted(actionData.ActionsExecuted).
			SetMnemonic(actionData.Mnemonic).
			SetFirstStartedMs(actionData.FirstStartedMs).
			SetLastEndedMs(actionData.LastEndedMs).
			SetSystemTime(actionData.SystemTime).
			SetUserTime(actionData.UserTime)
	}).Save(ctx)
}

func (act SaveActor) saveActionSummary(ctx context.Context, actionSummary summary.ActionSummary) (*ent.ActionSummary, error) {
	actionCacheStatistics, err := act.saveActionCacheStatistics(ctx, actionSummary.ActionCacheStatistics)
	if err != nil {
		return nil, err
	}
	runnerCounts, err := act.saveRunnerCounts(ctx, actionSummary.RunnerCount)
	if err != nil {
		return nil, err
	}
	actionDatas, err := act.saveActionDatas(ctx, actionSummary.ActionData)
	if err != nil {
		return nil, err
	}

	return act.db.ActionSummary.Create().
		SetActionsCreated(actionSummary.ActionsCreated).
		SetActionsCreatedNotIncludingAspects(actionSummary.ActionsCreatedNotIncludingAspects).
		SetActionsExecuted(actionSummary.ActionsExecuted).
		SetRemoteCacheHits(actionSummary.RemoteCacheHits).
		SetActionCacheStatistics(actionCacheStatistics).
		AddRunnerCount(runnerCounts...).
		AddActionData(actionDatas...).
		Save(ctx)
}

func (act SaveActor) saveBuildGraphMetrics(ctx context.Context, buildGraphMetrics summary.BuildGraphMetrics) (*ent.BuildGraphMetrics, error) {
	// TODO:implement EvalutionStats once they exist on the proto
	return act.db.BuildGraphMetrics.Create().
		SetActionLookupValueCount(buildGraphMetrics.ActionLookupValueCount).
		SetActionLookupValueCountNotIncludingAspects(buildGraphMetrics.ActionLookupValueCountNotIncludingAspects).
		SetActionCount(buildGraphMetrics.ActionCount).
		SetInputFileConfiguredTargetCount(buildGraphMetrics.InputFileConfiguredTargetCount).
		SetOutputFileConfiguredTargetCount(buildGraphMetrics.OutputFileConfiguredTargetCount).
		SetOtherConfiguredTargetCount(buildGraphMetrics.OtherConfiguredTargetCount).
		SetOutputArtifactCount(buildGraphMetrics.OutputArtifactCount).
		SetPostInvocationSkyframeNodeCount(buildGraphMetrics.PostInvocationSkyframeNodeCount).
		Save(ctx)
}

func (act SaveActor) saveGarbageMetrics(ctx context.Context, garbageMetrics []summary.GarbageMetrics) ([]*ent.GarbageMetrics, error) {
	return act.db.GarbageMetrics.MapCreateBulk(garbageMetrics, func(create *ent.GarbageMetricsCreate, i int) {
		garbageMetric := garbageMetrics[i]
		create.
			SetGarbageCollected(garbageMetric.GarbageCollected).
			SetType(garbageMetric.Type)
	}).Save(ctx)
}

func (act SaveActor) saveMemoryMetrics(ctx context.Context, memoryMetrics summary.MemoryMetrics) (*ent.MemoryMetrics, error) {
	garbageMetrics, err := act.saveGarbageMetrics(ctx, memoryMetrics.GarbageMetrics)
	if err != nil {
		return nil, err
	}
	return act.db.MemoryMetrics.Create().SetPeakPostGcHeapSize(memoryMetrics.PeakPostGcHeapSize).
		SetPeakPostGcTenuredSpaceHeapSize(memoryMetrics.PeakPostGcTenuredSpaceHeapSize).
		SetUsedHeapSizePostBuild(memoryMetrics.UsedHeapSizePostBuild).
		AddGarbageMetrics(garbageMetrics...).
		Save(ctx)
}

func (act SaveActor) saveTargetMetrics(ctx context.Context, targetMetrics summary.TargetMetrics) (*ent.TargetMetrics, error) {
	return act.db.TargetMetrics.Create().
		SetTargetsConfigured(targetMetrics.TargetsConfigured).
		SetTargetsConfiguredNotIncludingAspects(targetMetrics.TargetsConfiguredNotIncludingAspects).
		SetTargetsLoaded(targetMetrics.TargetsLoaded).
		Save(ctx)
}

func (act SaveActor) savePackageLoadMetrics(ctx context.Context, packageLoadMetrics []summary.PackageLoadMetrics) ([]*ent.PackageLoadMetrics, error) {
	return act.db.PackageLoadMetrics.MapCreateBulk(packageLoadMetrics, func(create *ent.PackageLoadMetricsCreate, i int) {
		packageLoadMetric := packageLoadMetrics[i]
		create.
			SetName(packageLoadMetric.Name).
			SetLoadDuration(packageLoadMetric.LoadDuration).
			SetNumTargets(packageLoadMetric.NumTargets).
			SetComputationSteps(packageLoadMetric.ComputationSteps).
			SetNumTransitiveLoads(packageLoadMetric.NumTransitiveLoads).
			SetPackageOverhead(packageLoadMetric.PackageOverhead)
	}).Save(ctx)
}

func (act SaveActor) savePackageMetrics(ctx context.Context, packageMetrics summary.PackageMetrics) (*ent.PackageMetrics, error) {
	packageLoadMetrics, err := act.savePackageLoadMetrics(ctx, packageMetrics.PackageLoadMetrics)
	if err != nil {
		return nil, err
	}
	return act.db.PackageMetrics.Create().
		SetPackagesLoaded(packageMetrics.PackagesLoaded).
		AddPackageLoadMetrics(packageLoadMetrics...).
		Save(ctx)
}

func (act SaveActor) saveCumulativeMetrics(ctx context.Context, cumulativeMetrics summary.CumulativeMetrics) (*ent.CumulativeMetrics, error) {
	return act.db.CumulativeMetrics.Create().
		SetNumAnalyses(cumulativeMetrics.NumAnalyses).
		SetNumBuilds(cumulativeMetrics.NumBuilds).
		Save(ctx)
}

func (act SaveActor) saveTimingMetrics(ctx context.Context, timingMetrics summary.TimingMetrics) (*ent.TimingMetrics, error) {
	// TODO: update when SetActionsExecutionStartInMs is added to and populated in the proto
	return act.db.TimingMetrics.Create().
		SetAnalysisPhaseTimeInMs(timingMetrics.AnalysisPhaseTimeInMs).
		SetCPUTimeInMs(timingMetrics.CPUTimeInMs).
		SetExecutionPhaseTimeInMs(timingMetrics.ExecutionPhaseTimeInMs).
		SetWallTimeInMs(timingMetrics.WallTimeInMs).
		Save(ctx)
}

func (act SaveActor) saveFilesMetric(ctx context.Context, filesMetric summary.FilesMetric) (*ent.FilesMetric, error) {
	return act.db.FilesMetric.Create().
		SetCount(filesMetric.Count).
		SetSizeInBytes(filesMetric.SizeInBytes).
		Save(ctx)
}

func (act SaveActor) saveArtifactMetrics(ctx context.Context, artifactMetrics summary.ArtifactMetrics) (*ent.ArtifactMetrics, error) {
	soureArtifactsRead, err := act.saveFilesMetric(ctx, artifactMetrics.SourceArtifactsRead)
	if err != nil {
		return nil, err
	}
	outputArtifactsSeen, err := act.saveFilesMetric(ctx, artifactMetrics.OutputArtifactsSeen)
	if err != nil {
		return nil, err
	}
	outputArtifactsFromActionCache, err := act.saveFilesMetric(ctx, artifactMetrics.OutputArtifactsFromActionCache)
	if err != nil {
		return nil, err
	}
	topLevelArtifacts, err := act.saveFilesMetric(ctx, artifactMetrics.TopLevelArtifacts)
	if err != nil {
		return nil, err
	}

	return act.db.ArtifactMetrics.Create().
		SetSourceArtifactsRead(soureArtifactsRead).
		SetOutputArtifactsSeen(outputArtifactsSeen).
		SetOutputArtifactsFromActionCache(outputArtifactsFromActionCache).
		SetTopLevelArtifacts(topLevelArtifacts).
		Save(ctx)
}

func (act SaveActor) saveNetworkMetrics(ctx context.Context, networkMetrics summary.NetworkMetrics) (*ent.NetworkMetrics, error) {
	systemNetworkStats, err := act.saveSystemNetworkStats(ctx, *networkMetrics.SystemNetworkStats)
	if err != nil {
		return nil, err
	}
	return act.db.NetworkMetrics.Create().
		SetSystemNetworkStats(systemNetworkStats).
		Save(ctx)
}

func (act SaveActor) saveSystemNetworkStats(ctx context.Context, systemNetworkStats summary.SystemNetworkStats) (*ent.SystemNetworkStats, error) {
	return act.db.SystemNetworkStats.Create().
		SetBytesRecv(systemNetworkStats.BytesRecv).
		SetBytesSent(systemNetworkStats.BytesSent).
		SetPacketsRecv(systemNetworkStats.PacketsRecv).
		SetPacketsSent(systemNetworkStats.PacketsSent).
		SetPeakBytesRecvPerSec(systemNetworkStats.PeakBytesRecvPerSec).
		SetPeakBytesSentPerSec(systemNetworkStats.PeakBytesSentPerSec).
		SetPeakPacketsRecvPerSec(systemNetworkStats.PeakPacketsRecvPerSec).
		SetPeakBytesSentPerSec(systemNetworkStats.PeakPacketsSentPerSec).
		Save(ctx)
}

func (act SaveActor) saveMetrics(ctx context.Context, metrics summary.Metrics) (*ent.Metrics, error) {
	actionSummary, err := act.saveActionSummary(ctx, metrics.ActionSummary)
	if err != nil {
		return nil, err
	}
	buildGraphMetrics, err := act.saveBuildGraphMetrics(ctx, metrics.BuildGraphMetrics)
	if err != nil {
		return nil, err
	}
	memoryMetrics, err := act.saveMemoryMetrics(ctx, metrics.MemoryMetrics)
	if err != nil {
		return nil, err
	}
	targetMetrics, err := act.saveTargetMetrics(ctx, metrics.TargetMetrics)
	if err != nil {
		return nil, err
	}
	packageMetrics, err := act.savePackageMetrics(ctx, metrics.PackageMetrics)
	if err != nil {
		return nil, err
	}
	cumulativeMetrics, err := act.saveCumulativeMetrics(ctx, metrics.CumulativeMetrics)
	if err != nil {
		return nil, err
	}
	timingMetrics, err := act.saveTimingMetrics(ctx, metrics.TimingMetrics)
	if err != nil {
		return nil, err
	}
	artifactMetrics, err := act.saveArtifactMetrics(ctx, metrics.ArtifactMetrics)
	if err != nil {
		return nil, err
	}
	create := act.db.Metrics.Create().
		SetActionSummary(actionSummary).
		SetBuildGraphMetrics(buildGraphMetrics).
		SetMemoryMetrics(memoryMetrics).
		SetTargetMetrics(targetMetrics).
		SetPackageMetrics(packageMetrics).
		SetCumulativeMetrics(cumulativeMetrics).
		SetTimingMetrics(timingMetrics).
		SetArtifactMetrics(artifactMetrics)

	if metrics.NetworkMetrics.SystemNetworkStats != nil {
		networkMetrics, err := act.saveNetworkMetrics(ctx, metrics.NetworkMetrics)
		if err != nil {
			return nil, err
		}
		create = create.SetNetworkMetrics(networkMetrics)
	}

	return create.Save(ctx)
}

func (act SaveActor) findOrCreateBuild(ctx context.Context, summary *summary.Summary) (*ent.Build, error) {
	var err error
	var buildRecord *ent.Build

	if summary.BuildURL == "" {
		return nil, nil
	}

	slog.Info("Querying for build", "url", summary.BuildURL, "uuid", summary.BuildUUID)
	buildRecord, err = act.db.Build.Query().
		Where(build.BuildUUID(summary.BuildUUID)).First(ctx)

	if ent.IsNotFound(err) {
		slog.Info("Creating build", "url", summary.BuildURL, "uuid", summary.BuildUUID)
		buildRecord, err = act.db.Build.Create().
			SetBuildURL(summary.BuildURL).
			SetBuildUUID(summary.BuildUUID).
			SetEnv(buildEnvVars(summary.EnvVars)).
			SetTimestamp(summary.StartedAt).
			Save(ctx)
	}

	if err != nil {
		return nil, fmt.Errorf("could not find or create build: %w", err)
	}
	return buildRecord, nil
}

func (act SaveActor) updateBlobRecord(ctx context.Context, b ent.Blob) {
	update := act.db.Blob.Update().Where(blob.URI(b.URI)).SetArchivingStatus(b.ArchivingStatus)
	if b.ArchiveURL != "" {
		update = update.SetArchiveURL(b.ArchiveURL)
	}
	if b.Reason != "" {
		update = update.SetReason(b.Reason)
	}
	if b.SizeBytes != 0 {
		update = update.SetSizeBytes(b.SizeBytes)
	}
	if _, err := update.Save(ctx); err != nil {
		slog.Error("failed to save archived blob", "uri", b.URI, "err", err)
	}
}

// buildEnvVars filters the input so it only contains well known environment
// variables injected into a CI build (e.g. a Jenkins build). These are well-known
// Jenkins, etc. environment variables and/or environment variables associated
// with plugins for GitHub, Gerrit, etc.
func buildEnvVars(env map[string]string) map[string]string {
	buildEnv := make(map[string]string)
	for k, v := range env {
		if !summary.IsBuildEnvKey(k) {
			continue
		}
		buildEnv[k] = v
	}

	return buildEnv
}
