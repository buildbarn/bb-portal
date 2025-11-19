"use client";

import { CodeFilled } from "@ant-design/icons";
import type React from "react";
import Content from "@/components/Content";
import OperationDetails from "@/components/OperationDetails";
import PageDisabled from "@/components/PageDisabled";
import PortalCard from "@/components/PortalCard";
import { FeatureType, isFeatureEnabled } from "@/utils/isFeatureEnabled";

interface PageParams {
  params: {
    slug: string;
  };
}

const Page: React.FC<PageParams> = ({ params }) => {
  if (!isFeatureEnabled(FeatureType.OPERATIONS)) {
    return <PageDisabled />;
  }

  const operationID = decodeURIComponent(params.slug);

  return (
    <Content
      content={
        <PortalCard
          icon={<CodeFilled />}
          titleBits={[<span key="title">{`Operation ${operationID}`}</span>]}
        >
          <OperationDetails operationID={operationID} />
        </PortalCard>
      }
    />
  );
};

export default Page;
