import { Flex } from "antd";
import { useCallback, useState } from "react";
import { Cell, Legend, type LegendType, Pie, PieChart } from "recharts";
import type { Payload } from "recharts/types/component/DefaultLegendContent";
import { renderActiveShapeCompact } from "../Utilities/renderShape";

// The `Legend` component uses `value`, `color`,
// and `type` by default to render the data.
export interface ActionsChartItem {
  key: React.Key;
  value: string;
  percent: string;
  color: string;
  count: number;
  type?: LegendType;
}

interface Props {
  items: ActionsChartItem[];
}

const INNER_RADIUS = 30;
const OUTER_RADIUS = 50;
const CHART_WIDTH = 600;

const ActionsPieChart: React.FC<Props> = ({ items }: Props) => {
  const renderLegendText = (value: string) => {
    const item = items.find((i) => i.value === value);
    if (value === items[activeIndexRunner].value) {
      return (
        <span>
          <u>
            <b>{item?.count ?? 0}</b> {value} ({item?.percent})
          </u>
        </span>
      );
    }
    return (
      <span>
        <b>{item?.count ?? 0}</b> {value} ({item?.percent})
      </span>
    );
  };

  const [activeIndexRunner, setActiveIndexRunner] = useState<number>(0);
  const onRunnerPieEnter = useCallback((_: Payload, index: number) => {
    setActiveIndexRunner(index);
  }, []);

  // Items are sorted to display the highest count first in the legend
  items.sort((a, b) => {
    return b.count - a.count;
  });
  return (
    <Flex vertical gap="middle">
      <PieChart height={OUTER_RADIUS * 3} width={CHART_WIDTH}>
        <Pie
          activeIndex={activeIndexRunner}
          activeShape={renderActiveShapeCompact}
          dataKey="count"
          data={items}
          innerRadius={INNER_RADIUS}
          outerRadius={OUTER_RADIUS}
          onMouseEnter={onRunnerPieEnter}
        >
          {items.map((value: ActionsChartItem) => {
            return <Cell key={value.key} fill={value.color} />;
          })}
        </Pie>
      </PieChart>
      <Legend
        payload={items}
        chartWidth={CHART_WIDTH}
        layout="vertical"
        align="center"
        wrapperStyle={{
          position: "static",
          columns: 2,
        }}
        formatter={renderLegendText}
        onMouseEnter={onRunnerPieEnter}
      />
    </Flex>
  );
};

export default ActionsPieChart;
