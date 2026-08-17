import { FieldTimeOutlined } from "@ant-design/icons";
import { Space } from "antd";
import type React from "react";
import type { BazelInvocationMetricsSystemNetworkStatsFragment } from "@/graphql/__generated__/graphql";
import { readableFileSize } from "@/utils/filesize";
import { MultiStatistic } from "../MultiStatistic";
import PortalCard from "../PortalCard";

interface Props {
  systemNetworkStats: BazelInvocationMetricsSystemNetworkStatsFragment;
  cardStyle?: React.CSSProperties;
}

export const SystemNetworkStatsDisplay: React.FC<Props> = ({
  systemNetworkStats,
  cardStyle,
}) => {
  return (
    <PortalCard
      type="inner"
      titleBits={["System Network Metrics"]}
      style={cardStyle}
      icon={<FieldTimeOutlined />}
    >
      <Space size="large">
        <MultiStatistic
          title="Packets Received"
          values={[
            {
              key: "count",
              value: (systemNetworkStats.packetsRecv ?? 0n).toString(),
            },
            {
              key: "size",
              value: readableFileSize(systemNetworkStats.bytesRecv ?? 0n),
            },
          ]}
        />
        <MultiStatistic
          title="Packets Sent"
          values={[
            {
              key: "count",
              value: (systemNetworkStats.packetsSent ?? 0n).toString(),
            },
            {
              key: "size",
              value: readableFileSize(systemNetworkStats.bytesSent ?? 0n),
            },
          ]}
        />
        <MultiStatistic
          title="Peak Packets Received(/s)"
          values={[
            {
              key: "count",
              value: (
                systemNetworkStats.peakPacketsRecvPerSec ?? 0n
              ).toString(),
            },
            {
              key: "size",
              value: readableFileSize(
                systemNetworkStats.peakBytesRecvPerSec ?? 0n,
              ),
            },
          ]}
        />
        <MultiStatistic
          title="Peak Packets Sent(/s)"
          values={[
            {
              key: "count",
              value: (
                systemNetworkStats.peakPacketsSentPerSec ?? 0n
              ).toString(),
            },
            {
              key: "size",
              value: readableFileSize(
                systemNetworkStats.peakBytesSentPerSec ?? 0n,
              ),
            },
          ]}
        />
      </Space>
    </PortalCard>
  );
};
