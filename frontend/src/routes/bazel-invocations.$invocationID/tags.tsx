import { createFileRoute, notFound } from "@tanstack/react-router";
import { apolloClient } from "@/components/ApolloWrapper";
import { InvocationTagTab } from "@/components/InvocationTagTab";
import { gql } from "@/graphql/__generated__";
import { generatePageTitle } from "@/utils/generatePageTitle";
import { parseGraphqlEdgeListWithFragment } from "@/utils/parseGraphqlEdgeList";

const GET_BAZEL_INVOCATION_TAGS = gql(/* GraphQL */ `
  query GetBazelInvocationTags($invocationID: UUID!) {
    getBazelInvocation(invocationID: $invocationID) {
      id
      tags(orderBy: { field: KEY, direction: ASC }) {
        edges {
          node {
            ...BazelInvocationTags
          }
        }
      }
    }
  }
`);

const BAZEL_INVOCATION_TAGS_FRAGMENT = gql(/* GraphQL */ `
  fragment BazelInvocationTags on InvocationTag {
    id
    key
    value
  }
`);

export const Route = createFileRoute("/bazel-invocations/$invocationID/tags")({
  component: RouteComponent,
  loader: async ({ params }) => {
    const { data } = await apolloClient.query({
      query: GET_BAZEL_INVOCATION_TAGS,
      variables: { invocationID: params.invocationID },
      fetchPolicy: "network-only",
    });

    if (!data?.getBazelInvocation?.tags) {
      throw notFound();
    }

    const tags = parseGraphqlEdgeListWithFragment(
      BAZEL_INVOCATION_TAGS_FRAGMENT,
      data.getBazelInvocation?.tags,
    );
    return { tags };
  },
  head: (_ctx) => ({
    meta: [
      {
        title: generatePageTitle([
          "Invocation",
          "Tags",
          _ctx.params.invocationID,
        ]),
      },
    ],
  }),
});

function RouteComponent() {
  const { tags } = Route.useLoaderData();
  return <InvocationTagTab tags={tags} />;
}
