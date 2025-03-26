import type { OperationState } from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import { digestFunctionValueToString } from "@/utils/digestFunctionUtils";

export const operationsStateToActionPageUrl = (
  operation: OperationState,
): string | undefined => {
  const instanceNamePrefix =
    operation.invocationName?.sizeClassQueueName?.platformQueueName
      ?.instanceNamePrefix;
  const instanceNameSuffix = operation.instanceNameSuffix;
  const digestFunction = digestFunctionValueToString(operation.digestFunction);
  const actionDigestHash = operation.actionDigest?.hash;
  const actionDigestSizeBytes = operation.actionDigest?.sizeBytes;

  if (!digestFunction || !actionDigestHash || !actionDigestSizeBytes) {
    return undefined;
  }

  let url = "/browser";
  if (instanceNamePrefix) url += `/${instanceNamePrefix}`;
  if (instanceNameSuffix) url += `/${instanceNameSuffix}`;
  url += `/blobs/${digestFunction}/action/${actionDigestHash}-${actionDigestSizeBytes}`;
  return url;
};
