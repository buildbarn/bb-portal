import React from "react";
import { Space, Table, Row, Col, Statistic } from 'antd';
import { DeploymentUnitOutlined, SearchOutlined } from '@ant-design/icons';
import type { StatisticProps, TableColumnsType } from "antd/lib";
import { Target, TargetMetrics } from "@/graphql/__generated__/graphql";
import PortalCard from "../PortalCard";
import { SearchFilterIcon, SearchWidget } from '@/components/SearchWidgets';
import NullBooleanTag from "../NullableBooleanTag";
import TargetAbortReasonTag, { AbortReasonsEnum } from "./targetAbortReasonTag";
import styles from "../../theme/theme.module.css"
import Link from "next/link";
import { readableDurationFromMilliseconds } from "@/utils/time";
interface TargetDataType {
    key: React.Key;
    name: string;           //label
    success: boolean;       //overall success/fail
    value: number;          //duration
    target_kind: string;    //target kind if available
    failure_reason: string  //failure reason if any
}

const TargetMetricsDisplay: React.FC<{
    targetMetrics: TargetMetrics | undefined | null,
    targetData: Target[] | undefined | null,
}> = ({
    targetMetrics,
    targetData,
}) => {

        var target_data: TargetDataType[] = []
        var count = 0;
        var all_types: string[] = []
        var targets_skipped: number = 0;
        var targets_built_successfully: number = 0;

        targetData?.map(x => {
            count++;
            var targetKind = x.targetKind ?? ""
            var failureReason = x.abortReason ?? ""

            if (failureReason == "SKIPPED") {
                targets_skipped++;
            }

            if (x.success == true) {
                targets_built_successfully++;
            }

            var row: TargetDataType = {
                key: "target_data_type" + count.toString(),
                name: x.label ?? "",
                success: x.success ?? false,
                value: x.durationInMs ?? 0,
                target_kind: targetKind,
                failure_reason: failureReason,
            }
            all_types.push(targetKind)
            target_data.push(row)

        })

        const targets_analyzed: number = targetData?.length ?? 0
        const type_filters: string[] = Array.from(new Set(all_types))

        const target_columns: TableColumnsType<TargetDataType> = [
            {

                title: "Label",
                dataIndex: "name",
                filterSearch: true,
                render: (_, record) => <Link href={"/targets/" + btoa(encodeURIComponent(record.name))}>{record.name}</Link>,
                filterDropdown: filterProps => (
                    <SearchWidget placeholder="Target Pattern..." {...filterProps} />
                ),
                filterIcon: filtered => <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />,
                onFilter: (value, record) => (record.name.includes(value.toString()) ? true : false)
            },
            {
                title: "Duration",
                dataIndex: "value",
                align: "right",
                render: (_, record) => <span className={styles.numberFormat}>{readableDurationFromMilliseconds(record.value, {smallestUnit: "ms"})}</span>,
                sorter: (a, b) => a.value - b.value,
            },
            {
                title: "Target Type",
                dataIndex: "target_kind",
                filters: type_filters.map(x => ({ text: x, value: x })),
                filterIcon: filtered => <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />,
                onFilter: (value, record) => (record.target_kind.includes(value.toString()) ? true : false),
                sorter: (a, b) => a.target_kind.localeCompare(b.target_kind),

            },
            {
                title: "Abort Reason",
                dataIndex: "failure_reason",
                filters: [
                    {
                        text: "Skipped",
                        value: "SKIPPED"
                    },
                    {
                        text: "User Interrupted",
                        value: "USER_INTERRUPTED"
                    },
                    {
                        text: "Time Out",
                        value: "TIME_OUT"
                    },
                    {
                        text: "Remote Environment Failure",
                        value: "REMOTE_ENVIRONMENT_FAILURE"
                    },
                    {
                        text: "Internal",
                        value: "INTERNAL"
                    },
                    {
                        text: "Loading Failure",
                        value: "LOADING_FAILURE"
                    },
                    {
                        text: "Analysis Failure",
                        value: "ANALYSIS_FAILURE"
                    },
                    {
                        text: "No Analyze",
                        value: "NO_ANALYZE"
                    },
                    {
                        text: "No Build",
                        value: "NO_BUILD"
                    },
                    {
                        text: "Incomplete",
                        value: "INCOMPLETE"
                    },
                    {
                        text: "Out of Memory",
                        value: "OUT_OF_MEMORY"
                    },
                ],
                render: (x) => <TargetAbortReasonTag key="failure_reason" reason={x as AbortReasonsEnum} />,
                filterIcon: filtered => <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />,
                onFilter: (value, record) => record.failure_reason == value,
                sorter: (a, b) => a.failure_reason.localeCompare(b.failure_reason),

            },
            {
                title: "Overall Success",
                dataIndex: "success",
                render: (x) => <NullBooleanTag key="success" status={x as boolean | null} />,
                sorter: (a, b) => Number(a.success) - Number(b.success),
                filters: [
                    {
                        text: "Yes",
                        value: true
                    },
                    {
                        text: "No",
                        value: false
                    }
                ],
                filterIcon: filtered => <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />,
                onFilter: (value, record) => record.success == value,
            },
        ]

        return (
            <Space direction="vertical" size="middle" style={{ display: 'flex' }} >
                <PortalCard type="inner" icon={<DeploymentUnitOutlined />} titleBits={["Targets"]}>
                    <Row>
                        <Space size="large">
                            <Statistic title="Targets Analyzed" value={targets_analyzed} />
                            <Statistic title="Targets Built Successfully" value={targets_built_successfully} valueStyle={{ color: "green" }} />
                            <Statistic title="Targets Skipped" value={targets_skipped} valueStyle={{ color: "purple" }} />
                            <Statistic title="Targets Configured" value={targetMetrics?.targetsConfigured ?? 0} />
                            <Statistic title="Targets Configured Not Including Aspects" value={targetMetrics?.targetsConfiguredNotIncludingAspects ?? 0} />
                        </Space>
                    </Row>
                    <Row justify="space-around" align="middle">
                        <Table
                            columns={target_columns}
                            dataSource={target_data}
                            showSorterTooltip={{ target: 'sorter-icon' }}
                        />
                    </Row>
                </PortalCard>
            </Space>
        )
    }

export default TargetMetricsDisplay;