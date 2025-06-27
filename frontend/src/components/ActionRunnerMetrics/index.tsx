import type { RunnerCount } from "@/graphql/__generated__/graphql";
import ActionsPieChart, { type ActionsChartItem } from "../ActionsPieChart";
import { nullPercent } from "../Utilities/nullPercent";

interface Props {
  runnerMetrics: RunnerCount[];
}

function colorSwitchOnExecStrat(exec: string) {
  switch (exec) {
    case "Remote":
      return "#49AA19";
    case "Local":
      return "#DC4446";
    default:
      return "#777777";
  }
}

const ActionRunnerMetrics: React.FC<Props> = ({ runnerMetrics }) => {
  const chartItems: ActionsChartItem[] = [];
  const totalCount: number =
    runnerMetrics.find((i) => i.name === "total")?.actionsExecuted ?? 0;

  runnerMetrics.forEach((item: RunnerCount, index: number) => {
    const chartItem: ActionsChartItem = {
      key: index,
      value: item.name ?? "",
      count: item.actionsExecuted ?? 0,
      percent: nullPercent(item.actionsExecuted, totalCount, 0),
      color: colorSwitchOnExecStrat(item.execKind ?? ""),
      type: "square",
    };
    if (chartItem.value !== "total") {
      chartItems.push(chartItem);
    }
  });

  return <ActionsPieChart items={chartItems} />;
};

export default ActionRunnerMetrics;
