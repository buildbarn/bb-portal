import { createFileRoute } from "@tanstack/react-router";
import { Divider, Typography } from "antd";
import { apolloClient } from "@/components/ApolloWrapper";
import { ExecutionRatioDonut } from "@/components/ExecutionRatioDonut";
import { InvocationOverviewDisplay } from "@/components/InvocationOverviewDisplay";
import { getFragmentData, gql } from "@/graphql/__generated__";
import type { RunnerCount } from "@/graphql/__generated__/graphql";
import { NotFoundError } from "@/main";
import { generatePageTitle } from "@/utils/generatePageTitle";

const GET_BAZEL_INVOCATION_OVERVIEW = gql(/* GraphQL */ `
  query GetBazelInvocationOverview($invocationID: UUID!) {
    getBazelInvocation(invocationID: $invocationID) {
      ...BazelInvocationOverview
      metrics {
        actionSummary {
          runnerCount {
            id
            name
            actionsExecuted
          }
        }
      }
    }
  }
`);

const BAZEL_INVOCATION_OVERVIEW_FRAGMENT = gql(/* GraphQL */ `
  fragment BazelInvocationOverview on BazelInvocation {
    id
    invocationID
    startedAt
    endedAt
    exitCodeName
    instanceName {
      id
      name
    }
    hostname
    connectionMetadata {
      id
      connectionLastOpenAt
      timeSinceLastConnectionMillis
    }
    originalCommandLine
    configurations {
      id
      cpu
      mnemonic
    }
    numFetches
    bazelVersion
  }
`);

export const Route = createFileRoute("/bazel-invocations/$invocationID/")({
  component: RouteComponent,
  loader: async ({ params }) => {
    const { data, error } = await apolloClient.query({
      errorPolicy: "all",
      query: GET_BAZEL_INVOCATION_OVERVIEW,
      variables: { invocationID: params.invocationID },
      fetchPolicy: "network-only",
    });

    if (!data?.getBazelInvocation) {
      throw new NotFoundError("invocation", error?.message);
    }

    const runnerCounts =
      (data.getBazelInvocation.metrics?.actionSummary?.runnerCount as
        | RunnerCount[]
        | null
        | undefined) ?? [];

    return {
      invocation: getFragmentData(
        BAZEL_INVOCATION_OVERVIEW_FRAGMENT,
        data.getBazelInvocation,
      ),
      runnerCounts,
    };
  },
  head: (_ctx) => ({
    meta: [
      {
        title: generatePageTitle([
          "Invocation",
          "Overview",
          _ctx.params.invocationID,
        ]),
      },
    ],
  }),
});

function RouteComponent() {
  const { invocation, runnerCounts } = Route.useLoaderData();
  return (
    <>
      <InvocationOverviewDisplay invocation={invocation} />
      {runnerCounts.length > 0 && (
        <>
          <Divider
            orientation="left"
            orientationMargin={0}
            style={{ marginTop: 24 }}
          >
            <Typography.Text type="secondary" style={{ fontSize: 13 }}>
              Execution Ratio
            </Typography.Text>
          </Divider>
          <ExecutionRatioDonut runnerCounts={runnerCounts} />
        </>
      )}
    </>
  );
}
