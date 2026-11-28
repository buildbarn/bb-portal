import { act } from "react";
import { createRoot } from "react-dom/client";
import { renderToStaticMarkup } from "react-dom/server";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { useSchedulerContext } from "@/schedulers/context";
import { SchedulerSelectionError } from "@/schedulers/registry";
import { FeatureDisabledError } from "@/utils/featureGuard";
import { SchedulerError, SchedulerLayout } from "./index";

const { navigate } = vi.hoisted(() => ({ navigate: vi.fn() }));

vi.mock("@/schedulers/context", () => ({ useSchedulerContext: vi.fn() }));
vi.mock("@tanstack/react-router", () => ({
  Outlet: () => <div>Scheduler view</div>,
  Link: ({ children }: { children: React.ReactNode }) => (
    <a href="/">{children}</a>
  ),
  useNavigate: () => navigate,
  useRouterState: () => "/operations/op-1",
}));
vi.mock("antd", async (importOriginal) => {
  const actual = await importOriginal<typeof import("antd")>();
  const React = await import("react");
  return {
    ...actual,
    Select: (props: {
      "aria-label"?: string;
      onChange?: (value: string) => void;
      options?: { value: string; label: string }[];
    }) =>
      React.createElement(
        "button",
        {
          "aria-label": props["aria-label"],
          onClick: () => {
            const selected = props.options?.[0];
            if (selected) props.onChange?.(selected.value);
          },
        },
        props.options?.[0]?.label,
      ),
  };
});
beforeEach(() => vi.mocked(useSchedulerContext).mockReturnValue({}));
afterEach(() => vi.restoreAllMocks());

const aws = { instanceNamePrefix: "ci/aws" };
const gcp = { instanceNamePrefix: "ci/gcp" };

describe("scheduler selector and errors", () => {
  it("renders single scheduler views without a selector", () => {
    vi.mocked(useSchedulerContext).mockReturnValue({
      scheduler: aws,
      schedulerRegistry: { schedulers: [aws] },
    });
    const html = renderToStaticMarkup(<SchedulerLayout />);
    expect(html).toContain("Scheduler view");
    expect(html).not.toContain('aria-label="Scheduler"');
  });
  it("renders a selector for multiple schedulers", () => {
    vi.mocked(useSchedulerContext).mockReturnValue({
      scheduler: gcp,
      schedulerRegistry: { schedulers: [aws, gcp] },
    });
    expect(renderToStaticMarkup(<SchedulerLayout />)).toContain(
      'aria-label="Scheduler"',
    );
  });
  it("asks to choose a scheduler, not reporting an error, when none is selected", () => {
    const error = new SchedulerSelectionError(
      "Select a scheduler",
      { schedulers: [aws, gcp] },
      "required",
    );
    const html = renderToStaticMarkup(<SchedulerError error={error} />);
    expect(html).toContain("Choose a scheduler");
    expect(html).not.toContain("Unable to open scheduler");
    expect(html).toContain('aria-label="Scheduler"');
  });
  it("preserves the disabled-feature page", () => {
    const html = renderToStaticMarkup(
      <SchedulerError error={new FeatureDisabledError()} />,
    );
    expect(html).toContain("This page is disabled.");
    expect(html).not.toContain("Unable to open scheduler");
  });
  it("lets users switch away from a scheduler whose operations loader failed", () => {
    vi.mocked(useSchedulerContext).mockReturnValue({
      scheduler: aws,
      schedulerRegistry: { schedulers: [aws, gcp] },
    });
    const html = renderToStaticMarkup(
      <SchedulerError
        error={new Error("dial tcp 10.23.45.67:4888 token=private-token")}
      />,
    );
    expect(html).toContain("Unable to load this scheduler");
    expect(html).not.toContain("10.23.45.67");
    expect(html).not.toContain("private-token");
    expect(html).toContain('aria-label="Scheduler"');
  });
  it("offers recovery for unknown selections even with only one valid scheduler", () => {
    const error = new SchedulerSelectionError(
      "Unknown scheduler: missing",
      { schedulers: [aws] },
      "unknown",
    );
    const html = renderToStaticMarkup(<SchedulerError error={error} />);
    expect(html).toContain("The selected scheduler is not configured.");
    expect(html).toContain('aria-label="Scheduler"');
  });
  it("reports when no scheduler is configured", () => {
    const error = new SchedulerSelectionError(
      "No schedulers are configured.",
      { schedulers: [] },
      "none-configured",
    );
    const html = renderToStaticMarkup(<SchedulerError error={error} />);
    expect(html).toContain("No schedulers are configured.");
    expect(html).not.toContain('aria-label="Scheduler"');
  });
  it("does not offer a selector for generic single-scheduler errors", () => {
    vi.mocked(useSchedulerContext).mockReturnValue({
      schedulerRegistry: { schedulers: [aws] },
    });
    const html = renderToStaticMarkup(
      <SchedulerError error={new Error("scheduler unavailable")} />,
    );
    expect(html).toContain("Unable to load this scheduler. Please try again.");
    expect(html).not.toContain("choose another scheduler");
    expect(html).not.toContain('aria-label="Scheduler"');
  });
  it("recovers from an unknown scheduler when one valid scheduler exists", async () => {
    Object.assign(globalThis, { IS_REACT_ACT_ENVIRONMENT: true });
    navigate.mockClear();
    const error = new SchedulerSelectionError(
      "Unknown scheduler: missing",
      { schedulers: [aws] },
      "unknown",
    );
    const container = document.createElement("div");
    const root = createRoot(container);
    try {
      await act(async () => root.render(<SchedulerError error={error} />));
      await act(async () => {
        container
          .querySelector("button")
          ?.dispatchEvent(new MouseEvent("click", { bubbles: true }));
      });
      expect(navigate).toHaveBeenCalledWith({
        to: "/operations",
        search: { scheduler: "ci/aws" },
      });
    } finally {
      await act(async () => root.unmount());
    }
  });
  it("preserves a deep link and its search when choosing a required scheduler", async () => {
    Object.assign(globalThis, { IS_REACT_ACT_ENVIRONMENT: true });
    navigate.mockClear();
    const error = new SchedulerSelectionError(
      "Select a scheduler",
      { schedulers: [aws, gcp] },
      "required",
    );
    const container = document.createElement("div");
    const root = createRoot(container);
    try {
      await act(async () => root.render(<SchedulerError error={error} />));
      await act(async () => {
        container
          .querySelector("button")
          ?.dispatchEvent(new MouseEvent("click", { bubbles: true }));
      });
      expect(navigate).toHaveBeenCalledTimes(1);
      const target = navigate.mock.calls[0][0] as {
        to: string;
        search: (previous: Record<string, unknown>) => Record<string, unknown>;
      };
      expect(target.to).toBe("/operations/op-1");
      expect(
        target.search({ filter: { toolInvocationId: "inv-1" }, page: "2" }),
      ).toEqual({
        filter: { toolInvocationId: "inv-1" },
        page: "2",
        scheduler: "ci/aws",
      });
    } finally {
      await act(async () => root.unmount());
    }
  });
});
