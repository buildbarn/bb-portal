export enum InvocationResult {
  // ExitCodes from the Build Event Stream
  SUCCESS = "SUCCESS",
  UNSTABLE = "UNSTABLE",
  PARSING_FAILURE = "PARSING_FAILURE",
  BUILD_FAILURE = "BUILD_FAILURE",
  TESTS_FAILED = "TESTS_FAILED",
  REMOTE_ERROR = "REMOTE_ERROR",
  NOT_BUILT = "NOT_BUILT",
  ABORTED = "ABORTED",
  INTERRUPTED = "INTERRUPTED",
  // Custom statuses
  UNKNOWN_EXIT_CODE = "UNKNOWN_EXIT_CODE",
  IN_PROGRESS = "IN_PROGRESS",
  DISCONNECTED = "DISCONNECTED",
}

export const InvocationExitCodes = [
  InvocationResult.SUCCESS.toString(),
  InvocationResult.UNSTABLE.toString(),
  InvocationResult.PARSING_FAILURE.toString(),
  InvocationResult.BUILD_FAILURE.toString(),
  InvocationResult.TESTS_FAILED.toString(),
  InvocationResult.REMOTE_ERROR.toString(),
  InvocationResult.NOT_BUILT.toString(),
  InvocationResult.ABORTED.toString(),
  InvocationResult.INTERRUPTED.toString(),
];

export const INVOCATION_IN_PROGRESS_TIMEOUT = 12 * 1000;

export const getInvocationResultTagEnum = (
  exitCodeName: string | undefined,
  timeSinceLastConnectionMillis: number | undefined,
): InvocationResult => {
  if (exitCodeName) {
    if (InvocationExitCodes.includes(exitCodeName)) {
      return exitCodeName as InvocationResult;
    }
    return InvocationResult.UNKNOWN_EXIT_CODE;
  }

  if (!timeSinceLastConnectionMillis) {
    return InvocationResult.DISCONNECTED;
  }

  if (timeSinceLastConnectionMillis <= INVOCATION_IN_PROGRESS_TIMEOUT) {
    return InvocationResult.IN_PROGRESS;
  }
  return InvocationResult.DISCONNECTED;
};
