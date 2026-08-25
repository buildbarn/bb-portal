import type { Maybe } from "graphql/jsutils/Maybe";
import type { ActionData } from "@/graphql/__generated__/graphql";
import SummaryPieChart, { type SummaryChartItem } from "../SummaryPieChart";
import { nullPercent } from "../Utilities/nullPercent";

interface Props {
  actionData?: Maybe<ActionData[]>;
  countField: "actionsCreated" | "actionsExecuted";
}

const ActionTypeMetrics: React.FC<Props> = ({ actionData, countField }) => {
  const actions: SummaryChartItem[] = [];
  const totalActions = actionData?.reduce(
    (accumulator, item) => accumulator + (item[countField] ?? 0),
    0,
  );

  if (actionData) {
    actionData.forEach((item: ActionData, index: number) => {
      const chartItem: SummaryChartItem = {
        key: index,
        value: item.mnemonic ?? "",
        percent: nullPercent(item[countField], totalActions, 0),
        count: item[countField] ?? 0,
      };
      actions.push(chartItem);
    });
  }

  return <SummaryPieChart items={actions} />;
};

export default ActionTypeMetrics;
