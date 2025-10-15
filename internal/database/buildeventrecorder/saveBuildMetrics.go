package buildeventrecorder

import (
	"context"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	besmetrics "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/packages/metrics"
	bescore "github.com/bazelbuild/bazel/src/main/protobuf"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/missdetail"
	"github.com/buildbarn/bb-storage/pkg/util"
)

func (r *BuildEventRecorder) saveMissDetails(ctx context.Context, tx *ent.Tx, missDetails []*bescore.ActionCacheStatistics_MissDetail, actionCacheStatisticsDbID int) error {
	if missDetails == nil {
		return nil
	}

	err := tx.MissDetail.MapCreateBulk(missDetails, func(create *ent.MissDetailCreate, i int) {
		missDetal := missDetails[i]
		create.
			SetCount(missDetal.Count).
			SetReason(missdetail.Reason(missDetal.Reason.String())).
			SetActionCacheStatisticsID(actionCacheStatisticsDbID)
	}).Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to save miss details to database")
	}
	return nil
}

func (r *BuildEventRecorder) saveActionCacheStatistics(ctx context.Context, tx *ent.Tx, actionCacheStastics *bescore.ActionCacheStatistics, actionSummaryDbID int) error {
	if actionCacheStastics == nil {
		return nil
	}

	actionCacheStasticsDb, err := tx.ActionCacheStatistics.Create().
		SetSizeInBytes(actionCacheStastics.SizeInBytes).
		SetSaveTimeInMs(actionCacheStastics.SaveTimeInMs).
		SetHits(actionCacheStastics.Hits).
		SetMisses(actionCacheStastics.Misses).
		SetActionSummaryID(actionSummaryDbID).
		Save(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to save action cache statistics to database")
	}

	err = r.saveMissDetails(ctx, tx, actionCacheStastics.MissDetails, actionCacheStasticsDb.ID)
	if err != nil {
		return util.StatusWrap(err, "Failed to save miss details to database")
	}
	return nil
}

func (r *BuildEventRecorder) saveRunnerCounts(ctx context.Context, tx *ent.Tx, runnerCounts []*bes.BuildMetrics_ActionSummary_RunnerCount, actionSummaryDbID int) error {
	if runnerCounts == nil {
		return nil
	}

	err := tx.RunnerCount.MapCreateBulk(runnerCounts, func(create *ent.RunnerCountCreate, i int) {
		runnerCount := runnerCounts[i]
		create.
			// TODO is there a better type for unsigned int?
			SetActionsExecuted(int64(runnerCount.Count)).
			SetName(runnerCount.Name).
			SetExecKind(runnerCount.ExecKind).
			SetActionSummaryID(actionSummaryDbID)
	}).Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to save runner counts to database")
	}
	return nil
}

func (r *BuildEventRecorder) saveActionDatas(ctx context.Context, tx *ent.Tx, actionDatas []*bes.BuildMetrics_ActionSummary_ActionData, actionSummaryDbID int) error {
	if actionDatas == nil {
		return nil
	}

	err := tx.ActionData.MapCreateBulk(actionDatas, func(create *ent.ActionDataCreate, i int) {
		actionData := actionDatas[i]
		ad := create.
			SetActionsExecuted(actionData.ActionsExecuted).
			SetMnemonic(actionData.Mnemonic).
			SetFirstStartedMs(actionData.FirstStartedMs).
			SetLastEndedMs(actionData.LastEndedMs).
			SetActionSummaryID(actionSummaryDbID)
		if actionData.SystemTime != nil {
			ad.SetSystemTime(actionData.SystemTime.AsDuration().Milliseconds())
		}
		if actionData.UserTime != nil {
			ad.SetUserTime(actionData.UserTime.AsDuration().Milliseconds())
		}
	}).Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to save action data to database")
	}
	return nil
}

func (r *BuildEventRecorder) saveActionSummary(ctx context.Context, tx *ent.Tx, actionSummary *bes.BuildMetrics_ActionSummary, metricsDbID int) error {
	if actionSummary == nil {
		return nil
	}

	actionSummaryDb, err := tx.ActionSummary.Create().
		SetActionsCreated(actionSummary.ActionsCreated).
		SetActionsCreatedNotIncludingAspects(actionSummary.ActionsCreatedNotIncludingAspects).
		SetActionsExecuted(actionSummary.ActionsExecuted).
		SetRemoteCacheHits(actionSummary.RemoteCacheHits).
		SetMetricsID(metricsDbID).
		Save(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to save action summary to database")
	}

	err = r.saveActionCacheStatistics(ctx, tx, actionSummary.ActionCacheStatistics, actionSummaryDb.ID)
	if err != nil {
		return util.StatusWrap(err, "Failed to save action summary")
	}
	err = r.saveRunnerCounts(ctx, tx, actionSummary.RunnerCount, actionSummaryDb.ID)
	if err != nil {
		return util.StatusWrap(err, "Failed to save action summary")
	}
	err = r.saveActionDatas(ctx, tx, actionSummary.ActionData, actionSummaryDb.ID)
	if err != nil {
		return util.StatusWrap(err, "Failed to save action summary")
	}
	return nil
}

