import { gql } from "@/graphql/__generated__";

export const LOAD_FULL_BAZEL_INVOCATION_DETAILS = gql(/* GraphQL */ `
  query LoadFullBazelInvocationDetails($invocationID: String!) {
    bazelInvocation(invocationId: $invocationID) {
      ...BazelInvocationInfo
    }
  }
`);

export const GET_PROBLEM_DETAILS = gql(/* GraphQL */ `
  query GetProblemDetails($invocationID: String!) {
    bazelInvocation(invocationId: $invocationID) {
      ...ProblemDetails
    }
  }
`);

export const PROBLEM_DETAILS_FRAGMENT = gql(/* GraphQL */`

  fragment ProblemDetails on BazelInvocation{
    problems {
        ...ProblemInfo
      }
  }

`)

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
      sourceArtifactsReadCount
      sourceArtifactsReadSizeInBytes
      outputArtifactsSeenCount
      outputArtifactsSeenSizeInBytes
      outputArtifactsFromActionCacheCount
      outputArtifactsFromActionCacheSizeInBytes
      topLevelArtifactsCount
      topLevelArtifactsSizeInBytes
    }
    cumulativeMetrics {
      id
      numBuilds
      numAnalyses
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
    id
    command
    executable
    residual
    explicitCmdLine
    cmdLine
    startupOptions
    explicitStartupOptions
  }
  id
  invocationID
  instanceName {
    name
  }
  authenticatedUser {
    displayName
    userUUID
  }
  bazelVersion
  build {
    id
    buildUUID
  }
  profile {
    id
    name
    digest
    sizeInBytes
    digestFunction
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
  user {
    Email
    LDAP
  }
  startedAt
  endedAt
  exitCodeName
  bepCompleted
  configurations {
    id
    cpu
    mnemonic
  }
  numFetches
  stepLabel
  hostname
  isCiWorker
  sourceControl {
    id
    provider
    instanceURL
    repo
    refs
    commitSha
    actor
    eventName
    workflow
    runID
    runNumber
    job
    action
    runnerName
    runnerArch
    runnerOs
  }
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
  ephemeralURL
}
`)


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