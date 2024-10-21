import React from "react";
import { Space, Table, Row, Col, Statistic } from 'antd';
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


    return (
        <Space direction="vertical" size="middle" style={{ display: 'flex' }} >
            <PortalCard type="inner" icon={<CodeOutlined />} titleBits={["Command Line"]}>
                <Row>
                    <Space size="large">
                        <strong>User Command Line:</strong>
                        <div>
                            {commandLineData?.executable} {commandLineData?.command} {commandLineData?.options} {commandLineData?.residual}
                        </div>
                    </Space>
                </Row>
                <Row>
                    <Space size="large">
                        <strong>Effective Command Line:</strong>
                        <div>
                            {commandLineData?.executable} {commandLineData?.command} {commandLineData?.options} {commandLineData?.residual}
                        </div>
                    </Space>
                </Row>
            </PortalCard>
        </Space>
    )
}

export default CommandLineDisplay;