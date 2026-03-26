import type {
  Digest,
  DigestFunction_Value,
} from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import { BrowserPageType } from "@/types/BrowserPageType";
import { digestFunctionValueToString } from "./digestFunctionUtils";

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
  return `${BACKEND_SERVE_FILE_URL}/${generateBrowserSplat(instanceName, digestFunction, digest, BrowserPageType.File)}/${fileName}`;
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
