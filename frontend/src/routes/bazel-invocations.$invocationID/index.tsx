import { createFileRoute } from "@tanstack/react-router";
import { apolloClient } from "@/components/ApolloWrapper";
import { InvocationOverviewDisplay } from "@/components/InvocationOverviewDisplay";
import { getFragmentData, gql } from "@/graphql/__generated__";
import { NotFoundError } from "@/main";
import { generatePageTitle } from "@/utils/generatePageTitle";

const GET_BAZEL_INVOCATION_OVERVIEW = gql(/* GraphQL */ `
  query GetBazelInvocationOverview($invocationID: UUID!) {
    getBazelInvocation(invocationID: $invocationID) {
      ...BazelInvocationOverview
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

    return {
      invocation: getFragmentData(
        BAZEL_INVOCATION_OVERVIEW_FRAGMENT,
        data.getBazelInvocation,
      ),
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
  const { invocation } = Route.useLoaderData();
  return <InvocationOverviewDisplay invocation={invocation} />;
}
