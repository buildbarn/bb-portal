import { ApartmentOutlined } from "@ant-design/icons";
import { Space, Table, type TableColumnsType, Typography } from "antd";
import type { BazelInvocationMetricsBuildGraphEvaluationMetricsFragment } from "@/graphql/__generated__/graphql";
import PortalCard from "../PortalCard";

interface Props {
  buildGraphMetrics: BazelInvocationMetricsBuildGraphEvaluationMetricsFragment;
}

interface EvaluationRow {
  key: string;
  skyfunctionName: string;
  dirtied: number;
  changed: number;
  built: number;
  cleaned: number;
  evaluated: number;
}

type CountField = Exclude<keyof EvaluationRow, "key" | "skyfunctionName">;

type RuleClassRow = NonNullable<
  BazelInvocationMetricsBuildGraphEvaluationMetricsFragment["ruleClassCounts"]
>[number];

type AspectRow = NonNullable<
  BazelInvocationMetricsBuildGraphEvaluationMetricsFragment["aspectCounts"]
>[number];

const OPERATION_FIELDS: Record<string, CountField> = {
  DIRTIED: "dirtied",
  CHANGED: "changed",
  BUILT: "built",
  CLEANED: "cleaned",
  EVALUATED: "evaluated",
};

const countColumn = (
  title: string,
  dataIndex: CountField,
): TableColumnsType<EvaluationRow>[number] => ({
  title,
  dataIndex,
  align: "right",
  sorter: (a, b) => a[dataIndex] - b[dataIndex],
});

const columns: TableColumnsType<EvaluationRow> = [
  {
    title: "Skyfunction",
    dataIndex: "skyfunctionName",
    sorter: (a, b) => a.skyfunctionName.localeCompare(b.skyfunctionName),
  },
  countColumn("Dirtied", "dirtied"),
  countColumn("Changed", "changed"),
  countColumn("Built", "built"),
  countColumn("Cleaned", "cleaned"),
  countColumn("Evaluated", "evaluated"),
];

const numericCountColumn = <
  T extends { count?: number | null; actionCount?: number | null },
>(
  title: string,
  dataIndex: "count" | "actionCount",
): TableColumnsType<T>[number] => ({
  title,
  dataIndex,
  align: "right",
  render: (value: number | null | undefined) => value ?? 0,
  sorter: (a, b) => Number(a[dataIndex] ?? 0) - Number(b[dataIndex] ?? 0),
});

const ruleClassColumns: TableColumnsType<RuleClassRow> = [
  {
    title: "Key",
    dataIndex: "key",
    sorter: (a, b) => (a.key ?? "").localeCompare(b.key ?? ""),
  },
  {
    title: "Rule Class",
    dataIndex: "ruleClass",
    sorter: (a, b) => (a.ruleClass ?? "").localeCompare(b.ruleClass ?? ""),
  },
  numericCountColumn<RuleClassRow>("Configured Targets", "count"),
  numericCountColumn<RuleClassRow>("Actions", "actionCount"),
];

const aspectColumns: TableColumnsType<AspectRow> = [
  {
    title: "Key",
    dataIndex: "key",
    sorter: (a, b) => (a.key ?? "").localeCompare(b.key ?? ""),
  },
  {
    title: "Aspect",
    dataIndex: "aspectName",
    sorter: (a, b) => (a.aspectName ?? "").localeCompare(b.aspectName ?? ""),
  },
  numericCountColumn<AspectRow>("Configured Targets", "count"),
  numericCountColumn<AspectRow>("Actions", "actionCount"),
];

export const BuildGraphEvaluationMetricsDisplay: React.FC<Props> = ({
  buildGraphMetrics,
}) => {
  const rowsBySkyfunction = new Map<string, EvaluationRow>();

  for (const stat of buildGraphMetrics.evaluationStats ?? []) {
    const skyfunctionName = stat.skyfunctionName ?? "Unknown";
    const row = rowsBySkyfunction.get(skyfunctionName) ?? {
      key: skyfunctionName,
      skyfunctionName,
      dirtied: 0,
      changed: 0,
      built: 0,
      cleaned: 0,
      evaluated: 0,
    };
    const countField = OPERATION_FIELDS[stat.operation ?? ""];
    if (countField) {
      row[countField] += stat.count ?? 0;
    }
    rowsBySkyfunction.set(skyfunctionName, row);
  }

  const rows = [...rowsBySkyfunction.values()].sort(
    (a, b) => b.evaluated - a.evaluated,
  );

  const ruleClassCounts = buildGraphMetrics.ruleClassCounts ?? [];
  const aspectCounts = buildGraphMetrics.aspectCounts ?? [];

  return (
    <PortalCard
      type="inner"
      icon={<ApartmentOutlined />}
      titleBits={["Build Graph Metrics"]}
    >
      <Space direction="vertical" size="middle" style={{ width: "100%" }}>
        {rows.length > 0 && (
          <div>
            <Typography.Title level={5}>Skyframe Evaluation</Typography.Title>
            <Table
              columns={columns}
              dataSource={rows}
              pagination={{ defaultPageSize: 20, hideOnSinglePage: true }}
              size="small"
            />
          </div>
        )}
        {ruleClassCounts.length > 0 && (
          <div>
            <Typography.Title level={5}>Rule Classes</Typography.Title>
            <Table
              columns={ruleClassColumns}
              dataSource={ruleClassCounts}
              pagination={{ defaultPageSize: 20, hideOnSinglePage: true }}
              rowKey="id"
              scroll={{ x: true }}
              size="small"
            />
          </div>
        )}
        {aspectCounts.length > 0 && (
          <div>
            <Typography.Title level={5}>Aspects</Typography.Title>
            <Table
              columns={aspectColumns}
              dataSource={aspectCounts}
              pagination={{ defaultPageSize: 20, hideOnSinglePage: true }}
              rowKey="id"
              scroll={{ x: true }}
              size="small"
            />
          </div>
        )}
      </Space>
    </PortalCard>
  );
};
