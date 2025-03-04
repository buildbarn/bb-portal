import type { OperationState } from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import { digestFunctionValueToString } from "@/utils/digestFunctionUtils";

export const operationsStateToActionPageUrl = (
  operation: OperationState,
): string | undefined => {
  const instanceNamePrefix =
    operation.invocationName?.sizeClassQueueName?.platformQueueName
      ?.instanceNamePrefix;
  const digestFunction = digestFunctionValueToString(operation.digestFunction);
  const actionDigestHash = operation.actionDigest?.hash;
  const actionDigestSizeBytes = operation.actionDigest?.sizeBytes;

  if (
    !instanceNamePrefix ||
    !digestFunction ||
    !actionDigestHash ||
    !actionDigestSizeBytes
  ) {
    return undefined;
  }

  return `/browser/${instanceNamePrefix}/blobs/${digestFunction}/action/${actionDigestHash}-${actionDigestSizeBytes}`;
};
