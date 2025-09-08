package auth

import (
	"context"
	"fmt"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/internal/api/common"
	"github.com/buildbarn/bb-storage/pkg/auth"
	auth_configuration "github.com/buildbarn/bb-storage/pkg/auth/configuration"
	bb_grpc "github.com/buildbarn/bb-storage/pkg/grpc"
	"github.com/buildbarn/bb-storage/pkg/program"
	auth_pb "github.com/buildbarn/bb-storage/pkg/proto/configuration/auth"
	"github.com/buildbarn/bb-storage/pkg/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func createInterceptor[T any](
	authorizer auth.Authorizer,
	shouldIncludeItem func(context.Context, auth.Authorizer, *T) bool,
) ent.InterceptFunc {
	return ent.InterceptFunc(func(next ent.Querier) ent.Querier {
		return ent.QuerierFunc(func(ctx context.Context, query ent.Query) (ent.Value, error) {
			value, err := next.Query(ctx, query)
			if err != nil {
				return nil, err
			}

			items, ok := value.([]*T)
			if !ok {
				return nil, fmt.Errorf("expected value to be of type %T, got %T", *new(T), value)
			}

			filteredItems := make([]*T, 0)
			for _, item := range items {
				if shouldIncludeItem(ctx, authorizer, item) {
					filteredItems = append(filteredItems, item)
				}
			}

			return filteredItems, nil
		})
	})
}

