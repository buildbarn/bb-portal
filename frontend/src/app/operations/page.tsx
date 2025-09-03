"use client";

import Content from "@/components/Content";
import OperationsGrid from "@/components/OperationsGrid";
import PageDisabled from "@/components/PageDisabled";
import PortalCard from "@/components/PortalCard";
import { FeatureType, isFeatureEnabled } from "@/utils/isFeatureEnabled";
import { CodeFilled } from "@ant-design/icons";
import type React from "react";

const Page: React.FC = () => {
  if (!isFeatureEnabled(FeatureType.OPERATIONS)) {
    return <PageDisabled />;
  }

  return (
    <Content
      content={
        <PortalCard
          icon={<CodeFilled />}
          titleBits={[<span key="title">Operations</span>]}
        >
          <OperationsGrid />
        </PortalCard>
      }
    />
  );
};

export default Page;
