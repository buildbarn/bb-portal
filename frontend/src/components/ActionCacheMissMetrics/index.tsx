import type {
  ActionCacheStatistics,
  Maybe,
  MissDetail,
} from "@/graphql/__generated__/graphql";
import SummaryPieChart, { type SummaryChartItem } from "../SummaryPieChart";
import { nullPercent } from "../Utilities/nullPercent";

interface Props {
  acStatistics?: Maybe<ActionCacheStatistics>;
}

const ActionCacheMissMetrics: React.FC<Props> = ({ acStatistics }) => {
  const chartItems: SummaryChartItem[] = [];

  if (acStatistics) {
    acStatistics?.missDetails?.forEach((item: MissDetail, index: number) => {
      const chartItem: SummaryChartItem = {
        key: index,
        count: item.count ?? 0,
        percent: nullPercent(item.count, acStatistics?.misses, 0),
        value: item.reason ?? "",
        type: "square",
      };
      chartItems.push(chartItem);
    });
  }

  return <SummaryPieChart items={chartItems} />;
};

export default ActionCacheMissMetrics;
