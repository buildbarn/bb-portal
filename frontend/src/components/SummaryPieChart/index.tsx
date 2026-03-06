import { Space, theme } from "antd";
import { Link, LinkOptions } from '@tanstack/react-router';
import { useState } from "react";
import { Legend, LegendPayload, Pie, PieChart, PieSectorDataItem, PieSectorShapeProps, Sector } from "recharts";
import { themeColor } from "./utils";
import React from "react";

// The `Legend` component uses `value`, `fill`,
// and `type` by default to render the data.
export interface SummaryChartItem {
  key: React.Key;
  value: string;
  percent: string;
  fill?: string;
  count: number;
  link?: LinkOptions;
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
  const [portalNode, setPortalNode] = useState<HTMLDivElement | null>(null);

  const coloredItems = [...items]
    .sort((a, b) => b.count - a.count)
    .map((item, index) => ({
      ...item,
      fill: item.fill ?? themeColor(token, index)
    }));

  const [hoverIndex, setHoverIndex] = React.useState<number | undefined>(undefined);
  const handleMouseEnter = (_payload: LegendPayload | PieSectorDataItem, index: number) => {
    setHoverIndex(index);
  };
  const handleMouseLeave = () => {
    setHoverIndex(undefined);
  };

  const renderLegendText = (_value: any, entry: LegendPayload, index: number) => {
    let item: SummaryChartItem | undefined
    if (entry.payload) {
      item = entry.payload as SummaryChartItem
    }

    let legendText = <>
      <b>{item?.count ?? 0}</b>{" "}
      {item?.link ? <Link {...item.link}>{item.value}</Link> : item?.value} (
      {item?.percent})
    </>

    if (index === hoverIndex) {
      legendText = <u>{legendText}</u>
    }

    return (
      <span style={{ wordBreak: "break-word" }}>
        {legendText}
      </span>
    );
  };

  const renderPieShape = ({
    cx,
    cy,
    innerRadius,
    outerRadius,
    startAngle,
    endAngle,
    fill,
  }: PieSectorShapeProps, index: number) => {
    return (
      <g>
        <Sector
          cx={cx}
          cy={cy}
          innerRadius={innerRadius}
          outerRadius={outerRadius}
          startAngle={startAngle}
          endAngle={endAngle}
          fill={fill}
          stroke="white"
        />
        {index == hoverIndex &&
          <Sector
            cx={cx}
            cy={cy}
            startAngle={startAngle}
            endAngle={endAngle}
            innerRadius={outerRadius && outerRadius + 6}
            outerRadius={outerRadius && outerRadius + 10}
            fill={fill}
            stroke="white"
          />
        }
      </g>
    );
  };

  return <Space direction="vertical" style={{ width: chartWidth }}>
    <PieChart height={OUTER_RADIUS * 3} width={chartWidth}>
      <Pie
        dataKey="count"
        data={coloredItems}
        innerRadius={INNER_RADIUS}
        outerRadius={OUTER_RADIUS}
        shape={renderPieShape}
        onMouseEnter={handleMouseEnter}
        onMouseLeave={handleMouseLeave}
      />
      <Legend
        layout="vertical"
        wrapperStyle={{
          position: 'static',
          columns: 2,
        }}
        portal={portalNode}
        formatter={renderLegendText}
        onMouseEnter={handleMouseEnter}
        onMouseLeave={handleMouseLeave}
      />
    </PieChart>
    <div style={{ breakInside: "avoid-column" }}>

      {/* Use a React Portal to move the legend to outside of the PieChart. This makes layout easier */}
      <div ref={setPortalNode} />
    </div>
  </Space>
};

export default SummaryPieChart;
