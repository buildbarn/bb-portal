import { ClientError, Status as GrpcStatus } from "nice-grpc-common";
import { Code } from "@/lib/grpc-client/google/rpc/code";
import { Status } from "@/lib/grpc-client/google/rpc/status";

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

export const isNotFoundGrpcError = (error: Error): boolean =>
  Status.fromJSON(error).code === Code.NOT_FOUND;

export const grpcErrorText = (error: Error): string =>
  error instanceof ClientError
    ? `${GrpcStatus[error.code]}: ${error.details}`
    : error.message || "Unknown error";
