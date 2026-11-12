import { afterEach, expect, test, vi } from "vitest";

// The guard resolves the scheduler registry once the feature check passes.
// Stub it so the enabled case exercises only the feature flag, not the network.
vi.mock("@/schedulers/context", () => ({
  schedulerContext: vi.fn().mockResolvedValue({}),
}));

afterEach(() => {
  Reflect.deleteProperty(window, "__env__");
  vi.resetModules();
});

test.each([
  {
    name: "null scheduler flag from backend JSON",
    configuration: '{"featureFlags":{"scheduler":null}}',
    disabled: true,
  },
  {
    name: "omitted scheduler flag",
    configuration: '{"featureFlags":{}}',
    disabled: true,
  },
  {
    name: "enabled scheduler flag",
    configuration: '{"featureFlags":{"scheduler":{}}}',
    disabled: false,
  },
])("scheduler route handles $name", async ({ configuration, disabled }) => {
  // Match the runtime configuration injected into the HTML by the backend.
  Object.assign(window, { __env__: JSON.parse(configuration) });
  const { Route } = await import("@/routes/scheduler");
  const beforeLoad = Route.options.beforeLoad;
  if (typeof beforeLoad !== "function") {
    throw new Error("Missing scheduler route guard");
  }

  // The guard is async because it also resolves the scheduler context, so a
  // disabled feature surfaces as a rejection rather than a synchronous throw.
  const loadDirectURL = async () => beforeLoad({ search: {} } as never);
  if (disabled) {
    await expect(loadDirectURL()).rejects.toThrow("FEATURE_DISABLED");
  } else {
    await expect(loadDirectURL()).resolves.toBeDefined();
  }
});
