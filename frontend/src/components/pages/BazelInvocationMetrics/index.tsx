import { Flex, Space } from "antd";
import ActionStatisticsDisplay from "@/components/ActionStatisticsDisplay";
import { ArtifactsMetricsDisplay } from "@/components/ArtifactsMetricsDisplay";
import { GarbageCollectionMetrics } from "@/components/GarbageCollectionMetrics";
import MemoryMetricsDisplay from "@/components/MemoryMetrics";
import { SystemNetworkStatsDisplay } from "@/components/SystemNetworkStatsDisplay";
import { TimingMetricsDisplay } from "@/components/TimingMetricsDisplay";
import { getFragmentData } from "@/graphql/__generated__";
import type { BazelInvocationMetricsFragment } from "@/graphql/__generated__/graphql";
import {
  BAZEL_INVOCATION_METRICS_ACTION_SUMMARY_FRAGMENT,
  BAZEL_INVOCATION_METRICS_ARTIFACT_METRICS_FRAGMENT,
  BAZEL_INVOCATION_METRICS_GARBAGE_METRICS_FRAGMENT,
  BAZEL_INVOCATION_METRICS_MEMORY_METRICS_FRAGMENT,
  BAZEL_INVOCATION_METRICS_SYSTEM_NETWORK_STATS_FRAGMENT,
  BAZEL_INVOCATION_METRICS_TIMING_METRICS_FRAGMENT,
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
    <Space direction="vertical" size="middle">
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
      </Flex>
      {actionSummary && (
        <ActionStatisticsDisplay actionSummary={actionSummary} />
      )}
      {garbageMetrics && (
        <GarbageCollectionMetrics garbageMetrics={garbageMetrics} />
      )}
    </Space>
  );
};
