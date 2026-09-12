import type { OperationState } from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import {
  type BrowserPageParams,
  parseBrowserPageSlug,
} from "@/types/BrowserPageParams";
import { BrowserPageType } from "@/types/BrowserPageType";
import { generateBrowserSplat } from "@/utils/urlGenerator";

export const historicalExecuteResponseDigestFromOperation = (
  operation: OperationState,
): BrowserPageParams | undefined => {
  if (
    !operation.completed?.message.startsWith(
      "Action details (uncached result):",
    )
  ) {
    return undefined;
  }
  const url = operation.completed.message.substring(34);
  const index = url.indexOf("/browser/");

  if (index === -1) {
    return undefined;
  }
  const urlSegments = url
    .substring(index + 9)
    .split("/")
    .filter((segment) => segment);
  return parseBrowserPageSlug(urlSegments);
};

const instanceNameFromOperationState = (operation: OperationState): string => {
  const instanceNamePrefix =
    operation.invocationName?.sizeClassQueueName?.platformQueueName
      ?.instanceNamePrefix;
  const instanceNameSuffix = operation.instanceNameSuffix;

  const instanceName = [];
  if (instanceNamePrefix) instanceName.push(instanceNamePrefix);
  if (instanceNameSuffix) instanceName.push(instanceNameSuffix);

  return instanceName.join("/");
};

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
    BrowserPageType.Action,
  );
};
