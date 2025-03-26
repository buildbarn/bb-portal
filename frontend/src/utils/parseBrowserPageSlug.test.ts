import { describe, expect, it } from "vitest";
import { parseBrowserPageSlug } from "./parseBrowserPageSlug";

describe("parseBrowserPageSlug", () => {
  it('should return undefined if "blobs" is not in the slug', () => {
    const result = parseBrowserPageSlug([
      "instance",
      "not-blobs",
      "sha256",
      "action",
      "hash-size",
      "other",
      "params",
    ]);
    expect(result).toBeUndefined();
  });

  it('should return undefined if "blobs" is at the end of the slug', () => {
    const result = parseBrowserPageSlug([
      "instance",
      "sha256",
      "action",
      "hash-size",
      "other",
      "params",
      "blobs",
    ]);
    expect(result).toBeUndefined();
  });

  it("should return undefined if browserPageType is undefined", () => {
    const result = parseBrowserPageSlug([
      "instance",
      "blobs",
      "sha256",
      "invalidType",
      "hash-size",
    ]);
    expect(result).toBeUndefined();
  });

  it("should return undefined if digest or sizeBytes is missing", () => {
    const result = parseBrowserPageSlug([
      "instance",
      "blobs",
      "sha256",
      "action",
      "digest",
    ]);
    expect(result).toBeUndefined();
  });

  it("should parse valid slug correctly", () => {
    const result = parseBrowserPageSlug([
      "instance",
      "blobs",
      "sha256",
      "action",
      "hash-size",
      "other",
      "params",
    ]);
    expect(result).toEqual({
      instanceName: "instance",
      digestFunction: 1,
      browserPageType: "action",
      digest: { hash: "hash", sizeBytes: "size" },
      otherParams: ["other", "params"],
    });
  });

  it("should parse valid slug without instance name", () => {
    const result = parseBrowserPageSlug([
      "blobs",
      "sha256",
      "action",
      "hash-size",
    ]);
    expect(result).toEqual({
      instanceName: "",
      digestFunction: 1,
      browserPageType: "action",
      digest: { hash: "hash", sizeBytes: "size" },
      otherParams: [],
    });
  });

  it("should parse valid slug with instance name with slashes correctly", () => {
    const result = parseBrowserPageSlug([
      "instance",
      "name",
      "with",
      "slashes",
      "blobs",
      "sha256",
      "action",
      "hash-size",
    ]);
    expect(result).toEqual({
      instanceName: "instance/name/with/slashes",
      digestFunction: 1,
      browserPageType: "action",
      digest: { hash: "hash", sizeBytes: "size" },
      otherParams: [],
    });
  });
});
