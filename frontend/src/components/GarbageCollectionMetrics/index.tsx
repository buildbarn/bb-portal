import { RestOutlined } from "@ant-design/icons";
import type React from "react";
import { PortalCard } from "@/components/PortalCard";
import type {
  BazelInvocationMetricsGarbageMetricsFragment,
  GarbageMetrics,
} from "@/graphql/__generated__/graphql";
import SummaryPieChart, { type SummaryChartItem } from "../SummaryPieChart";
import { nullPercent } from "../Utilities/nullPercent";

interface Props {
  garbageMetrics: BazelInvocationMetricsGarbageMetricsFragment[];
  cardStyle?: React.CSSProperties;
}

export const GarbageCollectionMetrics: React.FC<Props> = ({
  garbageMetrics,
  cardStyle,
}) => {
  if (garbageMetrics.length === 0) {
    return null;
  }

  const totalGarbageCollected = garbageMetrics.reduce(
    (acc, item) => acc + (item.garbageCollected ?? 0),
    0,
  );

  const chartItems: SummaryChartItem[] = [];
  garbageMetrics.forEach((item: GarbageMetrics, index) => {
    chartItems.push({
      key: index,
      value: item.type ?? "",
      percent: nullPercent(item.garbageCollected, totalGarbageCollected, 0),
      count: item.garbageCollected ?? 0,
    });
  });

  return (
    <PortalCard
      type="inner"
      icon={<RestOutlined />}
      style={cardStyle}
      titleBits={["Garbage Collection Breakdown"]}
    >
      <SummaryPieChart items={chartItems} />
    </PortalCard>
  );
};
