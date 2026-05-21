import { createFileRoute, notFound } from "@tanstack/react-router";
import { apolloClient } from "@/components/ApolloWrapper";
import { BazelInvocationDetailsPage } from "@/components/pages/BazelInvocationDetails";
import { getFragmentData, gql } from "@/graphql/__generated__";

const GET_BAZEL_INVOCATION_COMMON = gql(/* GraphQL */ `
  query GetBazelInvocationCommon($invocationID: UUID!) {
    getBazelInvocation(invocationID: $invocationID) {
      ...BazelInvocationCommon
    }
  }
`);

const BAZEL_INVOCATION_COMMON_FRAGMENT = gql(/* GraphQL */ `
  fragment BazelInvocationCommon on BazelInvocation {
    id
    invocationID
    startedAt
    endedAt
    exitCodeName
    instanceName {
      name
    }
    connectionMetadata {
      connectionLastOpenAt
      timeSinceLastConnectionMillis
    }
    username
    authenticatedUser {
      displayName
      userUUID
    }
    build {
      id
      buildUUID
    }
    metrics {
      id
    }
    actions {
      id
    }
    profile {
      id
      name
      digest
      sizeInBytes
      digestFunction
    }
    sourceControl {
      id
    }
  }
`);

export const Route = createFileRoute("/bazel-invocations/$invocationID")({
  component: RouteComponent,
  loader: async ({ params }) => {
    const { data } = await apolloClient.query({
      query: GET_BAZEL_INVOCATION_COMMON,
      variables: { invocationID: params.invocationID },
      fetchPolicy: "network-only",
    });

    if (!data?.getBazelInvocation) {
      throw notFound();
    }

    return {
      invocation: getFragmentData(
        BAZEL_INVOCATION_COMMON_FRAGMENT,
        data.getBazelInvocation,
      ),
    };
  },
});

function RouteComponent() {
  const { invocation } = Route.useLoaderData();
  return <BazelInvocationDetailsPage invocation={invocation} />;
}
