import type {
  ActionCacheStatistics,
  Maybe,
  MissDetail,
} from "@/graphql/__generated__/graphql";
import ActionsPieChart, { type ActionsChartItem } from "../ActionsPieChart";
import { chartColor } from "../ActionsPieChart/utils";
import { nullPercent } from "../Utilities/nullPercent";

interface Props {
  acStatistics?: Maybe<ActionCacheStatistics>;
}

const ActionCacheMissMetrics: React.FC<Props> = ({ acStatistics }) => {
  const chartItems: ActionsChartItem[] = [];

  if (acStatistics) {
    acStatistics?.missDetails?.forEach((item: MissDetail, index: number) => {
      const chartItem: ActionsChartItem = {
        key: index,
        count: item.count ?? 0,
        percent: nullPercent(item.count, acStatistics?.misses, 0),
        color: chartColor(index),
        value: item.reason ?? "",
        type: "square",
      };
      chartItems.push(chartItem);
    });
  }

  return <ActionsPieChart items={chartItems} />;
};

export default ActionCacheMissMetrics;
