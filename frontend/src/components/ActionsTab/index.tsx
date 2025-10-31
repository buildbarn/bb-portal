import { Collapse } from "antd";
import themeStyles from "@/theme/theme.module.css";
import PortalDuration from "../PortalDuration";
import { ActionDetails, type ActionDetailsData } from "./action";

const items = (
  instanceName: string,
  actions: ActionDetailsData[] | undefined | null,
) => {
  return actions?.map((action) => {
    return {
      key: action.id,
      label: action.label,
      extra: action.startTime && action.endTime && (
        <PortalDuration
          from={action.startTime}
          to={action.endTime}
          includePopover
          formatConfig={{ smallestUnit: "ms" }}
        />
      ),
      children: <ActionDetails instanceName={instanceName} action={action} />,
    };
  });
};

interface Props {
  instanceName: string;
  actions: ActionDetailsData[] | undefined | null;
}

export const ActionsTab: React.FC<Props> = ({ instanceName, actions }) => {
  return (
    <Collapse
      items={items(instanceName, actions)}
      bordered={true}
      defaultActiveKey={actions && actions.length === 1 ? [actions[0].id] : []}
      destroyInactivePanel
      className={themeStyles.collapse}
    />
  );
};
