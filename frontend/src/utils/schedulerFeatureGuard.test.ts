import { afterEach, expect, test, vi } from "vitest";

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

  const loadDirectURL = () => beforeLoad({} as never);
  if (disabled) {
    expect(loadDirectURL).toThrow("FEATURE_DISABLED");
  } else {
    expect(loadDirectURL).not.toThrow();
  }
});
