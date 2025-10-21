"use client";

import { DeploymentUnitOutlined } from "@ant-design/icons";
import { Space } from "antd";
import { useSearchParams } from "next/navigation";
import type React from "react";
import Content from "@/components/Content";
import PageDisabled from "@/components/PageDisabled";
import PortalCard from "@/components/PortalCard";
import { TargetDetails } from "@/components/Targets/TargetDetails";
import { FeatureType, isFeatureEnabled } from "@/utils/isFeatureEnabled";

const Page: React.FC = () => {
  const searchParams = useSearchParams();
  if (
    !isFeatureEnabled(FeatureType.BES) ||
    !isFeatureEnabled(FeatureType.BES_PAGE_TARGETS)
  ) {
    return <PageDisabled />;
  }

  const instanceName = decodeURIComponent(
    decodeURIComponent(searchParams.get("instanceName") || ""),
  );
  const label = decodeURIComponent(
    decodeURIComponent(searchParams.get("label") || ""),
  );
  const aspect = decodeURIComponent(
    decodeURIComponent(searchParams.get("aspect") || ""),
  );
  const targetKind = decodeURIComponent(
    decodeURIComponent(searchParams.get("targetKind") || ""),
  );

  return (
    <Content
      content={
        <Space direction="vertical" size="middle" style={{ display: "flex" }}>
          <PortalCard
            icon={<DeploymentUnitOutlined />}
            titleBits={[<span key="title">Target Details</span>]}
          >
            <TargetDetails
              instanceName={instanceName}
              label={label}
              aspect={aspect}
              targetKind={targetKind}
            />
          </PortalCard>
        </Space>
      }
    />
  );
};

export default Page;
