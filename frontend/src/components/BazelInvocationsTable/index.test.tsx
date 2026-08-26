import { beforeEach, describe, expect, it, vi } from "vitest";
import type { PortalFrontendConfiguration } from "@/lib/grpc-client/portal/frontend/frontend";

const mockEnv: Partial<PortalFrontendConfiguration> = {};

vi.mock("@/utils/env", () => ({
  get env() {
    return mockEnv;
  },
}));

// vitest hoists vi.mock, so getTableColumns must be imported after the mock is registered.
const { getTableColumns } = await import("./index");

describe("getTableColumns", () => {
  beforeEach(() => {
    for (const key of Object.keys(mockEnv)) {
      delete (mockEnv as Record<string, unknown>)[key];
    }
  });

  it("does not include the BES Instance column when the feature flag is unset", () => {
    const columns = getTableColumns();
    expect(columns.some((column) => column.key === "instanceName")).toBe(false);
  });

  it("does not include the BES Instance column when the feature flag is explicitly null", () => {
    // This mirrors what the backend actually sends: protojson serializes an
    // unset message-type feature flag as `null`, not as an absent key.
    mockEnv.featureFlags = {
      bes: { columnInstanceNameInvocations: null },
    } as unknown as PortalFrontendConfiguration["featureFlags"];

    const columns = getTableColumns();
    expect(columns.some((column) => column.key === "instanceName")).toBe(false);
  });

  it("includes the BES Instance column when the feature flag is enabled", () => {
    mockEnv.featureFlags = {
      bes: { columnInstanceNameInvocations: {} },
    } as unknown as PortalFrontendConfiguration["featureFlags"];

    const columns = getTableColumns();
    const instanceNameColumn = columns.find(
      (column) => column.key === "instanceName",
    );
    expect(instanceNameColumn).toBeDefined();
    expect(instanceNameColumn?.title).toBe("BES Instance");
  });
});
