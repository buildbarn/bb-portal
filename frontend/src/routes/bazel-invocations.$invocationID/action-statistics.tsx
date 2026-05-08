import { createFileRoute, notFound } from "@tanstack/react-router";
import ActionStatisticsDisplay from "@/components/ActionStatisticsDisplay";
import { apolloClient } from "@/components/ApolloWrapper";
import { getFragmentData, gql } from "@/graphql/__generated__";
import { generatePageTitle } from "@/utils/generatePageTitle";

const GET_BAZEL_INVOCATION_ACTION_STATISTICS = gql(/* GraphQL */ `
  query GetBazelInvocationActionStatistics($invocationID: UUID!) {
    getBazelInvocation(invocationID: $invocationID) {
      id
      metrics {
        id
        actionSummary {
          ...BazelInvocationActionSummary
          }
      }
    }
  }
`);

const BAZEL_INVOCATION_ACTION_SUMMARY_FRAGMENT = gql(/* GraphQL */ `
  fragment BazelInvocationActionSummary on ActionSummary {
    id
    actionsExecuted
    actionCacheStatistics {
      id
      loadTimeInMs
      saveTimeInMs
      hits
      misses
      sizeInBytes
      missDetails {
        id
        count
        reason
      }
    }
    runnerCount {
      id
      actionsExecuted
      name
      execKind
    }
    actionData {
      id
      mnemonic
      userTime
      systemTime
      actionsExecuted
    }
  }
`);

export const Route = createFileRoute(
  "/bazel-invocations/$invocationID/action-statistics",
)({
  component: RouteComponent,
  loader: async ({ params }) => {
    const { data } = await apolloClient.query({
      query: GET_BAZEL_INVOCATION_ACTION_STATISTICS,
      variables: { invocationID: params.invocationID },
      fetchPolicy: "network-only",
    });

    if (!data?.getBazelInvocation?.metrics?.actionSummary) {
      throw notFound();
    }

    const actionSummary = getFragmentData(
      BAZEL_INVOCATION_ACTION_SUMMARY_FRAGMENT,
      data.getBazelInvocation?.metrics?.actionSummary,
    );
    return { actionSummary };
  },
  head: (_ctx) => ({
    meta: [
      {
        title: generatePageTitle([
          "Invocation",
          "Action statistics",
          _ctx.params.invocationID,
        ]),
      },
    ],
  }),
});

function RouteComponent() {
  const { actionSummary } = Route.useLoaderData();
  return <ActionStatisticsDisplay actionSummary={actionSummary} />;
}
