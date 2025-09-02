import type { FindBuildFromUuidFragment } from "@/app/builds/[buildUUID]/[[...slugs]]/types";
import dayjs from "@/lib/dayjs";
import {
  readableDurationFromDates,
  readableDurationFromMilliseconds,
} from "@/utils/time";
import { theme } from "antd";
import Link from "next/link";
import { useMemo } from "react";
import {
  Bar,
  BarChart,
  CartesianGrid,
  Cell,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import type { InvocationInfo, TickProps } from "./types";
import { getBarColor } from "./utils";

interface Props {
  invocations: FindBuildFromUuidFragment[];
}

const BAR_HEIGHT = 20;
const CHART_PADDING = 40;

const InvocationTimeline: React.FC<Props> = ({ invocations }) => {
  const { token } = theme.useToken();

  const invocationsInfo: InvocationInfo[] = useMemo(
    () =>
      invocations.map((entry) => {
        return {
          invocationId: entry.invocationID,
          // Timestamp interval in milliseconds since UNIX epoch.
          timestamps: [
            dayjs(entry.startedAt).valueOf(),
            entry.endedAt ? dayjs(entry.endedAt).valueOf() : undefined,
          ],
          exitCode: entry.state.exitCode?.name,
        };
      }),
    [invocations],
  );

  // Place X-axis ticks at all defined timestamps.
  const ticks: number[] = useMemo(
    () =>
      invocationsInfo
        .flatMap((entry) => entry.timestamps)
        .filter((t): t is number => typeof t === "number"),
    [invocationsInfo],
  );
  const min = Math.min(...ticks);
  const max = Math.max(...ticks);

  const renderVerticalAxisTick = ({ x, y, payload }: TickProps) => {
    return (
      <g transform={`translate(${x},${y})`}>
        <Link href={`/bazel-invocations/${payload.value}`}>
          <text x={0} y={0} dy={8} textAnchor="end" fill={token.colorLink}>
            {`${payload.value.slice(0, 5)}...`}
          </text>
        </Link>
      </g>
    );
  };

  return (
    <ResponsiveContainer
      height={invocationsInfo.length * BAR_HEIGHT + CHART_PADDING}
      width="100%"
    >
      <BarChart layout="vertical" data={invocationsInfo}>
        <XAxis
          domain={[min, max]}
          type="number"
          ticks={ticks}
          tickFormatter={(value) => {
            return readableDurationFromDates(
              dayjs(min).toDate(),
              dayjs(value).toDate(),
              { precision: 1, smallestUnit: "s" },
            );
          }}
        />
        <YAxis
          dataKey="invocationId"
          type="category"
          tick={renderVerticalAxisTick}
          interval={0}
        />
        <CartesianGrid horizontal={false} syncWithTicks strokeDasharray="3 3" />
        <Tooltip
          // The label text turns white (against a white background, turning
          // it practically invisible) in dark mode unless color is set.
          labelStyle={{ color: "black" }}
          labelFormatter={(label, payload) => {
            const startedAt: number | undefined =
              payload[0]?.payload?.timestamps[0];
            const endedAt: number | undefined =
              payload[0]?.payload?.timestamps[1];
            return (
              <>
                <b>{label}</b>
                {startedAt && (
                  <p>Started at: {dayjs(startedAt).toISOString()}</p>
                )}
                {endedAt && <p>Ended at: {dayjs(endedAt).toISOString()}</p>}
              </>
            );
          }}
          formatter={(value: number[], name: string) => {
            // TODO: Handle in-progress invocations
            return [
              readableDurationFromMilliseconds(value[1] - value[0]),
              name,
            ];
          }}
        />
        <Bar
          dataKey="timestamps"
          name="Duration"
          minPointSize={5}
          barSize={BAR_HEIGHT}
        >
          {invocationsInfo.map((entry) => (
            <Cell key={entry.invocationId} fill={getBarColor(entry.exitCode)} />
          ))}
        </Bar>
      </BarChart>
    </ResponsiveContainer>
  );
};

export default InvocationTimeline;
