import { describe, expect, it } from "vitest";
import { generateFileUrlFromBepURI } from "./urlGenerator";

describe("generateFileUrlFromBepURI", () => {
  it("generates a file URL from a SHA-256 bytestream URI", () => {
    expect(
      generateFileUrlFromBepURI(
        "bytestream://cache.example.com/blobs/abcdef/42",
        "bazel-out/bin/output.txt",
      ),
    ).toBe("/api/v1/servefile/blobs/sha256/file/abcdef-42/output.txt");
  });

  it("preserves the instance name and explicit digest function", () => {
    expect(
      generateFileUrlFromBepURI(
        "bytestream://cache.example.com/projects/example/blobs/sha1/abcdef/42",
        "output.txt",
      ),
    ).toBe(
      "/api/v1/servefile/projects/example/blobs/sha1/file/abcdef-42/output.txt",
    );
  });

  it("does not generate links for local files", () => {
    expect(
      generateFileUrlFromBepURI("file:///tmp/output.txt", "output.txt"),
    ).toBeUndefined();
  });
});
