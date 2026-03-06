import type { OperationState } from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import { BrowserPageParams } from "@/types/BrowserPageType";
import { parseBrowserPageSlug } from "@/utils/parseBrowserPageSlug";

export const historicalExecuteResponseDigestFromOperation = (
  operation: OperationState,
): BrowserPageParams | undefined => {
  if (!operation.completed?.message.startsWith("Action details (uncached result):")) {
    return undefined
  }
  const url = operation.completed.message.substring(34);
  const index = url.indexOf("/browser/")

  if (index === -1) {
    return undefined
  }
  const urlSegments = url.substring(index + 9).split("/").filter(segment => segment)
  return parseBrowserPageSlug(urlSegments)
};
