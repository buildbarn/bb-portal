import type { FileDetailsFragment } from "@/graphql/__generated__/graphql";
import type {
  Digest,
  DigestFunction_Value,
} from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import { BrowserPageType } from "@/types/BrowserPageType";
import {
  digestFunctionValueFromString,
  digestFunctionValueToString,
} from "./digestFunctionUtils";

/////////////////////////////////////////////////////////////
// Frontend internal URLs
/////////////////////////////////////////////////////////////
export function generateBrowserSplat(
  instanceName: string | undefined,
  digestFunction: DigestFunction_Value,
  digest: Digest,
  pageType: BrowserPageType,
): string {
  return `${instanceName ? `${instanceName}/` : ""}blobs/${digestFunctionValueToString(
    digestFunction,
  )}/${pageType}/${digest.hash}-${digest.sizeBytes}`;
}

// TODO (isakstenstrom): Remove once the browser file component is rewritten
export function generateTreeUrl(
  instanceName: string | undefined,
  digestFunction: DigestFunction_Value,
  digest: Digest,
): string {
  return `/browser${instanceName ? `/${instanceName}` : ""}/blobs/${digestFunctionValueToString(
    digestFunction,
  )}/tree/${digest.hash}-${digest.sizeBytes}`;
}

// TODO (isakstenstrom): Remove once the browser file component is rewritten
export function generateDirectoryUrl(
  instanceName: string | undefined,
  digestFunction: DigestFunction_Value,
  digest: Digest,
): string {
  return `/browser${instanceName ? `/${instanceName}` : ""}/blobs/${digestFunctionValueToString(
    digestFunction,
  )}/directory/${digest.hash}-${digest.sizeBytes}`;
}

/////////////////////////////////////////////////////////////
// Backend URLs
/////////////////////////////////////////////////////////////

const BACKEND_SERVE_FILE_URL = "/api/v1/servefile";

export function generateFileUrl(
  instanceName: string | undefined,
  digestFunction: DigestFunction_Value,
  digest: Digest,
  fileName: string,
): string {
  // The link will not work if there are slashes in the filename. So we take
  // just the last segment from the file name.
  const lastFileNameSegment = fileName.split("/").pop() ?? fileName;
  return `${BACKEND_SERVE_FILE_URL}/${generateBrowserSplat(instanceName, digestFunction, digest, BrowserPageType.File)}/${lastFileNameSegment}`;
}

export function generateFileUrlFromGraphqlFile(
  file: FileDetailsFragment,
): string {
  return generateFileUrl(
    file.digest.rev2InstanceName,
    digestFunctionValueFromString(file.digest.digestFunction),
    {
      hash: file.digest.hash,
      sizeBytes: file.digest.sizeBytes.toString(),
    },
    file.filePath.path,
  );
}

export function generateCommandShellScriptUrl(
  instanceName: string | undefined,
  digestFunction: DigestFunction_Value,
  digest: Digest,
): string {
  return `${BACKEND_SERVE_FILE_URL}/${generateBrowserSplat(instanceName, digestFunction, digest, BrowserPageType.Command)}/?format=sh`;
}

export function generateDirectoryTarballUrl(
  instanceName: string | undefined,
  digestFunction: DigestFunction_Value,
  digest: Digest,
): string {
  return `${BACKEND_SERVE_FILE_URL}/${generateBrowserSplat(instanceName, digestFunction, digest, BrowserPageType.Directory)}/?format=tar`;
}
