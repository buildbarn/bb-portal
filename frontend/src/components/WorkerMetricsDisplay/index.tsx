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
    key: "workerIds",
    render: (_, worker) =>
      worker.workerIds?.map(({ workerID }) => workerID).join(", ") ?? "—",
  },
  {
    title: "Mnemonic",
    key: "mnemonic",
    render: (_, worker) => worker.mnemonic,
  },
  {
    title: "Process ID",
    key: "processID",
    align: "right",
    render: (_, worker) => worker.processID,
  },
  {
    title: "Worker Key",
    key: "workerKeyHash",
    align: "right",
    render: (_, worker) => worker.workerKeyHash,
  },
  {
    title: "Status",
    key: "status",
    render: (_, worker) =>
      [worker.workerStatus, worker.code].filter(Boolean).join(" / ") || "—",
  },
  {
    title: "Mode",
    key: "mode",
    render: (_, worker) => (worker.isMultiplex ? "Multiplex" : "Singleplex"),
  },
  {
    title: "Sandbox",
    key: "sandbox",
    render: (_, worker) => (worker.isSandbox ? "Yes" : "No"),
  },
  {
    title: "Measurable",
    key: "measurable",
    render: (_, worker) => (worker.isMeasurable ? "Yes" : "No"),
  },
  {
    title: "Actions",
    key: "actionsExecuted",
    align: "right",
    render: (_, worker) => worker.actionsExecuted,
  },
  {
    title: "Prior Actions",
    key: "priorActionsExecuted",
    align: "right",
    render: (_, worker) => worker.priorActionsExecuted,
  },
];

const workerStatsColumns: TableColumnsType<WorkerStatsRow> = [
  {
    title: "Worker",
    key: "worker",
    render: (_, workerStat) => workerStat.worker,
  },
  {
    title: "Collected",
    key: "collectTimeInMs",
    render: (_, workerStat) => displayTimestamp(workerStat.collectTimeInMs),
  },
  {
    title: "Memory",
    key: "workerMemoryInKB",
    render: (_, workerStat) =>
      readableFileSize(workerStat.workerMemoryInKB * 1024),
    align: "right",
  },
  {
    title: "Prior Memory",
    key: "priorWorkerMemoryInKB",
    render: (_, workerStat) =>
      readableFileSize(workerStat.priorWorkerMemoryInKB * 1024),
    align: "right",
  },
  {
    title: "Last Action Start",
    key: "lastActionStartTimeInMs",
    render: (_, workerStat) =>
      displayTimestamp(workerStat.lastActionStartTimeInMs),
  },
];

const workerPoolColumns: TableColumnsType<WorkerPoolRow> = [
  {
    title: "Hash",
    key: "hash",
    align: "right",
    render: (_, workerPoolStat) => workerPoolStat.hash,
  },
  {
    title: "Mnemonic",
    key: "mnemonic",
    render: (_, workerPoolStat) => workerPoolStat.mnemonic,
  },
  {
    title: "Created",
    key: "createdCount",
    align: "right",
    render: (_, workerPoolStat) => workerPoolStat.createdCount,
  },
  {
    title: "Destroyed",
    key: "destroyedCount",
    align: "right",
    render: (_, workerPoolStat) => workerPoolStat.destroyedCount,
  },
  {
    title: "Evicted",
    key: "evictedCount",
    align: "right",
    render: (_, workerPoolStat) => workerPoolStat.evictedCount,
  },
  {
    title: "Alive",
    key: "aliveCount",
    align: "right",
    render: (_, workerPoolStat) => workerPoolStat.aliveCount,
  },
  {
    title: "User Error",
    key: "userExecExceptionDestroyedCount",
    align: "right",
    render: (_, workerPoolStat) =>
      workerPoolStat.userExecExceptionDestroyedCount,
  },
  {
    title: "I/O Error",
    key: "ioExceptionDestroyedCount",
    align: "right",
    render: (_, workerPoolStat) => workerPoolStat.ioExceptionDestroyedCount,
  },
  {
    title: "Interrupted",
    key: "interruptedExceptionDestroyedCount",
    align: "right",
    render: (_, workerPoolStat) =>
      workerPoolStat.interruptedExceptionDestroyedCount,
  },
  {
    title: "Unknown",
    key: "unknownDestroyedCount",
    align: "right",
    render: (_, workerPoolStat) => workerPoolStat.unknownDestroyedCount,
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
      <Space orientation="vertical" size="middle" style={{ width: "100%" }}>
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
