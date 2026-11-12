import { describe, expect, it } from "vitest";
import {
  resolveScheduler,
  schedulerSearch,
  schedulerSwitchTarget,
} from "./registry";

const registry = {
  schedulers: [
    { id: "a", label: "A" },
    { id: "b", label: "B" },
  ],
  defaultSchedulerId: "a",
};

describe("scheduler selection", () => {
  it("resolves old links to the configured default and keeps explicit selections", () => {
    expect(resolveScheduler(registry).id).toBe("a");
    expect(resolveScheduler(registry, "b").id).toBe("b");
  });
  it("does not fall back for unknown, empty or malformed explicit selections", () => {
    for (const scheduler of [
      "unknown",
      "",
      "../a",
      schedulerSearch({ scheduler: 123 }).scheduler,
    ]) {
      expect(() => resolveScheduler(registry, scheduler)).toThrow();
    }
  });
  it("preserves malformed explicit input for recoverable selection errors", () => {
    for (const scheduler of [["a", "b"], 123, null, { id: "a" }]) {
      expect(schedulerSearch({ scheduler }).scheduler).toEqual(scheduler);
      expect(() => resolveScheduler(registry, scheduler)).toThrow();
    }
  });
  it("reports no configured scheduler", () => {
    expect(() =>
      resolveScheduler({ schedulers: [], defaultSchedulerId: "" }),
    ).toThrow("No schedulers are configured");
  });
  it("returns worker details to the queue list when switching schedulers", () => {
    expect(schedulerSwitchTarget("/scheduler/worker", "b")).toEqual({
      to: "/scheduler",
      search: { scheduler: "b" },
    });
  });
});
