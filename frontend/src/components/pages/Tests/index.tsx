import { ExperimentFilled } from "@ant-design/icons";
import { Alert, Space } from "antd";
import type React from "react";
import PortalCard from "@/components/PortalCard";
import TestGrid from "@/components/TestGrid";

type Props = React.ComponentProps<typeof TestGrid>;

export const TestsPage: React.FC<Props> = (props) => {
  return (
    <Space direction="vertical" size="middle" style={{ display: "flex" }}>
      <PortalCard
        icon={<ExperimentFilled />}
        extraBits={[
          <Alert
            key="search-by-label"
            showIcon
            message="Search by label and/or instance name to further refine your result"
            type="info"
          />,
        ]}
        titleBits={[<span key="title">Tests Overview</span>]}
      >
        <TestGrid {...props} />
      </PortalCard>
    </Space>
  );
};
