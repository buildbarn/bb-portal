import type {
  ActionCacheStatistics,
  Maybe,
} from "@/graphql/__generated__/graphql";
import { Space, Statistic, Tooltip } from "antd";
import { Bar, BarChart, LabelList, Legend } from "recharts";
import { nullPercent } from "../Utilities/nullPercent";

interface Props {
  acStatistics?: Maybe<ActionCacheStatistics>;
}

const HIT_COLOR = "#82ca9d";
const MISS_COLOR = "#8884d8";

const ActionCacheOverview: React.FC<Props> = ({ acStatistics }) => {
  const hitMissTotal = (acStatistics?.hits ?? 0) + (acStatistics?.misses ?? 0);

  const hitsData = [
    {
      key: "hitMissBarChart",
      Hit: acStatistics?.hits ?? 0,
      Miss: acStatistics?.misses ?? 0,
      hit_label: nullPercent(acStatistics?.hits, hitMissTotal, 0),
      miss_label: nullPercent(acStatistics?.misses, hitMissTotal, 0),
    },
  ];

  return (
    <Space size="middle" style={{ width: 600, height: 300 }}>
      <BarChart
        width={150}
        height={300}
        data={hitsData}
        margin={{ top: 0, left: 10, bottom: 10, right: 10 }}
      >
        <Bar dataKey="Miss" fill={MISS_COLOR} stackId="a">
          <LabelList dataKey="miss_label" position="center" />
        </Bar>
        <Bar dataKey="Hit" fill={HIT_COLOR} stackId="a">
          <LabelList dataKey="hit_label" position="center" />
        </Bar>
        <Tooltip />
        <Legend />
      </BarChart>
      <Statistic
        title="Hits"
        value={acStatistics?.hits ?? 0}
        valueStyle={{ color: HIT_COLOR }}
      />
      <Statistic
        title="Misses"
        value={acStatistics?.misses ?? 0}
        valueStyle={{ color: MISS_COLOR }}
      />
      <Statistic title="Size (bytes)" value={acStatistics?.sizeInBytes ?? 0} />
      <Statistic
        title="Save Time(ms)"
        value={acStatistics?.saveTimeInMs ?? 0}
      />
      <Statistic
        title="Load Time(ms)"
        value={acStatistics?.loadTimeInMs ?? 0}
      />
    </Space>
  );
};

export default ActionCacheOverview;