func (r *BuildEventRecorder) saveArtifactMetrics(ctx context.Context, tx *ent.Tx, artifactMetrics *bes.BuildMetrics_ArtifactMetrics, metricsDbID int) error {
	if artifactMetrics == nil {
		return nil
	}

	create := tx.ArtifactMetrics.Create().
		SetMetricsID(metricsDbID)

	if artifactMetrics.SourceArtifactsRead != nil {
		create.
			SetSourceArtifactsReadCount(artifactMetrics.SourceArtifactsRead.Count).
			SetSourceArtifactsReadSizeInBytes(artifactMetrics.SourceArtifactsRead.SizeInBytes)
	}

	if artifactMetrics.OutputArtifactsSeen != nil {
		create.
			SetOutputArtifactsSeenCount(artifactMetrics.OutputArtifactsSeen.Count).
			SetOutputArtifactsSeenSizeInBytes(artifactMetrics.OutputArtifactsSeen.SizeInBytes)
	}

	if artifactMetrics.OutputArtifactsFromActionCache != nil {
		create.
			SetOutputArtifactsFromActionCacheCount(artifactMetrics.OutputArtifactsFromActionCache.Count).
			SetOutputArtifactsFromActionCacheSizeInBytes(artifactMetrics.OutputArtifactsFromActionCache.SizeInBytes)
	}

	if artifactMetrics.TopLevelArtifacts != nil {
		create.
			SetTopLevelArtifactsCount(artifactMetrics.TopLevelArtifacts.Count).
			SetTopLevelArtifactsSizeInBytes(artifactMetrics.TopLevelArtifacts.SizeInBytes)
	}

	err := create.Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to save artifact metrics to database")
	}
	return nil
}

func (r *BuildEventRecorder) saveBuildGraphMetrics(ctx context.Context, tx *ent.Tx, buildGraphMetrics *bes.BuildMetrics_BuildGraphMetrics, metricsDbID int) error {
	if buildGraphMetrics == nil {
		return nil
	}

	// TODO:implement EvalutionStats once they exist on the proto
	err := tx.BuildGraphMetrics.Create().
		SetActionLookupValueCount(buildGraphMetrics.ActionLookupValueCount).
		SetActionLookupValueCountNotIncludingAspects(buildGraphMetrics.ActionLookupValueCountNotIncludingAspects).
		SetActionCount(buildGraphMetrics.ActionCount).
		SetInputFileConfiguredTargetCount(buildGraphMetrics.InputFileConfiguredTargetCount).
		SetOutputFileConfiguredTargetCount(buildGraphMetrics.OutputFileConfiguredTargetCount).
		SetOtherConfiguredTargetCount(buildGraphMetrics.OtherConfiguredTargetCount).
		SetOutputArtifactCount(buildGraphMetrics.OutputArtifactCount).
		SetPostInvocationSkyframeNodeCount(buildGraphMetrics.PostInvocationSkyframeNodeCount).
		SetMetricsID(metricsDbID).
		Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to save build graph metrics to database")
	}
	return nil
}

func (r *BuildEventRecorder) saveCumulativeMetrics(ctx context.Context, tx *ent.Tx, cumulativeMetrics *bes.BuildMetrics_CumulativeMetrics, metricsDbID int) error {
	if cumulativeMetrics == nil {
		return nil
	}

	err := tx.CumulativeMetrics.Create().
		SetNumAnalyses(cumulativeMetrics.NumAnalyses).
		SetNumBuilds(cumulativeMetrics.NumBuilds).
		SetMetricsID(metricsDbID).
		Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to save cumulative metrics to database")
	}
	return nil
}

