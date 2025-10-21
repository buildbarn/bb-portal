import { DeploymentUnitOutlined } from "@ant-design/icons";
import { Space } from "antd";
import type React from "react";
import PortalCard from "@/components/PortalCard";
import type { TargetMetrics } from "@/graphql/__generated__/graphql";
import { InvocationTargetsTable } from "../InvocationTargetsTable";

interface Props {
  invocationId: string;
  targetMetrics: TargetMetrics | undefined;
}

export const InvocationTargetsTab: React.FC<Props> = ({
  invocationId,
  targetMetrics,
}) => {
  return (
    <Space direction="vertical" size="middle" style={{ display: "flex" }}>
      <PortalCard
        type="inner"
        icon={<DeploymentUnitOutlined />}
        titleBits={["Targets"]}
      >
        <InvocationTargetsTable
          invocationId={invocationId}
          targetMetrics={targetMetrics}
        />
      </PortalCard>
    </Space>
  );
};
