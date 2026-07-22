import { ApartmentOutlined } from "@ant-design/icons";
import { Flex, Statistic } from "antd";
import type React from "react";
import type { BazelInvocationMetricsBuildGraphMetricsFragment } from "@/graphql/__generated__/graphql";
import { PortalCard } from "../PortalCard";

interface Props {
  buildGraphMetrics: BazelInvocationMetricsBuildGraphMetricsFragment;
  cardStyle?: React.CSSProperties;
}

export const BuildGraphMetricsDisplay: React.FC<Props> = ({
  buildGraphMetrics,
  cardStyle,
}) => {
  return (
    <PortalCard
      type="inner"
      icon={<ApartmentOutlined />}
      style={cardStyle}
      titleBits={["Build Graph Metrics"]}
    >
      <Flex gap="large" wrap={true}>
        <Statistic
          title="Action Lookup Values"
          value={buildGraphMetrics.actionLookupValueCount ?? 0}
        />
        <Statistic
          title="Action Lookup Values Without Aspects"
          value={
            buildGraphMetrics.actionLookupValueCountNotIncludingAspects ?? 0
          }
        />
        <Statistic title="Actions" value={buildGraphMetrics.actionCount ?? 0} />
        <Statistic
          title="Actions Without Aspects"
          value={buildGraphMetrics.actionCountNotIncludingAspects ?? 0}
        />
        <Statistic
          title="Input File Configured Targets"
          value={buildGraphMetrics.inputFileConfiguredTargetCount ?? 0}
        />
        <Statistic
          title="Output File Configured Targets"
          value={buildGraphMetrics.outputFileConfiguredTargetCount ?? 0}
        />
        <Statistic
          title="Other Configured Targets"
          value={buildGraphMetrics.otherConfiguredTargetCount ?? 0}
        />
        <Statistic
          title="Output Artifacts"
          value={buildGraphMetrics.outputArtifactCount ?? 0}
        />
        <Statistic
          title="Post-Invocation Skyframe Nodes"
          value={buildGraphMetrics.postInvocationSkyframeNodeCount ?? 0}
        />
      </Flex>
    </PortalCard>
  );
};