func (r *BuildEventRecorder) saveGarbageMetrics(ctx context.Context, tx *ent.Tx, garbageMetrics []*bes.BuildMetrics_MemoryMetrics_GarbageMetrics, memoryMetricsDbID int) error {
	if garbageMetrics == nil {
		return nil
	}

	err := tx.GarbageMetrics.MapCreateBulk(garbageMetrics, func(create *ent.GarbageMetricsCreate, i int) {
		garbageMetric := garbageMetrics[i]
		create.
			SetGarbageCollected(garbageMetric.GarbageCollected).
			SetType(garbageMetric.Type).
			SetMemoryMetricsID(memoryMetricsDbID)
	}).Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to save garbage metrics to database")
	}
	return nil
}

func (r *BuildEventRecorder) saveMemoryMetrics(ctx context.Context, tx *ent.Tx, memoryMetrics *bes.BuildMetrics_MemoryMetrics, metricsDbID int) error {
	if memoryMetrics == nil {
		return nil
	}

	memoryMetricsDb, err := tx.MemoryMetrics.Create().
		SetPeakPostGcHeapSize(memoryMetrics.PeakPostGcHeapSize).
		SetPeakPostGcTenuredSpaceHeapSize(memoryMetrics.PeakPostGcTenuredSpaceHeapSize).
		SetUsedHeapSizePostBuild(memoryMetrics.UsedHeapSizePostBuild).
		SetMetricsID(metricsDbID).
		Save(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to save memory metrics to database")
	}

	err = r.saveGarbageMetrics(ctx, tx, memoryMetrics.GarbageMetrics, memoryMetricsDb.ID)
	if err != nil {
		return util.StatusWrap(err, "Failed to save garbage metrics")
	}
	return nil
}

func (r *BuildEventRecorder) saveSystemNetworkStats(ctx context.Context, tx *ent.Tx, systemNetworkStats *bes.BuildMetrics_NetworkMetrics_SystemNetworkStats, networkMetricsDbID int) error {
	if systemNetworkStats == nil {
		return nil
	}

	err := tx.SystemNetworkStats.Create().
		SetBytesRecv(systemNetworkStats.BytesRecv).
		SetBytesSent(systemNetworkStats.BytesSent).
		SetPacketsRecv(systemNetworkStats.PacketsRecv).
		SetPacketsSent(systemNetworkStats.PacketsSent).
		SetPeakBytesRecvPerSec(systemNetworkStats.PeakBytesRecvPerSec).
		SetPeakBytesSentPerSec(systemNetworkStats.PeakBytesSentPerSec).
		SetPeakPacketsRecvPerSec(systemNetworkStats.PeakPacketsRecvPerSec).
		SetPeakBytesSentPerSec(systemNetworkStats.PeakPacketsSentPerSec).
		SetNetworkMetricsID(networkMetricsDbID).
		Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to save system network stats to database")
	}
	return nil
}

func (r *BuildEventRecorder) saveNetworkMetrics(ctx context.Context, tx *ent.Tx, networkMetrics *bes.BuildMetrics_NetworkMetrics, metricsDbID int) error {
	if networkMetrics == nil || networkMetrics.SystemNetworkStats == nil {
		return nil
	}

	networkMetricsDb, err := tx.NetworkMetrics.Create().
		SetMetricsID(metricsDbID).
		Save(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to save network metrics to database")
	}

	err = r.saveSystemNetworkStats(ctx, tx, networkMetrics.SystemNetworkStats, networkMetricsDb.ID)
	if err != nil {
		return util.StatusWrap(err, "Failed to save system network stats")
	}
	return nil
}

func (r *BuildEventRecorder) savePackageLoadMetrics(ctx context.Context, tx *ent.Tx, packageLoadMetrics []*besmetrics.PackageLoadMetrics, packageMetricsDbID int) error {
	if packageLoadMetrics == nil {
		return nil
	}

	err := tx.PackageLoadMetrics.MapCreateBulk(packageLoadMetrics, func(create *ent.PackageLoadMetricsCreate, i int) {
		packageLoadMetric := packageLoadMetrics[i]
		plm := create.
			SetPackageMetricsID(packageMetricsDbID).
			SetName(*packageLoadMetric.Name).
			SetNumTargets(*packageLoadMetric.NumTargets).
			SetComputationSteps(*packageLoadMetric.ComputationSteps).
			SetNumTransitiveLoads(*packageLoadMetric.NumTransitiveLoads).
			SetPackageOverhead(*packageLoadMetric.PackageOverhead)
		if packageLoadMetric.LoadDuration != nil {
			plm.SetLoadDuration(packageLoadMetric.LoadDuration.AsDuration().Milliseconds())
		}
	}).Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to save package load metrics to database")
	}
	return nil
}

