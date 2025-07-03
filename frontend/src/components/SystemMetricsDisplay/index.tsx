import React from "react";
import { Statistic, Space, Row } from 'antd';
import { FieldTimeOutlined, BuildOutlined, } from "@ant-design/icons";
import type { StatisticProps } from "antd/lib";
import { BuildGraphMetrics, NetworkMetrics, SystemNetworkStats, TimingMetrics } from "@/graphql/__generated__/graphql";
import PortalCard from "../PortalCard";

const SystemMetricsDisplay: React.FC<{
    buildGraphMetrics: BuildGraphMetrics | undefined,
    timingMetrics: TimingMetrics | undefined,
    systemNetworkStats: SystemNetworkStats | undefined
}> = ({
    timingMetrics,
    buildGraphMetrics,
    systemNetworkStats
}) => {
        return (
            <Space size={"large"} direction="vertical" style={{ display: 'flex' }}>
                <PortalCard type="inner" titleBits={["Timing Metrics"]} icon={<FieldTimeOutlined />}>
                    <Row>
                        <Space size={"large"}>
                            <Statistic title="Wall Time(ms)" value={timingMetrics?.wallTimeInMs ?? 0} />
                            <Statistic title="Analysis(ms)" value={timingMetrics?.analysisPhaseTimeInMs ?? 0} />
                            <Statistic title="CPU Time(ms)" value={timingMetrics?.cpuTimeInMs ?? 0} />
                            <Statistic title="Execuction(ms)" value={timingMetrics?.executionPhaseTimeInMs ?? 0} />
                            <Statistic title="Actions Execution Start" value={timingMetrics?.actionsExecutionStartInMs ?? 0} />
                        </Space>
                    </Row>
                </PortalCard>
                <PortalCard type="inner" titleBits={["Build Graph Metrics"]} icon={<BuildOutlined />}>
                    <Row>
                        <Space size={"large"}>
                            <Statistic title="Action Count" value={buildGraphMetrics?.actionCount ?? 0} />
                            <Statistic title="Action Lookup Value Count)" value={buildGraphMetrics?.actionLookupValueCount ?? 0} />
                            <Statistic title="Action Lookup Value Not Including Aspects" value={buildGraphMetrics?.actionLookupValueCountNotIncludingAspects ?? 0} />
                            <Statistic title="Input File Configured Target Count" value={buildGraphMetrics?.inputFileConfiguredTargetCount ?? 0} />
                            <Statistic title="Output Artifact Count" value={buildGraphMetrics?.outputArtifactCount ?? 0} />
                            <Statistic title="Output File Configured Target Count" value={buildGraphMetrics?.outputFileConfiguredTargetCount ?? 0} />
                            <Statistic title="Post Invocation Sky frameNode Count" value={buildGraphMetrics?.postInvocationSkyframeNodeCount ?? 0} />
                        </Space>
                    </Row>
                </PortalCard>
                <PortalCard type="inner" titleBits={["System Network Metrics"]} icon={<FieldTimeOutlined />}>
                    <Row>
                        <Space size="large">
                            <Statistic title="Bytes Recieved" value={systemNetworkStats?.bytesRecv ?? 0} />
                            <Statistic title="Bytes Sent" value={systemNetworkStats?.bytesSent ?? 0} />
                            <Statistic title="Packets Recieved" value={systemNetworkStats?.packetsRecv ?? 0} />
                            <Statistic title="Packets Sent" value={systemNetworkStats?.packetsSent ?? 0} />
                            <Statistic title="Peak Bytes Recieved(/s)" value={systemNetworkStats?.peakBytesRecvPerSec ?? 0} />
                            <Statistic title="Peak Bytes Sent(/s)" value={systemNetworkStats?.peakBytesSentPerSec ?? 0} />
                            <Statistic title="Peak Packets Recieved(/s)" value={systemNetworkStats?.peakPacketsRecvPerSec ?? 0} />
                            <Statistic title="Peak Packets Sent(/s)" value={systemNetworkStats?.peakPacketsSentPerSec ?? 0} />
                        </Space>
                    </Row>
                </PortalCard>
            </Space>
        )
    }

export default SystemMetricsDisplay;