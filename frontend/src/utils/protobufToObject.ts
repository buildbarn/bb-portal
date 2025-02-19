import { RequestMetadata } from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import { AuthenticationMetadata } from "@/lib/grpc-client/buildbarn/auth/auth";
import {
  FilePoolResourceUsage,
  InputRootResourceUsage,
  MonetaryResourceUsage,
  POSIXResourceUsage,
} from "@/lib/grpc-client/buildbarn/resourceusage/resourceusage";
import type { Any } from "@/lib/grpc-client/google/protobuf/any";
import { ProtobufTypeUrls } from "@/types/protobufTypeUrls";

export function protobufToObject<T>(
  objectType: {
    decode: (input: Uint8Array) => T;
    toJSON: (input: T) => unknown;
  },
  protobuf: Uint8Array,
  keepDefaultValues: boolean,
): T {
  if (keepDefaultValues) {
    return objectType.decode(protobuf);
  }
  return objectType.toJSON(objectType.decode(protobuf)) as T;
}

export function protobufToObjectWithTypeField(
  protobuf: Any,
  keepDefaultValues: boolean,
): unknown {
  const typeUrl = protobuf.typeUrl;
  const value = protobuf.value;

  switch (typeUrl) {
    case ProtobufTypeUrls.AUTHENTICATION_METADATA:
      return {
        "@type": typeUrl,
        ...protobufToObject(AuthenticationMetadata, value, keepDefaultValues),
      };
    case ProtobufTypeUrls.REQUEST_METADATA:
      return {
        "@type": typeUrl,
        ...protobufToObject(RequestMetadata, value, keepDefaultValues),
      };
    case ProtobufTypeUrls.POSIX_RESOURCE_USAGE:
      return {
        "@type": typeUrl,
        ...protobufToObject(POSIXResourceUsage, value, keepDefaultValues),
      };
    case ProtobufTypeUrls.FILE_POOL_RESOURCE_USAGE:
      return {
        "@type": typeUrl,
        ...protobufToObject(FilePoolResourceUsage, value, keepDefaultValues),
      };
    case ProtobufTypeUrls.INPUT_ROOT_RESOURCE_USAGE:
      return {
        "@type": typeUrl,
        ...protobufToObject(InputRootResourceUsage, value, keepDefaultValues),
      };
    case ProtobufTypeUrls.MONETARY_RESOURCE_USAGE:
      return {
        "@type": typeUrl,
        ...protobufToObject(MonetaryResourceUsage, value, keepDefaultValues),
      };
    default:
      console.error(`Unknown typeUrl: ${typeUrl}`);
      return {};
  }
}
