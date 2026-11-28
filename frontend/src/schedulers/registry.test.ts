import { describe, expect, it } from "vitest";
import {
  resolveScheduler,
  SchedulerSelectionError,
  schedulerLabel,
  schedulerSearch,
  schedulerSwitchTarget,
} from "./registry";

const registry = {
  schedulers: [
    { instanceNamePrefix: "ci/aws" },
    { instanceNamePrefix: "ci/gcp" },
  ],
};

function problem(run: () => unknown) {
  try {
    run();
  } catch (error) {
    if (error instanceof SchedulerSelectionError) return error.problem;
    throw error;
  }
  throw new Error("Expected a scheduler selection error");
}

describe("scheduler selection", () => {
  it("picks the only scheduler without a selection", () => {
    const single = { schedulers: [{ instanceNamePrefix: "" }] };
    expect(resolveScheduler(single).instanceNamePrefix).toBe("");
  });
  it("asks for a selection when several schedulers exist, instead of defaulting", () => {
    expect(problem(() => resolveScheduler(registry))).toBe("required");
  });
  it("keeps explicit selections, including prefixes with slashes", () => {
    expect(resolveScheduler(registry, "ci/gcp").instanceNamePrefix).toBe(
      "ci/gcp",
    );
  });
  it("selects the empty prefix when it is configured", () => {
    const withRoot = {
      schedulers: [{ instanceNamePrefix: "" }, { instanceNamePrefix: "ci" }],
    };
    expect(resolveScheduler(withRoot, "").instanceNamePrefix).toBe("");
  });
  it("does not fall back for unknown, nested or malformed explicit selections", () => {
    for (const scheduler of [
      "unknown",
      "",
      "ci",
      "ci/aws/mars",
      schedulerSearch({ scheduler: 123 }).scheduler,
    ]) {
      expect(problem(() => resolveScheduler(registry, scheduler))).toBe(
        "unknown",
      );
    }
  });
  it("preserves malformed explicit input for recoverable selection errors", () => {
    for (const scheduler of [["ci/aws", "ci/gcp"], 123, null, { id: "a" }]) {
      expect(schedulerSearch({ scheduler }).scheduler).toEqual(scheduler);
      expect(() => resolveScheduler(registry, scheduler)).toThrow();
    }
  });
  it("reports no configured scheduler", () => {
    expect(problem(() => resolveScheduler({ schedulers: [] }))).toBe(
      "none-configured",
    );
  });
  it("labels schedulers by prefix", () => {
    expect(schedulerLabel({ instanceNamePrefix: "ci/aws" })).toBe("ci/aws");
    expect(schedulerLabel({ instanceNamePrefix: "" })).toBe(
      "(all instance names)",
    );
  });
  it("returns worker details to the queue list when switching schedulers", () => {
    expect(schedulerSwitchTarget("/scheduler/worker", "ci/gcp")).toEqual({
      to: "/scheduler",
      search: { scheduler: "ci/gcp" },
    });
  });
});
