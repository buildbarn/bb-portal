import { renderToStaticMarkup } from "react-dom/server";
import { afterEach, describe, expect, it, vi } from "vitest";
import OperationsError from "@/components/OperationsError";
import { useSchedulerContext } from "@/schedulers/context";

vi.mock("@/schedulers/context", () => ({ useSchedulerContext: vi.fn() }));
vi.mock("@/components/PortalAlert", async () => {
  const React = await import("react");
  return {
    default: (props: { title?: string; description?: string }) =>
      React.createElement("div", null, props.title, props.description),
  };
});

afterEach(() => vi.restoreAllMocks());

describe("operations loader errors", () => {
  it("does not suggest switching schedulers in single-scheduler mode", () => {
    vi.mocked(useSchedulerContext).mockReturnValue({
      schedulerRegistry: { schedulers: [{ instanceNamePrefix: "ci/aws" }] },
    });
    const html = renderToStaticMarkup(
      <OperationsError error={new Error("scheduler unavailable")} />,
    );
    expect(html).toContain("Unable to load operations");
    expect(html).toContain("Please try again.");
    expect(html).not.toContain("choose another scheduler");
  });

  it("keeps the selector guidance when multiple schedulers are configured", () => {
    vi.mocked(useSchedulerContext).mockReturnValue({
      schedulerRegistry: {
        schedulers: [
          { instanceNamePrefix: "ci/aws" },
          { instanceNamePrefix: "ci/gcp" },
        ],
      },
    });
    const html = renderToStaticMarkup(
      <OperationsError error={new Error("scheduler unavailable")} />,
    );
    expect(html).toContain("choose another scheduler using the selector above");
  });
});
