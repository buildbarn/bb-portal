import { describe, expect, test } from "vitest";
import {
  type ArtifactFileNode,
  buildSetIndex,
  countSetFiles,
  type NamedSetNode,
  walkSetFiles,
} from "./graph";

function file(name: string): ArtifactFileNode {
  return { name };
}

describe("walkSetFiles", () => {
  test("walks nested named sets and yields every reachable file", () => {
    const sets: NamedSetNode[] = [
      { id: "root", files: [file("a")], childSetIds: ["child"] },
      { id: "child", files: [file("b"), file("c")], childSetIds: [] },
    ];
    const index = buildSetIndex(sets);
    const names = Array.from(walkSetFiles(index, ["root"])).map((f) => f.name);
    expect(names.sort()).toEqual(["a", "b", "c"]);
  });

  test("terminates on cycles", () => {
    const sets: NamedSetNode[] = [
      { id: "a", files: [file("fa")], childSetIds: ["b"] },
      { id: "b", files: [file("fb")], childSetIds: ["a"] },
    ];
    const index = buildSetIndex(sets);
    const names = Array.from(walkSetFiles(index, ["a"])).map((f) => f.name);
    expect(names.sort()).toEqual(["fa", "fb"]);
  });

  test("ignores unknown set ids", () => {
    const index = buildSetIndex([]);
    expect(Array.from(walkSetFiles(index, ["missing"]))).toEqual([]);
  });
});

describe("countSetFiles", () => {
  test("counts files reachable from the roots", () => {
    const sets: NamedSetNode[] = [
      { id: "root", files: [file("a"), file("b")], childSetIds: ["child"] },
      { id: "child", files: [file("c")], childSetIds: [] },
    ];
    const index = buildSetIndex(sets);
    expect(countSetFiles(index, ["root"])).toBe(3);
  });
});
