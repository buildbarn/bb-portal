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
  field: CountField,
): TableColumnsType<EvaluationRow>[number] => ({
  title,
  key: field,
  align: "right",
  render: (_, record) => record[field],
  sorter: (a, b) => a[field] - b[field],
});

const evaluationColumns: TableColumnsType<EvaluationRow> = [
  {
    title: "Skyfunction",
    key: "skyfunctionName",
    render: (_, record) => record.skyfunctionName,
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
  field: "count" | "actionCount",
): TableColumnsType<T>[number] => ({
  title,
  key: field,
  align: "right",
  render: (_, record) => record[field] ?? 0,
  sorter: (a, b) => Number(a[field] ?? 0) - Number(b[field] ?? 0),
});

const ruleClassColumns: TableColumnsType<RuleClassRow> = [
  {
    title: "Key",
    key: "key",
    render: (_, record) => record.key,
    sorter: (a, b) => (a.key ?? "").localeCompare(b.key ?? ""),
  },
  {
    title: "Rule Class",
    key: "ruleClass",
    render: (_, record) => record.ruleClass,
    sorter: (a, b) => (a.ruleClass ?? "").localeCompare(b.ruleClass ?? ""),
  },
  numericCountColumn<RuleClassRow>("Configured Targets", "count"),
  numericCountColumn<RuleClassRow>("Actions", "actionCount"),
];

const aspectColumns: TableColumnsType<AspectRow> = [
  {
    title: "Key",
    key: "key",
    render: (_, record) => record.key,
    sorter: (a, b) => (a.key ?? "").localeCompare(b.key ?? ""),
  },
  {
    title: "Aspect",
    key: "aspectName",
    render: (_, record) => record.aspectName,
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

  const evaluationRows = [...rowsBySkyfunction.values()].sort(
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
      <Space orientation="vertical" size="middle" style={{ width: "100%" }}>
        {evaluationRows.length > 0 && (
          <div>
            <Typography.Title level={5}>Skyframe Evaluation</Typography.Title>
            <Table
              columns={evaluationColumns}
              dataSource={evaluationRows}
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