// AddDatabaseAuthInterceptors adds interceptors to the ent.Client for filtering
// database queries based on the instanceNameAuthorizer.
func AddDatabaseAuthInterceptors(authorizerConfiguration *auth_pb.AuthorizerConfiguration, dbClient *ent.Client, dependenciesGroup program.Group, grpcClientFactory bb_grpc.ClientFactory) error {
	if authorizerConfiguration == nil {
		return status.Error(codes.NotFound, "No InstanceNameAuthorizer configured")
	}
	if _, ok := authorizerConfiguration.Policy.(*auth_pb.AuthorizerConfiguration_Allow); ok {
		// If the policy is "Allow", we don't need to add any interceptors. It
		// will just slow everything down.
		return nil
	}
	instanceNameAuthorizer, err := auth_configuration.DefaultAuthorizerFactory.NewAuthorizerFromConfiguration(authorizerConfiguration, dependenciesGroup, grpcClientFactory)
	if err != nil {
		return util.StatusWrap(err, "Failed to create InstanceNameAuthorizer")
	}

	dbClient.ActionCacheStatistics.Intercept(createInterceptor(instanceNameAuthorizer, isActionCacheStatisticsAllowed))
	dbClient.ActionData.Intercept(createInterceptor(instanceNameAuthorizer, isActionDataAllowed))
	dbClient.ActionSummary.Intercept(createInterceptor(instanceNameAuthorizer, isActionSummaryAllowed))
	dbClient.ArtifactMetrics.Intercept(createInterceptor(instanceNameAuthorizer, isArtifactMetricsAllowed))
	dbClient.BazelInvocation.Intercept(createInterceptor(instanceNameAuthorizer, isBazelInvocationAllowed))
	dbClient.BazelInvocationProblem.Intercept(createInterceptor(instanceNameAuthorizer, isBazelInvocationProblemAllowed))
	dbClient.Blob.Intercept(createInterceptor(instanceNameAuthorizer, isBlobAllowed))
	dbClient.Build.Intercept(createInterceptor(instanceNameAuthorizer, isBuildAllowed))
	dbClient.BuildGraphMetrics.Intercept(createInterceptor(instanceNameAuthorizer, isBuildGraphMetricsAllowed))
	dbClient.CumulativeMetrics.Intercept(createInterceptor(instanceNameAuthorizer, isCumulativeMetricsAllowed))
	dbClient.DynamicExecutionMetrics.Intercept(createInterceptor(instanceNameAuthorizer, isDynamicExecutionMetricsAllowed))
	dbClient.EvaluationStat.Intercept(createInterceptor(instanceNameAuthorizer, isEvaluationStatAllowed))
	dbClient.EventFile.Intercept(createInterceptor(instanceNameAuthorizer, isEventFileAllowed))
	dbClient.ExectionInfo.Intercept(createInterceptor(instanceNameAuthorizer, isExectionInfoAllowed))
	dbClient.FilesMetric.Intercept(createInterceptor(instanceNameAuthorizer, isFilesMetricAllowed))
	dbClient.GarbageMetrics.Intercept(createInterceptor(instanceNameAuthorizer, isGarbageMetricsAllowed))
	dbClient.MemoryMetrics.Intercept(createInterceptor(instanceNameAuthorizer, isMemoryMetricsAllowed))
	dbClient.Metrics.Intercept(createInterceptor(instanceNameAuthorizer, isMetricsAllowed))
	dbClient.MissDetail.Intercept(createInterceptor(instanceNameAuthorizer, isMissDetailAllowed))
	dbClient.NamedSetOfFiles.Intercept(createInterceptor(instanceNameAuthorizer, isNamedSetOfFilesAllowed))
	dbClient.NetworkMetrics.Intercept(createInterceptor(instanceNameAuthorizer, isNetworkMetricsAllowed))
	dbClient.OutputGroup.Intercept(createInterceptor(instanceNameAuthorizer, isOutputGroupAllowed))
	dbClient.PackageLoadMetrics.Intercept(createInterceptor(instanceNameAuthorizer, isPackageLoadMetricsAllowed))
	dbClient.PackageMetrics.Intercept(createInterceptor(instanceNameAuthorizer, isPackageMetricsAllowed))
	dbClient.RaceStatistics.Intercept(createInterceptor(instanceNameAuthorizer, isRaceStatisticsAllowed))
	dbClient.ResourceUsage.Intercept(createInterceptor(instanceNameAuthorizer, isResourceUsageAllowed))
	dbClient.RunnerCount.Intercept(createInterceptor(instanceNameAuthorizer, isRunnerCountAllowed))
	dbClient.SourceControl.Intercept(createInterceptor(instanceNameAuthorizer, isSourceControlAllowed))
	dbClient.SystemNetworkStats.Intercept(createInterceptor(instanceNameAuthorizer, isSystemNetworkStatsAllowed))
	dbClient.TargetComplete.Intercept(createInterceptor(instanceNameAuthorizer, isTargetCompleteAllowed))
	dbClient.TargetConfigured.Intercept(createInterceptor(instanceNameAuthorizer, isTargetConfiguredAllowed))
	dbClient.TargetMetrics.Intercept(createInterceptor(instanceNameAuthorizer, isTargetMetricsAllowed))
	dbClient.TargetPair.Intercept(createInterceptor(instanceNameAuthorizer, isTargetPairAllowed))
	dbClient.TestCollection.Intercept(createInterceptor(instanceNameAuthorizer, isTestCollectionAllowed))
	dbClient.TestFile.Intercept(createInterceptor(instanceNameAuthorizer, isTestFileAllowed))
	dbClient.TestResultBES.Intercept(createInterceptor(instanceNameAuthorizer, isTestResultBESAllowed))
	dbClient.TestSummary.Intercept(createInterceptor(instanceNameAuthorizer, isTestSummaryAllowed))
	dbClient.TimingBreakdown.Intercept(createInterceptor(instanceNameAuthorizer, isTimingBreakdownAllowed))
	dbClient.TimingChild.Intercept(createInterceptor(instanceNameAuthorizer, isTimingChildAllowed))
	dbClient.TimingMetrics.Intercept(createInterceptor(instanceNameAuthorizer, isTimingMetricsAllowed))
	return nil
}

func isActionCacheStatisticsAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, actionCacheStatistics *ent.ActionCacheStatistics) bool {
	actionSummary, err := actionCacheStatistics.ActionSummary(ctx)
	return err == nil && actionSummary != nil
}

func isActionDataAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, actionData *ent.ActionData) bool {
	actionSummary, err := actionData.ActionSummary(ctx)
	return err == nil && actionSummary != nil
}

func isActionSummaryAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, actionSummary *ent.ActionSummary) bool {
	metrics, err := actionSummary.Metrics(ctx)
	return err == nil && metrics != nil
}

func isArtifactMetricsAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, artifactMetrics *ent.ArtifactMetrics) bool {
	metrics, err := artifactMetrics.Metrics(ctx)
	return err == nil && metrics != nil
}

func isBazelInvocationAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, invocation *ent.BazelInvocation) bool {
	return common.IsInstanceNameAllowed(ctx, instanceNameAuthorizer, invocation.InstanceName)
}

func isBazelInvocationProblemAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, problem *ent.BazelInvocationProblem) bool {
	invocation, err := problem.BazelInvocation(ctx)
	return err == nil && invocation != nil
}

func isBlobAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, blob *ent.Blob) bool {
	return common.IsInstanceNameAllowed(ctx, instanceNameAuthorizer, blob.InstanceName)
}

func isBuildAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, build *ent.Build) bool {
	invocations, err := build.Invocations(ctx)
	return err == nil && invocations != nil && len(invocations) > 0
}

func isBuildGraphMetricsAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, buildGraphMetrics *ent.BuildGraphMetrics) bool {
	metrics, err := buildGraphMetrics.Metrics(ctx)
	return err == nil && metrics != nil
}

func isCumulativeMetricsAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, cumulativeMetrics *ent.CumulativeMetrics) bool {
	metrics, err := cumulativeMetrics.Metrics(ctx)
	return err == nil && metrics != nil
}

func isDynamicExecutionMetricsAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, dynamicExecutionMetrics *ent.DynamicExecutionMetrics) bool {
	metrics, err := dynamicExecutionMetrics.Metrics(ctx)
	return err == nil && metrics != nil
}

func isEvaluationStatAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, evaluationStat *ent.EvaluationStat) bool {
	buildGraphMetrics, err := evaluationStat.BuildGraphMetrics(ctx)
	return err == nil && buildGraphMetrics != nil
}

func isEventFileAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, eventFile *ent.EventFile) bool {
	invocation, err := eventFile.BazelInvocation(ctx)
	return err == nil && invocation != nil
}

func isExectionInfoAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, exectionInfo *ent.ExectionInfo) bool {
	testResult, err := exectionInfo.TestResult(ctx)
	return err == nil && testResult != nil
}

func isFilesMetricAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, filesMetric *ent.FilesMetric) bool {
	artifactMetrics, err := filesMetric.ArtifactMetrics(ctx)
	return err == nil && artifactMetrics != nil
}

func isGarbageMetricsAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, garbageMetrics *ent.GarbageMetrics) bool {
	memoryMetrics, err := garbageMetrics.MemoryMetrics(ctx)
	return err == nil && memoryMetrics != nil
}

func isMemoryMetricsAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, memoryMetrics *ent.MemoryMetrics) bool {
	metrics, err := memoryMetrics.Metrics(ctx)
	return err == nil && metrics != nil
}

func isMetricsAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, metrics *ent.Metrics) bool {
	invocation, err := metrics.BazelInvocation(ctx)
	return err == nil && invocation != nil
}

func isMissDetailAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, missDetail *ent.MissDetail) bool {
	actionCacheStatistics, err := missDetail.ActionCacheStatistics(ctx)
	return err == nil && actionCacheStatistics != nil
}

func isNamedSetOfFilesAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, namedSetOfFiles *ent.NamedSetOfFiles) bool {
	testFiles, err := namedSetOfFiles.Files(ctx)
	return err == nil && testFiles != nil && len(testFiles) > 0
}

func isNetworkMetricsAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, networkMetrics *ent.NetworkMetrics) bool {
	metrics, err := networkMetrics.Metrics(ctx)
	return err == nil && metrics != nil
}

func isOutputGroupAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, outputGroup *ent.OutputGroup) bool {
	targetComplete, err := outputGroup.TargetComplete(ctx)
	return err == nil && targetComplete != nil
}

func isPackageLoadMetricsAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, packageLoadMetrics *ent.PackageLoadMetrics) bool {
	packageMetrics, err := packageLoadMetrics.PackageMetrics(ctx)
	return err == nil && packageMetrics != nil
}

func isPackageMetricsAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, packageMetrics *ent.PackageMetrics) bool {
	metrics, err := packageMetrics.Metrics(ctx)
	return err == nil && metrics != nil
}

func isRaceStatisticsAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, raceStatistics *ent.RaceStatistics) bool {
	dynamicExecutionMetrics, err := raceStatistics.DynamicExecutionMetrics(ctx)
	return err == nil && dynamicExecutionMetrics != nil
}

func isResourceUsageAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, resourceUsage *ent.ResourceUsage) bool {
	executionInfo, err := resourceUsage.ExecutionInfo(ctx)
	return err == nil && executionInfo != nil
}

func isRunnerCountAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, runnerCount *ent.RunnerCount) bool {
	actionSummary, err := runnerCount.ActionSummary(ctx)
	return err == nil && actionSummary != nil
}

func isSourceControlAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, sourceControl *ent.SourceControl) bool {
	invocation, err := sourceControl.BazelInvocation(ctx)
	return err == nil && invocation != nil
}

func isSystemNetworkStatsAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, systemNetworkStats *ent.SystemNetworkStats) bool {
	networkMetrics, err := systemNetworkStats.NetworkMetrics(ctx)
	return err == nil && networkMetrics != nil
}

func isTargetCompleteAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, targetComplete *ent.TargetComplete) bool {
	targetPair, err := targetComplete.TargetPair(ctx)
	return err == nil && targetPair != nil
}

func isTargetConfiguredAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, targetConfigured *ent.TargetConfigured) bool {
	targetPair, err := targetConfigured.TargetPair(ctx)
	return err == nil && targetPair != nil
}

func isTargetMetricsAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, targetMetrics *ent.TargetMetrics) bool {
	metrics, err := targetMetrics.Metrics(ctx)
	return err == nil && metrics != nil
}

func isTargetPairAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, targetPair *ent.TargetPair) bool {
	invocation, err := targetPair.BazelInvocation(ctx)
	return err == nil && invocation != nil
}

func isTestCollectionAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, testCollection *ent.TestCollection) bool {
	invocation, err := testCollection.BazelInvocation(ctx)
	return err == nil && invocation != nil
}

func isTestFileAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, testFile *ent.TestFile) bool {
	testResult, err := testFile.TestResult(ctx)
	return err == nil && testResult != nil
}

func isTestResultBESAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, testResultBes *ent.TestResultBES) bool {
	testCollection, err := testResultBes.TestCollection(ctx)
	return err == nil && testCollection != nil
}

func isTestSummaryAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, testSummary *ent.TestSummary) bool {
	testCollection, err := testSummary.TestCollection(ctx)
	return err == nil && testCollection != nil
}

func isTimingBreakdownAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, timingBreakdown *ent.TimingBreakdown) bool {
	executionInfo, err := timingBreakdown.ExecutionInfo(ctx)
	return err == nil && executionInfo != nil
}

func isTimingChildAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, timingChild *ent.TimingChild) bool {
	timingBreakdown, err := timingChild.TimingBreakdown(ctx)
	return err == nil && timingBreakdown != nil
}

func isTimingMetricsAllowed(ctx context.Context, instanceNameAuthorizer auth.Authorizer, timingMetrics *ent.TimingMetrics) bool {
	metrics, err := timingMetrics.Metrics(ctx)
	return err == nil && metrics != nil
}
