const CONNECTION_TIME_REFETCH_CUTOFF = 60 * 1000;

interface ShouldPollInvocationType {
  exitCodeName?: string | undefined | null;
  connectionMetadata?:
    | {
        timeSinceLastConnectionMillis: number;
      }
    | undefined
    | null;
}

export const shouldPollInvocation = (
  invocation: ShouldPollInvocationType | undefined | null,
): boolean => {
  if (invocation?.exitCodeName) {
    return false;
  }
  const elapsedTime =
    invocation?.connectionMetadata?.timeSinceLastConnectionMillis;
  if (!elapsedTime) {
    return true;
  }
  if (elapsedTime <= CONNECTION_TIME_REFETCH_CUTOFF) {
    return true;
  }
  return false;
};
