import { createFileRoute, notFound } from "@tanstack/react-router";
import { apolloClient } from "@/components/ApolloWrapper";
import MemoryMetricsDisplay from "@/components/MemoryMetrics";
import { getFragmentData, gql } from "@/graphql/__generated__";
import { generatePageTitle } from "@/utils/generatePageTitle";

const GET_BAZEL_INVOCATION_MEMORY_METRICS = gql(/* GraphQL */ `
  query GetBazelInvocationMemoryMetrics($invocationID: UUID!) {
    getBazelInvocation(invocationID: $invocationID) {
      id
      metrics {
        id
        memoryMetrics {
          ...BazelInvocationMemoryMetrics
          }
      }
    }
  }
`);

const BAZEL_INVOCATION_MEMORY_METRICS_FRAGMENT = gql(/* GraphQL */ `
  fragment BazelInvocationMemoryMetrics on MemoryMetrics {
    id
    usedHeapSizePostBuild
    peakPostGcHeapSize
    peakPostGcTenuredSpaceHeapSize
    garbageMetrics {
      id
      garbageCollected
      type
    }
  }
`);

export const Route = createFileRoute("/bazel-invocations/$invocationID/memory")(
  {
    component: RouteComponent,
    loader: async ({ params }) => {
      const { data } = await apolloClient.query({
        query: GET_BAZEL_INVOCATION_MEMORY_METRICS,
        variables: { invocationID: params.invocationID },
        fetchPolicy: "network-only",
      });

      if (!data?.getBazelInvocation?.metrics?.memoryMetrics) {
        throw notFound();
      }

      const memoryMetrics = getFragmentData(
        BAZEL_INVOCATION_MEMORY_METRICS_FRAGMENT,
        data.getBazelInvocation?.metrics?.memoryMetrics,
      );
      return { memoryMetrics };
    },
    head: (_ctx) => ({
      meta: [
        {
          title: generatePageTitle([
            "Invocation",
            "Memroy metrics",
            _ctx.params.invocationID,
          ]),
        },
      ],
    }),
  },
);

function RouteComponent() {
  const { memoryMetrics } = Route.useLoaderData();
  return <MemoryMetricsDisplay memoryMetrics={memoryMetrics} />;
}
