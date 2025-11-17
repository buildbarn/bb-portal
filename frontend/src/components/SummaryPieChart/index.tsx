import { Flex, theme } from "antd";
import Link from "next/link";
import { useCallback, useState } from "react";
import { Cell, Legend, type LegendType, Pie, PieChart } from "recharts";
import type { Payload } from "recharts/types/component/DefaultLegendContent";
import { renderActiveShapeCompact } from "../Utilities/renderShape";
import styles from "./index.module.css";
import { themeColor } from "./utils";

// The `Legend` component uses `value`, `color`,
// and `type` by default to render the data.
export interface SummaryChartItem {
  key: React.Key;
  value: string;
  percent: string;
  color?: string;
  count: number;
  type?: LegendType;
  href?: string;
}

const { useToken } = theme;

interface Props {
  items: SummaryChartItem[];
  chartWidth?: number;
}

const INNER_RADIUS = 30;
const OUTER_RADIUS = 50;

const SummaryPieChart: React.FC<Props> = ({
  items,
  chartWidth = 600,
}: Props) => {
  const { token } = useToken();

  const renderLegendText = (value: string) => {
    const item = coloredItems.find((i) => i.value === value);
    if (value === coloredItems[activeIndexRunner].value) {
      return (
        <span>
          <u>
            <b>{item?.count ?? 0}</b>{" "}
            <span style={{ color: token.colorText }}>
              {item?.href ? <Link href={item.href}>{value}</Link> : value} (
              {item?.percent})
            </span>
          </u>
        </span>
      );
    }
    return (
      <span>
        <b>{item?.count ?? 0}</b>{" "}
        <span style={{ color: token.colorText }}>
          {item?.href ? <Link href={item.href}>{value}</Link> : value} (
          {item?.percent})
        </span>
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

  const coloredItems = items.map((item, index) => ({
    ...item,
    color: item?.color || themeColor(token, index),
  }));

  return (
    <Flex vertical gap="middle" style={{ width: chartWidth }}>
      <PieChart height={OUTER_RADIUS * 3} width={chartWidth}>
        <Pie
          activeIndex={activeIndexRunner}
          activeShape={renderActiveShapeCompact}
          dataKey="count"
          data={coloredItems}
          innerRadius={INNER_RADIUS}
          outerRadius={OUTER_RADIUS}
          onMouseEnter={onRunnerPieEnter}
        >
          {coloredItems.map((value: SummaryChartItem) => {
            return <Cell key={value.key} fill={value.color} />;
          })}
        </Pie>
      </PieChart>
      <div className={styles.summaryPieChartsWrapper}>
        <Legend
          payload={coloredItems}
          chartWidth={chartWidth}
          layout="vertical"
          align="center"
          wrapperStyle={{
            position: "static",
            columns: 2,
          }}
          formatter={renderLegendText}
          onMouseEnter={onRunnerPieEnter}
        />
      </div>
    </Flex>
  );
};

export default SummaryPieChart;
