import { describe, expect, it } from "vitest";
import { getActionCacheStatus } from "./cache";

describe("getActionCacheStatus", () => {
  it("identifies remote cache hits", () => {
    expect(getActionCacheStatus(true, "remote cache hit")).toMatchObject({
      label: "Remote hit",
      color: "success",
    });
  });

  it("identifies disk cache hits", () => {
    expect(getActionCacheStatus(true, "disk cache hit")).toMatchObject({
      label: "Disk hit",
      color: "success",
    });
  });

  it("preserves cache hits from unrecognized runners", () => {
    expect(getActionCacheStatus(true, "custom cache")).toMatchObject({
      label: "Hit",
      color: "success",
    });
  });

  it("distinguishes a reported non-hit from missing metadata", () => {
    expect(getActionCacheStatus(false, "remote").label).toBe("No hit");
    expect(getActionCacheStatus(undefined, "remote").label).toBe("Unknown");
  });
});
