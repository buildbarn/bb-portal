import { ClientError, Status as GrpcStatus } from "nice-grpc-common";
import { expect, test } from "vitest";
import { Code } from "@/lib/grpc-client/google/rpc/code";
import {
  grpcErrorText,
  isNotFoundGrpcError,
  isRetryableGrpcError,
} from "./grpcStatus";

const method =
  "/buildbarn.buildqueuestate.BuildQueueState/ListQueuedOperations";

test("isNotFoundGrpcError", () => {
  expect(
    isNotFoundGrpcError(
      new ClientError(method, GrpcStatus.NOT_FOUND, "Invocation not found"),
    ),
  ).toBe(true);
  expect(
    isNotFoundGrpcError(
      new ClientError(method, GrpcStatus.PERMISSION_DENIED, "Not allowed"),
    ),
  ).toBe(false);
  expect(isNotFoundGrpcError(new Error("Failed to fetch"))).toBe(false);
});

test("isRetryableGrpcError", () => {
  expect(
    isRetryableGrpcError({ code: Code.UNAVAILABLE, message: "", details: [] }),
  ).toBe(true);
  expect(
    isRetryableGrpcError({ code: Code.NOT_FOUND, message: "", details: [] }),
  ).toBe(false);
  expect(
    isRetryableGrpcError({ code: Code.OK, message: "", details: [] }),
  ).toBe(false);
});

test("grpcErrorText", () => {
  expect(
    grpcErrorText(
      new ClientError(method, GrpcStatus.PERMISSION_DENIED, "Not allowed"),
    ),
  ).toBe("PERMISSION_DENIED: Not allowed");
  expect(grpcErrorText(new Error("Failed to fetch"))).toBe("Failed to fetch");
  expect(grpcErrorText(new Error(""))).toBe("Unknown error");
});
