import { createFileRoute, notFound } from "@tanstack/react-router";
import { apolloClient } from "@/components/ApolloWrapper";
import SourceControlDisplay from "@/components/SourceControlDisplay";
import { getFragmentData, gql } from "@/graphql/__generated__";
import { generatePageTitle } from "@/utils/generatePageTitle";

const GET_BAZEL_INVOCATION_SOURCE_CONTROL = gql(/* GraphQL */ `
  query GetBazelInvocationSourceControl($invocationID: UUID!) {
    getBazelInvocation(invocationID: $invocationID) {
      id
      sourceControl {
        ...BazelInvocationSourceControl
      }
    }
  }
`);

const BAZEL_INVOCATION_SOURCE_CONTROL_FRAGMENT = gql(/* GraphQL */ `
  fragment BazelInvocationSourceControl on SourceControl {
    id
    repo
    repoURL
    ref
    refURL
    commit
    commitURL
  }
`);

export const Route = createFileRoute(
  "/bazel-invocations/$invocationID/source-control",
)({
  component: RouteComponent,
  loader: async ({ params }) => {
    const { data } = await apolloClient.query({
      query: GET_BAZEL_INVOCATION_SOURCE_CONTROL,
      variables: { invocationID: params.invocationID },
      fetchPolicy: "network-only",
    });

    if (!data?.getBazelInvocation?.sourceControl) {
      throw notFound();
    }

    const sourceControl = getFragmentData(
      BAZEL_INVOCATION_SOURCE_CONTROL_FRAGMENT,
      data.getBazelInvocation?.sourceControl,
    );
    return { sourceControl };
  },
  head: (_ctx) => ({
    meta: [
      {
        title: generatePageTitle([
          "Invocation",
          "Source control",
          _ctx.params.invocationID,
        ]),
      },
    ],
  }),
});

function RouteComponent() {
  const { sourceControl } = Route.useLoaderData();
  return <SourceControlDisplay sourceControl={sourceControl} />;
}
