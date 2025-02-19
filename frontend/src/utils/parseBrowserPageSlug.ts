import {
  type BrowserPageParams,
  getBrowserPageTypeFromString,
} from "@/types/BrowserPageType";
import { digestFunctionValueFromString } from "./digestFunctionUtils";

export const parseBrowserPageSlug = (
  slug: Array<string>,
): BrowserPageParams | undefined => {
  const blobIndex = slug.indexOf("blobs");
  if (blobIndex === -1 || blobIndex + 3 >= slug.length) {
    return undefined;
  }

  const instanceName = slug.slice(0, blobIndex).join("/");
  const digestFunction = digestFunctionValueFromString(slug[blobIndex + 1]);
  const browserPageType = getBrowserPageTypeFromString(slug[blobIndex + 2]);

  if (
    instanceName === "" ||
    digestFunction === undefined ||
    browserPageType === undefined
  ) {
    return undefined;
  }

  const hashAndSize = slug[blobIndex + 3];
  const [hash, sizeBytes] = hashAndSize.split("-");

  if (!hash || !sizeBytes) {
    return undefined;
  }

  const otherParams = slug.slice(blobIndex + 4);
  return {
    instanceName,
    digestFunction,
    browserPageType,
    digest: { hash, sizeBytes },
    otherParams,
  };
};
