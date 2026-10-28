import { ClientError, Status as GrpcStatus } from "nice-grpc-common";
import { expect, test } from "vitest";
import { Code } from "@/lib/grpc-client/google/rpc/code";
import { grpcErrorText, isRetryableGrpcError } from "./grpcStatus";

const method = "/buildbarn.buildqueuestate.BuildQueueState/AddDrain";

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
