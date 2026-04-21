import { gql } from "@/graphql/__generated__";

export const LOAD_FULL_BAZEL_INVOCATION_DETAILS = gql(/* GraphQL */ `
  query LoadFullBazelInvocationDetails($invocationID: UUID!) {
    getBazelInvocation(invocationID: $invocationID) {
      ...BazelInvocationInfo
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
      sourceArtifactsReadCount
      sourceArtifactsReadSizeInBytes
      outputArtifactsSeenCount
      outputArtifactsSeenSizeInBytes
      outputArtifactsFromActionCacheCount
      outputArtifactsFromActionCacheSizeInBytes
      topLevelArtifactsCount
      topLevelArtifactsSizeInBytes
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
  }
  canonicalCommandLine
  originalCommandLine
  optionsParsed
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
  actions {
    id
    label
    type
    success
    exitCode
    commandLine
    startTime
    endTime
    failureCode
    failureMessage
    stdoutHash
    stdoutSizeBytes
    stdoutHashFunction
    stderrHash
    stderrSizeBytes
    stderrHashFunction
    configuration {
      id
      configurationID
      mnemonic
      platformName
      cpu
      makeVariables
    }
  }
  profile {
    id
    name
    digest
    sizeInBytes
    digestFunction
  }
  username
  startedAt
  endedAt
  exitCodeName
  connectionMetadata {
    connectionLastOpenAt
    timeSinceLastConnectionMillis
  }
  configurations {
    id
    cpu
    mnemonic
  }
  numFetches
  hostname
  sourceControl {
    id
    repo
    repoURL
    ref
    refURL
    commit
    commitURL
  }
  tags(orderBy: { field: KEY, direction: ASC }) {
    edges {
      node {
        id
        key
        value
      }
    }
  }
}
`);
