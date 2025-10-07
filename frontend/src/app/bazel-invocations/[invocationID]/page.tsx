'use client';

import React, { useEffect } from 'react';
import Content from "@/components/Content";
import { ApolloError, NetworkStatus, useQuery } from "@apollo/client";
import {
  BazelInvocationInfoFragment,
  FullBazelInvocationDetailsFragment,
  LoadFullBazelInvocationDetailsQuery, ProblemInfoFragment
} from '@/graphql/__generated__/graphql';
import {
  BAZEL_INVOCATION_FRAGMENT,
  FULL_BAZEL_INVOCATION_DETAILS,
  LOAD_FULL_BAZEL_INVOCATION_DETAILS
} from "@/app/bazel-invocations/[invocationID]/index.graphql";
import { getFragmentData } from "@/graphql/__generated__";
import { Spin } from "antd";
import ErrorAlert from "@/components/ErrorAlert";
import BazelInvocation from "@/components/BazelInvocation";
import { isFeatureEnabled, FeatureType } from '@/utils/isFeatureEnabled';
import PageDisabled from '@/components/PageDisabled';

interface PageParams {
  params: {
    invocationID: string
  }
}

interface Props {
  loading: boolean
  error: ApolloError | undefined
  networkStatus: NetworkStatus
  invocationInfo: BazelInvocationInfoFragment
  problems: ProblemInfoFragment[]
}

const BazelInvocationsContent: React.FC<Props> = ({ loading, error, networkStatus, invocationInfo, problems }) => {
  if (loading && networkStatus !== NetworkStatus.poll && invocationInfo) {
    return (
      <Spin>
        <BazelInvocation invocationOverview={invocationInfo} />
      </Spin>
    );
  }
  if (loading && networkStatus !== NetworkStatus.poll) {
    return (
      <Spin />
    );
  }
  if (error && invocationInfo) {
    return (
      <>
        <BazelInvocation invocationOverview={invocationInfo} />
        <ErrorAlert error={new Error("REEEE")} />
      </>
    );
  }

  if (invocationInfo) {
    return <BazelInvocation invocationOverview={invocationInfo} />
  }

  return <></>
}

const shouldStopPolling = (invocation: FullBazelInvocationDetailsFragment | null | undefined): boolean => {
  return !!invocation;
}

const Page: React.FC<PageParams> = ({ params }) => {
  if (!isFeatureEnabled(FeatureType.BES) || !isFeatureEnabled(FeatureType.BES_PAGE_INVOCATIONS)) {
    return <PageDisabled />;
  }
  return <PageContent params={params}/>
}

const PageContent: React.FC<PageParams> = ({ params }) => {
  const { data, error, loading, stopPolling, networkStatus } = useQuery<LoadFullBazelInvocationDetailsQuery>(
    LOAD_FULL_BAZEL_INVOCATION_DETAILS,
    {
      variables: { invocationID: params.invocationID },
      fetchPolicy: 'cache-and-network',
      pollInterval: 5000,
      notifyOnNetworkStatusChange: true,
    },
  );

  const invocation = getFragmentData(FULL_BAZEL_INVOCATION_DETAILS, data?.bazelInvocation);
  const invocationOverview = getFragmentData(BAZEL_INVOCATION_FRAGMENT, invocation)


  const stop = shouldStopPolling(invocation);
  useEffect(() => {
    if (stop) {
      stopPolling();
    }
  }, [stop, stopPolling]);



  return (
    <Content
      content={<BazelInvocationsContent loading={loading} error={error} networkStatus={networkStatus} invocationInfo={invocationOverview as BazelInvocationInfoFragment} problems={[]} />}
    />
  );
}

export default Page;
