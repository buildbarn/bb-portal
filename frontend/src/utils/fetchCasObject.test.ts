import { describe, expect, it } from "vitest";
import { DigestFunction_Value } from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import type { ByteStreamClient } from "@/lib/grpc-client/google/bytestream/bytestream";
import { fetchCasObject } from "./fetchCasObject";

const clientReturning = (chunks: Uint8Array[]) =>
  ({
    read: async function* () {
      for (const chunk of chunks) {
        yield { data: chunk };
      }
    },
  }) as unknown as ByteStreamClient;

describe("fetchCasObject", () => {
  it("joins the streamed chunks in order", async () => {
    const data = await fetchCasObject(
      clientReturning(
        [[1, 2], [], [3], [4, 5, 6]].map((c) => new Uint8Array(c)),
      ),
      "fuse",
      DigestFunction_Value.SHA256,
      { hash: "00", sizeBytes: "6" },
    );
    expect(Array.from(data)).toEqual([1, 2, 3, 4, 5, 6]);
  });

  // The old number-array join took 13 s here, past the 5 s test timeout.
  it("joins a large blob", async () => {
    const chunk = Uint8Array.from({ length: 64 * 1024 }, (_, i) => i % 256);
    const data = await fetchCasObject(
      clientReturning(Array.from({ length: 256 }, () => chunk)),
      "",
      DigestFunction_Value.SHA256,
      { hash: "00", sizeBytes: String(256 * chunk.length) },
    );
    expect(data.length).toBe(256 * chunk.length);
    expect(data[chunk.length + 1]).toBe(1);
  });
});
