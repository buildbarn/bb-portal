import type { OperationState } from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import { BrowserPageType } from "@/types/BrowserPageType";
import { generateBrowserSplat } from "@/utils/urlGenerator";

export const instanceNameFromOperationState = (operation: OperationState): string => {
  const instanceNamePrefix =
    operation.invocationName?.sizeClassQueueName?.platformQueueName
      ?.instanceNamePrefix;
  const instanceNameSuffix = operation.instanceNameSuffix;

  const instanceName = []
  if (instanceNamePrefix) instanceName.push(instanceNamePrefix)
  if (instanceNameSuffix) instanceName.push(instanceNameSuffix)

  return instanceName.join("/")
}

export const operationsStateToBrowserSplat = (
  operation: OperationState,
): string | undefined => {
  if (operation.actionDigest === undefined) {
    return undefined;
  }
  return generateBrowserSplat(
    instanceNameFromOperationState(operation),
    operation.digestFunction,
    operation.actionDigest,
    BrowserPageType.Action
  )
};
