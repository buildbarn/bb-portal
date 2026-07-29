import { Space, Statistic } from "antd";
import type React from "react";
import type {
  BazelInvocationTargetCountsFragment,
  BazelInvocationTargetMetricsFragment,
} from "@/graphql/__generated__/graphql";

interface Props {
  targetMetrics: BazelInvocationTargetMetricsFragment | null | undefined;
  targetCounts: BazelInvocationTargetCountsFragment;
}

export const InvocationTargetsMetrics: React.FC<Props> = ({
  targetMetrics,
  targetCounts,
}) => {
  return (
    <Space size="large">
      {targetCounts.numTotal.totalCount !== undefined && (
        <Statistic
          title="Targets Analyzed"
          value={targetCounts.numTotal.totalCount}
        />
      )}
      {targetCounts.numSuccessful.totalCount !== undefined && (
        <Statistic
          title="Targets Built Successfully"
          value={targetCounts.numSuccessful.totalCount}
          valueStyle={{ color: "green" }}
        />
      )}
      {targetCounts.numSkipped.totalCount !== undefined && (
        <Statistic
          title="Targets Skipped"
          value={targetCounts.numSkipped.totalCount}
          valueStyle={{ color: "purple" }}
        />
      )}
      {targetMetrics?.targetsConfigured !== undefined &&
        targetMetrics?.targetsConfigured !== null && (
          <Statistic
            title="Targets Configured"
            value={targetMetrics.targetsConfigured}
          />
        )}
      {targetMetrics?.targetsConfiguredNotIncludingAspects !== undefined &&
        targetMetrics?.targetsConfiguredNotIncludingAspects !== null && (
          <Statistic
            title="Targets Configured Not Including Aspects"
            value={targetMetrics.targetsConfiguredNotIncludingAspects}
          />
        )}
    </Space>
  );
};
