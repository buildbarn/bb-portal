import { BuildOutlined } from "@ant-design/icons";
import { useQuery } from "@apollo/client/react";
import { Spin, Typography } from "antd";
import type React from "react";
import { useEffect } from "react";
import BazelInvocation from "@/components/BazelInvocation";
import Content from "@/components/Content";
import PortalAlert from "@/components/PortalAlert";
import PortalCard from "@/components/PortalCard";
import { getFragmentData } from "@/graphql/__generated__";
import type { LoadFullBazelInvocationDetailsQuery } from "@/graphql/__generated__/graphql";
import { shouldPollInvocation } from "@/utils/shouldPollInvocation";
import {
  BAZEL_INVOCATION_FRAGMENT,
  LOAD_FULL_BAZEL_INVOCATION_DETAILS,
} from "./index.graphql";

interface Params {
  invocationID: string;
}

export const BazelInvocationDetailsPage: React.FC<Params> = (params) => {
  return <Content content={<PageContent {...params} />} />;
};

const PageContent: React.FC<Params> = ({ invocationID }) => {
  const { data, error, loading, stopPolling } =
    useQuery<LoadFullBazelInvocationDetailsQuery>(
      LOAD_FULL_BAZEL_INVOCATION_DETAILS,
      {
        variables: { invocationID: invocationID },
        fetchPolicy: "cache-and-network",
        // nextFetchPolicy prevents unnecessary refetches if the logs are fetched
        nextFetchPolicy: "cache-and-network",
        pollInterval: 5000,
      },
    );
  const invocation = getFragmentData(
    BAZEL_INVOCATION_FRAGMENT,
    data?.getBazelInvocation,
  );

  // Poll continuously until the invocation is completed. Then we should stop.
  useEffect(() => {
    if (!shouldPollInvocation(invocation)) {
      stopPolling();
    }
  }, [invocation, stopPolling]);

  if (loading) {
    return (
      <PortalCard icon={<BuildOutlined />} titleBits={["Bazel Invocation"]}>
        <Spin />
      </PortalCard>
    );
  }
  if (error || !invocation) {
    return (
      <PortalCard icon={<BuildOutlined />} titleBits={["Bazel Invocation"]}>
        <PortalAlert
          showIcon
          type="error"
          message="Error fetching invocation details"
          description={
            <>
              {error?.message ? (
                <Typography.Paragraph>{error.message}</Typography.Paragraph>
              ) : (
                <Typography.Paragraph>
                  Unknown error occurred while fetching data from the server.
                </Typography.Paragraph>
              )}
              <Typography.Paragraph>
                It could be that the invocation is too old and has been removed,
                or that you don&quot;t have access to this invocation.
              </Typography.Paragraph>
            </>
          }
        />
      </PortalCard>
    );
  }

  return <BazelInvocation invocationOverview={invocation} />;
};
