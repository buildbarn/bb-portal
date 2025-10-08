import React, { useCallback, useState } from "react";
import { PieChart, Pie, Cell, Legend, BarChart, Bar, LabelList } from 'recharts';
import { Table, Row, Col, Statistic, Tooltip, Space } from 'antd';
import type { StatisticProps, TableColumnsType } from "antd/lib";
import { ActionCacheStatistics, ActionSummary, MissDetail } from "@/graphql/__generated__/graphql";
import PortalCard from "../PortalCard";
import { PieChartOutlined, DashboardOutlined, HddOutlined } from "@ant-design/icons";
import { renderActiveShape } from "../Utilities/renderShape"
import { nullPercent } from "../Utilities/nullPercent";
import "./index.module.css"
import MissDetailTag, { MissDetailEnum } from "./ActionCacheMissTag";
import styles from "../../theme/theme.module.css"
interface MissDetailDisplayDataType {
    key: React.Key;
    name: string;
    value: number;
    rate: string;
}

var ac_colors =
    [
        "grey",     //unknown
        "blue",     //different action key
        "pink",     //different deps
        "purple",   //different env
        "cyan",     //diff files
        "orange",   //corrupted cache entry
        "red",      //not cached
        "green"]    //unconditional execution

const ac_columns: TableColumnsType<MissDetailDisplayDataType> = [
    {
        title: "Miss Reason",
        dataIndex: "name",
        render: (x) => <MissDetailTag key="status" status={x as MissDetailEnum} />,
        // render: (idx, x) => <span className={"ac-miss-detail-" + x.name}>{x.name}</span>,
    },
    {
        title: "Count",
        dataIndex: "value",
        align: "right",
        defaultSortOrder: "descend",
        render: (_, record) => <span className={styles.numberFormat}>{record.value}</span>,
        sorter: (a, b) => a.value - b.value,
    },
    {
        title: "Rate (%)",
        dataIndex: "rate",
        align: "right",
        render: (_, record) => <span className={styles.numberFormat}>{record.rate}</span>,
        sorter: (a, b) => parseFloat(a.rate) - parseFloat(b.rate),
    }
]


const AcMetrics: React.FC<{ acMetrics: ActionSummary | undefined; }> = ({ acMetrics }) => {

    const acMetricsData: ActionCacheStatistics | undefined = acMetrics?.actionCacheStatistics ?? undefined

    var hitMissTotal: number = (acMetricsData?.misses ?? 0) + (acMetricsData?.hits ?? 0);

    const hits_data = [
        {
            key: "hitMissBarChart",
            Hit: acMetricsData?.hits,
            Miss: acMetricsData?.misses,
            hit_label: nullPercent(acMetricsData?.hits, hitMissTotal, 0),
            miss_label: nullPercent(acMetricsData?.misses, hitMissTotal, 0)
        },
    ]

    const ac_data: MissDetailDisplayDataType[] = [];
    var missTotal: number = acMetricsData?.misses ?? 0;

    acMetricsData?.missDetails?.map((item: MissDetail, index) => {
        var acd: MissDetailDisplayDataType = {
            key: index,
            name: item.reason ?? "",
            value: item.count ?? 0,
            rate: nullPercent(item.count, missTotal),

        }
        ac_data.push(acd)
    });


    const [activeIndexRunner, setActiveIndexRunner] = useState(0);
    const onRunnerPieEnter = useCallback(
        (_: any, runner_idx: any) => {
            setActiveIndexRunner(runner_idx);
        },
        [setActiveIndexRunner]
    );
    const acTitle: React.ReactNode[] = [<span key="label">Action Cache Statistics</span>];


    return (
        <Space direction="vertical" size="middle" style={{ display: 'flex' }} >
            <PortalCard type="inner" icon={<DashboardOutlined />} titleBits={acTitle} >
                <Row>
                    <Space size="large">
                        <BarChart width={170} height={150} data={hits_data} margin={{ top: 0, left: 10, bottom: 10, right: 10 }}>
                            <Bar dataKey="Miss" fill={"#8884d8"} stackId="a">
                                <LabelList dataKey="miss_label" position="center" />
                            </Bar>
                            <Bar dataKey="Hit" fill={"#82ca9d"} stackId="a">
                                <LabelList dataKey="hit_label" position="center" />
                            </Bar>
                            <Tooltip />
                            <Legend />
                        </BarChart>
                        <Statistic title="Hits" value={acMetricsData?.hits ?? 0} valueStyle={{ color: "#82ca9d" }} />
                        <Statistic title="Misses" value={acMetricsData?.misses ?? 0} valueStyle={{ color: "#8884d8" }} />
                        <Statistic title="Size (bytes)" value={acMetricsData?.sizeInBytes ?? 0} />
                        <Statistic title="Save Time(ms)" value={acMetricsData?.saveTimeInMs ?? 0} />
                        <Statistic title="Load Time(ms)" value={acMetricsData?.loadTimeInMs ?? 0} />

                    </Space>
                </Row>
                <Row justify="space-around" align="top" >
                    <Col span="12">
                        <PortalCard type="inner" icon={<PieChartOutlined />} titleBits={["Miss Detail Breakdown"]}>
                            <PieChart width={600} height={500}>

                                <Pie
                                    activeIndex={activeIndexRunner}
                                    activeShape={renderActiveShape}
                                    data={ac_data}
                                    dataKey="value"
                                    nameKey="name"
                                    cx="50%"
                                    cy="50%"
                                    innerRadius={70}
                                    outerRadius={90}
                                    onMouseEnter={onRunnerPieEnter}>
                                    {
                                        ac_data.map((entry, runner_index) => (
                                            <Cell key={`cell-${runner_index}`} fill={ac_colors[runner_index]} />
                                        ))
                                    }
                                </Pie>
                                <Legend layout="vertical" />
                            </PieChart>
                        </PortalCard>

                    </Col>
                    <Col span="12">
                        <PortalCard type="inner" icon={<HddOutlined />} titleBits={["Miss Detail Data"]}>
                            <Table
                                columns={ac_columns}
                                dataSource={ac_data}
                                showSorterTooltip={{ target: 'sorter-icon' }}
                                pagination={false}
                                rowClassName={(record, _) => ("ac-miss-detail-" + record.name)}
                            />
                        </PortalCard>
                    </Col>
                </Row >
            </PortalCard>
        </Space>
    )
}

export default AcMetrics;
