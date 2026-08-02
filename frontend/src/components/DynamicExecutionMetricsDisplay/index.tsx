import { ThunderboltOutlined } from "@ant-design/icons";
import { Table, type TableColumnsType } from "antd";
import type { BazelInvocationMetricsDynamicExecutionMetricsFragment } from "@/graphql/__generated__/graphql";
import PortalCard from "../PortalCard";

interface Props {
  dynamicExecutionMetrics: BazelInvocationMetricsDynamicExecutionMetricsFragment;
}

type RaceStatistic = NonNullable<
  BazelInvocationMetricsDynamicExecutionMetricsFragment["raceStatistics"]
>[number];

const columns: TableColumnsType<RaceStatistic> = [
  { title: "Mnemonic", dataIndex: "mnemonic" },
  { title: "Local Runner", dataIndex: "localRunner" },
  { title: "Remote Runner", dataIndex: "remoteRunner" },
  {
    title: "Local Wins",
    dataIndex: "localWins",
    align: "right",
    render: (value: number | null | undefined) => value ?? 0,
    sorter: (a, b) => (a.localWins ?? 0) - (b.localWins ?? 0),
  },
  {
    title: "Remote Wins",
    dataIndex: "remoteWins",
    align: "right",
    render: (value: number | null | undefined) => value ?? 0,
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
