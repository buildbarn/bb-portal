import { PieChartOutlined } from "@ant-design/icons";
import { Space, Statistic } from "antd";
import type React from "react";
import type { BazelInvocationMetricsMemoryMetricsFragment } from "@/graphql/__generated__/graphql";
import { readableFileSize } from "@/utils/filesize";
import PortalCard from "../PortalCard";

interface Props {
  memoryMetrics: BazelInvocationMetricsMemoryMetricsFragment;
  cardStyle?: React.CSSProperties;
}
const MemoryMetricsDisplay: React.FC<Props> = ({
  memoryMetrics,
  cardStyle,
}) => {
  return (
    <PortalCard
      type="inner"
      icon={<PieChartOutlined />}
      style={cardStyle}
      titleBits={["Memory Metrics"]}
    >
      <Space size="large">
        <Statistic
          title="Peak Post GC Heap Size"
          value={readableFileSize(memoryMetrics.peakPostGcHeapSize ?? 0)}
        />
        <Statistic
          title="Peak Post TC Tenured Space Heap Size"
          value={readableFileSize(
            memoryMetrics.peakPostGcTenuredSpaceHeapSize ?? 0,
          )}
        />
        <Statistic
          title="Used Heap Size Post Build"
          value={readableFileSize(memoryMetrics.usedHeapSizePostBuild ?? 0)}
        />
      </Space>
    </PortalCard>
  );
};

export default MemoryMetricsDisplay;
