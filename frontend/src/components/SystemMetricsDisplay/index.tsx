import React from "react";
import { Statistic, Space, Row } from 'antd';
import { FieldTimeOutlined, BuildOutlined, } from "@ant-design/icons";
import type { StatisticProps } from "antd/lib";
import { BuildGraphMetrics, NetworkMetrics, SystemNetworkStats, TimingMetrics } from "@/graphql/__generated__/graphql";
import PortalCard from "../PortalCard";
import { readableDurationFromMilliseconds } from "@/utils/time";
import { readableFileSize } from "@/utils/filesize";

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
                            <Statistic title="Wall Time" value={readableDurationFromMilliseconds(timingMetrics?.wallTimeInMs ?? 0, {smallestUnit: "ms"})} />
                            <Statistic title="Analysis" value={readableDurationFromMilliseconds(timingMetrics?.analysisPhaseTimeInMs ?? 0, {smallestUnit: "ms"})} />
                            <Statistic title="CPU Time" value={readableDurationFromMilliseconds(timingMetrics?.cpuTimeInMs ?? 0, {smallestUnit: "ms"})} />
                            <Statistic title="Execution" value={readableDurationFromMilliseconds(timingMetrics?.executionPhaseTimeInMs ?? 0, {smallestUnit: "ms"})} />
                            <Statistic title="Actions Execution Start" value={readableDurationFromMilliseconds(timingMetrics?.actionsExecutionStartInMs ?? 0, {smallestUnit: "ms"})} />
                        </Space>
                    </Row>
                </PortalCard>
                <PortalCard type="inner" titleBits={["Build Graph Metrics"]} icon={<BuildOutlined />}>
                    <Row>
                        <Space size={"large"}>
                            <Statistic title="Action Count" value={buildGraphMetrics?.actionCount ?? 0} />
                            <Statistic title="Action Lookup Value Count" value={buildGraphMetrics?.actionLookupValueCount ?? 0} />
                            <Statistic title="Action Lookup Value Not Including Aspects" value={buildGraphMetrics?.actionLookupValueCountNotIncludingAspects ?? 0} />
                            <Statistic title="Input File Configured Target Count" value={buildGraphMetrics?.inputFileConfiguredTargetCount ?? 0} />
                            <Statistic title="Output Artifact Count" value={buildGraphMetrics?.outputArtifactCount ?? 0} />
                            <Statistic title="Output File Configured Target Count" value={buildGraphMetrics?.outputFileConfiguredTargetCount ?? 0} />
                            <Statistic title="Post Invocation Skyframe Node Count" value={buildGraphMetrics?.postInvocationSkyframeNodeCount ?? 0} />
                        </Space>
                    </Row>
                </PortalCard>
                <PortalCard type="inner" titleBits={["System Network Metrics"]} icon={<FieldTimeOutlined />}>
                    <Row>
                        <Space size="large">
                            <Statistic title="Bytes Recieved" value={readableFileSize(systemNetworkStats?.bytesRecv ?? 0)} />
                            <Statistic title="Bytes Sent" value={readableFileSize(systemNetworkStats?.bytesSent ?? 0)} />
                            <Statistic title="Packets Recieved" value={systemNetworkStats?.packetsRecv ?? 0} />
                            <Statistic title="Packets Sent" value={systemNetworkStats?.packetsSent ?? 0} />
                            <Statistic title="Peak Bytes Recieved(/s)" value={readableFileSize(systemNetworkStats?.peakBytesRecvPerSec ?? 0)} />
                            <Statistic title="Peak Bytes Sent(/s)" value={readableFileSize(systemNetworkStats?.peakBytesSentPerSec ?? 0)} />
                            <Statistic title="Peak Packets Recieved(/s)" value={systemNetworkStats?.peakPacketsRecvPerSec ?? 0} />
                            <Statistic title="Peak Packets Sent(/s)" value={systemNetworkStats?.peakPacketsSentPerSec ?? 0} />
                        </Space>
                    </Row>
                </PortalCard>
            </Space>
        )
    }

export default SystemMetricsDisplay;