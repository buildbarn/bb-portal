import { describe, expect, it } from "vitest";
import { generateActionUrlFromGraphqlDigest } from "./urlGenerator";

describe("generateActionUrlFromGraphqlDigest", () => {
  it("generates a bb-browser Action Cache URL", () => {
    expect(
      generateActionUrlFromGraphqlDigest({
        rev2InstanceName: "projects/example",
        digestFunction: "SHA256",
        hash: "abcdef",
        sizeBytes: 145,
      }),
    ).toBe("/browser/projects/example/blobs/sha256/action/abcdef-145");
  });

  it("returns no URL without an Action digest", () => {
    expect(generateActionUrlFromGraphqlDigest(undefined)).toBeUndefined();
  });
});
