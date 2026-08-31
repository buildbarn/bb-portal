import { beforeEach, describe, expect, it, vi } from "vitest";
import type { PortalFrontendConfiguration } from "@/lib/grpc-client/portal/frontend/frontend";

const mockEnv: Partial<PortalFrontendConfiguration> = {};

vi.mock("@/utils/env", () => ({
  get env() {
    return mockEnv;
  },
}));

// vitest hoists vi.mock, so getColumns must be imported after the mock is registered.
const { getColumns } = await import("./Columns");

describe("getColumns", () => {
  beforeEach(() => {
    for (const key of Object.keys(mockEnv)) {
      delete (mockEnv as Record<string, unknown>)[key];
    }
    mockEnv.additionalBuildColumns = [];
  });

  it("includes the BES Instance column", () => {
    const columns = getColumns();
    const instanceNameColumn = columns.find(
      (column) => column.key === "instanceName",
    );
    expect(instanceNameColumn).toBeDefined();
    expect(instanceNameColumn?.title).toBe("BES Instance");
  });

  it("filters builds by BES instance name", () => {
    const columns = getColumns();
    const instanceNameColumn = columns.find(
      (column) => column.key === "instanceName",
    );
    expect(instanceNameColumn?.applyFilter?.(["my-instance"])).toEqual([
      { hasInstanceNameWith: [{ nameContainsFold: "my-instance" }] },
    ]);
    expect(instanceNameColumn?.applyFilter?.([])).toBeUndefined();
  });
});
