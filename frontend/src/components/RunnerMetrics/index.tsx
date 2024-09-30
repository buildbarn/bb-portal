import { RunnerCount } from "@/graphql/__generated__/graphql";
import React, { useCallback, useState } from "react";
import { PieChart, Pie, Cell, Legend } from 'recharts';
import { Table, Row, Col, Typography } from 'antd';
import type { TableColumnsType } from "antd/lib";
import { PieChartOutlined, PlayCircleOutlined } from "@ant-design/icons";
import { renderActiveShape } from "../Utilities/renderShape";
import { nullPercent } from "../Utilities/nullPercent";
import PortalCard from "../PortalCard";
import { BaseType } from "antd/es/typography/Base";

interface RunnerDataType {
    key: React.Key;
    name: string;
    exec: string;
    value: number;
    rate: string;
    color: string;
    text_type: string;
}

const runner_columns: TableColumnsType<RunnerDataType> = [
    {
        title: 'Runner Type',
        dataIndex: 'name',
        render: (_, x) => <Typography.Text type={x.text_type as BaseType}>{x.name}</Typography.Text>,
    },
    {
        title: 'Execution Type',
        dataIndex: 'exec',
        showSorterTooltip: { target: 'full-header' },
        filters: [
            {
                text: 'Remote',
                value: 'Remote',
            },
            {
                text: 'Local',
                value: 'Local',
            },
        ],

        onFilter: (value, record) => record.exec == value,
    },
    {
        title: 'Count',
        dataIndex: 'value',
        sorter: (a, b) => a.value - b.value,
    },
    {
        title: 'Rate (%)',
        dataIndex: 'rate',
        sorter: (a, b) => parseFloat(a.rate) - parseFloat(b.rate),
    },
];

function colorSwitchOnExecStrat(exec: string) {
    switch (exec) {
        case "Remote": return "#49AA19"
        case "Local": return "#DC4446"
        default: return "#777777"
    }
}

function getTextType(exec: string) {
    switch (exec) {
        case "Remote": return "success"
        case "Local": return "danger"
        default: return "secondary"
    }
}


const RunnerMetrics: React.FC<{ runnerMetrics: RunnerCount[]; }> = ({ runnerMetrics }) => {
    const runner_data: RunnerDataType[] = [];

    var totalCount: number = runnerMetrics.find(x => x.name == "total")?.actionsExecuted ?? 0

    runnerMetrics.map((item: RunnerCount, count: number) => {
        var rd: RunnerDataType = {
            key: count,
            name: item.name ?? "",
            value: item.actionsExecuted ?? 0,
            exec: item.execKind ?? "",
            rate: nullPercent(item.actionsExecuted, totalCount),
            color: colorSwitchOnExecStrat(item.execKind ?? ""),
            text_type: getTextType(item.execKind ?? ""),
        }
        count++;
        if (rd.name != "total") {
            runner_data.push(rd);
        }
    });

    runnerMetrics.sort((x, y) => {
        var a = x.execKind ?? ""
        var b = y.execKind ?? ""
        if (a < b) {
            return -1;
        }
        if (a > b) {
            return 1;
        }
        return 0;

    })

    const [activeIndexRunner, setActiveIndexRunner] = useState(0);
    const onRunnerPieEnter = useCallback(
        (_: any, runner_idx: any) => {
            setActiveIndexRunner(runner_idx);
        },
        [setActiveIndexRunner]
    );
    return (

        <Row justify="space-around" align="top">
            <Col span="10">
                <PortalCard icon={<PieChartOutlined />} titleBits={["Action Runners Breakdown"]}>
                    <PieChart width={500} height={500}>
                        <Pie
                            activeIndex={activeIndexRunner}
                            activeShape={renderActiveShape}
                            data={runner_data}
                            dataKey="value"
                            nameKey="name"
                            cx="50%"
                            cy="50%"
                            innerRadius={70}
                            outerRadius={90}
                            onMouseEnter={onRunnerPieEnter}>
                            {
                                runner_data.map((entry, runner_index) => (
                                    <Cell key={`cell-${runner_index}`} fill={entry.color} />
                                ))
                            }
                        </Pie>
                        <Legend layout="vertical" />
                    </PieChart>
                </PortalCard>
            </Col>
            <Col span="12">
                <PortalCard icon={<PlayCircleOutlined />} titleBits={["Action Runner Data"]}>
                    <Table
                        columns={runner_columns}
                        dataSource={runner_data}
                        showSorterTooltip={{ target: 'sorter-icon' }}
                        pagination={false} />

                </PortalCard>
            </Col>
            <Col span="2" />
        </Row>
    )
}

export default RunnerMetrics;