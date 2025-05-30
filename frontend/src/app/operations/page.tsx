"use client";

import Content from "@/components/Content";
import OperationsGrid from "@/components/OperationsGrid";
import PortalCard from "@/components/PortalCard";
import { FeatureType, isFeatureEnabled } from "@/utils/isFeatureEnabled";
import { CodeFilled } from "@ant-design/icons";
import { notFound } from "next/navigation";
import type React from "react";

const Page: React.FC = () => {
  if (!isFeatureEnabled(FeatureType.SCHEDULER)) {
    return notFound();
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
