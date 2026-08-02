package buildeventrecorder

import (
	"context"
	"fmt"
	"log/slog"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	bescore "github.com/bazelbuild/bazel/src/main/protobuf"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-storage/pkg/util"
)

func (r *buildEventRecorder) saveMissDetails(ctx context.Context, tx *ent.Client, missDetails []*bescore.ActionCacheStatistics_MissDetail, actionCacheStatisticsDbID int64) error {
	if missDetails == nil {
		return nil
	}

	err := tx.MissDetail.MapCreateBulk(missDetails, func(create *ent.MissDetailCreate, i int) {
		missDetail := missDetails[i]

		create.
			SetCount(missDetail.Count).
			SetActionCacheStatisticsID(actionCacheStatisticsDbID)

		if value, ok := bescore.ActionCacheStatistics_MissReason_name[int32(*missDetail.Reason.Enum())]; ok {
			create.SetReason(value)
		} else {
			create.SetReason("UNKOWN")
			slog.Warn(fmt.Sprintf("Unknown ActionCacheStatistic Miss reason enum value: %d. This is probably because a new Bazel verison has a new enum value that bb-portal doens't implement.", *missDetail.Reason.Enum()))
		}
	}).Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to save miss details to database")
	}
	return nil
}

func (r *buildEventRecorder) saveActionCacheStatistics(ctx context.Context, tx *ent.Client, actionCacheStastics *bescore.ActionCacheStatistics, actionSummaryDbID int64) error {
	if actionCacheStastics == nil {
		return nil
	}

	actionCacheStasticsDb, err := tx.ActionCacheStatistics.Create().
		SetSizeInBytes(actionCacheStastics.SizeInBytes).
		SetSaveTimeInMs(actionCacheStastics.SaveTimeInMs).
		SetCacheCheckSemaphoreWaitTimeInMs(actionCacheStastics.CacheCheckSemaphoreWaitTimeInMs).
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

func (r *buildEventRecorder) saveRunnerCounts(ctx context.Context, tx *ent.Client, runnerCounts []*bes.BuildMetrics_ActionSummary_RunnerCount, actionSummaryDbID int64) error {
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

func (r *buildEventRecorder) saveActionDatas(ctx context.Context, tx *ent.Client, actionDatas []*bes.BuildMetrics_ActionSummary_ActionData, actionSummaryDbID int64) error {
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

func (r *buildEventRecorder) saveActionSummary(ctx context.Context, tx *ent.Client, actionSummary *bes.BuildMetrics_ActionSummary, metricsDbID int64) error {
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

func (r *buildEventRecorder) saveArtifactMetrics(ctx context.Context, tx *ent.Client, artifactMetrics *bes.BuildMetrics_ArtifactMetrics, metricsDbID int64) error {
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

func (r *buildEventRecorder) saveBuildGraphMetrics(ctx context.Context, tx *ent.Client, buildGraphMetrics *bes.BuildMetrics_BuildGraphMetrics, metricsDbID int64) error {
	if buildGraphMetrics == nil {
		return nil
	}

	buildGraphMetricsDb, err := tx.BuildGraphMetrics.Create().
		SetActionLookupValueCount(buildGraphMetrics.ActionLookupValueCount).
		SetActionLookupValueCountNotIncludingAspects(buildGraphMetrics.ActionLookupValueCountNotIncludingAspects).
		SetActionCount(buildGraphMetrics.ActionCount).
		SetInputFileConfiguredTargetCount(buildGraphMetrics.InputFileConfiguredTargetCount).
		SetOutputFileConfiguredTargetCount(buildGraphMetrics.OutputFileConfiguredTargetCount).
		SetOtherConfiguredTargetCount(buildGraphMetrics.OtherConfiguredTargetCount).
		SetOutputArtifactCount(buildGraphMetrics.OutputArtifactCount).
		SetPostInvocationSkyframeNodeCount(buildGraphMetrics.PostInvocationSkyframeNodeCount).
		SetMetricsID(metricsDbID).
		Save(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to save build graph metrics to database")
	}

	err = r.saveBuildGraphEvaluationStats(ctx, tx, buildGraphMetrics, buildGraphMetricsDb.ID)
	if err != nil {
		return util.StatusWrap(err, "Failed to save build graph evaluation stats")
	}
	return nil
}

type buildGraphEvaluationStat struct {
	operation string
	stat      *bes.BuildMetrics_EvaluationStat
}

func (r *buildEventRecorder) saveBuildGraphEvaluationStats(ctx context.Context, tx *ent.Client, buildGraphMetrics *bes.BuildMetrics_BuildGraphMetrics, buildGraphMetricsDbID int64) error {
	evaluationStats := make([]buildGraphEvaluationStat, 0,
		len(buildGraphMetrics.DirtiedValues)+
			len(buildGraphMetrics.ChangedValues)+
			len(buildGraphMetrics.BuiltValues)+
			len(buildGraphMetrics.CleanedValues)+
			len(buildGraphMetrics.EvaluatedValues))

	appendStats := func(operation string, stats []*bes.BuildMetrics_EvaluationStat) {
		for _, stat := range stats {
			if stat != nil {
				evaluationStats = append(evaluationStats, buildGraphEvaluationStat{
					operation: operation,
					stat:      stat,
				})
			}
		}
	}

	appendStats("DIRTIED", buildGraphMetrics.DirtiedValues)
	appendStats("CHANGED", buildGraphMetrics.ChangedValues)
	appendStats("BUILT", buildGraphMetrics.BuiltValues)
	appendStats("CLEANED", buildGraphMetrics.CleanedValues)
	appendStats("EVALUATED", buildGraphMetrics.EvaluatedValues)

	if len(evaluationStats) == 0 {
		return nil
	}

	err := tx.BuildGraphEvaluationStat.MapCreateBulk(evaluationStats, func(create *ent.BuildGraphEvaluationStatCreate, i int) {
		evaluationStat := evaluationStats[i]
		create.
			SetOperation(evaluationStat.operation).
			SetSkyfunctionName(evaluationStat.stat.SkyfunctionName).
			SetCount(evaluationStat.stat.Count).
			SetBuildGraphMetricsID(buildGraphMetricsDbID)
	}).Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to save build graph evaluation stats to database")
	}
	return nil
}

func (r *buildEventRecorder) saveGarbageMetrics(ctx context.Context, tx *ent.Client, garbageMetrics []*bes.BuildMetrics_MemoryMetrics_GarbageMetrics, memoryMetricsDbID int64) error {
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

func (r *buildEventRecorder) saveMemoryMetrics(ctx context.Context, tx *ent.Client, memoryMetrics *bes.BuildMetrics_MemoryMetrics, metricsDbID int64) error {
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

func (r *buildEventRecorder) saveSystemNetworkStats(ctx context.Context, tx *ent.Client, systemNetworkStats *bes.BuildMetrics_NetworkMetrics_SystemNetworkStats, networkMetricsDbID int64) error {
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
		SetPeakPacketsSentPerSec(systemNetworkStats.PeakPacketsSentPerSec).
		SetNetworkMetricsID(networkMetricsDbID).
		Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to save system network stats to database")
	}
	return nil
}

func (r *buildEventRecorder) saveNetworkMetrics(ctx context.Context, tx *ent.Client, networkMetrics *bes.BuildMetrics_NetworkMetrics, metricsDbID int64) error {
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

func (r *buildEventRecorder) saveTargetMetrics(ctx context.Context, tx *ent.Client, targetMetrics *bes.BuildMetrics_TargetMetrics, metricsDbID int64) error {
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

func (r *buildEventRecorder) savePackageMetrics(ctx context.Context, tx *ent.Client, packageMetrics *bes.BuildMetrics_PackageMetrics, metricsDbID int64) error {
	if packageMetrics == nil {
		return nil
	}

	err := tx.PackageMetrics.Create().
		SetPackagesLoaded(packageMetrics.PackagesLoaded).
		SetMetricsID(metricsDbID).
		Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to save package metrics to database")
	}
	return nil
}

func (r *buildEventRecorder) saveCumulativeMetrics(ctx context.Context, tx *ent.Client, cumulativeMetrics *bes.BuildMetrics_CumulativeMetrics, metricsDbID int64) error {
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

func (r *buildEventRecorder) saveTimingMetrics(ctx context.Context, tx *ent.Client, timingMetrics *bes.BuildMetrics_TimingMetrics, metricsDbID int64) error {
	if timingMetrics == nil {
		return nil
	}

	// TODO: update when SetActionsExecutionStartInMs is added to and populated in the proto
	create := tx.TimingMetrics.Create().
		SetMetricsID(metricsDbID).
		SetAnalysisPhaseTimeInMs(timingMetrics.AnalysisPhaseTimeInMs).
		SetCPUTimeInMs(timingMetrics.CpuTimeInMs).
		SetExecutionPhaseTimeInMs(timingMetrics.ExecutionPhaseTimeInMs).
		SetWallTimeInMs(timingMetrics.WallTimeInMs)

	if timingMetrics.CriticalPathTime != nil {
		create.SetCriticalPathTimeInMs(timingMetrics.CriticalPathTime.AsDuration().Milliseconds())
	}

	err := create.Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to save timing metrics to database")
	}
	return nil
}

func (r *buildEventRecorder) saveWorkerStats(ctx context.Context, tx *ent.Client, workerStats []*bes.BuildMetrics_WorkerMetrics_WorkerStats, workerMetricsDbID int64) error {
	nonNilWorkerStats := make([]*bes.BuildMetrics_WorkerMetrics_WorkerStats, 0, len(workerStats))
	for _, workerStat := range workerStats {
		if workerStat != nil {
			nonNilWorkerStats = append(nonNilWorkerStats, workerStat)
		}
	}
	if len(nonNilWorkerStats) == 0 {
		return nil
	}

	err := tx.WorkerStats.MapCreateBulk(nonNilWorkerStats, func(create *ent.WorkerStatsCreate, i int) {
		workerStat := nonNilWorkerStats[i]
		create.
			SetCollectTimeInMs(workerStat.CollectTimeInMs).
			SetWorkerMemoryInKB(workerStat.WorkerMemoryInKb).
			SetPriorWorkerMemoryInKB(workerStat.PriorWorkerMemoryInKb).
			SetLastActionStartTimeInMs(workerStat.LastActionStartTimeInMs).
			SetWorkerMetricsID(workerMetricsDbID)
	}).Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to save worker stats to database")
	}
	return nil
}

func (r *buildEventRecorder) saveWorkerIDs(ctx context.Context, tx *ent.Client, workerMetrics *bes.BuildMetrics_WorkerMetrics, workerMetricsDbID int64) error {
	workerIDs := workerMetrics.WorkerIds
	if len(workerIDs) == 0 && workerMetrics.WorkerId > 0 {
		// Bazel versions before multiplex worker metrics used the singular ID.
		workerIDs = []uint32{uint32(workerMetrics.WorkerId)}
	}
	if len(workerIDs) == 0 {
		return nil
	}

	err := tx.WorkerID.MapCreateBulk(workerIDs, func(create *ent.WorkerIDCreate, i int) {
		create.
			SetWorkerID(workerIDs[i]).
			SetWorkerMetricsID(workerMetricsDbID)
	}).Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to save worker IDs to database")
	}
	return nil
}

func (r *buildEventRecorder) saveWorkerMetrics(ctx context.Context, tx *ent.Client, workerMetrics []*bes.BuildMetrics_WorkerMetrics, metricsDbID int64) error {
	for _, workerMetric := range workerMetrics {
		if workerMetric == nil {
			continue
		}

		create := tx.WorkerMetrics.Create().
			SetProcessID(workerMetric.ProcessId).
			SetMnemonic(workerMetric.Mnemonic).
			SetIsMultiplex(workerMetric.IsMultiplex).
			SetIsSandbox(workerMetric.IsSandbox).
			SetIsMeasurable(workerMetric.IsMeasurable).
			SetWorkerKeyHash(workerMetric.WorkerKeyHash).
			SetWorkerStatus(workerMetric.WorkerStatus.String()).
			SetActionsExecuted(workerMetric.ActionsExecuted).
			SetPriorActionsExecuted(workerMetric.PriorActionsExecuted).
			SetMetricsID(metricsDbID)

		if workerMetric.Code != nil {
			create.SetCode(workerMetric.Code.String())
		}

		workerMetricsDb, err := create.Save(ctx)
		if err != nil {
			return util.StatusWrap(err, "Failed to save worker metrics to database")
		}

		if err := r.saveWorkerIDs(ctx, tx, workerMetric, workerMetricsDb.ID); err != nil {
			return util.StatusWrap(err, "Failed to save worker IDs")
		}
		if err := r.saveWorkerStats(ctx, tx, workerMetric.WorkerStats, workerMetricsDb.ID); err != nil {
			return util.StatusWrap(err, "Failed to save worker stats")
		}
	}
	return nil
}

func (r *buildEventRecorder) saveWorkerPoolMetrics(ctx context.Context, tx *ent.Client, workerPoolMetrics *bes.BuildMetrics_WorkerPoolMetrics, metricsDbID int64) error {
	if workerPoolMetrics == nil {
		return nil
	}

	workerPoolMetricsDb, err := tx.WorkerPoolMetrics.Create().
		SetMetricsID(metricsDbID).
		Save(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to save worker pool metrics to database")
	}

	workerPoolStats := make([]*bes.BuildMetrics_WorkerPoolMetrics_WorkerPoolStats, 0, len(workerPoolMetrics.WorkerPoolStats))
	for _, workerPoolStat := range workerPoolMetrics.WorkerPoolStats {
		if workerPoolStat != nil {
			workerPoolStats = append(workerPoolStats, workerPoolStat)
		}
	}
	if len(workerPoolStats) == 0 {
		return nil
	}

	err = tx.WorkerPoolStats.MapCreateBulk(workerPoolStats, func(create *ent.WorkerPoolStatsCreate, i int) {
		workerPoolStat := workerPoolStats[i]
		create.
			SetHash(workerPoolStat.Hash).
			SetMnemonic(workerPoolStat.Mnemonic).
			SetCreatedCount(workerPoolStat.CreatedCount).
			SetDestroyedCount(workerPoolStat.DestroyedCount).
			SetEvictedCount(workerPoolStat.EvictedCount).
			SetUserExecExceptionDestroyedCount(workerPoolStat.UserExecExceptionDestroyedCount).
			SetIoExceptionDestroyedCount(workerPoolStat.IoExceptionDestroyedCount).
			SetInterruptedExceptionDestroyedCount(workerPoolStat.InterruptedExceptionDestroyedCount).
			SetUnknownDestroyedCount(workerPoolStat.UnknownDestroyedCount).
			SetAliveCount(workerPoolStat.AliveCount).
			SetWorkerPoolMetricsID(workerPoolMetricsDb.ID)
	}).Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to save worker pool stats to database")
	}
	return nil
}

func (r *buildEventRecorder) saveBuildMetrics(ctx context.Context, tx *ent.Client, metrics *bes.BuildMetrics) error {
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
	err = r.saveMemoryMetrics(ctx, tx, metrics.MemoryMetrics, metricsDb.ID)
	if err != nil {
		return util.StatusWrap(err, "Failed to save memory metrics")
	}
	err = r.saveNetworkMetrics(ctx, tx, metrics.NetworkMetrics, metricsDb.ID)
	if err != nil {
		return util.StatusWrap(err, "Failed to save network metrics")
	}
	err = r.saveTargetMetrics(ctx, tx, metrics.TargetMetrics, metricsDb.ID)
	if err != nil {
		return util.StatusWrap(err, "Failed to save target metrics")
	}
	err = r.savePackageMetrics(ctx, tx, metrics.PackageMetrics, metricsDb.ID)
	if err != nil {
		return util.StatusWrap(err, "Failed to save package metrics")
	}
	err = r.saveCumulativeMetrics(ctx, tx, metrics.CumulativeMetrics, metricsDb.ID)
	if err != nil {
		return util.StatusWrap(err, "Failed to save cumulative metrics")
	}
	err = r.saveTimingMetrics(ctx, tx, metrics.TimingMetrics, metricsDb.ID)
	if err != nil {
		return util.StatusWrap(err, "Failed to save timing metrics")
	}
	err = r.saveWorkerMetrics(ctx, tx, metrics.WorkerMetrics, metricsDb.ID)
	if err != nil {
		return util.StatusWrap(err, "Failed to save worker metrics")
	}
	err = r.saveWorkerPoolMetrics(ctx, tx, metrics.WorkerPoolMetrics, metricsDb.ID)
	if err != nil {
		return util.StatusWrap(err, "Failed to save worker pool metrics")
	}
	return nil
}
