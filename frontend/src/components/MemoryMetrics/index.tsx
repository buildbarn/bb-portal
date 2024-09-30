import React, { useCallback, useState } from "react";
import { PieChart, Pie, Cell, Legend } from 'recharts';
import { Table, Row, Col, Statistic, Space } from 'antd';
import type { StatisticProps, TableColumnsType } from "antd/lib";
import CountUp from 'react-countup';
import { MemoryMetrics, GarbageMetrics } from "@/graphql/__generated__/graphql";
import PortalCard from "../PortalCard";
import { PieChartOutlined, HddOutlined } from "@ant-design/icons";
import { renderActiveShape, newColorFind } from "../Utilities/renderShape";

interface GarbageMetricDetailDisplayType {
    key: React.Key;
    name: string;
    value: number;
    color: string;
    //    rate: string;
}

const formatter: StatisticProps['formatter'] = (value) => (
    <CountUp end={value as number} separator="," />
);

const garbage_columns: TableColumnsType<GarbageMetricDetailDisplayType> = [
    {
        title: "Type",
        dataIndex: "name",
    },
    {
        title: "Garbage Collected",
        dataIndex: "value",
        sorter: (a, b) => a.value - b.value,
    },
]

const MemoryMetricsDisplay: React.FC<{ memoryMetrics: MemoryMetrics | undefined; }> = ({ memoryMetrics }) => {

    const garbage_data: GarbageMetricDetailDisplayType[] = [];
    memoryMetrics?.garbageMetrics?.map((item: GarbageMetrics, index) => {
        var gm: GarbageMetricDetailDisplayType = {
            key: index,
            name: item.type ?? "",
            value: item.garbageCollected ?? 0,
            color: newColorFind(index) ?? "#333333"
        }
        garbage_data.push(gm)
    });

    const [activeIndexRunner, setActiveIndexRunner] = useState(0);
    const onRunnerPieEnter = useCallback(
        (_: any, runner_idx: any) => {
            setActiveIndexRunner(runner_idx);
        },
        [setActiveIndexRunner]
    );

    return (
        <Space direction="vertical" size="middle" style={{ display: 'flex' }} >
            <PortalCard icon={<PieChartOutlined />} titleBits={["Memory Metrics"]} >
                <Row>
                    <Space size="large">
                        <Statistic title="Peak Post GC Heap Size" value={memoryMetrics?.peakPostGcHeapSize ?? 0} formatter={formatter} />
                        <Statistic title="Peak Post TC Tenured Space Heap Size" value={memoryMetrics?.peakPostGcTenuredSpaceHeapSize ?? 0} formatter={formatter} />
                        <Statistic title="Used Heap Size Post Build" value={memoryMetrics?.usedHeapSizePostBuild ?? 0} formatter={formatter} />
                    </Space>
                </Row>
                <Row justify="space-around" align="top">
                    <Col span="12">
                        <PortalCard icon={<PieChartOutlined />} titleBits={["Garbage Collection Breakdown"]}>
                            <PieChart width={500} height={500}>

                                <Pie
                                    activeIndex={activeIndexRunner}
                                    activeShape={renderActiveShape}
                                    data={garbage_data}
                                    dataKey="value"
                                    nameKey="name"
                                    cx="50%"
                                    cy="50%"
                                    innerRadius={70}
                                    outerRadius={90}
                                    onMouseEnter={onRunnerPieEnter}>
                                    {
                                        garbage_data.map((entry, runner_index) => (
                                            <Cell key={`cell-${runner_index}`} fill={entry.color} />
                                        ))
                                    }
                                </Pie>
                                <Legend layout="vertical" />
                            </PieChart>
                        </PortalCard>
                    </Col>
                    <Col span="12">
                        <PortalCard icon={<HddOutlined />} titleBits={["Gargage Collection Data"]}>
                            <Table
                                columns={garbage_columns}
                                dataSource={garbage_data}
                                showSorterTooltip={{ target: 'sorter-icon' }}
                            />
                        </PortalCard>
                    </Col>
                </Row >
            </PortalCard>
        </Space>
    )
}

export default MemoryMetricsDisplay;