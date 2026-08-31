import { describe, expect, it } from "vitest";
import { getTableColumns } from "./index";

describe("getTableColumns", () => {
  it("includes the BES Instance column", () => {
    const columns = getTableColumns();
    const instanceNameColumn = columns.find(
      (column) => column.key === "instanceName",
    );
    expect(instanceNameColumn).toBeDefined();
    expect(instanceNameColumn?.title).toBe("BES Instance");
  });
});
