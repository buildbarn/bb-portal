import { createChannel, createClient } from "nice-grpc-web";
import { describe, expect, it, vi } from "vitest";
import {
  getBuildQueueStateClient,
  SCHEDULER_SELECTION_HEADER,
} from "./buildQueueStateClient";

vi.mock("nice-grpc-web", () => ({
  Metadata: (init: Record<string, string>) => ({ init }),
  createChannel: vi.fn((url: string) => ({ url })),
  createClient: vi.fn(
    (_definition: unknown, channel: unknown, defaults: unknown) => ({
      channel,
      defaults,
    }),
  ),
}));

describe("scheduler gRPC-Web clients", () => {
  it("shares one endpoint and scopes each client by its instance name prefix", () => {
    const aws = getBuildQueueStateClient("client-test/aws");
    const root = getBuildQueueStateClient("");
    expect(aws).not.toBe(root);
    expect(getBuildQueueStateClient("client-test/aws")).toBe(aws);
    expect(vi.mocked(createChannel).mock.calls).toEqual([["/api/v1/grpcweb"]]);
    expect(
      vi.mocked(createClient).mock.calls.map(([, , defaults]) => defaults),
    ).toEqual([
      {
        "*": {
          metadata: {
            init: { [SCHEDULER_SELECTION_HEADER]: "client-test/aws" },
          },
        },
      },
      { "*": { metadata: { init: { [SCHEDULER_SELECTION_HEADER]: "" } } } },
    ]);
  });
});
