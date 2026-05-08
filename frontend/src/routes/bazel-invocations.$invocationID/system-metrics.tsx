import { createFileRoute } from "@tanstack/react-router";
import { apolloClient } from "@/components/ApolloWrapper";
import SystemMetricsDisplay from "@/components/SystemMetricsDisplay";
import { getFragmentData, gql } from "@/graphql/__generated__";
import { generatePageTitle } from "@/utils/generatePageTitle";

const GET_BAZEL_INVOCATION_SYSTEM_METRICS = gql(/* GraphQL */ `
  query GetBazelInvocationSystemMetrics($invocationID: UUID!) {
    getBazelInvocation(invocationID: $invocationID) {
      id
      metrics {
        id
        timingMetrics {
          ...BazelInvocationTimingMetrics
        }
        networkMetrics {
          id
          systemNetworkStats {
            ...BazelInvocationSystemNetworkStats
          }
        }
      }
    }
  }
`);

const BAZEL_INVOCATION_TIMING_METRICS_FRAGMENT = gql(/* GraphQL */ `
  fragment BazelInvocationTimingMetrics on TimingMetrics {
    id
    cpuTimeInMs
    wallTimeInMs
    analysisPhaseTimeInMs
    executionPhaseTimeInMs
    actionsExecutionStartInMs
  }
`);

const BAZEL_INVOCATION_SYSTEM_NETWORK_STATS_FRAGMENT = gql(/* GraphQL */ `
  fragment BazelInvocationSystemNetworkStats on SystemNetworkStats {
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
  "/bazel-invocations/$invocationID/system-metrics",
)({
  component: RouteComponent,
  loader: async ({ params }) => {
    const { data } = await apolloClient.query({
      query: GET_BAZEL_INVOCATION_SYSTEM_METRICS,
      variables: { invocationID: params.invocationID },
      fetchPolicy: "network-only",
    });

    const timingMetrics = getFragmentData(
      BAZEL_INVOCATION_TIMING_METRICS_FRAGMENT,
      data?.getBazelInvocation?.metrics?.timingMetrics,
    );
    const systemNetworkStats = getFragmentData(
      BAZEL_INVOCATION_SYSTEM_NETWORK_STATS_FRAGMENT,
      data?.getBazelInvocation?.metrics?.networkMetrics?.systemNetworkStats,
    );
    return { timingMetrics, systemNetworkStats };
  },
  head: (_ctx) => ({
    meta: [
      {
        title: generatePageTitle([
          "Invocation",
          "System metrics",
          _ctx.params.invocationID,
        ]),
      },
    ],
  }),
});

function RouteComponent() {
  const { timingMetrics, systemNetworkStats } = Route.useLoaderData();
  return (
    <SystemMetricsDisplay
      timingMetrics={timingMetrics}
      systemNetworkStats={systemNetworkStats}
    />
  );
}
