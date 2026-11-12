import { createChannel } from "nice-grpc-web";
import { describe, expect, it, vi } from "vitest";
import { getBuildQueueStateClient } from "./buildQueueStateClient";

vi.mock("nice-grpc-web", () => ({
  createChannel: vi.fn((url: string) => ({ url })),
  createClient: vi.fn((_definition: unknown, channel: unknown) => ({
    channel,
  })),
}));

describe("scheduler gRPC-Web clients", () => {
  it("uses the scoped endpoint and reuses only the matching scheduler client", () => {
    const a = getBuildQueueStateClient("client-test-a");
    const b = getBuildQueueStateClient("client-test-b");
    expect(a).not.toBe(b);
    expect(getBuildQueueStateClient("client-test-a")).toBe(a);
    expect(vi.mocked(createChannel).mock.calls.map(([url]) => url)).toEqual([
      "/api/v1/schedulers/client-test-a/grpcweb",
      "/api/v1/schedulers/client-test-b/grpcweb",
    ]);
  });
});
