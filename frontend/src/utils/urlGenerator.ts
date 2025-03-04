import {
  Digest,
  DigestFunction_Value,
} from '@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution';
import { env } from 'next-runtime-env';
import { digestFunctionValueToString } from './digestFunctionUtils';

export function generateFileUrl(
  instanceName: string,
  digestFunction: DigestFunction_Value,
  digest: Digest,
  fileName: string,
): string {
  return `${env(
    'NEXT_PUBLIC_BES_BACKEND_URL',
  )}/api/servefile/${instanceName}/blobs/${digestFunctionValueToString(
    digestFunction,
  )}/file/${digest.hash}-${digest.sizeBytes}/${fileName}`;
}

export function generateCommandShellScriptUrl(
  instanceName: string,
  digestFunction: DigestFunction_Value,
  digest: Digest,
): string {
  return `${env(
    'NEXT_PUBLIC_BES_BACKEND_URL',
  )}/api/servefile/${instanceName}/blobs/${digestFunctionValueToString(
    digestFunction,
  )}/command/${digest.hash}-${digest.sizeBytes}/?format=sh`;
}

export function generateDirectoryTarballUrl(
  instanceName: string,
  digestFunction: DigestFunction_Value,
  digest: Digest,
): string {
  return `${env(
    'NEXT_PUBLIC_BES_BACKEND_URL',
  )}/api/servefile/${instanceName}/blobs/${digestFunctionValueToString(
    digestFunction,
  )}/directory/${digest.hash}-${digest.sizeBytes}/?format=tar`;
}
