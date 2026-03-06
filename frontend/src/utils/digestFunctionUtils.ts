import {
  Action,
  Digest,
  DigestFunction_Value,
  type Platform,
  digestFunction_ValueFromJSON,
  digestFunction_ValueToJSON,
} from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";

export async function sha256hex(data: BufferSource): Promise<string> {
  const hashBuffer = await window.crypto.subtle.digest('SHA-256', data);
  const hashArray = Array.from(new Uint8Array(hashBuffer));
  return hashArray.map((b) => b.toString(16).padStart(2, '0')).join('');
}

export const digestFunctionValueFromString = (
  string: string,
): DigestFunction_Value => {
  return digestFunction_ValueFromJSON(string.toUpperCase());
};

export const digestFunctionValueToString = (
  value: DigestFunction_Value,
): string | undefined => {
  return digestFunction_ValueToJSON(value).toLowerCase();
};

export const includeDigestFunctionInCasFetch = (
  digestFunction: DigestFunction_Value,
): boolean => {
  return ![
    DigestFunction_Value.MD5,
    DigestFunction_Value.MURMUR3,
    DigestFunction_Value.SHA1,
    DigestFunction_Value.SHA256,
    DigestFunction_Value.SHA384,
    DigestFunction_Value.SHA512,
    DigestFunction_Value.VSO,
  ].includes(digestFunction);
};

// Currently we only support SHA256, as some of the other
// algorithms are difficult to implement in Node.
// TODO: Handle different types of algorithms.
export const getReducedActionDigest_SHA256 = async (
  commandDigest: Digest,
  platform: Platform,
): Promise<Digest> => {
  const encodedReducedAction = Action.encode(
    Action.fromPartial({
      commandDigest: commandDigest,
      platform: platform,
    }),
  ).finish();

  return Digest.create({
    hash: await sha256hex(encodedReducedAction),
    sizeBytes: encodedReducedAction.length.toString(),
  });
};
