import { useQuery } from "@tanstack/react-query";
import { theme } from "antd";
import type React from "react";
import {
  ResponsiveContainer,
  Scatter,
  ScatterChart,
  type ScatterPointItem,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import type { FileDetailsFragment } from "@/graphql/__generated__/graphql";
import { readableDurationFromMilliseconds } from "@/utils/time";
import { type CriticalEvent, getCriticalPath } from "./utils";

interface Props {
  profile: FileDetailsFragment;
  hideTinyActions?: boolean;
}

export const CriticalPathDisplay: React.FC<Props> = ({
  profile,
  hideTinyActions,
}) => {
  const { token } = theme.useToken();
  const {
    data: criticalPath,
    isLoading,
    isError,
  } = useQuery({
    queryKey: ["criticalPath", profile.digest.hash],
    queryFn: () => getCriticalPath(profile),
  });
  if (!criticalPath || isLoading || isError || criticalPath.length === 0) {
    return null;
  }
  let chartItems: ChartItem[] = [];
  let totalDuration = 0;

  criticalPath.forEach((item) => {
    if (totalDuration < item.ts + item.dur) {
      totalDuration = item.ts + item.dur;
    }
    chartItems.push({ ...item, xCenter: item.ts + item.dur / 2, y: 1 });
  });
  if (hideTinyActions) {
    chartItems = chartItems.filter((item) => item.dur / totalDuration > 0.01);
  }

  return (
    <div
      style={{
        position: "relative",
        display: "block",
        width: "100%",
        height: 40,
        minWidth: 0,
        minHeight: 40,
      }}
    >
      <ResponsiveContainer width="100%" height={80}>
        <ScatterChart margin={{ top: 0, right: 0, left: 0, bottom: 20 }}>
          <XAxis
            type="number"
            dataKey="xCenter"
            domain={[0, totalDuration]}
            tickCount={10}
            stroke={token.colorTextSecondary}
            tickFormatter={(value) =>
              value === 0
                ? value
                : readableDurationFromMilliseconds(value / 1000)
            }
          />
          <YAxis type="number" dataKey="y" domain={[0, 2]} hide />

          <Tooltip
            content={<SegmentTooltip />}
            cursor={{ strokeDasharray: "3 3" }}
          />

          <Scatter data={chartItems} shape={<TimelineBar />} />
        </ScatterChart>
      </ResponsiveContainer>
    </div>
  );
};

interface ChartItem extends CriticalEvent {
  xCenter: number;
  y: number;
}

const CHAR_HEIGHT = 12;
const CHAR_WIDTH = CHAR_HEIGHT * 0.6;
const MIN_NAME_LENGTH = 20; // When we stop trunkating the name and instead hide it entierly, in characters
const TEXT_OVERFLOW = "[…]";

interface TimelineBarProps extends Partial<ScatterPointItem> {
  payload?: ChartItem;
}

const TimelineBar: React.FC<TimelineBarProps> = ({ cx, cy, payload }) => {
  const { token } = theme.useToken();
  if (!payload || cx === undefined || cy === undefined) return null;
  let width = 0;
  if (payload.xCenter > 0) {
    // Calculate how many pixels equal 1 unit of time
    const pixelsPerUnit = cx / payload.xCenter;

    // Scale duration into pixel width
    width = payload.dur * pixelsPerUnit;
  }

  const barWidth = Math.max(width, 2); // Ensure the bar is at least 2px wide
  const barHeight = 20;
  const yPos = cy - barHeight / 2;
  var displayName = payload.name;
  var textWidth = displayName.length * CHAR_WIDTH;
  if (
    textWidth > barWidth &&
    displayName.length > MIN_NAME_LENGTH &&
    MIN_NAME_LENGTH * CHAR_WIDTH < barWidth
  ) {
    const availableCharCount = barWidth / CHAR_WIDTH;
    displayName = `${displayName.substring(0, availableCharCount - TEXT_OVERFLOW.length)}${TEXT_OVERFLOW}`;
    textWidth = displayName.length * CHAR_WIDTH;
  }

  return (
    <g>
      <rect
        x={cx - width / 2}
        y={yPos}
        width={barWidth}
        height={barHeight}
        fill={token.colorPrimary}
        rx={4}
        ry={4}
      />

      {barWidth >= textWidth && (
        <text
          x={cx}
          y={cy}
          fill={token.colorTextLightSolid}
          textAnchor="middle"
          dominantBaseline="central"
          fontSize={CHAR_HEIGHT}
          style={{ pointerEvents: "none" }}
          textLength={textWidth}
          lengthAdjust="spacingAndGlyphs"
        >
          {payload.name}
        </text>
      )}
    </g>
  );
};

interface CustomTooltipProps {
  active?: boolean;
  payload?: Array<{ payload: CriticalEvent }>;
}

const SegmentTooltip: React.FC<CustomTooltipProps> = ({ active, payload }) => {
  const { token } = theme.useToken();
  if (active && payload && payload.length > 0 && payload[0].payload) {
    const data = payload[0].payload;
    return (
      <div
        style={{
          background: token.colorBgElevated,
          border: `1px solid ${token.colorBgTextActive}`,
          padding: "8px",
          borderRadius: "4px",
        }}
      >
        <p style={{ margin: "0 0 4px 0", fontWeight: "bold" }}>{data.name}</p>
        <p style={{ margin: 0, fontSize: "13px", color: token.colorTextLabel }}>
          Start: {readableDurationFromMilliseconds(data.ts / 1000)}
        </p>
        <p style={{ margin: 0, fontSize: "13px", color: token.colorTextLabel }}>
          Duration: {readableDurationFromMilliseconds(data.dur / 1000)}
        </p>
      </div>
    );
  }
  return null;
};
