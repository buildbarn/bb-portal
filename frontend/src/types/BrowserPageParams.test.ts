import { describe, expect, it } from "vitest";
import {
  parseBrowserPageSlug,
  validateBrowserPageSlug,
} from "./BrowserPageParams";

const validHash =
  "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef";

describe("validateBrowserPageSlug", () => {
  const validate = (path: string) => validateBrowserPageSlug(path.split("/"));

  it("should accept a valid slug", () => {
    expect(validate(`fuse/blobs/sha256/action/${validHash}-308`)).toEqual({
      success: true,
      data: {
        instanceName: "fuse",
        digestFunction: 1,
        browserPageType: "action",
        digest: { hash: validHash, sizeBytes: "308" },
      },
    });
  });

  it("should accept instance names containing slashes", () => {
    const result = validate(`a/b/c/blobs/sha256/tree/${validHash}-308/sub`);
    expect(result.success && result.data.instanceName).toBe("a/b/c");
  });

  it.each([
    [
      "fuse",
      'Missing "blobs" segment',
      "/browser/fuse/blobs/<digest function>/<page type>/<hash>-<sizeBytes>",
    ],
    [
      "a/b",
      'Missing "blobs" segment',
      "/browser/a/b/blobs/<digest function>/<page type>/<hash>-<sizeBytes>",
    ],
    [
      "fuse/blobs",
      "Missing digest function",
      "/browser/fuse/blobs/<digest function>/<page type>/<hash>-<sizeBytes>",
    ],
    [
      "fuse/blobs/",
      "Missing digest function",
      "/browser/fuse/blobs/<digest function>/<page type>/<hash>-<sizeBytes>",
    ],
    [
      "fuse/blobs/sha256",
      "Missing page type",
      "/browser/fuse/blobs/sha256/<page type>/<hash>-<sizeBytes>",
    ],
    [
      "a/b/blobs/sha256/action",
      "Missing digest",
      "/browser/a/b/blobs/sha256/action/<hash>-<sizeBytes>",
    ],
  ])("should report %s as incomplete", (path, error, suggestion) => {
    expect(validate(path)).toEqual({
      success: false,
      incomplete: true,
      error,
      suggestion,
    });
  });

  it.each([
    [
      "fuse/blobs/xyz/action",
      expect.stringMatching(
        /^Unknown digest function "xyz", expected one of: sha256, sha1, md5, .*blake3/,
      ),
    ],
    [
      "fuse/blobs/sha256/xyz",
      'Unknown page type "xyz", expected one of: action, command, directory, file, historical_execute_response, previous_execution_stats, tree',
    ],
    [
      "fuse/blobs/sha256/action/abc",
      'Malformed digest "abc", expected "<hash>-<sizeBytes>"',
    ],
    [
      "fuse/blobs/sha256/action/abc-1",
      'Malformed digest "abc-1": hash must only contain hex digits and resemble a hash in length',
    ],
    [
      `fuse/blobs/sha256/action/${validHash}-x`,
      `Malformed digest "${validHash}-x": sizeBytes must only contain digits`,
    ],
  ])("should reject %s", (path, error) => {
    expect(validate(path)).toEqual({
      success: false,
      incomplete: false,
      error,
    });
  });
});

describe("parseBrowserPageSlug", () => {
  it('should return undefined if "blobs" is not in the slug', () => {
    const result = parseBrowserPageSlug([
      "instance",
      "not-blobs",
      "sha256",
      "action",
      "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef-308",
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
      "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef-308",
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
      "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef-308",
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
      "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef-308",
    ]);
    expect(result).toEqual({
      instanceName: "instance",
      digestFunction: 1,
      browserPageType: "action",
      digest: {
        hash: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
        sizeBytes: "308",
      },
    });
  });

  it("should parse valid slug without instance name", () => {
    const result = parseBrowserPageSlug([
      "blobs",
      "sha256",
      "action",
      "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef-308",
    ]);
    expect(result).toEqual({
      instanceName: "",
      digestFunction: 1,
      browserPageType: "action",
      digest: {
        hash: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
        sizeBytes: "308",
      },
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
      "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef-308",
    ]);
    expect(result).toEqual({
      instanceName: "instance/name/with/slashes",
      digestFunction: 1,
      browserPageType: "action",
      digest: {
        hash: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
        sizeBytes: "308",
      },
    });
  });
});
