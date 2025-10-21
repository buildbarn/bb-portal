import { Space, Statistic } from "antd";
import type React from "react";
import type { TargetMetrics } from "@/graphql/__generated__/graphql";

interface Props {
  targetMetrics: TargetMetrics | undefined;
  invocationTargetsCount: number | undefined;
  invocationTargetsBuiltSuccessfully: number | undefined;
  invocationTargetsSkipped: number | undefined;
}

export const InvocationTargetsMetrics: React.FC<Props> = ({
  targetMetrics,
  invocationTargetsCount,
  invocationTargetsBuiltSuccessfully,
  invocationTargetsSkipped,
}) => {
  return (
    <Space size="large">
      {invocationTargetsCount !== undefined && (
        <Statistic title="Targets Analyzed" value={invocationTargetsCount} />
      )}
      {invocationTargetsBuiltSuccessfully !== undefined && (
        <Statistic
          title="Targets Built Successfully"
          value={invocationTargetsBuiltSuccessfully}
          valueStyle={{ color: "green" }}
        />
      )}
      {invocationTargetsSkipped !== undefined && (
        <Statistic
          title="Targets Skipped"
          value={invocationTargetsSkipped}
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
