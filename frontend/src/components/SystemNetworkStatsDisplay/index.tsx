import { FieldTimeOutlined } from "@ant-design/icons";
import { Row, Space, Statistic } from "antd";
import type React from "react";
import type { BazelInvocationMetricsSystemNetworkStatsFragment } from "@/graphql/__generated__/graphql";
import { readableFileSize } from "@/utils/filesize";
import PortalCard from "../PortalCard";

interface Props {
  systemNetworkStats:
    | BazelInvocationMetricsSystemNetworkStatsFragment
    | undefined
    | null;
}

export const SystemNetworkStatsDisplay: React.FC<Props> = ({
  systemNetworkStats,
}) => {
  return (
    <PortalCard
      type="inner"
      titleBits={["System Network Metrics"]}
      icon={<FieldTimeOutlined />}
    >
      <Row>
        <Space size="large">
          <Statistic
            title="Bytes Recieved"
            value={readableFileSize(systemNetworkStats?.bytesRecv ?? 0)}
          />
          <Statistic
            title="Bytes Sent"
            value={readableFileSize(systemNetworkStats?.bytesSent ?? 0)}
          />
          <Statistic
            title="Packets Recieved"
            value={systemNetworkStats?.packetsRecv ?? 0}
          />
          <Statistic
            title="Packets Sent"
            value={systemNetworkStats?.packetsSent ?? 0}
          />
          <Statistic
            title="Peak Bytes Recieved(/s)"
            value={readableFileSize(
              systemNetworkStats?.peakBytesRecvPerSec ?? 0,
            )}
          />
          <Statistic
            title="Peak Bytes Sent(/s)"
            value={readableFileSize(
              systemNetworkStats?.peakBytesSentPerSec ?? 0,
            )}
          />
          <Statistic
            title="Peak Packets Recieved(/s)"
            value={systemNetworkStats?.peakPacketsRecvPerSec ?? 0}
          />
          <Statistic
            title="Peak Packets Sent(/s)"
            value={systemNetworkStats?.peakPacketsSentPerSec ?? 0}
          />
        </Space>
      </Row>
    </PortalCard>
  );
};
