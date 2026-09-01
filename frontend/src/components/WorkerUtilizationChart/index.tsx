import { useQuery } from "@tanstack/react-query";
import { Select, Skeleton, Space, Typography } from "antd";
import type React from "react";
import { useMemo, useState } from "react";
import {
  Area,
  AreaChart,
  CartesianGrid,
  Legend,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import { buildQueueStateClient } from "@/grpc/buildQueueStateClient";
import { env } from "@/utils/env";
import PortalAlert from "../PortalAlert";

const FETCHING_COLOR = "#FA8C16";
const RUNNING_COLOR = "#49AA19";
const UPLOADING_COLOR = "#722ED1";
const IDLE_COLOR = "#8C8C8C";

const RANGE_OPTIONS = [
  { label: "Last 15 min", value: 900, step: 15, window: "1m" },
  { label: "Last 1 hour", value: 3600, step: 60, window: "5m" },
  { label: "Last 6 hours", value: 21600, step: 300, window: "15m" },
];

interface DataPoint {
  time: number;
  Fetching: number;
  Executing: number;
  Uploading: number;
  Idle: number;
}

interface PrometheusMatrix {
  metric: { stage?: string };
  values: [number, string][];
}

function safeFloat(s: string): number {
  const v = parseFloat(s);
  return Number.isFinite(v) ? v : 0;
}

async function fetchPhaseRates(
  rangeSeconds: number,
  step: number,
  window: string,
): Promise<PrometheusMatrix[]> {
  const end = Math.floor(Date.now() / 1000);
  const start = end - rangeSeconds;
  // Rate window scales with the step so short-lived spikes remain visible.
  const query = `sum by (stage) (rate(buildbarn_builder_build_executor_duration_seconds_sum[${window}]))`;
  // In production the Go backend proxies /api/v1/prometheus/* → Prometheus.
  // In dev the Vite config proxies /api/v1/prometheus/* → localhost:9090.
  const url = `/api/v1/prometheus/api/v1/query_range?query=${encodeURIComponent(query)}&start=${start}&end=${end}&step=${step}`;
  const res = await fetch(url);
  if (!res.ok) throw new Error(`Prometheus error: ${res.status} ${res.statusText}`);
  const json = await res.json();
  if (json.status !== "success") throw new Error(json.error ?? "Prometheus query failed");
  return json.data.result as PrometheusMatrix[];
}

const WorkerUtilizationChart: React.FC = () => {
  const [rangeIndex, setRangeIndex] = useState(1);
  const { value: rangeSeconds, step, window } = RANGE_OPTIONS[rangeIndex];

  const prometheusEnabled = Boolean(env.prometheusUrl);

  const { data: queueData } = useQuery({
    queryKey: ["listPlatformQueuesForChart"],
    queryFn: async () => buildQueueStateClient.listPlatformQueues({}),
    enabled: prometheusEnabled,
    refetchInterval: 30_000,
  });

  const totalWorkers = useMemo(
    () =>
      // biome-ignore lint/suspicious/noExplicitAny: queueData may be PlatformQueueTableState[] from cache
      ((queueData as any)?.platformQueues ?? []).reduce(
        (sum: number, q: any) =>
          sum + (q.sizeClassQueues ?? []).reduce((s: number, scq: any) => s + (scq.workersCount ?? 0), 0),
        0,
      ),
    [queueData],
  );

  const {
    data: phaseData,
    isLoading,
    isError,
    error,
  } = useQuery({
    queryKey: ["prometheusWorkerPhases", rangeSeconds, step],
    queryFn: () => fetchPhaseRates(rangeSeconds, step, window),
    refetchInterval: 60_000,
    enabled: prometheusEnabled,
  });

  // Compute chart data only when both datasets change.
  const chartData = useMemo<DataPoint[]>(() => {
    if (!phaseData) return [];
    const byStage: Record<string, Map<number, number>> = {};
    for (const series of phaseData) {
      const stage = series.metric.stage ?? "unknown";
      const map = new Map<number, number>();
      for (const [ts, val] of series.values) map.set(ts, safeFloat(val));
      byStage[stage] = map;
    }
    const allTs = new Set<number>();
    for (const m of Object.values(byStage)) m.forEach((_, ts) => allTs.add(ts));
    return [...allTs].sort((a, b) => a - b).map((ts) => {
      const fetching = byStage["FetchingInputs"]?.get(ts) ?? 0;
      const executing = byStage["Running"]?.get(ts) ?? 0;
      const uploading = byStage["UploadingOutputs"]?.get(ts) ?? 0;
      const active = fetching + executing + uploading;
      // Clamp idle to [0, totalWorkers] — metric rates can transiently exceed the
      // actual worker count during bursts or counter resets.
      const idle = totalWorkers > 0 ? Math.max(0, Math.min(totalWorkers, totalWorkers - active)) : 0;
      return {
        time: ts * 1000,
        Fetching: parseFloat(fetching.toFixed(2)),
        Executing: parseFloat(executing.toFixed(2)),
        Uploading: parseFloat(uploading.toFixed(2)),
        Idle: parseFloat(idle.toFixed(2)),
      };
    });
  }, [phaseData, totalWorkers]);

  if (!prometheusEnabled) {
    return (
      <Typography.Text type="secondary">
        Worker utilization chart requires Prometheus. Set{" "}
        <code>prometheus_url</code> in the bb-portal configuration.
      </Typography.Text>
    );
  }

  if (isLoading) {
    return <Skeleton active paragraph={{ rows: 4 }} />;
  }

  if (isError) {
    return (
      <PortalAlert
        showIcon
        type="error"
        message="Error fetching worker phase metrics"
        description={
          (error as Error).message ||
          "Could not reach Prometheus. Ensure the cluster is running and port-forwarded."
        }
      />
    );
  }

  const fmt = (ts: number) => {
    const d = new Date(ts);
    return rangeSeconds <= 900
      ? d.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit", second: "2-digit" })
      : d.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
  };

  return (
    <Space direction="vertical" style={{ width: "100%" }}>
      <Space>
        <Typography.Text type="secondary" style={{ fontSize: 12 }}>
          Time range:
        </Typography.Text>
        <Select
          size="small"
          value={rangeIndex}
          onChange={setRangeIndex}
          options={RANGE_OPTIONS.map((o, i) => ({ label: o.label, value: i }))}
          style={{ width: 130 }}
        />
        {totalWorkers > 0 && (
          <Typography.Text type="secondary" style={{ fontSize: 12 }}>
            {totalWorkers} workers total
          </Typography.Text>
        )}
      </Space>
      {chartData.length === 0 ? (
        <Typography.Text type="secondary">
          No data — Prometheus may have no metrics yet. Run a build first.
        </Typography.Text>
      ) : (
        <ResponsiveContainer width="100%" height={240}>
          <AreaChart data={chartData} margin={{ top: 4, right: 16, left: 0, bottom: 4 }}>
            <CartesianGrid strokeDasharray="3 3" vertical={false} />
            <XAxis
              dataKey="time"
              type="number"
              scale="time"
              domain={["dataMin", "dataMax"]}
              tickFormatter={fmt}
              tick={{ fontSize: 10 }}
              tickCount={6}
            />
            <YAxis
              tick={{ fontSize: 10 }}
              label={{ value: "workers", angle: -90, position: "insideLeft", fontSize: 10, offset: 10 }}
              width={52}
            />
            <Tooltip
              labelFormatter={(v) => fmt(v as number)}
              formatter={(v: number, name: string) => [
                Number.isFinite(v) ? `${v.toFixed(2)} workers` : "—",
                name,
              ]}
            />
            <Legend />
            <Area type="monotone" dataKey="Idle" stackId="w" fill={IDLE_COLOR} stroke={IDLE_COLOR} fillOpacity={0.6} isAnimationActive={false} />
            <Area type="monotone" dataKey="Fetching" stackId="w" fill={FETCHING_COLOR} stroke={FETCHING_COLOR} fillOpacity={0.8} isAnimationActive={false} />
            <Area type="monotone" dataKey="Executing" stackId="w" fill={RUNNING_COLOR} stroke={RUNNING_COLOR} fillOpacity={0.8} isAnimationActive={false} />
            <Area type="monotone" dataKey="Uploading" stackId="w" fill={UPLOADING_COLOR} stroke={UPLOADING_COLOR} fillOpacity={0.8} isAnimationActive={false} />
          </AreaChart>
        </ResponsiveContainer>
      )}
    </Space>
  );
};

export default WorkerUtilizationChart;
