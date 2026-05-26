import { DeploymentUnitOutlined } from "@ant-design/icons";
import { Alert, Space } from "antd";
import type React from "react";
import PortalCard from "@/components/PortalCard";
import TargetGrid from "@/components/Targets/TargetGrid";

type Props = React.ComponentProps<typeof TargetGrid>;

export const TargetsPage: React.FC<Props> = (props) => {
  return (
    <Space direction="vertical" size="middle" style={{ display: "flex" }}>
      <PortalCard
        icon={<DeploymentUnitOutlined />}
        extraBits={[
          <Alert
            key="search-by-label"
            showIcon
            message="Search by label to further refine your result"
            type="info"
          />,
        ]}
        titleBits={[<span key="title">Targets Overview </span>]}
      >
        <TargetGrid {...props} />
      </PortalCard>
    </Space>
  );
};
