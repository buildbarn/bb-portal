import type {
  Digest,
  DigestFunction_Value,
} from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import type { ByteStreamClient } from "@/lib/grpc-client/google/bytestream/bytestream";
import {
  digestFunctionValueToString,
  includeDigestFunctionInCasFetch,
} from "@/utils/digestFunctionUtils";
import { protobufToObject } from "@/utils/protobufToObject";

export const fetchCasObject = async (
  casByteStreamClient: ByteStreamClient,
  instanceName: string | undefined,
  digestFunction: DigestFunction_Value,
  digest: Digest,
): Promise<Uint8Array> => {
  const resourceName = [
    instanceName ? `/${instanceName}` : "",
    "/blobs",
    includeDigestFunctionInCasFetch(digestFunction)
      ? `/${digestFunctionValueToString(digestFunction)}`
      : "",
    `/${digest.hash}/${digest.sizeBytes}`,
  ].join("");

  const responseStream = casByteStreamClient.read({
    resourceName,
    readOffset: "0",
    readLimit: "0",
  });

  const chunks: Uint8Array[] = [];
  for await (const chunk of responseStream) {
    chunks.push(chunk.data);
  }

  return new Uint8Array(
    chunks.reduce(
      (acc: number[], chunk) => acc.concat(Array.from(chunk)),
      [] as number[],
    ),
  );
};

export const fetchCasObjectAndParse = async <T>(
  casByteStreamClient: ByteStreamClient,
  instanceName: string | undefined,
  digestFunction: DigestFunction_Value,
  digest: Digest,
  objectType: {
    decode: (input: Uint8Array) => T;
    toJSON: (input: T) => unknown;
  },
): Promise<T> => {
  const combinedChunks = await fetchCasObject(
    casByteStreamClient,
    instanceName,
    digestFunction,
    digest,
  );

  return protobufToObject(objectType, combinedChunks, true);
};
