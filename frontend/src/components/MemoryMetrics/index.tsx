import React from "react";
import { Table, Row, Col, Statistic, Space } from "antd";
import type { TableColumnsType } from "antd/lib";
import { MemoryMetrics, GarbageMetrics } from "@/graphql/__generated__/graphql";
import PortalCard from "../PortalCard";
import { PieChartOutlined, HddOutlined } from "@ant-design/icons";
import styles from "../../theme/theme.module.css";
import { readableFileSize } from "@/utils/filesize";
import SummaryPieChart, { SummaryChartItem } from "../SummaryPieChart";
import { nullPercent } from "../Utilities/nullPercent";

interface GarbageMetricDetailDisplayType {
  key: React.Key;
  name: string;
  value: number;
}

const garbage_columns: TableColumnsType<GarbageMetricDetailDisplayType> = [
  {
    title: "Type",
    dataIndex: "name",
  },
  {
    title: "Garbage Collected",
    dataIndex: "value",
    sorter: (a, b) => a.value - b.value,
    align: "right",
    render: (_, record) => (
      <span className={styles.numberFormat}>{record.value}</span>
    ),
  },
];

const MemoryMetricsDisplay: React.FC<{
  memoryMetrics: MemoryMetrics | undefined;
}> = ({ memoryMetrics }) => {
  const garbage_data: GarbageMetricDetailDisplayType[] = [];
  memoryMetrics?.garbageMetrics?.map((item: GarbageMetrics, index) => {
    var gm: GarbageMetricDetailDisplayType = {
      key: index,
      name: item.type ?? "",
      value: item.garbageCollected ?? 0,
    };
    garbage_data.push(gm);
  });

  const chartItems: SummaryChartItem[] = [];
  const totalGarbageCollected = memoryMetrics?.garbageMetrics?.reduce(
    (acc, item) => acc + (item.garbageCollected ?? 0),
    0,
  );

  memoryMetrics?.garbageMetrics?.forEach((item: GarbageMetrics, index) => {
    const chartItem: SummaryChartItem = {
      key: index,
      value: item.type ?? "",
      percent: nullPercent(item.garbageCollected, totalGarbageCollected, 0),
      count: item.garbageCollected ?? 0,
    };
    chartItems.push(chartItem);
  });

  return (
    <Space direction="vertical" size="middle" style={{ display: "flex" }}>
      <PortalCard
        type="inner"
        icon={<PieChartOutlined />}
        titleBits={["Memory Metrics"]}
      >
        <Row>
          <Space size="large">
            <Statistic
              title="Peak Post GC Heap Size"
              value={readableFileSize(memoryMetrics?.peakPostGcHeapSize ?? 0)}
            />
            <Statistic
              title="Peak Post TC Tenured Space Heap Size"
              value={readableFileSize(memoryMetrics?.peakPostGcTenuredSpaceHeapSize ?? 0)}
            />
            <Statistic
              title="Used Heap Size Post Build"
              value={readableFileSize(memoryMetrics?.usedHeapSizePostBuild ?? 0)}
            />
          </Space>
        </Row>
        <Row justify="space-around" align="top">
          <Col span="12">
            <PortalCard
              type="inner"
              icon={<PieChartOutlined />}
              titleBits={["Garbage Collection Breakdown"]}
            >
              <SummaryPieChart items={chartItems} />
            </PortalCard>
          </Col>
          <Col span="12">
            <PortalCard
              type="inner"
              icon={<HddOutlined />}
              titleBits={["Gargage Collection Data"]}
            >
              <Table
                columns={garbage_columns}
                dataSource={garbage_data}
                showSorterTooltip={{ target: "sorter-icon" }}
              />
            </PortalCard>
          </Col>
        </Row>
      </PortalCard>
    </Space>
  );
};

export default MemoryMetricsDisplay;
