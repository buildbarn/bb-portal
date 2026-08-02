import { ToolOutlined } from "@ant-design/icons";
import { Space, Table, type TableColumnsType } from "antd";
import type {
  BazelInvocationMetricsWorkerMetricsFragment,
  BazelInvocationMetricsWorkerPoolMetricsFragment,
} from "@/graphql/__generated__/graphql";
import { readableFileSize } from "@/utils/filesize";
import PortalCard from "../PortalCard";

interface Props {
  workerMetrics: BazelInvocationMetricsWorkerMetricsFragment[];
  workerPoolMetrics?: BazelInvocationMetricsWorkerPoolMetricsFragment | null;
}

type WorkerRow = BazelInvocationMetricsWorkerMetricsFragment;
type WorkerPoolRow = NonNullable<
  BazelInvocationMetricsWorkerPoolMetricsFragment["workerPoolStats"]
>[number];

interface WorkerStatsRow {
  key: string;
  worker: string;
  collectTimeInMs: number;
  workerMemoryInKB: number;
  priorWorkerMemoryInKB: number;
  lastActionStartTimeInMs: number;
}

const displayTimestamp = (timestampInMs: number): string =>
  timestampInMs > 0 ? new Date(timestampInMs).toLocaleString() : "—";

const workerColumns: TableColumnsType<WorkerRow> = [
  {
    title: "Worker IDs",
    render: (_, worker) =>
      worker.workerIds?.map(({ workerID }) => workerID).join(", ") ?? "—",
  },
  { title: "Mnemonic", dataIndex: "mnemonic" },
  { title: "Process ID", dataIndex: "processID", align: "right" },
  { title: "Worker Key", dataIndex: "workerKeyHash", align: "right" },
  {
    title: "Status",
    render: (_, worker) =>
      [worker.workerStatus, worker.code].filter(Boolean).join(" / ") || "—",
  },
  {
    title: "Mode",
    render: (_, worker) => (worker.isMultiplex ? "Multiplex" : "Singleplex"),
  },
  {
    title: "Sandbox",
    render: (_, worker) => (worker.isSandbox ? "Yes" : "No"),
  },
  {
    title: "Measurable",
    render: (_, worker) => (worker.isMeasurable ? "Yes" : "No"),
  },
  { title: "Actions", dataIndex: "actionsExecuted", align: "right" },
  {
    title: "Prior Actions",
    dataIndex: "priorActionsExecuted",
    align: "right",
  },
];

const workerStatsColumns: TableColumnsType<WorkerStatsRow> = [
  { title: "Worker", dataIndex: "worker" },
  {
    title: "Collected",
    dataIndex: "collectTimeInMs",
    render: displayTimestamp,
  },
  {
    title: "Memory",
    dataIndex: "workerMemoryInKB",
    render: (value: number) => readableFileSize(value * 1024),
    align: "right",
  },
  {
    title: "Prior Memory",
    dataIndex: "priorWorkerMemoryInKB",
    render: (value: number) => readableFileSize(value * 1024),
    align: "right",
  },
  {
    title: "Last Action Start",
    dataIndex: "lastActionStartTimeInMs",
    render: displayTimestamp,
  },
];

const workerPoolColumns: TableColumnsType<WorkerPoolRow> = [
  { title: "Hash", dataIndex: "hash", align: "right" },
  { title: "Mnemonic", dataIndex: "mnemonic" },
  { title: "Created", dataIndex: "createdCount", align: "right" },
  { title: "Destroyed", dataIndex: "destroyedCount", align: "right" },
  { title: "Evicted", dataIndex: "evictedCount", align: "right" },
  { title: "Alive", dataIndex: "aliveCount", align: "right" },
  {
    title: "User Error",
    dataIndex: "userExecExceptionDestroyedCount",
    align: "right",
  },
  {
    title: "I/O Error",
    dataIndex: "ioExceptionDestroyedCount",
    align: "right",
  },
  {
    title: "Interrupted",
    dataIndex: "interruptedExceptionDestroyedCount",
    align: "right",
  },
  {
    title: "Unknown",
    dataIndex: "unknownDestroyedCount",
    align: "right",
  },
];

export const WorkerMetricsDisplay: React.FC<Props> = ({
  workerMetrics,
  workerPoolMetrics,
}) => {
  const workerStats: WorkerStatsRow[] = workerMetrics.flatMap((worker) => {
    const workerName =
      worker.workerIds?.map(({ workerID }) => workerID).join(", ") ||
      worker.mnemonic ||
      "Unknown";
    return (worker.workerStats ?? []).map((stat) => ({
      key: stat.id,
      worker: workerName,
      collectTimeInMs: stat.collectTimeInMs ?? 0,
      workerMemoryInKB: stat.workerMemoryInKB ?? 0,
      priorWorkerMemoryInKB: stat.priorWorkerMemoryInKB ?? 0,
      lastActionStartTimeInMs: stat.lastActionStartTimeInMs ?? 0,
    }));
  });
  const workerPoolStats = workerPoolMetrics?.workerPoolStats ?? [];

  return (
    <PortalCard
      type="inner"
      icon={<ToolOutlined />}
      titleBits={["Persistent Worker Metrics"]}
    >
      <Space direction="vertical" size="middle" style={{ width: "100%" }}>
        {workerMetrics.length > 0 && (
          <Table
            columns={workerColumns}
            dataSource={workerMetrics}
            pagination={false}
            rowKey="id"
            scroll={{ x: true }}
            size="small"
          />
        )}
        {workerStats.length > 0 && (
          <Table
            columns={workerStatsColumns}
            dataSource={workerStats}
            pagination={false}
            size="small"
          />
        )}
        {workerPoolStats.length > 0 && (
          <Table
            columns={workerPoolColumns}
            dataSource={workerPoolStats}
            pagination={false}
            rowKey="id"
            scroll={{ x: true }}
            size="small"
          />
        )}
      </Space>
    </PortalCard>
  );
};
