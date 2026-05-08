import { createFileRoute, notFound } from "@tanstack/react-router";
import { apolloClient } from "@/components/ApolloWrapper";
import { InvocationTargetsTab } from "@/components/InvocationTargets/InvocationTargetsTab";
import { getFragmentData, gql } from "@/graphql/__generated__";
import { env } from "@/utils/env";
import { requireFeature } from "@/utils/featureGuard";
import { generatePageTitle } from "@/utils/generatePageTitle";

const GET_BAZEL_INVOCATION_TARGETS = gql(/* GraphQL */ `
  query GetBazelInvocationTargets($invocationID: UUID!) {
    getBazelInvocation(invocationID: $invocationID) {
      id
      metrics {
        id
        targetMetrics {
          ...BazelInvocationTargetMetrics
          }
      }
    }
  }
`);

const BAZEL_INVOCATION_TARGET_METRICS_FRAGMENT = gql(/* GraphQL */ `
  fragment BazelInvocationTargetMetrics on TargetMetrics {
    id
    targetsLoaded
    targetsConfigured
    targetsConfiguredNotIncludingAspects
  }
`);

export const Route = createFileRoute(
  "/bazel-invocations/$invocationID/targets",
)({
  component: RouteComponent,
  beforeLoad: requireFeature(env.featureFlags?.bes?.pageTargets),
  loader: async ({ params }) => {
    const { data } = await apolloClient.query({
      query: GET_BAZEL_INVOCATION_TARGETS,
      variables: { invocationID: params.invocationID },
      fetchPolicy: "network-only",
    });

    if (!data?.getBazelInvocation) {
      throw notFound();
    }

    const targetMetrics = getFragmentData(
      BAZEL_INVOCATION_TARGET_METRICS_FRAGMENT,
      data.getBazelInvocation?.metrics?.targetMetrics ?? undefined,
    );
    return { targetMetrics };
  },
  head: (_ctx) => ({
    meta: [
      {
        title: generatePageTitle([
          "Invocation",
          "Targets",
          _ctx.params.invocationID,
        ]),
      },
    ],
  }),
});

function RouteComponent() {
  const { invocationID } = Route.useParams();
  const { targetMetrics } = Route.useLoaderData();
  // TODO (isakstenstrom): Fetch data in the data loader
  return (
    <InvocationTargetsTab
      invocationId={invocationID}
      targetMetrics={targetMetrics}
    />
  );
}
