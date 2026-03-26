import { useNavigate } from "@tanstack/react-router";
import { theme } from "antd";
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
import type { GetBuildInvocationFragment } from "@/graphql/__generated__/graphql";
import dayjs from "@/lib/dayjs";
import {
  readableDurationFromDates,
  readableDurationFromMilliseconds,
} from "@/utils/time";
import CommandLinePreview from "../CommandLinePreview";
import PortalAlert from "../PortalAlert";
import type { InvocationInfo, TickProps } from "./types";
import { getInvocationResultTagColor } from "./utils";

interface Props {
  invocations: GetBuildInvocationFragment[];
}

const BAR_HEIGHT = 20;
const CHART_PADDING = 40;

const InvocationTimeline: React.FC<Props> = ({ invocations }) => {
  const navigate = useNavigate();
  const { token } = theme.useToken();

  const invocationsInfo: InvocationInfo[] = useMemo(
    () =>
      invocations
        .filter((entry) => !!entry.startedAt)
        .map((entry) => {
          const startTime = entry.startedAt;
          let endTime = entry.endedAt;
          if (!endTime) {
            endTime = entry.connectionMetadata?.connectionLastOpenAt;
          }
          if (!endTime) {
            endTime = new Date();
          }
          return {
            invocationId: entry.invocationID,
            // Timestamp interval in milliseconds since UNIX epoch.
            timestamps: [dayjs(startTime).valueOf(), dayjs(endTime).valueOf()],
            exitCodeName: entry.exitCodeName || undefined,
            timeSinceLastConnectionMillis:
              entry.connectionMetadata?.timeSinceLastConnectionMillis ||
              undefined,
            command: entry.originalCommandLine,
            job: entry.sourceControl?.job,
            workflow: entry.sourceControl?.workflow,
            action: entry.sourceControl?.action,
          };
        }),
    [invocations],
  );

  // Place X-axis ticks at all defined timestamps.
  const ticks: number[] = useMemo(
    () =>
      invocationsInfo
        .flatMap((entry) => entry.timestamps)
        // Provide sort function because otherwise JS converts the numbers to strings and sorts
        .sort((a, b) => a - b)
        .filter(
          // Remove duplicates, which cause issues with the rendering.
          (timestamp, index, array) => !index || timestamp !== array[index - 1],
        ),
    [invocationsInfo],
  );

  if (invocationsInfo.length < 1)
    return (
      <PortalAlert
        showIcon
        type="warning"
        message="The provided invocations list was empty"
      />
    );

  const min = ticks[0];
  const max = ticks[ticks.length - 1];

  const renderVerticalAxisTick = ({ x, y, payload }: TickProps) => {
    return (
      <g transform={`translate(${x},${y})`}>
        <text x={0} y={0} dy={8} textAnchor="end" fill={token.colorText}>
          {`${payload.value.slice(0, 5)}...`}
        </text>
      </g>
    );
  };

  return (
    <ResponsiveContainer
      height={invocationsInfo.length * BAR_HEIGHT + CHART_PADDING}
      width="100%"
    >
      <BarChart
        layout="vertical"
        data={invocationsInfo}
        onClick={(state) => {
          if (state.activeLabel !== undefined) {
            navigate({
              to: "/bazel-invocations/$invocationID",
              params: { invocationID: `${state.activeLabel}` },
            });
          }
        }}
      >
        <XAxis
          domain={[min, max]}
          type="number"
          interval={"preserveStartEnd"}
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
          contentStyle={{
            backgroundColor: token.colorBgContainer,
            borderColor: token.colorBgTextActive,
          }}
          wrapperStyle={{ maxWidth: "50vw", zIndex: 999 }}
          labelFormatter={(label, payload) => {
            const invocationEntry = payload[0]?.payload;
            return (
              // The labels are wrapped in a span with `display: block` to
              // simulate a div for text formatting purposes. Using divs
              // directly would cause hydration errors as the label
              // formatter wraps the elements below in a <p> tag.
              <>
                <b>Invocation ID:</b> <code>{label}</code>
                {invocationEntry?.workflow && (
                  <span style={{ display: "block" }}>
                    <b>Workflow: </b> <code>{invocationEntry?.workflow}</code>
                  </span>
                )}
                {invocationEntry?.job && (
                  <span style={{ display: "block" }}>
                    <b>Job: </b> <code>{invocationEntry?.job}</code>
                  </span>
                )}
                {invocationEntry?.action && (
                  <span style={{ display: "block" }}>
                    <b>Action: </b> <code>{invocationEntry?.action}</code>
                  </span>
                )}
                {invocationEntry?.timestamps[0] && (
                  <span style={{ display: "block" }}>
                    <b>Duration: </b>
                    {readableDurationFromMilliseconds(
                      payload[0].payload?.timestamps[1] -
                        payload[0].payload?.timestamps[0],
                    )}
                  </span>
                )}
                {invocationEntry?.command && (
                  <>
                    <b>Bazel command: </b>{" "}
                    <CommandLinePreview
                      codeBlockWrapper
                      command={invocationEntry.command}
                    />
                  </>
                )}
              </>
            );
          }}
          formatter={() => []}
        />
        <Bar
          dataKey="timestamps"
          name="Duration"
          minPointSize={5}
          barSize={BAR_HEIGHT}
        >
          {invocationsInfo.map((entry) => (
            <Cell
              key={entry.invocationId}
              fill={getInvocationResultTagColor(
                entry.exitCodeName,
                entry.timeSinceLastConnectionMillis,
              )}
            />
          ))}
        </Bar>
      </BarChart>
    </ResponsiveContainer>
  );
};

export default InvocationTimeline;
