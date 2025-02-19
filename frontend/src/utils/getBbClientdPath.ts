import type {
  Digest,
  DigestFunction_Value,
} from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import { digestFunctionValueToString } from "./digestFunctionUtils";

export function getBBClientdPath(
  instanceName: string,
  digestFunction: DigestFunction_Value,
  digest: Digest,
  blobType: string,
): string {
  return `~/bb_clientd/cas/${instanceName}/blobs/${digestFunctionValueToString(
    digestFunction,
  )}/${blobType}/${digest.hash}-${digest.sizeBytes}`;
}
