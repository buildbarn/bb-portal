import { FolderOpenOutlined } from "@ant-design/icons";
import { Table, type TableColumnsType } from "antd";
import type { BazelInvocationMetricsPackageMetricsFragment } from "@/graphql/__generated__/graphql";
import { readableDurationFromMilliseconds } from "@/utils/time";
import { PortalCard } from "../PortalCard";

interface Props {
  packageMetrics: BazelInvocationMetricsPackageMetricsFragment;
}

type PackageLoadRow = NonNullable<
  BazelInvocationMetricsPackageMetricsFragment["packageLoadMetrics"]
>[number];

type NumericPackageLoadField =
  | "numTargets"
  | "computationSteps"
  | "numTransitiveLoads"
  | "packageOverhead"
  | "globFilesystemOperationCost";

const numericColumn = (
  title: string,
  field: NumericPackageLoadField,
): TableColumnsType<PackageLoadRow>[number] => ({
  title,
  key: field,
  align: "right",
  render: (_, record) => record[field] ?? 0,
  sorter: (a, b) => Number(a[field] ?? 0) - Number(b[field] ?? 0),
});

const columns: TableColumnsType<PackageLoadRow> = [
  {
    title: "Package",
    key: "name",
    render: (_, record) => record.name || "—",
    sorter: (a, b) => (a.name ?? "").localeCompare(b.name ?? ""),
  },
  {
    title: "Load Duration",
    key: "loadDurationInNs",
    align: "right",
    render: (_, record) =>
      readableDurationFromMilliseconds(
        (record.loadDurationInNs ?? 0) / 1_000_000,
      ),
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
