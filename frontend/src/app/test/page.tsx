"use client";

import { ExperimentFilled } from "@ant-design/icons";
import { Space } from "antd";
import { useSearchParams } from "next/navigation";
import type React from "react";
import Content from "@/components/Content";
import PageDisabled from "@/components/PageDisabled";
import PortalCard from "@/components/PortalCard";
import TestDetails from "@/components/TestDetails";
import { FeatureType, isFeatureEnabled } from "@/utils/isFeatureEnabled";

const Page: React.FC = () => {
  const searchParams = useSearchParams();
  if (
    !isFeatureEnabled(FeatureType.BES) ||
    !isFeatureEnabled(FeatureType.BES_PAGE_TESTS)
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
            icon={<ExperimentFilled />}
            titleBits={[<span key="title">Test Details</span>]}
          >
            <TestDetails
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
