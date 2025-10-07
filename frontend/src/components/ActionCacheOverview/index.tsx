import type {
  ActionCacheStatistics,
  Maybe,
} from "@/graphql/__generated__/graphql";
import { Space, Statistic } from "antd";
import { readableFileSize } from "@/utils/filesize";
import { readableDurationFromMilliseconds } from "@/utils/time";
import { Bar, BarChart, LabelList, Legend } from "recharts";
import { YAxis } from "recharts";
import { nullPercent } from "../Utilities/nullPercent";

interface Props {
  acStatistics?: Maybe<ActionCacheStatistics>;
}

const MISS_COLOR = "#777777";
const HIT_COLOR = "#49AA19";
const LABBEL_COLOR = "#000000";
const LABEL_THRESHOLD = 0.05;

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
          {(acStatistics?.misses ?? 0) / hitMissTotal >= LABEL_THRESHOLD && (
            <LabelList
              dataKey="miss_label"
              position="center"
              fill={LABBEL_COLOR}
            />
          )}
        </Bar>
        <Bar dataKey="Hit" fill={HIT_COLOR} stackId="a">
          {(acStatistics?.hits ?? 0) / hitMissTotal >= LABEL_THRESHOLD && (
            <LabelList
              dataKey="hit_label"
              position="center"
              fill={LABBEL_COLOR}
            />
          )}
        </Bar>
        {/* YAxis is hidden, but needs to exist to force the correct domain for the chart */}
        <YAxis domain={[0, hitMissTotal]} hide={true} />
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
      <Statistic
        title="Size"
        value={readableFileSize(acStatistics?.sizeInBytes ?? 0)}
      />
      <Statistic
        title="Save Time"
        value={readableDurationFromMilliseconds(
          acStatistics?.saveTimeInMs ?? 0,
          { smallestUnit: "ms" },
        )}
      />
      <Statistic
        title="Load Time"
        value={readableDurationFromMilliseconds(
          acStatistics?.loadTimeInMs ?? 0,
          { smallestUnit: "ms" },
        )}
      />
    </Space>
  );
};

export default ActionCacheOverview;
