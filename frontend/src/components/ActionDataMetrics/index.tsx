import React, { useCallback, useState } from "react";
import { PieChart, Pie, Cell } from 'recharts';
import { Table, Row, Col, Space, Statistic } from 'antd';
import { BuildOutlined, PieChartOutlined } from "@ant-design/icons";
import type { StatisticProps, TableColumnsType } from "antd/lib";
import { ActionSummary, ActionData } from "@/graphql/__generated__/graphql";
import PortalCard from "../PortalCard";
import { renderActiveShape, newColorFind } from "../Utilities/renderShape"
import styles from "../../theme/theme.module.css"
import { millisecondsToTime } from "../Utilities/time";
interface ActionDataGraphDisplayType {
    key: React.Key;
    name: string;
    value: number;
    color: string;
}

interface ActionDataColumnType {
    key: React.Key;
    mnemonic: string;
    actionsExecuted: number;
    actionsCreated: number;
    firstStartedMs: number;
    lastEndedMs: number;
    systemTime: number;
    userTime: number;
}

const ad_columns: TableColumnsType<ActionDataColumnType> = [
    {
        title: "Mnemonic",
        dataIndex: "mnemonic"
    },
    {
        title: "Actions Executed",
        dataIndex: "actionsExecuted",
        align: "right",
        render: (_, record) => <span className={styles.numberFormat}>{record.actionsExecuted}</span>,
        sorter: (a, b) => (a.actionsExecuted ?? 0) - (b.actionsExecuted ?? 0),
    },
    {
        title: "Actions Created",
        dataIndex: "actionsCreated",
        align: "right",
        render: (_, record) => <span className={styles.numberFormat}>{record.actionsCreated}</span>,
        sorter: (a, b) => (a.actionsCreated ?? 0) - (b.actionsCreated ?? 0),
    },
    {
        title: "System Time",
        dataIndex: "systemTime",
        align: "right",
        render: (_, record) => <span className={styles.numberFormat}>{millisecondsToTime(record.systemTime)}</span>,
        sorter: (a, b) => (a.systemTime ?? 0) - (b.systemTime ?? 0),
    },
    {
        title: "User Time",
        dataIndex: "userTime",
        align: "right",
        render: (_, record) => <span className={styles.numberFormat}>{millisecondsToTime(record.userTime)}</span>,
        sorter: (a, b) => (a.userTime ?? 0) - (b.userTime ?? 0),
    },
]

const ActionDataMetrics: React.FC<{ acMetrics: ActionSummary | undefined; }> = ({ acMetrics }) => {

    const actions_data: ActionDataColumnType[] = [];
    const actions_graph_data: ActionDataGraphDisplayType[] = [];
    acMetrics?.actionData?.map((ad: ActionData, idx) => {
        actions_data.push({
            key: "action_data_key" + ad.id,
            mnemonic: ad.mnemonic ?? "",
            actionsExecuted: ad.actionsExecuted ?? 0,
            actionsCreated: ad.actionsCreated ?? 0,
            firstStartedMs: ad.firstStartedMs ?? 0,
            lastEndedMs: ad.lastEndedMs ?? 0,
            systemTime: ad.systemTime ?? 0,
            userTime: ad.userTime ?? 0
        })
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
            <PortalCard icon={<PieChart />} type="inner" titleBits={["Actions"]}>
                <Row>
                    <Space size={"large"}>
                        <Statistic title="Actions Executed" value={totalActionsExecuted} />
                        <Statistic title="Actions Created" value={totalActionsCreated} />
                        <Statistic title="Total User Time(ms)" value={totalUserTime} />
                        <Statistic title="Total System Time(ms)" value={totalSystemTime} />
                    </Space>
                </Row>
                <Row justify="space-around" align="top">
                    <Col span="14">
                        <Table
                            columns={ad_columns}
                            dataSource={actions_data}
                            showSorterTooltip={{ target: 'sorter-icon' }}
                        />
                        {/* <PortalCard type="inner" icon={<BuildOutlined />} titleBits={["Actions Data"]}>
                        </PortalCard> */}
                    </Col>
                    <Col span="10">
                        <PortalCard type="inner" icon={<PieChartOutlined />} titleBits={["User Time Breakdown"]} hidden={totalUserTime == 0}>
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