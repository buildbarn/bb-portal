import { describe, expect, it } from "vitest";
import type { FileDetailsFragment } from "@/graphql/__generated__/graphql";
import {
  createPerfettoTraceMessage,
  waitForPerfettoToLoad,
} from "./index";

describe("createPerfettoTraceMessage", () => {
  it("sends the file path as the Perfetto file name", () => {
    const profile = {
      id: "file-id",
      filePath: { id: "path-id", path: "profiles/profile.pftrace" },
      digest: {
        id: "digest-id",
        rev2InstanceName: "",
        digestFunction: "SHA256",
        hash: "hash",
        sizeBytes: 123,
      },
    } satisfies FileDetailsFragment;
    const buffer = new ArrayBuffer(4);

    expect(createPerfettoTraceMessage(buffer, profile, "invocation-id")).toEqual(
      {
        perfetto: {
          buffer,
          title: "invocation-id",
          fileName: "profiles/profile.pftrace",
        },
      },
    );
  });
});

describe("waitForPerfettoToLoad", () => {
  it("ignores PONG messages from another window", async () => {
    const handle = {} as Window;
    const otherWindow = {} as Window;
    const waitPromise = waitForPerfettoToLoad(handle);
    let resolved = false;
    waitPromise.then(() => {
      resolved = true;
    });

    const dispatchPong = (source: Window) => {
      const event = new MessageEvent("message", { data: "PONG" });
      Object.defineProperty(event, "source", { value: source });
      window.dispatchEvent(event);
    };

    dispatchPong(otherWindow);
    await Promise.resolve();
    expect(resolved).toBe(false);

    dispatchPong(handle);
    await waitPromise;
  });
});
