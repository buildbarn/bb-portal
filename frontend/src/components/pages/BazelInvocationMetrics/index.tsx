import { Flex, Space } from "antd";
import ActionStatisticsDisplay from "@/components/ActionStatisticsDisplay";
import { ArtifactsMetricsDisplay } from "@/components/ArtifactsMetricsDisplay";
import { BazelServerMetricsDisplay } from "@/components/BazelServerMetricsDisplay";
import { BuildGraphEvaluationMetricsDisplay } from "@/components/BuildGraphEvaluationMetricsDisplay";
import { BuildGraphMetricsDisplay } from "@/components/BuildGraphMetricsDisplay";
import { DynamicExecutionMetricsDisplay } from "@/components/DynamicExecutionMetricsDisplay";
import { GarbageCollectionMetrics } from "@/components/GarbageCollectionMetrics";
import MemoryMetricsDisplay from "@/components/MemoryMetrics";
import { PackageLoadMetricsDisplay } from "@/components/PackageLoadMetricsDisplay";
import { SystemNetworkStatsDisplay } from "@/components/SystemNetworkStatsDisplay";
import { TimingMetricsDisplay } from "@/components/TimingMetricsDisplay";
import { WorkerMetricsDisplay } from "@/components/WorkerMetricsDisplay";
import { getFragmentData } from "@/graphql/__generated__";
import type { BazelInvocationMetricsFragment } from "@/graphql/__generated__/graphql";
import {
  BAZEL_INVOCATION_METRICS_ACTION_SUMMARY_FRAGMENT,
  BAZEL_INVOCATION_METRICS_ARTIFACT_METRICS_FRAGMENT,
  BAZEL_INVOCATION_METRICS_BUILD_GRAPH_EVALUATION_METRICS_FRAGMENT,
  BAZEL_INVOCATION_METRICS_BUILD_GRAPH_METRICS_FRAGMENT,
  BAZEL_INVOCATION_METRICS_CUMULATIVE_METRICS_FRAGMENT,
  BAZEL_INVOCATION_METRICS_DYNAMIC_EXECUTION_METRICS_FRAGMENT,
  BAZEL_INVOCATION_METRICS_GARBAGE_METRICS_FRAGMENT,
  BAZEL_INVOCATION_METRICS_MEMORY_METRICS_FRAGMENT,
  BAZEL_INVOCATION_METRICS_PACKAGE_METRICS_FRAGMENT,
  BAZEL_INVOCATION_METRICS_SYSTEM_NETWORK_STATS_FRAGMENT,
  BAZEL_INVOCATION_METRICS_TIMING_METRICS_FRAGMENT,
  BAZEL_INVOCATION_METRICS_WORKER_METRICS_FRAGMENT,
  BAZEL_INVOCATION_METRICS_WORKER_POOL_METRICS_FRAGMENT,
} from "@/routes/bazel-invocations.$invocationID/metrics";

const CARD_STYLE: React.CSSProperties = {
  width: "750px",
};

interface Props {
  metrics: BazelInvocationMetricsFragment;
}

