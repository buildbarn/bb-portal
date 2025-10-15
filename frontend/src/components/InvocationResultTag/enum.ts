export enum InvocationResultTagEnum {
  SUCCESS = "SUCCESS",
  UNSTABLE = "UNSTABLE",
  PARSING_FAILURE = "PARSING_FAILURE",
  BUILD_FAILURE = "BUILD_FAILURE",
  TESTS_FAILED = "TESTS_FAILED",
  REMOTE_ERROR = "REMOTE_ERROR",
  NOT_BUILT = "NOT_BUILT",
  ABORTED = "ABORTED",
  INTERRUPTED = "INTERRUPTED",
  UNKNOWN = "UNKNOWN",
  IN_PROGRESS = "IN_PROGRESS",
  BEP_UPLOAD_ABORTED = "BEP_UPLOAD_ABORTED",
}

export const InvocationExitCodes = [
  InvocationResultTagEnum.SUCCESS.toString(),
  InvocationResultTagEnum.UNSTABLE.toString(),
  InvocationResultTagEnum.PARSING_FAILURE.toString(),
  InvocationResultTagEnum.BUILD_FAILURE.toString(),
  InvocationResultTagEnum.TESTS_FAILED.toString(),
  InvocationResultTagEnum.REMOTE_ERROR.toString(),
  InvocationResultTagEnum.NOT_BUILT.toString(),
  InvocationResultTagEnum.ABORTED.toString(),
  InvocationResultTagEnum.INTERRUPTED.toString(),
];

export const getInvocationResultTagEnum = (
  exitCodeName: string | undefined,
  bepCompleted: boolean,
): InvocationResultTagEnum => {
  if (exitCodeName === undefined || exitCodeName === "") {
    if (!bepCompleted) {
      return InvocationResultTagEnum.IN_PROGRESS;
    } else {
      return InvocationResultTagEnum.BEP_UPLOAD_ABORTED;
    }
  }
  if (InvocationExitCodes.includes(exitCodeName)) {
    return exitCodeName as InvocationResultTagEnum;
  }
  return InvocationResultTagEnum.UNKNOWN;
};
