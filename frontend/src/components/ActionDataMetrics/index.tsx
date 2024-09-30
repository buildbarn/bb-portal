import React, { useCallback, useState } from "react";
import { PieChart, Pie, Cell } from 'recharts';
import { Table, Row, Col, Space, Statistic } from 'antd';
import { BuildOutlined, PieChartOutlined } from "@ant-design/icons";
import type { StatisticProps, TableColumnsType } from "antd/lib";
import CountUp from 'react-countup';
import { ActionSummary, ActionData } from "@/graphql/__generated__/graphql";
import PortalCard from "../PortalCard";
import { renderActiveShape, newColorFind } from "../Utilities/renderShape"
interface ActionDataGraphDisplayType {
    key: React.Key;
    name: string;
    value: number;
    color: string;
}

const formatter: StatisticProps['formatter'] = (value) => (
    <CountUp end={value as number} separator="," />
);

const ad_columns: TableColumnsType<ActionData> = [
    {
        title: "Mnemonic",
        dataIndex: "mnemonic"
    },
    {
        title: "Actions Executed",
        dataIndex: "actionsExecuted",
        sorter: (a, b) => (a.actionsExecuted ?? 0) - (b.actionsExecuted ?? 0),
    },
    {
        title: "Actions Created",
        dataIndex: "actionsCreated",
        sorter: (a, b) => (a.actionsCreated ?? 0) - (b.actionsCreated ?? 0),
    },
    {
        title: "First Started(ms)",
        dataIndex: "firstStartedMs",
        sorter: (a, b) => (a.firstStartedMs ?? 0) - (b.firstStartedMs ?? 0),
    },
    {
        title: "Last Ended(ms)",
        dataIndex: "lastEndedMs",
        sorter: (a, b) => (a.lastEndedMs ?? 0) - (b.lastEndedMs ?? 0),
    },
    {
        title: "System Time(ms)",
        dataIndex: "systemTime",
        sorter: (a, b) => (a.systemTime ?? 0) - (b.systemTime ?? 0),
    },
    {
        title: "User Time(ms)",
        dataIndex: "userTime",
        sorter: (a, b) => (a.userTime ?? 0) - (b.userTime ?? 0),
    },
]

const ActionDataMetrics: React.FC<{ acMetrics: ActionSummary | undefined; }> = ({ acMetrics }) => {

    const actions_data: ActionData[] = [];
    const actions_graph_data: ActionDataGraphDisplayType[] = [];
    acMetrics?.actionData?.map((ad: ActionData, idx) => {
        actions_data.push(ad)
        var agd: ActionDataGraphDisplayType = {
            key: "actiondatagraphdisplaytype-" + String(idx),
            name: ad.mnemonic ?? "",
            value: ad.userTime ?? 0,
            color: newColorFind(idx) ?? "#333333"
        }
        actions_graph_data.push(agd)
    });

    const [activeIndexRunner, setActiveIndexRunner] = useState(0);
    const onRunnerPieEnter = useCallback(
        (_: any, runner_idx: any) => {
            setActiveIndexRunner(runner_idx);
        },
        [setActiveIndexRunner]
    );
    var totalUserTime = actions_data.reduce((accumulator, item) => accumulator + (item.userTime ?? 0), 0);
    var totalSystemTime = actions_data.reduce((accumulator, item) => accumulator + (item.systemTime ?? 0), 0);
    var totalActionsExecuted = actions_data.reduce((accumulator, item) => accumulator + (item.actionsExecuted ?? 0), 0);
    var totalActionsCreated = actions_data.reduce((accumulator, item) => accumulator + (item.actionsCreated ?? 0), 0);

    return (
        <Space direction="vertical" size="middle" style={{ display: 'flex' }} >
            <PortalCard icon={<PieChart />} titleBits={["Actions"]}>
                <Row>
                    <Space size={"large"}>
                        <Statistic title="Actions Executed" value={totalActionsExecuted} formatter={formatter} />
                        <Statistic title="Actions Created" value={totalActionsCreated} formatter={formatter} />
                        <Statistic title="Total User Time(ms)" value={totalUserTime} formatter={formatter} />
                        <Statistic title="Total System Time(ms)" value={totalSystemTime} formatter={formatter} />
                    </Space>
                </Row>
                <Row justify="space-around" align="top">
                    <Col span="14">
                        <PortalCard icon={<BuildOutlined />} titleBits={["Actions Data"]}>
                            <Table
                                columns={ad_columns}
                                dataSource={actions_data}
                                showSorterTooltip={{ target: 'sorter-icon' }}
                            />
                        </PortalCard>
                    </Col>
                    <Col span="10">
                        <PortalCard icon={<PieChartOutlined />} titleBits={["User Time(ms)"]}>
                            <PieChart width={600} height={556}>
                                <Pie
                                    activeIndex={activeIndexRunner}
                                    activeShape={renderActiveShape}
                                    data={actions_graph_data}
                                    dataKey="value"
                                    cx="50%"
                                    cy="50%"
                                    innerRadius={50}
                                    outerRadius={90}
                                    onMouseEnter={onRunnerPieEnter}>
                                    {
                                        actions_graph_data.map((entry, actions_index) => (
                                            <Cell key={`cell-${actions_index}`} fill={entry.color} />
                                        ))
                                    }
                                </Pie>
                            </PieChart>
                        </PortalCard>
                    </Col>
                </Row>
            </PortalCard>
        </Space>
    )
}

export default ActionDataMetrics;