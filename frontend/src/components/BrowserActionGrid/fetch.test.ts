import { describe, expect, it, vi } from "vitest";
import {
  Action,
  type ActionCacheClient,
  Command,
  type Digest,
  DigestFunction_Value,
} from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import type { FileSystemAccessCacheClient } from "@/lib/grpc-client/buildbarn/fsac/fsac";
import type { InitialSizeClassCacheClient } from "@/lib/grpc-client/buildbarn/iscc/iscc";
import type { ByteStreamClient } from "@/lib/grpc-client/google/bytestream/bytestream";
import { BrowserPageType } from "@/types/BrowserPageType";
import { fetchBrowserActionGrid } from "./fetch";

const digest = (hash: string, sizeBytes: string): Digest => ({
  hash,
  sizeBytes,
});

describe("fetchBrowserActionGrid", () => {
  it("renders action data without eagerly fetching an unavailable input root", async () => {
    const actionDigest = digest("action", "10");
    const commandDigest = digest("command", "20");
    const inputRootDigest = digest("missing-input-root", "30");
    const actionBytes = Action.encode(
      Action.fromPartial({ commandDigest, inputRootDigest }),
    ).finish();
    const commandBytes = Command.encode(
      Command.fromPartial({ arguments: ["tool", "--flag"] }),
    ).finish();

    const read = vi.fn(({ resourceName }: { resourceName?: string }) => {
      const data = resourceName?.includes("/action/")
        ? actionBytes
        : resourceName?.includes("/command/")
          ? commandBytes
          : undefined;
      if (!data) {
        throw new Error(`Unexpected CAS fetch: ${resourceName}`);
      }
      return (async function* () {
        yield { data };
      })();
    });

    const result = await fetchBrowserActionGrid(
      {
        instanceName: "",
        digestFunction: DigestFunction_Value.SHA256,
        browserPageType: BrowserPageType.Action,
        digest: actionDigest,
        otherParams: [],
      },
      {
        getActionResult: vi.fn().mockResolvedValue({ exitCode: 0 }),
      } as unknown as ActionCacheClient,
      { read } as unknown as ByteStreamClient,
      {} as InitialSizeClassCacheClient,
      {} as FileSystemAccessCacheClient,
    );

    expect(result.action.inputRootDigest).toEqual(inputRootDigest);
    expect(result.casCommand?.arguments).toEqual(["tool", "--flag"]);
    expect(read).toHaveBeenCalledTimes(2);
    expect(
      read.mock.calls.some(([request]) =>
        request.resourceName?.includes("missing-input-root"),
      ),
    ).toBe(false);
  });
});
