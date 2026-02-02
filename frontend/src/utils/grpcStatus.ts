import { Code } from "@/lib/grpc-client/google/rpc/code";
import type { Status } from "@/lib/grpc-client/google/rpc/status";

export const isRetryableGrpcError = (status: Status): boolean => {
  switch (status.code) {
    case Code.CANCELLED:
    case Code.UNKNOWN:
    case Code.DEADLINE_EXCEEDED:
    case Code.RESOURCE_EXHAUSTED:
    case Code.ABORTED:
    case Code.INTERNAL:
    case Code.UNAVAILABLE:
      return true;
    default:
      return false;
  }
};
