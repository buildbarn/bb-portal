import type { PreviousExecutionStats } from "@/lib/grpc-client/buildbarn/iscc/iscc";
import { formatDurationFromSeconds } from "@/utils/formatValues";
import {
  Legend,
  ReferenceArea,
  Scatter,
  ScatterChart,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import { durationToSeconds } from "../Utilities/time";

interface Props {
  prevStats: PreviousExecutionStats;
}

interface PlotDataPoint {
  x: number;
  y: number;
  sizeClass: number;
}

const PADDING_FACTOR = 4;
const REFERENCE_AREA_WIDTH = 0.4;

const PreviousExecutionsPlot: React.FC<Props> = ({ prevStats }) => {
  const succeeded: PlotDataPoint[] = [];
  const timedOut: PlotDataPoint[] = [];
  const sizeClasses: number[] = [];
  const entries = Object.entries(prevStats.sizeClasses);

  for (let i = 0; i < entries.length; ++i) {
    const sizeClass = Number.parseInt(entries[i][0]);
    sizeClasses.push(sizeClass);
    for (const prevExec of entries[i][1].previousExecutions) {
      // TODO: Make random scatter deterministic for each data point
      const xValue = i + (Math.random() - 0.5) / 3;
      if (prevExec.succeeded) {
        const time = durationToSeconds(prevExec.succeeded);
        succeeded.push({
          x: xValue,
          y: time,
          sizeClass: sizeClass,
        });
      }
      if (prevExec.timedOut) {
        const time = durationToSeconds(prevExec.timedOut);
        timedOut.push({
          x: xValue,
          y: time,
          sizeClass: sizeClass,
        });
        // `prevExec.failed` has no time information,
        // so we cannot visualize them in the graph
      }
    }
  }

  return (
    <ScatterChart
      width={750}
      height={500}
      margin={{
        top: 10,
        bottom: 10,
        left: 20,
        right: 20,
      }}
    >
      <XAxis
        dataKey="x"
        type="number"
        name="Size class"
        label={{
          value: "Size class",
          position: "insideBottom",
          offset: 0,
        }}
        ticks={Array.from(Array(entries.length).keys())}
        tickFormatter={(_, index) => {
          return entries[index][0];
        }}
        domain={() => {
          const len = entries.length;
          return [-PADDING_FACTOR / len, len - 1 + PADDING_FACTOR / len];
        }}
      />

      <YAxis
        dataKey="y"
        type="number"
        name="Execution time"
        tickLine={false}
        label={{
          value: "Execution time (s)",
          angle: -90,
          position: "insideLeft",
        }}
      />
      <Tooltip
        cursor={{ stroke: "3 3" }}
        formatter={(value, name, props) => {
          switch (name) {
            case "Size class": {
              return [props.payload.sizeClass, name];
            }
            case "Execution time": {
              return [formatDurationFromSeconds(props.payload.y, 10), name];
            }
            default: {
              return [value, name];
            }
          }
        }}
      />
      {entries.map((value, index) => {
        return (
          <ReferenceArea
            key={value[0]}
            x1={index - REFERENCE_AREA_WIDTH}
            x2={index + REFERENCE_AREA_WIDTH}
            y1={0}
            fill="gray"
            ifOverflow="extendDomain"
          />
        );
      })}
      {succeeded.length > 0 && (
        <Scatter shape="cross" name="Succeeded" data={succeeded} fill="green" />
      )}
      {timedOut.length > 0 && (
        <Scatter shape="cross" name="Timed out" data={timedOut} fill="orange" />
      )}
      <Legend verticalAlign="top" align="right" />
    </ScatterChart>
  );
};

export default PreviousExecutionsPlot;
