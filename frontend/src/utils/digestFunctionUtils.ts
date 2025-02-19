import {
  DigestFunction_Value, digestFunction_ValueFromJSON,
  digestFunction_ValueToJSON
} from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";

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
