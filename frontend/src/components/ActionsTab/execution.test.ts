import { describe, expect, it } from "vitest";
import { getActionExecutionKind } from "./execution";

describe("getActionExecutionKind", () => {
  it.each([
    "remote",
    "remote cache hit",
  ])("classifies %s as remote", (runner) => {
    expect(getActionExecutionKind(runner)).toBe("Remote");
  });

  it.each([
    "darwin-sandbox",
    "linux-sandbox",
    "worker",
    "disk cache hit",
  ])("classifies %s as local", (runner) => {
    expect(getActionExecutionKind(runner)).toBe("Local");
  });

  it("classifies non-spawn actions as internal", () => {
    expect(getActionExecutionKind("internal")).toBe("Internal");
  });

  it("classifies actions without execution-log metadata as unknown", () => {
    expect(getActionExecutionKind(undefined)).toBe("Unknown");
  });
});
