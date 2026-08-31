import { Space, Statistic, Tooltip } from "antd";
import type React from "react";
import type { ActionTimingMetrics } from "@/graphql/__generated__/graphql";
import { readableDurationFromMilliseconds } from "@/utils/time";

interface Props {
  actionTimingMetrics: Pick<
    ActionTimingMetrics,
    | "totalExpectedTimeInMs"
    | "timeSavedByCacheHitsInMs"
    | "totalActions"
    | "timedActions"
    | "cacheHitActions"
    | "timedCacheHitActions"
  >;
  executionPhaseTimeInMs: number | null | undefined;
}

const duration = (milliseconds: number | null | undefined) =>
  readableDurationFromMilliseconds(milliseconds, { smallestUnit: "ms" });

export const ActionsMetrics: React.FC<Props> = ({
  actionTimingMetrics,
  executionPhaseTimeInMs,
}) => (
  <Space size="large">
    <Statistic
      title={
        <Tooltip
          title={`Sum of the reported durations for ${actionTimingMetrics.timedActions} of ${actionTimingMetrics.totalActions} actions. This is the time they would take if run serially.`}
        >
          <span>Total Expected Time</span>
        </Tooltip>
      }
      value={duration(actionTimingMetrics.totalExpectedTimeInMs)}
    />
    <Statistic
      title={
        <Tooltip title="Bazel's execution-phase wall time, including actions that ran in parallel.">
          <span>Real Execution Time</span>
        </Tooltip>
      }
      value={duration(executionPhaseTimeInMs)}
    />
    <Statistic
      title={
        <Tooltip
          title={`Sum of the reported durations for ${actionTimingMetrics.timedCacheHitActions} of ${actionTimingMetrics.cacheHitActions} cache-hit actions.`}
        >
          <span>Time Saved by Cache Hits</span>
        </Tooltip>
      }
      value={duration(actionTimingMetrics.timeSavedByCacheHitsInMs)}
      valueStyle={{ color: "green" }}
    />
  </Space>
);
