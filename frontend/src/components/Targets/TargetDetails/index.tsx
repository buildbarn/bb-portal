import React, { useState } from 'react';
import { Space, Row, Statistic } from 'antd';
import { TestStatusEnum } from '../../TestStatusTag';
import type { StatisticProps } from "antd/lib";
import CountUp from 'react-countup';
import { useQuery } from '@apollo/client';
import { FindTargetsQueryVariables } from '@/graphql/__generated__/graphql';
import PortalAlert from '../../PortalAlert';
import { AreaChart, Area, CartesianGrid, XAxis, YAxis, Tooltip } from 'recharts';
import PortalCard from '../../PortalCard';
import { FieldTimeOutlined, BorderInnerOutlined } from '@ant-design/icons/lib/icons';
import TargetGridRow from '../TargetGridRow';
import { FIND_TARGETS } from '@/app/targets/graphql';

interface Props {
    label: string
}

const formatter: StatisticProps['formatter'] = (value) => (
    <CountUp end={value as number} separator="," />
);

export interface TargetStatusType {
    label: string
    invocationId: string,
    status: TestStatusEnum
}
const PAGE_SIZE = 10
interface GraphDataPoint {
    name: string
    duration: number
}


const TestDetails: React.FC<Props> = ({ label }) => {

    const [variables, setVariables] = useState<FindTargetsQueryVariables>({ first: 1000, where: { label: label } })
    const { loading: labelLoading, data: labelData, previousData: labelPreviousData, error: labelError } = useQuery(FIND_TARGETS, {
        variables: variables,
        fetchPolicy: 'network-only',
        //pollInterval: 120000,
    });


    const data = labelLoading ? labelPreviousData : labelData;
    var result: GraphDataPoint[] = []
    var totalCnt: number = 0
    var total_duration: number = 0

    if (labelError) {
        <PortalAlert className="error" message="There was a problem communicating w/the backend server." />
    } else {
        totalCnt = data?.findTargets.totalCount ?? 0
        data?.findTargets.edges?.map(edge => {
            var row = edge?.node
            result.push({
                name: row?.bazelInvocation?.invocationID ?? "",
                duration: row?.durationInMs ?? 0,
            })
            total_duration += row?.durationInMs ?? 0
        });
    }

    return (
        <Space direction="vertical" size="middle" style={{ display: 'flex' }}>
            <h1>{label}</h1>
            <Row>
                <Space size="large">
                    <Statistic title="Average Duration" value={total_duration / totalCnt} formatter={formatter} />
                    <Statistic title="Total Runs" value={totalCnt} formatter={formatter} />
                </Space>
            </Row>
            <PortalCard icon={<FieldTimeOutlined />} titleBits={["Target Duration Over Time"]} >
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
                <PortalCard icon={<BorderInnerOutlined />} titleBits={["Target Pass/Fail Grid"]}>
                    <TargetGridRow rowLabel={label} first={1000} reverseOrder={true} />
                </PortalCard>
            </Row>
        </Space>
    );
}
export default TestDetails