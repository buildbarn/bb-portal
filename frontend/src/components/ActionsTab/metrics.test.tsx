import { renderToStaticMarkup } from "react-dom/server";
import { describe, expect, it } from "vitest";
import { ActionsMetrics } from "./metrics";

const timingMetrics = {
  totalExpectedTimeInMs: 3000,
  timeSavedByCacheHitsInMs: 1250,
  totalActions: 4,
  timedActions: 4,
  cacheHitActions: 1,
  timedCacheHitActions: 1,
};

describe("ActionsMetrics", () => {
  it("renders invocation-wide expected, real, and cache-saved timing", () => {
    const metrics = renderToStaticMarkup(
      <ActionsMetrics
        actionTimingMetrics={timingMetrics}
        executionPhaseTimeInMs={2250}
      />,
    );

    expect(metrics).toContain("Total Expected Time");
    expect(metrics).toContain("Real Execution Time");
    expect(metrics).toContain("Time Saved by Cache Hits");
    expect(metrics).toContain("3.00s");
    expect(metrics).toContain("2.25s");
    expect(metrics).toContain("1.25s");
  });

  it("shows an unavailable real execution time when Bazel did not report one", () => {
    const metrics = renderToStaticMarkup(
      <ActionsMetrics
        actionTimingMetrics={timingMetrics}
        executionPhaseTimeInMs={null}
      />,
    );

    expect(metrics).toContain(">-</span>");
  });
});
