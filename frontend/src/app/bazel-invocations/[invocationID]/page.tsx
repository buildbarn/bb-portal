'use client';

import React from 'react';
import Content from "@/components/Content";
import { useQuery } from "@apollo/client";
import {
  LoadFullBazelInvocationDetailsQuery
} from '@/graphql/__generated__/graphql';
import {
  BAZEL_INVOCATION_FRAGMENT,
  LOAD_FULL_BAZEL_INVOCATION_DETAILS
} from "@/app/bazel-invocations/[invocationID]/index.graphql";
import { getFragmentData } from "@/graphql/__generated__";
import { Spin, Typography } from "antd";
import BazelInvocation from "@/components/BazelInvocation";
import { isFeatureEnabled, FeatureType } from '@/utils/isFeatureEnabled';
import PageDisabled from '@/components/PageDisabled';
import PortalCard from '@/components/PortalCard';
import { BuildOutlined } from '@ant-design/icons';
import PortalAlert from '@/components/PortalAlert';

interface PageParams {
  params: {
    invocationID: string
  }
}

const Page: React.FC<PageParams> = ({ params }) => {
  if (!isFeatureEnabled(FeatureType.BES) || !isFeatureEnabled(FeatureType.BES_PAGE_INVOCATIONS)) {
    return <PageDisabled />;
  }

  return (
    <Content
      content={
        <PageContent params={params} />
      }
    />
  );
}

const PageContent: React.FC<PageParams> = ({ params }) => {
  const { data, error, loading } = useQuery<LoadFullBazelInvocationDetailsQuery>(
    LOAD_FULL_BAZEL_INVOCATION_DETAILS,
    {
      variables: { invocationID: params.invocationID },
      fetchPolicy: 'cache-and-network',
      // nextFetchPolicy prevents unnecessary refetches if the logs are fetched
      nextFetchPolicy: 'cache-only',
      notifyOnNetworkStatusChange: true,
    },
  );
  const invocation = getFragmentData(BAZEL_INVOCATION_FRAGMENT, data?.bazelInvocation);

  if (loading) {
    return (
      <PortalCard
        icon={<BuildOutlined />}
        titleBits={["Bazel Invocation"]}
      >
        <Spin />
      </PortalCard>
    )
  }

  if (error || !invocation) {
    return (
      <PortalCard
        icon={<BuildOutlined />}
        titleBits={["Bazel Invocation"]}
      >
        <PortalAlert
          showIcon
          type="error"
          message="Error fetching invocation details"
          description={
            <>
              {error?.message ? <Typography.Paragraph>{error.message}</Typography.Paragraph> :
                  <Typography.Paragraph>Unknown error occurred while fetching data from the server.</Typography.Paragraph>}
              <Typography.Paragraph>It could be that the invocation is too old and has been removed, or that you don&quot;t have access to this invocation.</Typography.Paragraph>
            </>
          }
        />
      </PortalCard>
    )
  }

  return <BazelInvocation invocationOverview={invocation} />
}

export default Page;
