import { useQuery } from "@apollo/client/react";
import { createFileRoute } from "@tanstack/react-router";
import { apolloClient } from "@/components/ApolloWrapper";
import BuildLogsDisplay from "@/components/BuildLogsDisplay";
import { getFragmentData, gql } from "@/graphql/__generated__";
import { NotFoundError } from "@/main";
import { commandLineDataToString } from "@/utils/commandLineDataToString";
import { generatePageTitle } from "@/utils/generatePageTitle";
import { shouldPollInvocation } from "@/utils/shouldPollInvocation";

const INVOCATION_POLL_INTERVAL_MS = 5000;

const GET_BAZEL_INVOCATION_LOG = gql(/* GraphQL */ `
  query GetBazelInvocationLog($invocationID: UUID!) {
    getBazelInvocation(invocationID: $invocationID) {
      ...BazelInvocationLog
    }
  }
`);

const BAZEL_INVOCATION_LOG_FRAGMENT = gql(/* GraphQL */ `
  fragment BazelInvocationLog on BazelInvocation {
    id
    invocationID
    originalCommandLine
    exitCodeName
    connectionMetadata {
      id
      timeSinceLastConnectionMillis
    }
  }
`);

export const Route = createFileRoute("/bazel-invocations/$invocationID/log")({
  component: RouteComponent,
  loader: async ({ params }) => {
    // TODO (isakstenstrom): Move the log fetching to this loader instead
    const { data, error } = await apolloClient.query({
      errorPolicy: "all",
      query: GET_BAZEL_INVOCATION_LOG,
      variables: { invocationID: params.invocationID },
      fetchPolicy: "network-only",
    });

    if (!data?.getBazelInvocation) {
      throw new NotFoundError("invocation", error?.message);
    }

    const invocation = getFragmentData(
      BAZEL_INVOCATION_LOG_FRAGMENT,
      data.getBazelInvocation,
    );

    return {
      invocationID: invocation.invocationID,
      command: commandLineDataToString(invocation.originalCommandLine),
      isLiveAtLoad: shouldPollInvocation(invocation),
    };
  },
  head: (_ctx) => ({
    meta: [
      {
        title: generatePageTitle([
          "Invocation",
          "Log",
          _ctx.params.invocationID,
        ]),
      },
    ],
  }),
});

function RouteComponent() {
  const { invocationID, command, isLiveAtLoad } = Route.useLoaderData();

  // Keep the invocation's connection state fresh, so that the log stops
  // tailing once the invocation is done.
  const { data } = useQuery(GET_BAZEL_INVOCATION_LOG, {
    variables: { invocationID },
    pollInterval: INVOCATION_POLL_INTERVAL_MS,
    skip: !isLiveAtLoad,
  });
  const invocation = getFragmentData(
    BAZEL_INVOCATION_LOG_FRAGMENT,
    data?.getBazelInvocation,
  );

  return (
    <BuildLogsDisplay
      invocationId={invocationID}
      rawCommand={command}
      isLive={isLiveAtLoad && shouldPollInvocation(invocation)}
    />
  );
}
