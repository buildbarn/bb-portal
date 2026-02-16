"use client";

import React from "react";
import PortalCard from "@/components/PortalCard";
import { Spin, Typography } from "antd";
import {
  DeploymentUnitOutlined
} from "@ant-design/icons";
import { FIND_BUILD_BY_UUID_QUERY } from "./index.graphql";
import { useQuery } from "@apollo/client";
import PortalAlert from "@/components/PortalAlert";
import { isFeatureEnabled, FeatureType } from '@/utils/isFeatureEnabled';
import PageDisabled from "@/components/PageDisabled";
import { BuildDetails } from "@/components/BuildDetails";
import Content from "@/components/Content";

interface PageParams {
  params: {
    buildUUID: string;
  }
}

const Page: React.FC<PageParams> = ({ params }) => {
  if (!isFeatureEnabled(FeatureType.BES) || !isFeatureEnabled(FeatureType.BES_PAGE_BUILDS)) {
    return <PageDisabled />;
  }
  return <Content content={
    <PageContent params={params} />
  } />;
}

const PageContent: React.FC<PageParams> = ({ params }) => {
  const { data, loading, error } = useQuery(FIND_BUILD_BY_UUID_QUERY, {
    variables: { uuid: params.buildUUID },
  });

  if (loading) {
    return (
      <PortalCard
        icon={<DeploymentUnitOutlined />}
        titleBits={["Build"]}
      >
        <Spin />
      </PortalCard>
    )
  }

  if (error || !data) {
    return (
      <PortalCard
        icon={<DeploymentUnitOutlined />}
        titleBits={["Build"]}
      >
        <PortalAlert
          showIcon
          type="error"
          message="Error fetching build details"
          description={
            <>
              {error?.message ? <Typography.Paragraph>{error.message}</Typography.Paragraph> :
                <Typography.Paragraph>Unknown error occurred while fetching data from the server.</Typography.Paragraph>}
              <Typography.Paragraph>It could be that the build is too old and has been removed, or that you don&quot;t have access to this build.</Typography.Paragraph>
            </>
          }
        />
      </PortalCard>
    )
  }

  return <BuildDetails data={data} />
}

export default Page;