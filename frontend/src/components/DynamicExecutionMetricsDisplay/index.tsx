import { ThunderboltOutlined } from "@ant-design/icons";
import { Table, type TableColumnsType } from "antd";
import type { BazelInvocationMetricsDynamicExecutionMetricsFragment } from "@/graphql/__generated__/graphql";
import { PortalCard } from "../PortalCard";

interface Props {
  dynamicExecutionMetrics: BazelInvocationMetricsDynamicExecutionMetricsFragment;
}

type RaceStatistic = NonNullable<
  BazelInvocationMetricsDynamicExecutionMetricsFragment["raceStatistics"]
>[number];

const columns: TableColumnsType<RaceStatistic> = [
  {
    title: "Mnemonic",
    key: "mnemonic",
    render: (_, record) => record.mnemonic,
  },
  {
    title: "Local Runner",
    key: "localRunner",
    render: (_, record) => record.localRunner,
  },
  {
    title: "Remote Runner",
    key: "remoteRunner",
    render: (_, record) => record.remoteRunner,
  },
  {
    title: "Local Wins",
    key: "localWins",
    align: "right",
    render: (_, record) => record.localWins ?? 0,
    sorter: (a, b) => (a.localWins ?? 0) - (b.localWins ?? 0),
  },
  {
    title: "Remote Wins",
    key: "remoteWins",
    align: "right",
    render: (_, record) => record.remoteWins ?? 0,
    sorter: (a, b) => (a.remoteWins ?? 0) - (b.remoteWins ?? 0),
  },
];

export const DynamicExecutionMetricsDisplay: React.FC<Props> = ({
  dynamicExecutionMetrics,
}) => (
  <PortalCard
    type="inner"
    icon={<ThunderboltOutlined />}
    titleBits={["Dynamic Execution Metrics"]}
  >
    <Table
      columns={columns}
      dataSource={dynamicExecutionMetrics.raceStatistics ?? []}
      pagination={{ defaultPageSize: 20, hideOnSinglePage: true }}
      rowKey="id"
      scroll={{ x: true }}
      size="small"
    />
  </PortalCard>
);
