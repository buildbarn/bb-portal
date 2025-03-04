import React from "react";
import { Statistic, Space, Row } from 'antd';
import { FieldTimeOutlined } from "@ant-design/icons";
import type { StatisticProps } from "antd/lib";
import { NetworkMetrics, SystemNetworkStats } from "@/graphql/__generated__/graphql";
import PortalCard from "../PortalCard";

const NetworkMetricsDisplay: React.FC<{ networkMetrics: NetworkMetrics | undefined }> = ({
    networkMetrics: networkMetrics
}) => {
    const systemNetworkStats: SystemNetworkStats | undefined = networkMetrics?.systemNetworkStats ?? undefined
    return (
        <Space direction="vertical" size="middle" style={{ display: 'flex' }} >
            <PortalCard titleBits={["System Network Metrics"]} icon={<FieldTimeOutlined />}>
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

export default NetworkMetricsDisplay