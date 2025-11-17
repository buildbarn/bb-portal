import type { Maybe } from "graphql/jsutils/Maybe";
import type { ActionData } from "@/graphql/__generated__/graphql";
import SummaryPieChart, { type SummaryChartItem } from "../SummaryPieChart";
import { nullPercent } from "../Utilities/nullPercent";

interface Props {
  actionData?: Maybe<ActionData[]>;
}

const ActionTypeMetrics: React.FC<Props> = ({ actionData }) => {
  const actions: SummaryChartItem[] = [];
  const totalActionsExecuted = actionData?.reduce(
    (accumulator, item) => accumulator + (item.actionsExecuted ?? 0),
    0,
  );

  if (actionData) {
    actionData.forEach((item: ActionData, index: number) => {
      const chartItem: SummaryChartItem = {
        key: index,
        value: item.mnemonic ?? "",
        percent: nullPercent(item.actionsExecuted, totalActionsExecuted, 0),
        count: item.actionsExecuted ?? 0,
        type: "square",
      };
      actions.push(chartItem);
    });
  }

  return <SummaryPieChart items={actions} />;
};

export default ActionTypeMetrics;
