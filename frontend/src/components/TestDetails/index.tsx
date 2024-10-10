import React, { useState } from 'react';
import { Space, Row, Statistic } from 'antd';
import { TestStatusEnum } from '../TestStatusTag';
import type { StatisticProps } from "antd/lib";
import CountUp from 'react-countup';
import { useQuery } from '@apollo/client';
import { FindTestsQueryVariables } from '@/graphql/__generated__/graphql';
import TestGridRow from '../TestGridRow';
import PortalAlert from '../PortalAlert';
import { AreaChart, Area, CartesianGrid, XAxis, YAxis, Tooltip } from 'recharts';
import { FIND_TESTS_WITH_CACHE } from './graphql';
import PortalCard from '../PortalCard';
import { FieldTimeOutlined, BorderInnerOutlined } from '@ant-design/icons/lib/icons';

interface Props {
    label: string
}

const formatter: StatisticProps['formatter'] = (value) => (
    <CountUp end={value as number} separator="," />
);

export interface TestStatusType {
    label: string
    invocationId: string,
    status: TestStatusEnum
}
const PAGE_SIZE = 10
interface GraphDataPoint {
    name: string
    duration: number
    local: boolean
    remote: boolean
}


const TestDetails: React.FC<Props> = ({ label }) => {

    const [variables, setVariables] = useState<FindTestsQueryVariables>({ first: 1000, where: { label: label } })
    const { loading: labelLoading, data: labelData, previousData: labelPreviousData, error: labelError } = useQuery(FIND_TESTS_WITH_CACHE, {
        variables: variables,
        fetchPolicy: 'cache-and-network',
        //pollInterval: 120000,
    });


    const data = labelLoading ? labelPreviousData : labelData;
    var result: GraphDataPoint[] = []
    var totalCnt: number = 0
    var local_cached: number = 0
    var remote_cached: number = 0
    var total_duration: number = 0

    if (labelError) {
        <PortalAlert className="error" message="There was a problem communicating w/the backend server." />
    } else {
        totalCnt = data?.findTests.totalCount ?? 0
        data?.findTests.edges?.map(edge => {
            var row = edge?.node
            result.push({
                name: row?.bazelInvocation?.invocationID ?? "",
                duration: row?.durationMs ?? 0,
                local: row?.cachedLocally ?? false,
                remote: row?.cachedRemotely ?? false
            })
            if (row?.cachedLocally) {
                local_cached++
            }
            if (row?.cachedRemotely) {
                remote_cached++
            }
            total_duration += row?.durationMs ?? 0
        });
    }

    return (
        <Space direction="vertical" size="middle" style={{ display: 'flex' }}>
            <h1>{label}</h1>
            <Row>
                <Space size="large">
                    <Statistic title="Average Duration" value={total_duration / totalCnt} formatter={formatter} />
                    <Statistic title="Total Runs" value={totalCnt} formatter={formatter} />
                    <Statistic title="Cached Locally" value={local_cached} formatter={formatter} />
                    <Statistic title="Cached Remotely" value={remote_cached} formatter={formatter} valueStyle={{ color: "#82ca9d" }} />
                </Space>
            </Row>
            <PortalCard icon={<FieldTimeOutlined />} titleBits={["Test Duration Over Time"]} >
                <AreaChart width={1500} height={250} data={result}
                    margin={{ top: 10, right: 30, left: 0, bottom: 0 }}>
                    <defs>
                        <linearGradient id="colorUv" x1="0" y1="0" x2="0" y2="1">
                            <stop offset="5%" stopColor="#8884d8" stopOpacity={0.8} />
                            <stop offset="95%" stopColor="#8884d8" stopOpacity={0} />
                        </linearGradient>
                    </defs>
                    <XAxis />
                    <YAxis />
                    <CartesianGrid strokeDasharray="3 3" />
                    <Tooltip />
                    <Area type="monotone" dataKey="duration" stroke="#8884d8" fillOpacity={1} fill="url(#colorUv)" />
                </AreaChart>
            </PortalCard>
            <Row>
                <PortalCard icon={<BorderInnerOutlined />} titleBits={["Test Pass/Fail Grid"]}>
                    <TestGridRow rowLabel={label} first={1000} reverseOrder={true} />
                </PortalCard>
            </Row>
        </Space>
    );
}
export default TestDetails