import { createFileRoute, notFound } from "@tanstack/react-router";
import { apolloClient } from "@/components/ApolloWrapper";
import { BazelInvocationMetrics } from "@/components/pages/BazelInvocationMetrics";
import { getFragmentData, gql } from "@/graphql/__generated__";
import { generatePageTitle } from "@/utils/generatePageTitle";

const GET_BAZEL_INVOCATION_METRICS = gql(/* GraphQL */ `
  query GetBazelInvocationMetrics($invocationID: UUID!) {
    getBazelInvocation(invocationID: $invocationID) {
      id
      metrics {
        ...BazelInvocationMetrics
      }
    }
  }
`);

const BAZEL_INVOCATION_METRICS_FRAGMENT = gql(/* GraphQL */ `
  fragment BazelInvocationMetrics on Metrics {
    id
    actionSummary {
      ...BazelInvocationMetricsActionSummary
    }
    artifactMetrics {
      ...BazelInvocationMetricsArtifactMetrics
    }
    memoryMetrics {
      ...BazelInvocationMetricsMemoryMetrics
    }
    timingMetrics {
      ...BazelInvocationMetricsTimingMetrics
    }
    networkMetrics {
      id
      systemNetworkStats {
        ...BazelInvocationMetricsSystemNetworkStats
      }
    }
  }
`);

export const BAZEL_INVOCATION_METRICS_ACTION_SUMMARY_FRAGMENT =
  gql(/* GraphQL */ `
  fragment BazelInvocationMetricsActionSummary on ActionSummary {
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

export const BAZEL_INVOCATION_METRICS_ARTIFACT_METRICS_FRAGMENT =
  gql(/* GraphQL */ `
  fragment BazelInvocationMetricsArtifactMetrics on ArtifactMetrics {
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

export const BAZEL_INVOCATION_METRICS_MEMORY_METRICS_FRAGMENT =
  gql(/* GraphQL */ `
  fragment BazelInvocationMetricsMemoryMetrics on MemoryMetrics {
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

export const BAZEL_INVOCATION_METRICS_TIMING_METRICS_FRAGMENT =
  gql(/* GraphQL */ `
  fragment BazelInvocationMetricsTimingMetrics on TimingMetrics {
    id
    cpuTimeInMs
    wallTimeInMs
    analysisPhaseTimeInMs
    executionPhaseTimeInMs
    actionsExecutionStartInMs
  }
`);

export const BAZEL_INVOCATION_METRICS_SYSTEM_NETWORK_STATS_FRAGMENT =
  gql(/* GraphQL */ `
  fragment BazelInvocationMetricsSystemNetworkStats on SystemNetworkStats {
    id
    bytesSent
    bytesRecv
    packetsSent
    packetsRecv
    peakBytesSentPerSec
    peakBytesRecvPerSec
    peakPacketsSentPerSec
    peakPacketsRecvPerSec
  }
`);

export const Route = createFileRoute(
  "/bazel-invocations/$invocationID/metrics",
)({
  component: RouteComponent,
  loader: async ({ params }) => {
    const { data } = await apolloClient.query({
      query: GET_BAZEL_INVOCATION_METRICS,
      variables: { invocationID: params.invocationID },
      fetchPolicy: "network-only",
    });

    if (!data?.getBazelInvocation?.metrics) {
      throw notFound();
    }

    const metrics = getFragmentData(
      BAZEL_INVOCATION_METRICS_FRAGMENT,
      data?.getBazelInvocation?.metrics,
    );

    return { metrics };
  },
  head: (_ctx) => ({
    meta: [
      {
        title: generatePageTitle([
          "Invocation",
          "Metrics",
          _ctx.params.invocationID,
        ]),
      },
    ],
  }),
});

function RouteComponent() {
  const { metrics } = Route.useLoaderData();
  return <BazelInvocationMetrics metrics={metrics} />;
}
