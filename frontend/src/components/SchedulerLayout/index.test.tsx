import { renderToStaticMarkup } from "react-dom/server";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { useSchedulerContext } from "@/schedulers/context";
import { SchedulerSelectionError } from "@/schedulers/registry";
import { FeatureDisabledError } from "@/utils/featureGuard";
import { SchedulerError, SchedulerLayout } from "./index";

vi.mock("@/schedulers/context", () => ({ useSchedulerContext: vi.fn() }));
vi.mock("@tanstack/react-router", () => ({
  Outlet: () => <div>Scheduler view</div>,
  Link: ({ children }: { children: React.ReactNode }) => (
    <a href="/">{children}</a>
  ),
  useNavigate: () => vi.fn(),
  useRouterState: () => "/scheduler",
}));
beforeEach(() => vi.mocked(useSchedulerContext).mockReturnValue({}));
afterEach(() => vi.restoreAllMocks());

const a = { id: "a", label: "A" };
const b = { id: "b", label: "B" };

describe("scheduler selector and errors", () => {
  it("renders single scheduler views without a selector", () => {
    vi.mocked(useSchedulerContext).mockReturnValue({
      scheduler: a,
      schedulerRegistry: { schedulers: [a], defaultSchedulerId: "a" },
    });
    const html = renderToStaticMarkup(<SchedulerLayout />);
    expect(html).toContain("Scheduler view");
    expect(html).not.toContain('aria-label="Scheduler"');
  });
  it("renders a selector for multiple schedulers", () => {
    vi.mocked(useSchedulerContext).mockReturnValue({
      scheduler: b,
      schedulerRegistry: { schedulers: [a, b], defaultSchedulerId: "a" },
    });
    expect(renderToStaticMarkup(<SchedulerLayout />)).toContain(
      'aria-label="Scheduler"',
    );
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
      scheduler: a,
      schedulerRegistry: { schedulers: [a, b], defaultSchedulerId: "a" },
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
    const error = new SchedulerSelectionError("Unknown scheduler: missing", {
      schedulers: [a],
      defaultSchedulerId: "a",
    });
    const html = renderToStaticMarkup(<SchedulerError error={error} />);
    expect(html).toContain("The selected scheduler is not configured.");
    expect(html).toContain('aria-label="Scheduler"');
  });
});
