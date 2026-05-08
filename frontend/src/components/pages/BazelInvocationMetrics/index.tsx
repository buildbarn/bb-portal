import { Space } from "antd";
import ActionStatisticsDisplay from "@/components/ActionStatisticsDisplay";
import { ArtifactsMetricsDisplay } from "@/components/ArtifactsMetricsDisplay";
import MemoryMetricsDisplay from "@/components/MemoryMetrics";
import { SystemNetworkStatsDisplay } from "@/components/SystemNetworkStatsDisplay";
import { TimingMetricsDisplay } from "@/components/TimingMetricsDisplay";
import { getFragmentData } from "@/graphql/__generated__";
import type { BazelInvocationMetricsFragment } from "@/graphql/__generated__/graphql";
import {
  BAZEL_INVOCATION_METRICS_ACTION_SUMMARY_FRAGMENT,
  BAZEL_INVOCATION_METRICS_ARTIFACT_METRICS_FRAGMENT,
  BAZEL_INVOCATION_METRICS_MEMORY_METRICS_FRAGMENT,
  BAZEL_INVOCATION_METRICS_SYSTEM_NETWORK_STATS_FRAGMENT,
  BAZEL_INVOCATION_METRICS_TIMING_METRICS_FRAGMENT,
} from "@/routes/bazel-invocations.$invocationID/metrics";

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
  const timingMetrics = getFragmentData(
    BAZEL_INVOCATION_METRICS_TIMING_METRICS_FRAGMENT,
    metrics.timingMetrics,
  );
  const systemNetworkStats = getFragmentData(
    BAZEL_INVOCATION_METRICS_SYSTEM_NETWORK_STATS_FRAGMENT,
    metrics.networkMetrics?.systemNetworkStats,
  );

  return (
    <Space direction="vertical" size="middle">
      {actionSummary && (
        <ActionStatisticsDisplay actionSummary={actionSummary} />
      )}
      {artifactMetrics && (
        <ArtifactsMetricsDisplay artifactMetrics={artifactMetrics} />
      )}
      {memoryMetrics && <MemoryMetricsDisplay memoryMetrics={memoryMetrics} />}
      {timingMetrics && <TimingMetricsDisplay timingMetrics={timingMetrics} />}
      {systemNetworkStats && (
        <SystemNetworkStatsDisplay systemNetworkStats={systemNetworkStats} />
      )}
    </Space>
  );
};
