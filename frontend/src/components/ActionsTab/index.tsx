import { Collapse } from "antd";
import { useMemo } from "react";
import type { BazelInvocationActionsFragment } from "@/graphql/__generated__/graphql";
import themeStyles from "@/theme/theme.module.css";
import PortalDuration from "../PortalDuration";
import { ActionDetails } from "./action";

const getCollapseItems = (
  instanceName: string,
  actions: BazelInvocationActionsFragment[],
) => {
  return actions?.map((action) => {
    return {
      key: action.id,
      label: action.label,
      extra: action.startTime && action.endTime && (
        <PortalDuration
          from={action.startTime || undefined}
          to={action.endTime || undefined}
          formatConfig={{ smallestUnit: "ms" }}
        />
      ),
      children: <ActionDetails instanceName={instanceName} action={action} />,
    };
  });
};

interface Props {
  instanceName: string;
  actions: BazelInvocationActionsFragment[];
}

export const ActionsTab: React.FC<Props> = ({ instanceName, actions }) => {
  const items = useMemo(
    () => getCollapseItems(instanceName, actions),
    [instanceName, actions],
  );
  return (
    <Collapse
      items={items}
      bordered={true}
      defaultActiveKey={actions && actions.length === 1 ? [actions[0].id] : []}
      className={themeStyles.collapse}
    />
  );
};
