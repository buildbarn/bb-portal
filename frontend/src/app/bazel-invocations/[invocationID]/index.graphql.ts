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
