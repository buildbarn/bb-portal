import type { ActionData } from "@/graphql/__generated__/graphql";
import type { Maybe } from "graphql/jsutils/Maybe";
import ActionsPieChart, { type ActionsChartItem } from "../ActionsPieChart";
import { chartColor } from "../ActionsPieChart/utils";
import { nullPercent } from "../Utilities/nullPercent";

interface Props {
  actionData?: Maybe<ActionData[]>;
}

const ActionTypeMetrics: React.FC<Props> = ({ actionData }) => {
  const actions: ActionsChartItem[] = [];
  const totalActionsExecuted = actionData?.reduce(
    (accumulator, item) => accumulator + (item.actionsExecuted ?? 0),
    0,
  );

  if (actionData) {
    actionData.forEach((item: ActionData, index: number) => {
      const chartItem: ActionsChartItem = {
        key: index,
        value: item.mnemonic ?? "",
        percent: nullPercent(item.actionsExecuted, totalActionsExecuted, 0),
        count: item.actionsExecuted ?? 0,
        color: chartColor(index),
        type: "square",
      };
      actions.push(chartItem);
    });
  }

  return <ActionsPieChart items={actions} />;
};

export default ActionTypeMetrics;
