import React from "react";
import { Space, Table, Row, Col, Statistic, List } from 'antd';
import { CodeOutlined, DeploymentUnitOutlined, SearchOutlined } from '@ant-design/icons';
import type { StatisticProps, TableColumnsType } from "antd/lib";
import { BazelCommand, TargetMetrics, TargetPair } from "@/graphql/__generated__/graphql";
import PortalCard from "../PortalCard";
import { SearchFilterIcon, SearchWidget } from '@/components/SearchWidgets';
import NullBooleanTag from "../NullableBooleanTag";
import styles from "../../theme/theme.module.css"



const CommandLineDisplay: React.FC<{ commandLineData: BazelCommand | undefined | null }> = ({
    commandLineData: commandLineData
}) => {

    var commandLineOptions: string[] = []
    commandLineData?.cmdLine?.forEach(x => commandLineOptions.push(x ?? ""))
    var cmdLine = [commandLineData?.executable, commandLineData?.command, commandLineData?.residual, commandLineData?.explicitCmdLine].join(" ")

    return (

        <Space direction="vertical" size="middle" style={{ display: 'flex' }} >
            <PortalCard type="inner" icon={<CodeOutlined />} titleBits={["Explicit Command Line:", cmdLine]}>
                <Row>
                    <List
                        bordered
                        size="small"
                        header={<div><strong>Raw Command Line Options:</strong></div>}
                        dataSource={commandLineData?.cmdLine?.filter(x => x !== undefined) as string[]}
                        renderItem={(item) => <List.Item>{item}</List.Item>}
                    />
                </Row>
                <Row>
                    <List
                        bordered
                        size="small"
                        header={<div><strong>Explicit Startup Options:</strong></div>}
                        dataSource={commandLineData?.explicitStartupOptions?.filter(x => x !== undefined) as string[]}
                        renderItem={(item) => <List.Item>{item}</List.Item>}
                    />
                </Row>
                <Row>
                    <List
                        bordered
                        size="small"
                        header={<div><strong>Effective Startup Options:</strong></div>}
                        dataSource={commandLineData?.startupOptions?.filter(x => x !== undefined) as string[]}
                        renderItem={(item) => <List.Item>{item}</List.Item>}
                    />
                </Row>
            </PortalCard>
        </Space>
    )
}

export default CommandLineDisplay;