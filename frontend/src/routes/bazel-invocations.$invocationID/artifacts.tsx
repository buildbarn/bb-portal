import { createFileRoute, notFound } from "@tanstack/react-router";
import { apolloClient } from "@/components/ApolloWrapper";
import ArtifactsDataMetrics from "@/components/Artifacts";
import { getFragmentData, gql } from "@/graphql/__generated__";
import { generatePageTitle } from "@/utils/generatePageTitle";

const GET_BAZEL_INVOCATION_ARTIFACT_METRICS = gql(/* GraphQL */ `
  query GetBazelInvocationArtifactMetrics($invocationID: UUID!) {
    getBazelInvocation(invocationID: $invocationID) {
      id
      metrics {
        id
        artifactMetrics {
          ...BazelInvocationArtifactMetrics
          }
      }
    }
  }
`);

const BAZEL_INVOCATION_ARTIFACT_METRICS_FRAGMENT = gql(/* GraphQL */ `
  fragment BazelInvocationArtifactMetrics on ArtifactMetrics {
    id
    sourceArtifactsReadCount
    sourceArtifactsReadSizeInBytes
    outputArtifactsSeenCount
    outputArtifactsSeenSizeInBytes
    outputArtifactsFromActionCacheCount
    outputArtifactsFromActionCacheSizeInBytes
    topLevelArtifactsCount
    topLevelArtifactsSizeInBytes
  }
`);

export const Route = createFileRoute(
  "/bazel-invocations/$invocationID/artifacts",
)({
  component: RouteComponent,
  loader: async ({ params }) => {
    const { data } = await apolloClient.query({
      query: GET_BAZEL_INVOCATION_ARTIFACT_METRICS,
      variables: { invocationID: params.invocationID },
      fetchPolicy: "network-only",
    });

    if (!data?.getBazelInvocation?.metrics?.artifactMetrics) {
      throw notFound();
    }

    const artifactMetrics = getFragmentData(
      BAZEL_INVOCATION_ARTIFACT_METRICS_FRAGMENT,
      data.getBazelInvocation?.metrics?.artifactMetrics,
    );
    return { artifactMetrics };
  },
  head: (_ctx) => ({
    meta: [
      {
        title: generatePageTitle([
          "Invocation",
          "Artifacts",
          _ctx.params.invocationID,
        ]),
      },
    ],
  }),
});

function RouteComponent() {
  const { artifactMetrics } = Route.useLoaderData();
  return <ArtifactsDataMetrics artifactMetrics={artifactMetrics} />;
}
