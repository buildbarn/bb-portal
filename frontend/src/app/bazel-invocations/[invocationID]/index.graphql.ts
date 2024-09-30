import { gql } from "@/graphql/__generated__";

export const LOAD_FULL_BAZEL_INVOCATION_DETAILS = gql(/* GraphQL */ `
  query LoadFullBazelInvocationDetails($invocationID: String!) {
    bazelInvocation(invocationId: $invocationID) {
      ...FullBazelInvocationDetails
    }
  }
`);

export const BAZEL_INVOCATION_FRAGMENT = gql(/* GraphQL */ `
fragment BazelInvocationInfo on BazelInvocation {
  metrics {
    id
    actionSummary {
      id
      actionsCreated
      actionsExecuted
      actionsCreatedNotIncludingAspects
      remoteCacheHits
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
        lastEndedMs
        actionsCreated
        actionsExecuted
        firstStartedMs
      }
    }
    artifactMetrics {
      id
      sourceArtifactsRead {
        id
        sizeInBytes
        count
      }
      outputArtifactsSeen {
        id
        sizeInBytes
        count
      }
      outputArtifactsFromActionCache {
        id
        sizeInBytes
        count
      }
      topLevelArtifacts {
        id
        sizeInBytes
        count
      }
    }
    cumulativeMetrics {
      id
      numBuilds
      numAnalyses
    }
    dynamicExecutionMetrics {
      id
      raceStatistics {
        id
        localWins
        mnemonic
        renoteWins
        localRunner
        remoteRunner
      }
    }
    buildGraphMetrics {
      id
      actionLookupValueCount
      actionLookupValueCountNotIncludingAspects
      actionCount
      inputFileConfiguredTargetCount
      outputFileConfiguredTargetCount
      otherConfiguredTargetCount
      outputArtifactCount
      postInvocationSkyframeNodeCount
    }
    memoryMetrics {
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
    targetMetrics {
      id
      targetsLoaded
      targetsConfigured
      targetsConfiguredNotIncludingAspects
    }
    timingMetrics {
      id
      cpuTimeInMs
      wallTimeInMs
      analysisPhaseTimeInMs
      executionPhaseTimeInMs
      actionsExecutionStartInMs
    }
    networkMetrics {
      id
      systemNetworkStats {
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
    }
    packageMetrics {
      id
      packagesLoaded
      packageLoadMetrics {
        id
        name
        numTargets
        loadDuration
        packageOverhead
        computationSteps
        numTransitiveLoads
      }
    }
  }
  bazelCommand {
    command
    executable
    id
    buildOptions: options
    residual
  }
  id
  invocationID
  build {
    id
    buildUUID
  }
  targets {
    id
    label
    success
    testSize
    targetKind
    durationInMs
    abortReason
  }
  testCollection {
    id
    label
    strategy
    durationMs
    overallStatus
    cachedLocally
    cachedRemotely
  }
  relatedFiles {
    name
    url
  }
  user {
    Email
    LDAP
  }
  startedAt
  endedAt
  state {
    bepCompleted
    buildEndTime
    buildStartTime
    exitCode {
      code
      id
      name
    }
    id
  }
  stepLabel
}
`);
export const PROBLEM_INFO_FRAGMENT = gql(/* GraphQL */ `
 fragment ProblemInfo on Problem {
  id
  label
  __typename
  ... on ActionProblem {
    __typename
    id
    label
    type
    stdout {
      ...BlobReferenceInfo
    }
    stderr {
      ...BlobReferenceInfo
    }
  }
  ... on TestProblem {
    __typename
    id
    label
    status
    results {
      __typename
      id
      run
      shard
      attempt
      status
      actionLogOutput {
        ...BlobReferenceInfo
      }
      undeclaredTestOutputs {
        ...BlobReferenceInfo
      }
    }
  }
  ... on TargetProblem {
    __typename
    id
    label
  }
  ... on ProgressProblem {
    __typename
    id
    output
    label
  }
}
`);

export const BLOB_REFERENCE_INFO_FRAGMENT = gql(/* GraphQL */ `
fragment BlobReferenceInfo on BlobReference {
  availabilityStatus
  name
  sizeInBytes
  downloadURL
}
`)


export const FULL_BAZEL_INVOCATION_DETAILS = gql(/* GraphQL */ `
    fragment FullBazelInvocationDetails on BazelInvocation {
      problems {
        ...ProblemInfo
      }
      ...BazelInvocationInfo
    }
`);




export const GET_ACTION_PROBLEM = gql(/* GraphQL */ `
  query GetActionProblem($id: ID!) {
    node(id: $id) {
      id
      ... on ActionProblem {
        label
        stdout {
          ...BlobReferenceInfo
        }
        stderr {
          ...BlobReferenceInfo
        }
      }
    }
  }
`);

export const TEST_RESULT_FRAGMENT = gql(/* GraphQL */`
fragment TestResultInfo on TestResult {
      actionLogOutput {
  ...BlobReferenceInfo
  }
  attempt
  run
  shard
  status
  undeclaredTestOutputs {
    ...BlobReferenceInfo
  }
}`)