func (r *BuildEventRecorder) savePackageMetrics(ctx context.Context, tx *ent.Tx, packageMetrics *bes.BuildMetrics_PackageMetrics, metricsDbID int) error {
	if packageMetrics == nil {
		return nil
	}

	packageMetricsDb, err := tx.PackageMetrics.Create().
		SetPackagesLoaded(packageMetrics.PackagesLoaded).
		SetMetricsID(metricsDbID).
		Save(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to save package metrics to database")
	}

	err = r.savePackageLoadMetrics(ctx, tx, packageMetrics.PackageLoadMetrics, packageMetricsDb.ID)
	if err != nil {
		return util.StatusWrap(err, "Failed to save package load metrics")
	}
	return nil
}

func (r *BuildEventRecorder) saveTargetMetrics(ctx context.Context, tx *ent.Tx, targetMetrics *bes.BuildMetrics_TargetMetrics, metricsDbID int) error {
	if targetMetrics == nil {
		return nil
	}

	err := tx.TargetMetrics.Create().
		SetTargetsConfigured(targetMetrics.TargetsConfigured).
		SetTargetsConfiguredNotIncludingAspects(targetMetrics.TargetsConfiguredNotIncludingAspects).
		SetTargetsLoaded(targetMetrics.TargetsLoaded).
		SetMetricsID(metricsDbID).
		Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to save target metrics to database")
	}
	return nil
}

func (r *BuildEventRecorder) saveTimingMetrics(ctx context.Context, tx *ent.Tx, timingMetrics *bes.BuildMetrics_TimingMetrics, metricsDbID int) error {
	if timingMetrics == nil {
		return nil
	}

	// TODO: update when SetActionsExecutionStartInMs is added to and populated in the proto
	err := tx.TimingMetrics.Create().
		SetMetricsID(metricsDbID).
		SetAnalysisPhaseTimeInMs(timingMetrics.AnalysisPhaseTimeInMs).
		SetCPUTimeInMs(timingMetrics.CpuTimeInMs).
		SetExecutionPhaseTimeInMs(timingMetrics.ExecutionPhaseTimeInMs).
		SetWallTimeInMs(timingMetrics.WallTimeInMs).
		Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to save timing metrics to database")
	}
	return nil
}

func (r *BuildEventRecorder) saveBuildMetrics(ctx context.Context, tx *ent.Tx, metrics *bes.BuildMetrics) error {
	if metrics == nil {
		return nil
	}

	metricsDb, err := tx.Metrics.
		Create().
		SetBazelInvocationID(r.InvocationDbID).
		Save(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to save build metrics to database")
	}

	err = r.saveActionSummary(ctx, tx, metrics.ActionSummary, metricsDb.ID)
	if err != nil {
		return util.StatusWrap(err, "Failed to save action summary")
	}
	err = r.saveArtifactMetrics(ctx, tx, metrics.ArtifactMetrics, metricsDb.ID)
	if err != nil {
		return util.StatusWrap(err, "Failed to save artifact metrics")
	}
	err = r.saveBuildGraphMetrics(ctx, tx, metrics.BuildGraphMetrics, metricsDb.ID)
	if err != nil {
		return util.StatusWrap(err, "Failed to save build graph metrics")
	}
	err = r.saveCumulativeMetrics(ctx, tx, metrics.CumulativeMetrics, metricsDb.ID)
	if err != nil {
		return util.StatusWrap(err, "Failed to save cumulative metrics")
	}
	err = r.saveMemoryMetrics(ctx, tx, metrics.MemoryMetrics, metricsDb.ID)
	if err != nil {
		return util.StatusWrap(err, "Failed to save memory metrics")
	}
	err = r.savePackageMetrics(ctx, tx, metrics.PackageMetrics, metricsDb.ID)
	if err != nil {
		return util.StatusWrap(err, "Failed to save package metrics")
	}
	err = r.saveNetworkMetrics(ctx, tx, metrics.NetworkMetrics, metricsDb.ID)
	if err != nil {
		return util.StatusWrap(err, "Failed to save network metrics")
	}
	err = r.saveTargetMetrics(ctx, tx, metrics.TargetMetrics, metricsDb.ID)
	if err != nil {
		return util.StatusWrap(err, "Failed to save target metrics")
	}
	err = r.saveTimingMetrics(ctx, tx, metrics.TimingMetrics, metricsDb.ID)
	if err != nil {
		return util.StatusWrap(err, "Failed to save timing metrics")
	}
	return nil
}
