import { RadiusUprightOutlined } from "@ant-design/icons";
import { Space } from "antd";
import type React from "react";
import type { BazelInvocationMetricsArtifactMetricsFragment } from "@/graphql/__generated__/graphql";
import { readableFileSize } from "@/utils/filesize";
import { MultiStatistic } from "../MultiStatistic";
import PortalCard from "../PortalCard";

interface Props {
  artifactMetrics: BazelInvocationMetricsArtifactMetricsFragment;
  cardStyle?: React.CSSProperties;
}

export const ArtifactsMetricsDisplay: React.FC<Props> = ({
  artifactMetrics,
  cardStyle,
}) => {
  return (
    <PortalCard
      type="inner"
      icon={<RadiusUprightOutlined />}
      style={cardStyle}
      titleBits={["Artifacts Metrics"]}
    >
      <Space size="large">
        <MultiStatistic
          title="Source Artifacts Read"
          values={[
            {
              key: "count",
              value: artifactMetrics.sourceArtifactsReadCount ?? 0,
            },
            {
              key: "size",
              value: readableFileSize(
                artifactMetrics.sourceArtifactsReadSizeInBytes ?? 0,
              ),
            },
          ]}
        />
        <MultiStatistic
          title="Output Artifacts From Action Cache"
          values={[
            {
              key: "count",
              value: artifactMetrics.outputArtifactsFromActionCacheCount ?? 0,
            },
            {
              key: "size",
              value: readableFileSize(
                artifactMetrics.outputArtifactsFromActionCacheSizeInBytes ?? 0,
              ),
            },
          ]}
        />
        <MultiStatistic
          title="Output Artifacts Seen"
          values={[
            {
              key: "count",
              value: artifactMetrics.outputArtifactsSeenCount ?? 0,
            },
            {
              key: "size",
              value: readableFileSize(
                artifactMetrics.outputArtifactsSeenSizeInBytes ?? 0,
              ),
            },
          ]}
        />
        <MultiStatistic
          title="Top Level Artifacts"
          values={[
            {
              key: "count",
              value: artifactMetrics.topLevelArtifactsCount ?? 0,
            },
            {
              key: "size",
              value: readableFileSize(
                artifactMetrics.topLevelArtifactsSizeInBytes ?? 0,
              ),
            },
          ]}
        />
      </Space>
    </PortalCard>
  );
};
