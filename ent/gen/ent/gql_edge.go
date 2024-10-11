// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
)

func (acs *ActionCacheStatistics) ActionSummary(ctx context.Context) (*ActionSummary, error) {
	result, err := acs.Edges.ActionSummaryOrErr()
	if IsNotLoaded(err) {
		result, err = acs.QueryActionSummary().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (acs *ActionCacheStatistics) MissDetails(ctx context.Context) (result []*MissDetail, err error) {
	if fc := graphql.GetFieldContext(ctx); fc != nil && fc.Field.Alias != "" {
		result, err = acs.NamedMissDetails(graphql.GetFieldContext(ctx).Field.Alias)
	} else {
		result, err = acs.Edges.MissDetailsOrErr()
	}
	if IsNotLoaded(err) {
		result, err = acs.QueryMissDetails().All(ctx)
	}
	return result, err
}

func (ad *ActionData) ActionSummary(ctx context.Context) (*ActionSummary, error) {
	result, err := ad.Edges.ActionSummaryOrErr()
	if IsNotLoaded(err) {
		result, err = ad.QueryActionSummary().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (as *ActionSummary) Metrics(ctx context.Context) (*Metrics, error) {
	result, err := as.Edges.MetricsOrErr()
	if IsNotLoaded(err) {
		result, err = as.QueryMetrics().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (as *ActionSummary) ActionData(ctx context.Context) (result []*ActionData, err error) {
	if fc := graphql.GetFieldContext(ctx); fc != nil && fc.Field.Alias != "" {
		result, err = as.NamedActionData(graphql.GetFieldContext(ctx).Field.Alias)
	} else {
		result, err = as.Edges.ActionDataOrErr()
	}
	if IsNotLoaded(err) {
		result, err = as.QueryActionData().All(ctx)
	}
	return result, err
}

func (as *ActionSummary) RunnerCount(ctx context.Context) (result []*RunnerCount, err error) {
	if fc := graphql.GetFieldContext(ctx); fc != nil && fc.Field.Alias != "" {
		result, err = as.NamedRunnerCount(graphql.GetFieldContext(ctx).Field.Alias)
	} else {
		result, err = as.Edges.RunnerCountOrErr()
	}
	if IsNotLoaded(err) {
		result, err = as.QueryRunnerCount().All(ctx)
	}
	return result, err
}

func (as *ActionSummary) ActionCacheStatistics(ctx context.Context) (*ActionCacheStatistics, error) {
	result, err := as.Edges.ActionCacheStatisticsOrErr()
	if IsNotLoaded(err) {
		result, err = as.QueryActionCacheStatistics().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (am *ArtifactMetrics) Metrics(ctx context.Context) (*Metrics, error) {
	result, err := am.Edges.MetricsOrErr()
	if IsNotLoaded(err) {
		result, err = am.QueryMetrics().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (am *ArtifactMetrics) SourceArtifactsRead(ctx context.Context) (*FilesMetric, error) {
	result, err := am.Edges.SourceArtifactsReadOrErr()
	if IsNotLoaded(err) {
		result, err = am.QuerySourceArtifactsRead().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (am *ArtifactMetrics) OutputArtifactsSeen(ctx context.Context) (*FilesMetric, error) {
	result, err := am.Edges.OutputArtifactsSeenOrErr()
	if IsNotLoaded(err) {
		result, err = am.QueryOutputArtifactsSeen().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (am *ArtifactMetrics) OutputArtifactsFromActionCache(ctx context.Context) (*FilesMetric, error) {
	result, err := am.Edges.OutputArtifactsFromActionCacheOrErr()
	if IsNotLoaded(err) {
		result, err = am.QueryOutputArtifactsFromActionCache().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (am *ArtifactMetrics) TopLevelArtifacts(ctx context.Context) (*FilesMetric, error) {
	result, err := am.Edges.TopLevelArtifactsOrErr()
	if IsNotLoaded(err) {
		result, err = am.QueryTopLevelArtifacts().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (bi *BazelInvocation) EventFile(ctx context.Context) (*EventFile, error) {
	result, err := bi.Edges.EventFileOrErr()
	if IsNotLoaded(err) {
		result, err = bi.QueryEventFile().Only(ctx)
	}
	return result, err
}

func (bi *BazelInvocation) Build(ctx context.Context) (*Build, error) {
	result, err := bi.Edges.BuildOrErr()
	if IsNotLoaded(err) {
		result, err = bi.QueryBuild().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (bi *BazelInvocation) Metrics(ctx context.Context) (*Metrics, error) {
	result, err := bi.Edges.MetricsOrErr()
	if IsNotLoaded(err) {
		result, err = bi.QueryMetrics().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (bi *BazelInvocation) TestCollection(ctx context.Context) (result []*TestCollection, err error) {
	if fc := graphql.GetFieldContext(ctx); fc != nil && fc.Field.Alias != "" {
		result, err = bi.NamedTestCollection(graphql.GetFieldContext(ctx).Field.Alias)
	} else {
		result, err = bi.Edges.TestCollectionOrErr()
	}
	if IsNotLoaded(err) {
		result, err = bi.QueryTestCollection().All(ctx)
	}
	return result, err
}

func (bi *BazelInvocation) Targets(ctx context.Context) (result []*TargetPair, err error) {
	if fc := graphql.GetFieldContext(ctx); fc != nil && fc.Field.Alias != "" {
		result, err = bi.NamedTargets(graphql.GetFieldContext(ctx).Field.Alias)
	} else {
		result, err = bi.Edges.TargetsOrErr()
	}
	if IsNotLoaded(err) {
		result, err = bi.QueryTargets().All(ctx)
	}
	return result, err
}

func (bip *BazelInvocationProblem) BazelInvocation(ctx context.Context) (*BazelInvocation, error) {
	result, err := bip.Edges.BazelInvocationOrErr()
	if IsNotLoaded(err) {
		result, err = bip.QueryBazelInvocation().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (b *Build) Invocations(ctx context.Context) (result []*BazelInvocation, err error) {
	if fc := graphql.GetFieldContext(ctx); fc != nil && fc.Field.Alias != "" {
		result, err = b.NamedInvocations(graphql.GetFieldContext(ctx).Field.Alias)
	} else {
		result, err = b.Edges.InvocationsOrErr()
	}
	if IsNotLoaded(err) {
		result, err = b.QueryInvocations().All(ctx)
	}
	return result, err
}

func (bgm *BuildGraphMetrics) Metrics(ctx context.Context) (*Metrics, error) {
	result, err := bgm.Edges.MetricsOrErr()
	if IsNotLoaded(err) {
		result, err = bgm.QueryMetrics().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (bgm *BuildGraphMetrics) DirtiedValues(ctx context.Context) (*EvaluationStat, error) {
	result, err := bgm.Edges.DirtiedValuesOrErr()
	if IsNotLoaded(err) {
		result, err = bgm.QueryDirtiedValues().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (bgm *BuildGraphMetrics) ChangedValues(ctx context.Context) (*EvaluationStat, error) {
	result, err := bgm.Edges.ChangedValuesOrErr()
	if IsNotLoaded(err) {
		result, err = bgm.QueryChangedValues().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (bgm *BuildGraphMetrics) BuiltValues(ctx context.Context) (*EvaluationStat, error) {
	result, err := bgm.Edges.BuiltValuesOrErr()
	if IsNotLoaded(err) {
		result, err = bgm.QueryBuiltValues().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (bgm *BuildGraphMetrics) CleanedValues(ctx context.Context) (*EvaluationStat, error) {
	result, err := bgm.Edges.CleanedValuesOrErr()
	if IsNotLoaded(err) {
		result, err = bgm.QueryCleanedValues().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (bgm *BuildGraphMetrics) EvaluatedValues(ctx context.Context) (*EvaluationStat, error) {
	result, err := bgm.Edges.EvaluatedValuesOrErr()
	if IsNotLoaded(err) {
		result, err = bgm.QueryEvaluatedValues().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (cm *CumulativeMetrics) Metrics(ctx context.Context) (*Metrics, error) {
	result, err := cm.Edges.MetricsOrErr()
	if IsNotLoaded(err) {
		result, err = cm.QueryMetrics().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (dem *DynamicExecutionMetrics) Metrics(ctx context.Context) (*Metrics, error) {
	result, err := dem.Edges.MetricsOrErr()
	if IsNotLoaded(err) {
		result, err = dem.QueryMetrics().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (dem *DynamicExecutionMetrics) RaceStatistics(ctx context.Context) (result []*RaceStatistics, err error) {
	if fc := graphql.GetFieldContext(ctx); fc != nil && fc.Field.Alias != "" {
		result, err = dem.NamedRaceStatistics(graphql.GetFieldContext(ctx).Field.Alias)
	} else {
		result, err = dem.Edges.RaceStatisticsOrErr()
	}
	if IsNotLoaded(err) {
		result, err = dem.QueryRaceStatistics().All(ctx)
	}
	return result, err
}

func (es *EvaluationStat) BuildGraphMetrics(ctx context.Context) (*BuildGraphMetrics, error) {
	result, err := es.Edges.BuildGraphMetricsOrErr()
	if IsNotLoaded(err) {
		result, err = es.QueryBuildGraphMetrics().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (ef *EventFile) BazelInvocation(ctx context.Context) (*BazelInvocation, error) {
	result, err := ef.Edges.BazelInvocationOrErr()
	if IsNotLoaded(err) {
		result, err = ef.QueryBazelInvocation().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (ei *ExectionInfo) TestResult(ctx context.Context) (*TestResultBES, error) {
	result, err := ei.Edges.TestResultOrErr()
	if IsNotLoaded(err) {
		result, err = ei.QueryTestResult().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (ei *ExectionInfo) TimingBreakdown(ctx context.Context) (*TimingBreakdown, error) {
	result, err := ei.Edges.TimingBreakdownOrErr()
	if IsNotLoaded(err) {
		result, err = ei.QueryTimingBreakdown().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (ei *ExectionInfo) ResourceUsage(ctx context.Context) (result []*ResourceUsage, err error) {
	if fc := graphql.GetFieldContext(ctx); fc != nil && fc.Field.Alias != "" {
		result, err = ei.NamedResourceUsage(graphql.GetFieldContext(ctx).Field.Alias)
	} else {
		result, err = ei.Edges.ResourceUsageOrErr()
	}
	if IsNotLoaded(err) {
		result, err = ei.QueryResourceUsage().All(ctx)
	}
	return result, err
}

func (fm *FilesMetric) ArtifactMetrics(ctx context.Context) (*ArtifactMetrics, error) {
	result, err := fm.Edges.ArtifactMetricsOrErr()
	if IsNotLoaded(err) {
		result, err = fm.QueryArtifactMetrics().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (gm *GarbageMetrics) MemoryMetrics(ctx context.Context) (*MemoryMetrics, error) {
	result, err := gm.Edges.MemoryMetricsOrErr()
	if IsNotLoaded(err) {
		result, err = gm.QueryMemoryMetrics().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (mm *MemoryMetrics) Metrics(ctx context.Context) (*Metrics, error) {
	result, err := mm.Edges.MetricsOrErr()
	if IsNotLoaded(err) {
		result, err = mm.QueryMetrics().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (mm *MemoryMetrics) GarbageMetrics(ctx context.Context) (result []*GarbageMetrics, err error) {
	if fc := graphql.GetFieldContext(ctx); fc != nil && fc.Field.Alias != "" {
		result, err = mm.NamedGarbageMetrics(graphql.GetFieldContext(ctx).Field.Alias)
	} else {
		result, err = mm.Edges.GarbageMetricsOrErr()
	}
	if IsNotLoaded(err) {
		result, err = mm.QueryGarbageMetrics().All(ctx)
	}
	return result, err
}

func (m *Metrics) BazelInvocation(ctx context.Context) (*BazelInvocation, error) {
	result, err := m.Edges.BazelInvocationOrErr()
	if IsNotLoaded(err) {
		result, err = m.QueryBazelInvocation().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (m *Metrics) ActionSummary(ctx context.Context) (*ActionSummary, error) {
	result, err := m.Edges.ActionSummaryOrErr()
	if IsNotLoaded(err) {
		result, err = m.QueryActionSummary().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (m *Metrics) MemoryMetrics(ctx context.Context) (*MemoryMetrics, error) {
	result, err := m.Edges.MemoryMetricsOrErr()
	if IsNotLoaded(err) {
		result, err = m.QueryMemoryMetrics().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (m *Metrics) TargetMetrics(ctx context.Context) (*TargetMetrics, error) {
	result, err := m.Edges.TargetMetricsOrErr()
	if IsNotLoaded(err) {
		result, err = m.QueryTargetMetrics().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (m *Metrics) PackageMetrics(ctx context.Context) (*PackageMetrics, error) {
	result, err := m.Edges.PackageMetricsOrErr()
	if IsNotLoaded(err) {
		result, err = m.QueryPackageMetrics().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (m *Metrics) TimingMetrics(ctx context.Context) (*TimingMetrics, error) {
	result, err := m.Edges.TimingMetricsOrErr()
	if IsNotLoaded(err) {
		result, err = m.QueryTimingMetrics().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (m *Metrics) CumulativeMetrics(ctx context.Context) (*CumulativeMetrics, error) {
	result, err := m.Edges.CumulativeMetricsOrErr()
	if IsNotLoaded(err) {
		result, err = m.QueryCumulativeMetrics().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (m *Metrics) ArtifactMetrics(ctx context.Context) (*ArtifactMetrics, error) {
	result, err := m.Edges.ArtifactMetricsOrErr()
	if IsNotLoaded(err) {
		result, err = m.QueryArtifactMetrics().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (m *Metrics) NetworkMetrics(ctx context.Context) (*NetworkMetrics, error) {
	result, err := m.Edges.NetworkMetricsOrErr()
	if IsNotLoaded(err) {
		result, err = m.QueryNetworkMetrics().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (m *Metrics) DynamicExecutionMetrics(ctx context.Context) (*DynamicExecutionMetrics, error) {
	result, err := m.Edges.DynamicExecutionMetricsOrErr()
	if IsNotLoaded(err) {
		result, err = m.QueryDynamicExecutionMetrics().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (m *Metrics) BuildGraphMetrics(ctx context.Context) (*BuildGraphMetrics, error) {
	result, err := m.Edges.BuildGraphMetricsOrErr()
	if IsNotLoaded(err) {
		result, err = m.QueryBuildGraphMetrics().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (md *MissDetail) ActionCacheStatistics(ctx context.Context) (*ActionCacheStatistics, error) {
	result, err := md.Edges.ActionCacheStatisticsOrErr()
	if IsNotLoaded(err) {
		result, err = md.QueryActionCacheStatistics().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (nsof *NamedSetOfFiles) OutputGroup(ctx context.Context) (*OutputGroup, error) {
	result, err := nsof.Edges.OutputGroupOrErr()
	if IsNotLoaded(err) {
		result, err = nsof.QueryOutputGroup().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (nsof *NamedSetOfFiles) Files(ctx context.Context) (result []*TestFile, err error) {
	if fc := graphql.GetFieldContext(ctx); fc != nil && fc.Field.Alias != "" {
		result, err = nsof.NamedFiles(graphql.GetFieldContext(ctx).Field.Alias)
	} else {
		result, err = nsof.Edges.FilesOrErr()
	}
	if IsNotLoaded(err) {
		result, err = nsof.QueryFiles().All(ctx)
	}
	return result, err
}

func (nsof *NamedSetOfFiles) FileSets(ctx context.Context) (*NamedSetOfFiles, error) {
	result, err := nsof.Edges.FileSetsOrErr()
	if IsNotLoaded(err) {
		result, err = nsof.QueryFileSets().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (nm *NetworkMetrics) Metrics(ctx context.Context) (*Metrics, error) {
	result, err := nm.Edges.MetricsOrErr()
	if IsNotLoaded(err) {
		result, err = nm.QueryMetrics().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (nm *NetworkMetrics) SystemNetworkStats(ctx context.Context) (*SystemNetworkStats, error) {
	result, err := nm.Edges.SystemNetworkStatsOrErr()
	if IsNotLoaded(err) {
		result, err = nm.QuerySystemNetworkStats().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (og *OutputGroup) TargetComplete(ctx context.Context) (*TargetComplete, error) {
	result, err := og.Edges.TargetCompleteOrErr()
	if IsNotLoaded(err) {
		result, err = og.QueryTargetComplete().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (og *OutputGroup) InlineFiles(ctx context.Context) (result []*TestFile, err error) {
	if fc := graphql.GetFieldContext(ctx); fc != nil && fc.Field.Alias != "" {
		result, err = og.NamedInlineFiles(graphql.GetFieldContext(ctx).Field.Alias)
	} else {
		result, err = og.Edges.InlineFilesOrErr()
	}
	if IsNotLoaded(err) {
		result, err = og.QueryInlineFiles().All(ctx)
	}
	return result, err
}

func (og *OutputGroup) FileSets(ctx context.Context) (*NamedSetOfFiles, error) {
	result, err := og.Edges.FileSetsOrErr()
	if IsNotLoaded(err) {
		result, err = og.QueryFileSets().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (plm *PackageLoadMetrics) PackageMetrics(ctx context.Context) (*PackageMetrics, error) {
	result, err := plm.Edges.PackageMetricsOrErr()
	if IsNotLoaded(err) {
		result, err = plm.QueryPackageMetrics().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (pm *PackageMetrics) Metrics(ctx context.Context) (*Metrics, error) {
	result, err := pm.Edges.MetricsOrErr()
	if IsNotLoaded(err) {
		result, err = pm.QueryMetrics().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (pm *PackageMetrics) PackageLoadMetrics(ctx context.Context) (result []*PackageLoadMetrics, err error) {
	if fc := graphql.GetFieldContext(ctx); fc != nil && fc.Field.Alias != "" {
		result, err = pm.NamedPackageLoadMetrics(graphql.GetFieldContext(ctx).Field.Alias)
	} else {
		result, err = pm.Edges.PackageLoadMetricsOrErr()
	}
	if IsNotLoaded(err) {
		result, err = pm.QueryPackageLoadMetrics().All(ctx)
	}
	return result, err
}

func (rs *RaceStatistics) DynamicExecutionMetrics(ctx context.Context) (*DynamicExecutionMetrics, error) {
	result, err := rs.Edges.DynamicExecutionMetricsOrErr()
	if IsNotLoaded(err) {
		result, err = rs.QueryDynamicExecutionMetrics().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (ru *ResourceUsage) ExecutionInfo(ctx context.Context) (*ExectionInfo, error) {
	result, err := ru.Edges.ExecutionInfoOrErr()
	if IsNotLoaded(err) {
		result, err = ru.QueryExecutionInfo().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (rc *RunnerCount) ActionSummary(ctx context.Context) (*ActionSummary, error) {
	result, err := rc.Edges.ActionSummaryOrErr()
	if IsNotLoaded(err) {
		result, err = rc.QueryActionSummary().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (sns *SystemNetworkStats) NetworkMetrics(ctx context.Context) (*NetworkMetrics, error) {
	result, err := sns.Edges.NetworkMetricsOrErr()
	if IsNotLoaded(err) {
		result, err = sns.QueryNetworkMetrics().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (tc *TargetComplete) TargetPair(ctx context.Context) (*TargetPair, error) {
	result, err := tc.Edges.TargetPairOrErr()
	if IsNotLoaded(err) {
		result, err = tc.QueryTargetPair().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (tc *TargetComplete) ImportantOutput(ctx context.Context) (result []*TestFile, err error) {
	if fc := graphql.GetFieldContext(ctx); fc != nil && fc.Field.Alias != "" {
		result, err = tc.NamedImportantOutput(graphql.GetFieldContext(ctx).Field.Alias)
	} else {
		result, err = tc.Edges.ImportantOutputOrErr()
	}
	if IsNotLoaded(err) {
		result, err = tc.QueryImportantOutput().All(ctx)
	}
	return result, err
}

func (tc *TargetComplete) DirectoryOutput(ctx context.Context) (result []*TestFile, err error) {
	if fc := graphql.GetFieldContext(ctx); fc != nil && fc.Field.Alias != "" {
		result, err = tc.NamedDirectoryOutput(graphql.GetFieldContext(ctx).Field.Alias)
	} else {
		result, err = tc.Edges.DirectoryOutputOrErr()
	}
	if IsNotLoaded(err) {
		result, err = tc.QueryDirectoryOutput().All(ctx)
	}
	return result, err
}

func (tc *TargetComplete) OutputGroup(ctx context.Context) (*OutputGroup, error) {
	result, err := tc.Edges.OutputGroupOrErr()
	if IsNotLoaded(err) {
		result, err = tc.QueryOutputGroup().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (tc *TargetConfigured) TargetPair(ctx context.Context) (*TargetPair, error) {
	result, err := tc.Edges.TargetPairOrErr()
	if IsNotLoaded(err) {
		result, err = tc.QueryTargetPair().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (tm *TargetMetrics) Metrics(ctx context.Context) (*Metrics, error) {
	result, err := tm.Edges.MetricsOrErr()
	if IsNotLoaded(err) {
		result, err = tm.QueryMetrics().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (tp *TargetPair) BazelInvocation(ctx context.Context) (*BazelInvocation, error) {
	result, err := tp.Edges.BazelInvocationOrErr()
	if IsNotLoaded(err) {
		result, err = tp.QueryBazelInvocation().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (tp *TargetPair) Configuration(ctx context.Context) (*TargetConfigured, error) {
	result, err := tp.Edges.ConfigurationOrErr()
	if IsNotLoaded(err) {
		result, err = tp.QueryConfiguration().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (tp *TargetPair) Completion(ctx context.Context) (*TargetComplete, error) {
	result, err := tp.Edges.CompletionOrErr()
	if IsNotLoaded(err) {
		result, err = tp.QueryCompletion().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (tc *TestCollection) BazelInvocation(ctx context.Context) (*BazelInvocation, error) {
	result, err := tc.Edges.BazelInvocationOrErr()
	if IsNotLoaded(err) {
		result, err = tc.QueryBazelInvocation().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (tc *TestCollection) TestSummary(ctx context.Context) (*TestSummary, error) {
	result, err := tc.Edges.TestSummaryOrErr()
	if IsNotLoaded(err) {
		result, err = tc.QueryTestSummary().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (tc *TestCollection) TestResults(ctx context.Context) (result []*TestResultBES, err error) {
	if fc := graphql.GetFieldContext(ctx); fc != nil && fc.Field.Alias != "" {
		result, err = tc.NamedTestResults(graphql.GetFieldContext(ctx).Field.Alias)
	} else {
		result, err = tc.Edges.TestResultsOrErr()
	}
	if IsNotLoaded(err) {
		result, err = tc.QueryTestResults().All(ctx)
	}
	return result, err
}

func (tf *TestFile) TestResult(ctx context.Context) (*TestResultBES, error) {
	result, err := tf.Edges.TestResultOrErr()
	if IsNotLoaded(err) {
		result, err = tf.QueryTestResult().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (trb *TestResultBES) TestCollection(ctx context.Context) (*TestCollection, error) {
	result, err := trb.Edges.TestCollectionOrErr()
	if IsNotLoaded(err) {
		result, err = trb.QueryTestCollection().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (trb *TestResultBES) TestActionOutput(ctx context.Context) (result []*TestFile, err error) {
	if fc := graphql.GetFieldContext(ctx); fc != nil && fc.Field.Alias != "" {
		result, err = trb.NamedTestActionOutput(graphql.GetFieldContext(ctx).Field.Alias)
	} else {
		result, err = trb.Edges.TestActionOutputOrErr()
	}
	if IsNotLoaded(err) {
		result, err = trb.QueryTestActionOutput().All(ctx)
	}
	return result, err
}

func (trb *TestResultBES) ExecutionInfo(ctx context.Context) (*ExectionInfo, error) {
	result, err := trb.Edges.ExecutionInfoOrErr()
	if IsNotLoaded(err) {
		result, err = trb.QueryExecutionInfo().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (ts *TestSummary) TestCollection(ctx context.Context) (*TestCollection, error) {
	result, err := ts.Edges.TestCollectionOrErr()
	if IsNotLoaded(err) {
		result, err = ts.QueryTestCollection().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (ts *TestSummary) Passed(ctx context.Context) (result []*TestFile, err error) {
	if fc := graphql.GetFieldContext(ctx); fc != nil && fc.Field.Alias != "" {
		result, err = ts.NamedPassed(graphql.GetFieldContext(ctx).Field.Alias)
	} else {
		result, err = ts.Edges.PassedOrErr()
	}
	if IsNotLoaded(err) {
		result, err = ts.QueryPassed().All(ctx)
	}
	return result, err
}

func (ts *TestSummary) Failed(ctx context.Context) (result []*TestFile, err error) {
	if fc := graphql.GetFieldContext(ctx); fc != nil && fc.Field.Alias != "" {
		result, err = ts.NamedFailed(graphql.GetFieldContext(ctx).Field.Alias)
	} else {
		result, err = ts.Edges.FailedOrErr()
	}
	if IsNotLoaded(err) {
		result, err = ts.QueryFailed().All(ctx)
	}
	return result, err
}

func (tb *TimingBreakdown) ExecutionInfo(ctx context.Context) (*ExectionInfo, error) {
	result, err := tb.Edges.ExecutionInfoOrErr()
	if IsNotLoaded(err) {
		result, err = tb.QueryExecutionInfo().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (tb *TimingBreakdown) Child(ctx context.Context) (result []*TimingChild, err error) {
	if fc := graphql.GetFieldContext(ctx); fc != nil && fc.Field.Alias != "" {
		result, err = tb.NamedChild(graphql.GetFieldContext(ctx).Field.Alias)
	} else {
		result, err = tb.Edges.ChildOrErr()
	}
	if IsNotLoaded(err) {
		result, err = tb.QueryChild().All(ctx)
	}
	return result, err
}

func (tc *TimingChild) TimingBreakdown(ctx context.Context) (*TimingBreakdown, error) {
	result, err := tc.Edges.TimingBreakdownOrErr()
	if IsNotLoaded(err) {
		result, err = tc.QueryTimingBreakdown().Only(ctx)
	}
	return result, MaskNotFound(err)
}

func (tm *TimingMetrics) Metrics(ctx context.Context) (*Metrics, error) {
	result, err := tm.Edges.MetricsOrErr()
	if IsNotLoaded(err) {
		result, err = tm.QueryMetrics().Only(ctx)
	}
	return result, MaskNotFound(err)
}
