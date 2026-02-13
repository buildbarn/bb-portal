"use client";

import type React from "react";
import { BuildDetails } from "@/components/BuildDetails";
import Content from "@/components/Content";
import PageDisabled from "@/components/PageDisabled";
import { FeatureType, isFeatureEnabled } from "@/utils/isFeatureEnabled";

interface PageParams {
  params: {
    buildUUID: string;
  };
}

const Page: React.FC<PageParams> = ({ params }) => {
  if (
    !isFeatureEnabled(FeatureType.BES) ||
    !isFeatureEnabled(FeatureType.BES_PAGE_BUILDS)
  ) {
    return <PageDisabled />;
  }
  return <Content content={<BuildDetails buildUUID={params.buildUUID} />} />;
};

export default Page;
