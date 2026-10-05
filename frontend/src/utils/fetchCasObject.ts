import { useQuery } from "@tanstack/react-query";
import { casClient } from "@/grpc/casClient";
import type {
  Digest,
  DigestFunction_Value,
  FindMissingBlobsRequest,
} from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import { protobufToObject } from "@/utils/protobufToObject";
import { generateFileUrl } from "@/utils/urlGenerator";

export const fetchCasBytes = async (
  instanceName: string | undefined,
  digestFunction: DigestFunction_Value,
  digest: Digest,
): Promise<Uint8Array> => {
  const res = await fetch(
    generateFileUrl(instanceName, digestFunction, digest, ""),
  );
  if (!res.ok) {
    throw new Error(`Failed to fetch CAS object: ${res.statusText}`);
  }
  return new Uint8Array(await res.arrayBuffer());
};

export const fetchCasObjectAndParse = async <T>(
  instanceName: string | undefined,
  digestFunction: DigestFunction_Value,
  digest: Digest,
  objectType: {
    decode: (input: Uint8Array) => T;
    toJSON: (input: T) => unknown;
  },
): Promise<T> => {
  const bytes = await fetchCasBytes(instanceName, digestFunction, digest);

  return protobufToObject(objectType, bytes, true);
};

async function checkDataExists(
  instanceName: string,
  digests: Digest[],
  digestFunction: DigestFunction_Value,
): Promise<boolean> {
  const requestBody: FindMissingBlobsRequest = {
    instanceName: instanceName,
    blobDigests: digests,
    digestFunction: digestFunction,
  };
  const response = await casClient.findMissingBlobs(requestBody);
  return response.missingBlobDigests.length === 0;
}

export function useCheckDataExists(
  instanceName: string,
  digests: Digest[],
  digestFunction: DigestFunction_Value,
  enabled: boolean = true,
) {
  const { data: exists, isLoading } = useQuery({
    queryKey: ["checkDataExists", instanceName, digests, digestFunction],
    queryFn: () => checkDataExists(instanceName, digests, digestFunction),
    enabled: enabled,
  });

  return {
    exists: exists ?? false,
    isLoading: isLoading,
  };
}
