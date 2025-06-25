import type { OperationState } from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";

export const historicalExecuteResponseUrlFromOperation = (
  operation: OperationState,
): string | undefined => {
  if (
    operation.completed?.message.startsWith("Action details (uncached result):")
  ) {
    return operation.completed.message.substring(34);
  }
  return undefined;
};

export const historicalExecuteResponseDigestFromUrl = (
  url: string | undefined,
): string | undefined => {
  if (!url) return undefined;
  const match = url.match(/([a-f0-9]{64})-[0-9]+/);
  return match ? match[0] : undefined;
};
