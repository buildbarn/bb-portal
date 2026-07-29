import { Collapse } from "antd";
import { useMemo } from "react";
import type { BazelInvocationActionsFragment } from "@/graphql/__generated__/graphql";
import themeStyles from "@/theme/theme.module.css";
import PortalDuration from "../PortalDuration";
import { ActionDetails } from "./action";

const getCollapseItems = (actions: BazelInvocationActionsFragment[]) => {
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
      children: <ActionDetails action={action} />,
    };
  });
};

interface Props {
  actions: BazelInvocationActionsFragment[];
}

export const ActionsTab: React.FC<Props> = ({ actions }) => {
  const items = useMemo(() => getCollapseItems(actions), [actions]);
  return (
    <Collapse
      items={items}
      bordered={true}
      defaultActiveKey={actions && actions.length === 1 ? [actions[0].id] : []}
      className={themeStyles.collapse}
    />
  );
};
