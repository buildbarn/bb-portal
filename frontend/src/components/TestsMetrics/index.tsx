import React from "react";
import { Table, Row, Statistic, Space } from 'antd';
import type { StatisticProps, TableColumnsType } from "antd/lib";
import TestStatusTag from "../TestStatusTag";
import { TestCollection } from "@/graphql/__generated__/graphql";
import { TestStatusEnum } from "../TestStatusTag";
import NullBooleanTag from "../NullableBooleanTag";
import PortalCard from "../PortalCard";
import { SearchFilterIcon, SearchWidget } from '@/components/SearchWidgets';
import { SearchOutlined, ExperimentOutlined, } from "@ant-design/icons";
import Link from "next/link";
import { millisecondsToTime } from "../Utilities/time";
import styles from "../../theme/theme.module.css"
interface TestDataType {
    key: React.Key;
    status: string;
    name: string;
    value: number;
    strategy: string;
    cached_local: boolean | null;
    cached_remote: boolean | null;
    duration: number;
}

const test_columns: TableColumnsType<TestDataType> = [
    {
        title: "Status",
        dataIndex: "status",
        render: (x) => <TestStatusTag displayText={true} key="status" status={x as TestStatusEnum} />,
        showSorterTooltip: { target: 'full-header' },
        filters: [
            {
                text: 'No Status',
                value: 'NO_STATUS',
            },
            {
                text: 'Passed',
                value: 'PASSED',
            },
            {
                text: "Flaky",
                value: "FLAKY"
            },
            {
                text: "Timeout",
                value: "TIMEOUT"
            },
            {
                text: "Failed",
                value: "FAILED"
            },
            {
                text: "Incomplete",
                value: "INCOMPLETE"
            },
            {
                text: "Remote Failure",
                value: "REMOTE_FAILURE"
            },
            {
                text: "Failed to Build",
                value: "FAILED_TO_BUILD"
            },
            {
                text: "Tool Halted Before Testing",
                value: "TOOL_HALTED_BEFORE_TESTING"
            },
        ],
        onFilter: (value, record) => record.status == value,
    },
    {
        title: "Label",
        dataIndex: "name",
        render: (_, record) => <Link href={"/tests/" + btoa(encodeURIComponent(record.name))}>{record.name}</Link>,
        filterSearch: true,
        filterDropdown: filterProps => (
            <SearchWidget placeholder="Target Pattern..." {...filterProps} />
        ),
        filterIcon: filtered => <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />,
        onFilter: (value, record) => (record.name.includes(value.toString()) ? true : false)
    },
    {
        title: "Strategy",
        dataIndex: "strategy",
        sorter: (a, b) => a.strategy.localeCompare(b.strategy),
        filters: [
            {
                text: "Remote Cache Hit",
                value: "remote cache hit"
            },
            {
                text: "Remote",
                value: "remote"
            },
            {
                text: "Linux Sandbox",
                value: "linux-sandbox"
            },
            {
                text: "Disk Cache Hit",
                value: "disk cache hit"
            },
            {
                text: "None",
                value: ""
            },
        ],
        filterIcon: filtered => <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />,
        onFilter: (value, record) => record.strategy == value
    },
    {
        title: "Cached Locally",
        dataIndex: "cached_local",
        render: (x) => <NullBooleanTag key="cached_locally" status={x as boolean | null} />,
        sorter: (a, b) => Number(a.cached_local) - Number(b.cached_local),
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
        onFilter: (value, record) => record.cached_local == value,
    },
    {
        title: "Cached Remotely",
        dataIndex: "cached_remote",
        render: (x) => <NullBooleanTag key="cached_remotely" status={x as boolean | null} />,
        sorter: (a, b) => Number(a.cached_local) - Number(b.cached_local),
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
        onFilter: (value, record) => record.cached_remote == value
    },
    {
        title: "Duration",
        dataIndex: "value",
        render: (_, record) => <span className={styles.numberFormat}>{millisecondsToTime(record.duration)}</span>,
        align: "right",
        sorter: (a, b) => a.value - b.value,
    },

]

const TestMetricsDisplay: React.FC<{
    testMetrics: TestCollection[] | undefined | null,
    targetTimes: Map<string, number>,
}> = ({
    testMetrics,
    targetTimes,
}) => {
        const totalTests: number = testMetrics?.length ?? 0
        const test_data: TestDataType[] = []

        testMetrics?.map((item: TestCollection, index) => {
            var label = item.label ?? "NO_TARGET_LABEL"

            var row: TestDataType = {
                key: "test-data-type-row-" + index,
                name: item.label ?? "",
                value: item.durationMs ?? 0,
                strategy: item.strategy ?? "",
                cached_local: item.cachedLocally ?? null,
                cached_remote: item.cachedRemotely ?? null,
                duration: (item.durationMs ?? 0) + (targetTimes.get(item.label ?? "") ?? 0),
                status: item.overallStatus ?? ""
            }
            test_data.push(row);
        })

        var numPassed = test_data.filter(x => x.status == "PASSED").length
        var numFlaky = test_data.filter(x => x.status == "FLAKY").length
        var numFailed = test_data.filter(x => x.status == "FAILED").length
        var numExecutedLocally = test_data.filter(x => x.strategy == "linux-sandbox").length
        var numExecutedRemotely = test_data.filter(x => x.strategy == "remote").length
        var localCacheHit = test_data.filter(x => x.strategy in ["disk cache hit", ""] || x.cached_local == true).length
        var remoteCacheHit = test_data.filter(x => x.strategy == "remote cache hit" || x.cached_remote == true).length

        return (
            <Space direction="vertical" size="middle" style={{ display: 'flex' }} >
                <PortalCard type="inner" icon={<ExperimentOutlined />} titleBits={["Tests"]}>
                    <Row>
                        <Space size="large">
                            <Statistic title="Tests Completed" value={totalTests} />
                            <Statistic title="Passed" value={numPassed} />
                            <Statistic title="Flaky" value={numFlaky} />
                            <Statistic title="Failed" value={numFailed} />
                            <Statistic title="Executed Locally" value={numExecutedLocally} />
                            <Statistic title="Executed Remotely" value={numExecutedRemotely} />
                            <Statistic title="Local Cache Hit" value={localCacheHit} />
                            <Statistic title="Remote Cache Hit" value={remoteCacheHit} />
                        </Space>
                    </Row>
                    <Row justify="space-around" align="middle">
                        <Table
                            columns={test_columns}
                            dataSource={test_data}
                            showSorterTooltip={{ target: 'sorter-icon' }}
                        />
                    </Row>
                </PortalCard>
            </Space>
        )
    }

export default TestMetricsDisplay;