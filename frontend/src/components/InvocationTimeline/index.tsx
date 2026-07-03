import { useQuery } from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";
import { theme } from "antd";
import { useMemo, useState } from "react";
import {
  Bar,
  BarChart,
  type BarShapeProps,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import { getFragmentData } from "@/graphql/__generated__";
import type { GetBuildInvocationFragment } from "@/graphql/__generated__/graphql";
import dayjs from "@/lib/dayjs";
import { FILE_DETAILS_FRAGMENT } from "@/types/GraphqlFileFragment";
import { env } from "@/utils/env";
import { parseGraphqlEdgeList } from "@/utils/parseGraphqlEdgeList";
import {
  readableDurationFromDates,
  readableDurationFromMilliseconds,
} from "@/utils/time";
import CommandLinePreview from "../CommandLinePreview";
import { getCriticalPath } from "../CriticalPath/utils";
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
  const [hoveredCriticalAction, setHoveredCriticalAction] = useState<
    string | null
  >(null);

  const { data: criticalPaths = [] } = useQuery({
    queryKey: ["criticalPaths", invocations?.map((inv) => inv.id)],
    queryFn: async () => {
      const promises = invocations.map(async (invocation) => {
        if (!invocation.profile) {
          return undefined;
        }
        return getCriticalPath(
          getFragmentData(FILE_DETAILS_FRAGMENT, invocation.profile),
        );
      });
      return await Promise.all(promises);
    },
    staleTime: Infinity,
  });

  const invocationsInfo: InvocationInfo[] = useMemo(
    () =>
      [...invocations]
        .filter((entry) => !!entry.startedAt)
        .map((entry, index) => {
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
            criticalTraceEvents: criticalPaths[index],
          };
        })
        .sort((a, b) => {
          for (const columnName of env.additionalBuildInvocationColumns) {
            const aValue =
              a.tags.find((tag) => tag.key === columnName.valueKey)?.value ??
              "";
            const bValue =
              b.tags.find((tag) => tag.key === columnName.valueKey)?.value ??
              "";
            if (aValue.localeCompare(bValue) !== 0) {
              return aValue.localeCompare(bValue);
            }
          }
          return a.timestamps[0] - b.timestamps[0];
        }),
    [invocations, criticalPaths],
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
        title="The provided invocations list was empty"
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
                {hoveredCriticalAction && (
                  <span style={{ display: "block" }}>
                    <b>Critical path: </b> <code>{hoveredCriticalAction}</code>
                  </span>
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
          shape={(props: BarShapeProps) => (
            <CustomBar {...props} onHoverAction={setHoveredCriticalAction} />
          )}
        />
      </BarChart>
    </ResponsiveContainer>
  );
};

interface BarProps extends BarShapeProps {
  payload?: InvocationInfo;
  onHoverAction?: (action: string | null) => void;
}

const CHAR_HEIGHT = 12;
const CHAR_WIDTH = CHAR_HEIGHT * 0.6;
const MIN_NAME_LENGTH = 20; // When we stop trunkating the name and instead hide it entierly, in characters
const TEXT_OVERFLOW = "[…]";

const CustomBar: React.FC<BarProps> = ({
  x,
  y,
  width,
  height,
  payload,
  onHoverAction,
}) => {
  const { token } = theme.useToken();
  if (!x || !y || !width || !height || !payload) {
    return null;
  }
  const durationInMs = payload.timestamps[1] - payload.timestamps[0];
  const xPixelsFromMicroSeconds = (time: number) => {
    return (time / (durationInMs * 1000)) * width;
  };
  const margin = 2;

  return (
    <g key={payload.invocationId}>
      {/* Full event bar */}
      <rect
        x={x}
        y={y}
        width={width}
        height={height}
        fill={
          INVOCATION_RESULT_TAGS[
            payload?.invocationStatus || InvocationResult.UNKNOWN_EXIT_CODE
          ].color
        }
      />
      {/* Critical path components */}
      {payload.criticalTraceEvents?.map((event) => {
        const eventWidth = xPixelsFromMicroSeconds(event.dur);
        var displayName = event.name;
        var textWidth = displayName.length * CHAR_WIDTH;
        if (
          textWidth > eventWidth &&
          displayName.length > MIN_NAME_LENGTH &&
          MIN_NAME_LENGTH * CHAR_WIDTH < eventWidth
        ) {
          const availableCharCount = eventWidth / CHAR_WIDTH;
          displayName = `${displayName.substring(0, availableCharCount - TEXT_OVERFLOW.length)}${TEXT_OVERFLOW}`;
          textWidth = displayName.length * CHAR_WIDTH;
        }
        return (
          <>
            {/* biome-ignore lint/a11y/noInteractiveElementToNoninteractiveRole: There's no good non-static tags or role usable inside an SVG to represent a region that alters the displayed tooltip. */}
            <g
              key={event.name}
              role="graphics-symbol"
              aria-label={`Critical path action: ${event.name}`}
              tabIndex={0}
              onMouseEnter={() => onHoverAction?.(event.name)}
              onMouseLeave={() => onHoverAction?.(null)}
            >
              <rect
                x={x + xPixelsFromMicroSeconds(event.ts)}
                y={y + margin}
                width={eventWidth}
                height={height - margin * 2}
                fill={token.colorPrimaryBg}
              />
              {eventWidth > textWidth && (
                <text
                  x={x + xPixelsFromMicroSeconds(event.ts) + eventWidth / 2}
                  y={y + height * 0.5}
                  dy=".35em"
                  textAnchor="middle"
                  fontSize={CHAR_HEIGHT}
                  fontFamily="monospace"
                  fill={token.colorText}
                  style={{ pointerEvents: "none" }} // Ensure text doesn't interfere with mouse events
                  textLength={textWidth}
                  lengthAdjust="spacingAndGlyphs"
                >
                  {displayName}
                </text>
              )}
            </g>
          </>
        );
      })}
    </g>
  );
};

export default InvocationTimeline;
