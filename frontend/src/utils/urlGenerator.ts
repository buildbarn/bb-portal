import { env } from "next-runtime-env";
import type {
  Digest,
  DigestFunction_Value,
} from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import { digestFunctionValueToString } from "./digestFunctionUtils";

function getServeFileUrl(): string {
  return `${env("NEXT_PUBLIC_BES_BACKEND_URL")}/api/v1/servefile`;
}

export function generateFileUrl(
  instanceName: string | undefined,
  digestFunction: DigestFunction_Value,
  digest: Digest,
  fileName: string,
): string {
  return `${getServeFileUrl()}${instanceName ? `/${instanceName}` : ""}/blobs/${digestFunctionValueToString(
    digestFunction,
  )}/file/${digest.hash}-${digest.sizeBytes}/${fileName}`;
}

export function generateCommandShellScriptUrl(
  instanceName: string | undefined,
  digestFunction: DigestFunction_Value,
  digest: Digest,
): string {
  return `${getServeFileUrl()}${instanceName ? `/${instanceName}` : ""}/blobs/${digestFunctionValueToString(
    digestFunction,
  )}/command/${digest.hash}-${digest.sizeBytes}/?format=sh`;
}

export function generateDirectoryUrl(
  instanceName: string | undefined,
  digestFunction: DigestFunction_Value,
  digest: Digest,
): string {
  return `/browser${instanceName ? `/${instanceName}` : ""}/blobs/${digestFunctionValueToString(
    digestFunction,
  )}/directory/${digest.hash}-${digest.sizeBytes}`;
}

export function generateDirectoryTarballUrl(
  instanceName: string | undefined,
  digestFunction: DigestFunction_Value,
  digest: Digest,
): string {
  return `${getServeFileUrl()}${instanceName ? `/${instanceName}` : ""}/blobs/${digestFunctionValueToString(
    digestFunction,
  )}/directory/${digest.hash}-${digest.sizeBytes}/?format=tar`;
}

export function generateUrlFromEphemeralUrl(
  instanceName: string | undefined,
  ephemeralUrl: string,
): string {
  return `${getServeFileUrl()}${instanceName ? `/${instanceName}` : ""}${ephemeralUrl}`;
}

export function generateLinkToTargetsPage(
  instanceName: string,
  label: string,
  aspect: string,
  targetKind: string,
): string {
  const params = new URLSearchParams({
    instanceName: encodeURIComponent(encodeURIComponent(instanceName)),
    label: encodeURIComponent(encodeURIComponent(label)),
    aspect: encodeURIComponent(encodeURIComponent(aspect)),
    targetKind: encodeURIComponent(encodeURIComponent(targetKind)),
  });
  return `/target?${params.toString()}`;
}
