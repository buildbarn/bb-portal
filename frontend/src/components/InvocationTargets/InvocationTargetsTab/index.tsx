import { DeploymentUnitOutlined } from "@ant-design/icons";
import type React from "react";
import PortalCard from "@/components/PortalCard";
import { InvocationTargetsTable } from "../InvocationTargetsTable";

type Props = React.ComponentProps<typeof InvocationTargetsTable>;

export const InvocationTargetsTab: React.FC<Props> = (params: Props) => {
  return (
    <PortalCard
      type="inner"
      icon={<DeploymentUnitOutlined />}
      titleBits={["Targets"]}
    >
      <InvocationTargetsTable {...params} />
    </PortalCard>
  );
};
