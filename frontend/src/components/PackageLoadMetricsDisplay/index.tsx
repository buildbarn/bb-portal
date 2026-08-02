import { FolderOpenOutlined } from "@ant-design/icons";
import { Table, type TableColumnsType } from "antd";
import type { BazelInvocationMetricsPackageMetricsFragment } from "@/graphql/__generated__/graphql";
import { readableDurationFromMilliseconds } from "@/utils/time";
import PortalCard from "../PortalCard";

interface Props {
  packageMetrics: BazelInvocationMetricsPackageMetricsFragment;
}

type PackageLoadRow = NonNullable<
  BazelInvocationMetricsPackageMetricsFragment["packageLoadMetrics"]
>[number];

const numericColumn = (
  title: string,
  dataIndex: keyof PackageLoadRow,
): TableColumnsType<PackageLoadRow>[number] => ({
  title,
  dataIndex,
  align: "right",
  render: (value: number | null | undefined) => value ?? 0,
  sorter: (a, b) => Number(a[dataIndex] ?? 0) - Number(b[dataIndex] ?? 0),
});

const columns: TableColumnsType<PackageLoadRow> = [
  {
    title: "Package",
    dataIndex: "name",
    render: (value: string | null | undefined) => value || "—",
    sorter: (a, b) => (a.name ?? "").localeCompare(b.name ?? ""),
  },
  {
    title: "Load Duration",
    dataIndex: "loadDurationInNs",
    align: "right",
    render: (value: number | null | undefined) =>
      readableDurationFromMilliseconds((value ?? 0) / 1_000_000),
    sorter: (a, b) => (a.loadDurationInNs ?? 0) - (b.loadDurationInNs ?? 0),
  },
  numericColumn("Targets", "numTargets"),
  numericColumn("Computation Steps", "computationSteps"),
  numericColumn("Transitive Loads", "numTransitiveLoads"),
  numericColumn("Package Overhead", "packageOverhead"),
  numericColumn("Glob Filesystem Cost", "globFilesystemOperationCost"),
];

export const PackageLoadMetricsDisplay: React.FC<Props> = ({
  packageMetrics,
}) => (
  <PortalCard
    type="inner"
    icon={<FolderOpenOutlined />}
    titleBits={["Package Load Metrics"]}
  >
    <Table
      columns={columns}
      dataSource={packageMetrics.packageLoadMetrics ?? []}
      pagination={{ defaultPageSize: 20, hideOnSinglePage: true }}
      rowKey="id"
      scroll={{ x: true }}
      size="small"
    />
  </PortalCard>
);
