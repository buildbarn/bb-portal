import { describe, expect, it } from "vitest";
import { getColumns } from "./columns";

describe("ActionsTab selectable filters", () => {
  const columns = getColumns(
    ["CppCompile", "GoLink"],
    ["k8-fastbuild", "k8-opt"],
  );

  it("filters by all selected action mnemonics", () => {
    const column = columns.find(({ key }) => key === "type");

    expect(column?.filterSearch).toBe(true);
    expect(column?.filterMultiple).toBe(true);
    expect(column?.filters).toEqual([
      { text: "CppCompile", value: "CppCompile" },
      { text: "GoLink", value: "GoLink" },
    ]);
    expect(column?.applyFilter?.(["CppCompile", "GoLink"])).toEqual([
      { typeIn: ["CppCompile", "GoLink"] },
    ]);
  });

  it("filters by all selected configurations", () => {
    const column = columns.find(({ key }) => key === "configurationMnemonic");

    expect(column?.filterSearch).toBe(true);
    expect(column?.filterMultiple).toBe(true);
    expect(column?.applyFilter?.(["k8-fastbuild", "k8-opt"])).toEqual([
      {
        hasConfigurationWith: [{ mnemonicIn: ["k8-fastbuild", "k8-opt"] }],
      },
    ]);
  });
});
