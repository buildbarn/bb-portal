import React from "react";
import { Statistic, Space, Row } from 'antd';
import { FieldTimeOutlined, BuildOutlined, } from "@ant-design/icons";
import type { StatisticProps } from "antd/lib";
import CountUp from 'react-countup';
import { BuildGraphMetrics, TimingMetrics } from "@/graphql/__generated__/graphql";
import PortalCard from "../PortalCard";


const formatter: StatisticProps['formatter'] = (value) => (
    <CountUp end={value as number} separator="," />
);

const TimingMetricsDisplay: React.FC<{
    buildGraphMetrics: BuildGraphMetrics | undefined,
    timingMetrics: TimingMetrics | undefined
}> = ({
    timingMetrics,
    buildGraphMetrics
}) => {
        return (

            <Space size={"large"} direction="vertical" style={{ display: 'flex' }}>
                <PortalCard titleBits={["Timing Metrics"]} icon={<FieldTimeOutlined />}>
                    <Row>
                        <Space size={"large"}>
                            <Statistic title="Wall Time(ms)" value={timingMetrics?.wallTimeInMs ?? 0} formatter={formatter} />
                            <Statistic title="Analysis(ms)" value={timingMetrics?.analysisPhaseTimeInMs ?? 0} formatter={formatter} />
                            <Statistic title="CPU Time(ms)" value={timingMetrics?.cpuTimeInMs ?? 0} formatter={formatter} />
                            <Statistic title="Execuction(ms)" value={timingMetrics?.executionPhaseTimeInMs ?? 0} formatter={formatter} />
                            <Statistic title="Actions Execution Start" value={timingMetrics?.actionsExecutionStartInMs ?? 0} formatter={formatter} />

                        </Space>
                    </Row>
                </PortalCard>
                <PortalCard titleBits={["Build Graph Metrics"]} icon={<BuildOutlined />}>
                    <Row>
                        <Space size={"large"}>
                            <Statistic title="Action Count" value={buildGraphMetrics?.actionCount ?? 0} formatter={formatter} />
                            <Statistic title="Action Lookup Value Count)" value={buildGraphMetrics?.actionLookupValueCount ?? 0} formatter={formatter} />
                            <Statistic title="Action Lookup Value Not Including Aspects" value={buildGraphMetrics?.actionLookupValueCountNotIncludingAspects ?? 0} formatter={formatter} />
                            <Statistic title="Input File Configured Target Count" value={buildGraphMetrics?.inputFileConfiguredTargetCount ?? 0} formatter={formatter} />
                            <Statistic title="Output Artifact Count" value={buildGraphMetrics?.outputArtifactCount ?? 0} formatter={formatter} />
                            <Statistic title="Output File Configured Target Count" value={buildGraphMetrics?.outputFileConfiguredTargetCount ?? 0} formatter={formatter} />
                            <Statistic title="Post Invocation Sky frameNode Count" value={buildGraphMetrics?.postInvocationSkyframeNodeCount ?? 0} formatter={formatter} />
                        </Space>
                    </Row>
                </PortalCard>
            </Space>

        )
    }

export default TimingMetricsDisplay