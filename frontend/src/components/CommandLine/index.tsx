import React from "react";
import { Space, Table, Row, Col, Statistic, List } from 'antd';
import { CodeOutlined, DeploymentUnitOutlined, SearchOutlined } from '@ant-design/icons';
import type { StatisticProps, TableColumnsType } from "antd/lib";
import CountUp from 'react-countup';
import { BazelCommand, TargetMetrics, TargetPair } from "@/graphql/__generated__/graphql";
import PortalCard from "../PortalCard";
import { SearchFilterIcon, SearchWidget } from '@/components/SearchWidgets';
import NullBooleanTag from "../NullableBooleanTag";
import styles from "../../theme/theme.module.css"
import { millisecondsToTime } from "../Utilities/time";



const CommandLineDisplay: React.FC<{ commandLineData: BazelCommand | undefined | null }> = ({
    commandLineData: commandLineData
}) => {

    const createUnorderedList = (items: string[]): JSX.Element => {
        return (
            <ul>
                {items.map((item, index) => (
                    <li key={index}>{item}</li>
                ))}
            </ul>
        );
    };

    var commandLineOptions: string[] = []
    commandLineData?.cmdLine?.forEach(x => commandLineOptions.push(x ?? ""))

    return (

        <Space direction="vertical" size="middle" style={{ display: 'flex' }} >
            <PortalCard type="inner" icon={<CodeOutlined />} titleBits={["Command Line"]}>
                <Row>
                    <Space size="large">
                        <strong>Explicit Command Line:</strong>
                        <div>
                            {commandLineData?.executable} {commandLineData?.command} {commandLineData?.residual} {commandLineData?.explicitCmdLine}
                        </div>
                    </Space>
                </Row>
                <Row>
                    <Space size="large">
                        <div>
                            <List
                                bordered
                                header={<div><strong>Effective Command Line Options:</strong></div>}
                                dataSource={commandLineData?.cmdLine?.filter(x => x !== undefined).toSorted() as string[]}
                                renderItem={(item) => <List.Item>{item}</List.Item>}
                            />

                            <List
                                bordered
                                header={<div><strong>Explicit Startup Options:</strong></div>}
                                dataSource={commandLineData?.explicitStartupOptions?.filter(x => x !== undefined) as string[]}
                                renderItem={(item) => <List.Item>{item}</List.Item>}
                            />
                            <List
                                bordered
                                header={<div><strong>Effective Startup Options:</strong></div>}
                                dataSource={commandLineData?.startupOptions?.filter(x => x !== undefined) as string[]}
                                renderItem={(item) => <List.Item>{item}</List.Item>}
                            />
                        </div>
                    </Space>
                </Row>
            </PortalCard>
        </Space>
    )
}

export default CommandLineDisplay;