export const BazelInvocationMetrics: React.FC<Props> = ({ metrics }) => {
  const actionSummary = getFragmentData(
    BAZEL_INVOCATION_METRICS_ACTION_SUMMARY_FRAGMENT,
    metrics.actionSummary,
  );
  const artifactMetrics = getFragmentData(
    BAZEL_INVOCATION_METRICS_ARTIFACT_METRICS_FRAGMENT,
    metrics.artifactMetrics,
  );
  const memoryMetrics = getFragmentData(
    BAZEL_INVOCATION_METRICS_MEMORY_METRICS_FRAGMENT,
    metrics.memoryMetrics,
  );
  const packageMetrics = getFragmentData(
    BAZEL_INVOCATION_METRICS_PACKAGE_METRICS_FRAGMENT,
    metrics.packageMetrics,
  );
  const cumulativeMetrics = getFragmentData(
    BAZEL_INVOCATION_METRICS_CUMULATIVE_METRICS_FRAGMENT,
    metrics.cumulativeMetrics,
  );
  const buildGraphEvaluationMetrics = getFragmentData(
    BAZEL_INVOCATION_METRICS_BUILD_GRAPH_EVALUATION_METRICS_FRAGMENT,
    metrics.buildGraphMetrics,
  );
  const buildGraphMetrics = getFragmentData(
    BAZEL_INVOCATION_METRICS_BUILD_GRAPH_METRICS_FRAGMENT,
    metrics.buildGraphMetrics,
  );
  const dynamicExecutionMetrics = getFragmentData(
    BAZEL_INVOCATION_METRICS_DYNAMIC_EXECUTION_METRICS_FRAGMENT,
    metrics.dynamicExecutionMetrics,
  );
  const workerMetrics = getFragmentData(
    BAZEL_INVOCATION_METRICS_WORKER_METRICS_FRAGMENT,
    metrics.workerMetrics,
  );
  const workerPoolMetrics = getFragmentData(
    BAZEL_INVOCATION_METRICS_WORKER_POOL_METRICS_FRAGMENT,
    metrics.workerPoolMetrics,
  );
  const buildGraphEvaluationStats =
    buildGraphEvaluationMetrics?.evaluationStats ?? [];
  const packageLoadMetrics = packageMetrics?.packageLoadMetrics ?? [];
  const buildGraphRuleClassCounts =
    buildGraphEvaluationMetrics?.ruleClassCounts ?? [];
  const buildGraphAspectCounts =
    buildGraphEvaluationMetrics?.aspectCounts ?? [];
  const dynamicExecutionRaceStatistics =
    dynamicExecutionMetrics?.raceStatistics ?? [];
  const workers = workerMetrics ?? [];
  const garbageMetrics = getFragmentData(
    BAZEL_INVOCATION_METRICS_GARBAGE_METRICS_FRAGMENT,
    memoryMetrics?.garbageMetrics,
  );
  const timingMetrics = getFragmentData(
    BAZEL_INVOCATION_METRICS_TIMING_METRICS_FRAGMENT,
    metrics.timingMetrics,
  );
  const systemNetworkStats = getFragmentData(
    BAZEL_INVOCATION_METRICS_SYSTEM_NETWORK_STATS_FRAGMENT,
    metrics.networkMetrics?.systemNetworkStats,
  );

  return (
    <Space orientation="vertical" size="middle">
      <Flex vertical={false} gap="small" wrap={true}>
        {timingMetrics && (
          <TimingMetricsDisplay
            timingMetrics={timingMetrics}
            cardStyle={CARD_STYLE}
          />
        )}
        {systemNetworkStats && (
          <SystemNetworkStatsDisplay
            systemNetworkStats={systemNetworkStats}
            cardStyle={CARD_STYLE}
          />
        )}
        {artifactMetrics && (
          <ArtifactsMetricsDisplay
            artifactMetrics={artifactMetrics}
            cardStyle={CARD_STYLE}
          />
        )}
        {memoryMetrics && (
          <MemoryMetricsDisplay
            memoryMetrics={memoryMetrics}
            cardStyle={CARD_STYLE}
          />
        )}
        {buildGraphMetrics && (
          <BuildGraphMetricsDisplay
            buildGraphMetrics={buildGraphMetrics}
            cardStyle={CARD_STYLE}
          />
        )}
      </Flex>
      {actionSummary && (
        <ActionStatisticsDisplay actionSummary={actionSummary} />
      )}
      <Flex vertical={false} gap="small" wrap={true}>
        {(packageMetrics || cumulativeMetrics) && (
          <BazelServerMetricsDisplay
            packageMetrics={packageMetrics}
            cumulativeMetrics={cumulativeMetrics}
            cardStyle={CARD_STYLE}
          />
        )}
        {garbageMetrics && (
          <GarbageCollectionMetrics
            garbageMetrics={garbageMetrics}
            cardStyle={CARD_STYLE}
          />
        )}
      </Flex>
      {packageMetrics && packageLoadMetrics.length > 0 && (
        <PackageLoadMetricsDisplay packageMetrics={packageMetrics} />
      )}
      {buildGraphEvaluationMetrics &&
        (buildGraphEvaluationStats.length > 0 ||
          buildGraphRuleClassCounts.length > 0 ||
          buildGraphAspectCounts.length > 0) && (
          <BuildGraphEvaluationMetricsDisplay
            buildGraphMetrics={buildGraphEvaluationMetrics}
          />
        )}
      {dynamicExecutionMetrics && dynamicExecutionRaceStatistics.length > 0 && (
        <DynamicExecutionMetricsDisplay
          dynamicExecutionMetrics={dynamicExecutionMetrics}
        />
      )}
      {(workers.length > 0 ||
        (workerPoolMetrics?.workerPoolStats?.length ?? 0) > 0) && (
        <WorkerMetricsDisplay
          workerMetrics={workers}
          workerPoolMetrics={workerPoolMetrics}
        />
      )}
    </Space>
  );
};
