import { createFileRoute } from "@tanstack/react-router";
import { apolloClient } from "@/components/ApolloWrapper";
import { BazelInvocationMetrics } from "@/components/pages/BazelInvocationMetrics";
import { InvocationDataNotFoundAlert } from "@/components/pages/InvocationDataNotFoundAlert";
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
    packageMetrics {
      ...BazelInvocationMetricsPackageMetrics
    }
    cumulativeMetrics {
      ...BazelInvocationMetricsCumulativeMetrics
    }
    buildGraphMetrics {
      ...BazelInvocationMetricsBuildGraphEvaluationMetrics
    }
    dynamicExecutionMetrics {
      ...BazelInvocationMetricsDynamicExecutionMetrics
    }
    workerMetrics {
      ...BazelInvocationMetricsWorkerMetrics
    }
    workerPoolMetrics {
      ...BazelInvocationMetricsWorkerPoolMetrics
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
      cacheCheckSemaphoreWaitTimeInMs
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
      ...BazelInvocationMetricsGarbageMetrics
    }
  }
`);

export const BAZEL_INVOCATION_METRICS_PACKAGE_METRICS_FRAGMENT =
  gql(/* GraphQL */ `
  fragment BazelInvocationMetricsPackageMetrics on PackageMetrics {
    id
    packagesLoaded
    packageLoadMetrics {
      id
      name
      loadDurationInNs
      numTargets
      computationSteps
      numTransitiveLoads
      packageOverhead
      globFilesystemOperationCost
    }
  }
`);

export const BAZEL_INVOCATION_METRICS_CUMULATIVE_METRICS_FRAGMENT =
  gql(/* GraphQL */ `
  fragment BazelInvocationMetricsCumulativeMetrics on CumulativeMetrics {
    id
    numAnalyses
    numBuilds
  }
`);

export const BAZEL_INVOCATION_METRICS_BUILD_GRAPH_EVALUATION_METRICS_FRAGMENT =
  gql(/* GraphQL */ `
  fragment BazelInvocationMetricsBuildGraphEvaluationMetrics on BuildGraphMetrics {
    id
    evaluationStats {
      id
      operation
      skyfunctionName
      count
    }
    ruleClassCounts {
      id
      key
      ruleClass
      count
      actionCount
    }
    aspectCounts {
      id
      key
      aspectName
      count
      actionCount
    }
  }
`);

export const BAZEL_INVOCATION_METRICS_DYNAMIC_EXECUTION_METRICS_FRAGMENT =
  gql(/* GraphQL */ `
  fragment BazelInvocationMetricsDynamicExecutionMetrics on DynamicExecutionMetrics {
    id
    raceStatistics {
      id
      mnemonic
      localRunner
      remoteRunner
      localWins
      remoteWins
    }
  }
`);

export const BAZEL_INVOCATION_METRICS_WORKER_METRICS_FRAGMENT =
  gql(/* GraphQL */ `
  fragment BazelInvocationMetricsWorkerMetrics on WorkerMetrics {
    id
    processID
    mnemonic
    isMultiplex
    isSandbox
    isMeasurable
    workerKeyHash
    workerStatus
    code
    actionsExecuted
    priorActionsExecuted
    workerIds {
      id
      workerID
    }
    workerStats {
      id
      collectTimeInMs
      workerMemoryInKB
      priorWorkerMemoryInKB
      lastActionStartTimeInMs
    }
  }
`);

export const BAZEL_INVOCATION_METRICS_WORKER_POOL_METRICS_FRAGMENT =
  gql(/* GraphQL */ `
  fragment BazelInvocationMetricsWorkerPoolMetrics on WorkerPoolMetrics {
    id
    workerPoolStats {
      id
      hash
      mnemonic
      createdCount
      destroyedCount
      evictedCount
      userExecExceptionDestroyedCount
      ioExceptionDestroyedCount
      interruptedExceptionDestroyedCount
      unknownDestroyedCount
      aliveCount
    }
  }
`);

export const BAZEL_INVOCATION_METRICS_GARBAGE_METRICS_FRAGMENT =
  gql(/* GraphQL */ `
  fragment BazelInvocationMetricsGarbageMetrics on GarbageMetrics {
    id
    garbageCollected
    type
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
    criticalPathTimeInMs
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
      return { metrics: undefined };
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

  if (metrics === undefined) {
    return <InvocationDataNotFoundAlert type="metrics" />;
  }

  return <BazelInvocationMetrics metrics={metrics} />;
}
