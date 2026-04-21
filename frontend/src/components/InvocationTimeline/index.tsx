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
import { env } from "@/utils/env";
import { parseGraphqlEdgeList } from "@/utils/parseGraphqlEdgeList";
import {
  readableDurationFromDates,
  readableDurationFromMilliseconds,
} from "@/utils/time";
import CommandLinePreview from "../CommandLinePreview";
import { INVOCATION_RESULT_TAGS } from "../InvocationResultTag";
import {
  getInvocationResultTagEnum,
  InvocationResult,
} from "../InvocationResultTag/enum";
import PortalAlert from "../PortalAlert";
import type { InvocationInfo, TickProps } from "./types";

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
          const invocationStatus = getInvocationResultTagEnum(
            entry.exitCodeName || undefined,
            entry.connectionMetadata?.timeSinceLastConnectionMillis,
          );
          let endTime = entry.endedAt;
          if (!endTime && invocationStatus !== InvocationResult.IN_PROGRESS) {
            endTime = entry.connectionMetadata?.connectionLastOpenAt;
          }
          if (!endTime) {
            endTime = new Date();
          }
          return {
            invocationId: entry.invocationID,
            // Timestamp interval in milliseconds since UNIX epoch.
            timestamps: [
              dayjs(entry.startedAt).valueOf(),
              dayjs(endTime).valueOf(),
            ],
            invocationStatus,
            command: entry.originalCommandLine,
            tags: parseGraphqlEdgeList(entry.tags),
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
            const columns = env.additionalBuildInvocationColumns;
            const invocationEntry = payload[0]?.payload as
              | InvocationInfo
              | undefined;
            return (
              // The labels are wrapped in a span with `display: block` to
              // simulate a div for text formatting purposes. Using divs
              // directly would cause hydration errors as the label
              // formatter wraps the elements below in a <p> tag.
              <>
                <b>Invocation ID:</b> <code>{label}</code>
                {invocationEntry &&
                  columns.map((column) => (
                    <span key={column.valueKey} style={{ display: "block" }}>
                      <b>{column.title}:</b>{" "}
                      <code>
                        {invocationEntry.tags.find(
                          (tag) => tag.key === column.valueKey,
                        )?.value || "-"}
                      </code>
                    </span>
                  ))}
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
              fill={INVOCATION_RESULT_TAGS[entry.invocationStatus].color}
            />
          ))}
        </Bar>
      </BarChart>
    </ResponsiveContainer>
  );
};

export default InvocationTimeline;
