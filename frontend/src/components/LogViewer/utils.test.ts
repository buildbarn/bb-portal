import { describe, expect, it } from "vitest";
import { findMatchIndices } from "./utils";

const LINES = ["first line", "second line", "third line with line twice"];

describe("findMatchIndices", () => {
  it("finds no matches for an empty query", () => {
    expect(findMatchIndices(LINES, "", 10000)).toStrictEqual([]);
  });

  it("returns the same empty result for repeated empty queries", () => {
    // A new array each call would re-run every effect keyed on the match list.
    expect(findMatchIndices(LINES, "", 10000)).toBe(
      findMatchIndices([], "", 10000),
    );
  });

  it("reports one entry per occurrence", () => {
    expect(findMatchIndices(LINES, "line", 10000)).toStrictEqual([0, 1, 2, 2]);
  });

  it("matches case insensitively", () => {
    expect(findMatchIndices(LINES, "FIRST", 10000)).toStrictEqual([0]);
  });

  it("treats the query as plain text", () => {
    expect(findMatchIndices(["a+b", "aab"], "a+b", 10000)).toStrictEqual([0]);
  });

  it("ignores ansi escape codes in the log", () => {
    expect(
      findMatchIndices(["\x1B[31merror\x1B[0m: failed"], "error:", 10000),
    ).toStrictEqual([0]);
  });

  it("stops at the match limit", () => {
    expect(findMatchIndices(LINES, "line", 2)).toStrictEqual([0, 1]);
  });
});
