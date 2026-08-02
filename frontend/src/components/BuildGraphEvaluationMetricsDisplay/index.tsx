import { ApartmentOutlined } from "@ant-design/icons";
import { Table, type TableColumnsType } from "antd";
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

  return (
    <PortalCard
      type="inner"
      icon={<ApartmentOutlined />}
      titleBits={["Skyframe Evaluation Metrics"]}
    >
      <Table
        columns={columns}
        dataSource={rows}
        pagination={{ defaultPageSize: 20, hideOnSinglePage: true }}
        size="small"
      />
    </PortalCard>
  );
};
