import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import {
  createMemoryHistory,
  Outlet,
  RouterProvider,
} from "@tanstack/react-router";
import { act, type ReactNode } from "react";
import { createRoot } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { getBuildQueueStateClient } from "@/grpc/buildQueueStateClient";
import { OperationState } from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import { createPortalRouter } from "@/router";
import { schedulerSearch } from "./registry";
import { parsePortalSearch, stringifyPortalSearch } from "./urlSearch";

const state = vi.hoisted(() => ({
  cache: undefined as QueryClient | undefined,
}));
vi.mock("@/main", () => import("@/utils/NotFoundError"));
vi.mock("@/utils/env", () => ({
  env: { featureFlags: { scheduler: {}, browser: {} }, companyName: "Test" },
}));
vi.mock("@/grpc/buildQueueStateClient", () => ({
  getBuildQueueStateClient: vi.fn(),
}));
vi.mock("@/components/pages/RootLayout", () => ({
  RootLayout: () => {
    if (!state.cache) throw new Error("Missing test cache");
    return (
      <QueryClientProvider client={state.cache}>
        <Outlet />
      </QueryClientProvider>
    );
  },
}));
// Native controls exercise the production onChange/onClick callbacks without
// coupling routing and authorization checks to Ant Design's dropdown markup.
vi.mock("antd", async (importOriginal) => ({
  ...(await importOriginal<typeof import("antd")>()),
  Select: ({
    value,
    options,
    onChange,
  }: {
    value?: string;
    options: { value: string; label: string }[];
    onChange: (value: string) => void;
  }) => (
    <select
      aria-label="Scheduler"
      value={value ?? ""}
      onChange={(event) => onChange(event.currentTarget.value)}
    >
      {options.map((option) => (
        <option key={option.value} value={option.value}>
          {option.label}
        </option>
      ))}
    </select>
  ),
  Button: ({
    children,
    disabled,
    onClick,
  }: {
    children: ReactNode;
    disabled?: boolean;
    onClick?: () => void;
  }) => (
    <button type="button" disabled={disabled} onClick={onClick}>
      {children}
    </button>
  ),
}));

const scalarIds = [
  "1",
  "2024",
  "true",
  "false",
  "null",
  "001",
  "1e3",
  "999999999999999999999999999999999999999999999999999999999999999",
];
const registry = {
  schedulers: ["a", "b", ...scalarIds].map((id) => ({
    id,
    label: id.toUpperCase(),
  })),
  defaultSchedulerId: "a",
};
const operation = (targetId: string) =>
  OperationState.fromPartial({ name: "same-id", targetId });
const clients: Record<
  string,
  {
    listOperations: ReturnType<typeof vi.fn>;
    getOperation: ReturnType<typeof vi.fn>;
    killOperations: ReturnType<typeof vi.fn>;
  }
> = Object.fromEntries(
  registry.schedulers.map(({ id }) => [
    id,
    { listOperations: vi.fn(), getOperation: vi.fn(), killOperations: vi.fn() },
  ]),
);
let permission: (scheduler: string) => Promise<Response>;
const reply = (body: unknown, status = 200) =>
  new Response(JSON.stringify(body), {
    status,
    headers: { "Content-Type": "application/json" },
  });

beforeEach(() => {
  Object.assign(globalThis, { IS_REACT_ACT_ENVIRONMENT: true });
  state.cache = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  vi.stubGlobal("scrollTo", vi.fn());
  vi.stubGlobal(
    "matchMedia",
    vi.fn((media: string) => ({
      matches: false,
      media,
      addListener: vi.fn(),
      removeListener: vi.fn(),
    })),
  );
  permission = async () => reply({ allowed: false });
  vi.stubGlobal(
    "fetch",
    vi.fn(async (url: string) => {
      if (url === "/api/v1/schedulers") return reply(registry);
      const id = /\/schedulers\/([^/]+)\/checkPermissions\//.exec(url)?.[1];
      if (!id) throw new Error(`Unexpected HTTP request: ${url}`);
      return permission(id);
    }),
  );
  for (const [id, client] of Object.entries(clients)) {
    client.listOperations
      .mockReset()
      .mockResolvedValue({ operations: [operation(`${id}-list`)] });
    client.getOperation
      .mockReset()
      .mockResolvedValue({ operation: operation(`${id}-detail`) });
    client.killOperations.mockReset().mockResolvedValue({});
  }
  vi.mocked(getBuildQueueStateClient).mockImplementation((id) => {
    if (!clients[id]) throw new Error(`Unexpected scheduler: ${id}`);
    return clients[id] as unknown as ReturnType<
      typeof getBuildQueueStateClient
    >;
  });
});
afterEach(() => {
  state.cache?.clear();
  vi.restoreAllMocks();
  vi.unstubAllGlobals();
});

async function mount(url: string) {
  const router = createPortalRouter(
    createMemoryHistory({ initialEntries: [url] }),
  );
  const container = document.createElement("div");
  document.body.append(container);
  const root = createRoot(container);
  await act(async () => {
    root.render(<RouterProvider router={router} />);
    await router.load();
  });
  return {
    router,
    container,
    close: async () => {
      await act(async () => root.unmount());
      container.remove();
    },
  };
}
async function waitFor(check: () => void) {
  await act(async () => {
    await vi.waitFor(check);
  });
}
async function choose(container: HTMLElement, id: string) {
  const select = container.querySelector<HTMLSelectElement>(
    'select[aria-label="Scheduler"]',
  );
  if (!select) throw new Error("Missing production scheduler selector");
  await act(async () => {
    select.value = id;
    select.dispatchEvent(new Event("change", { bubbles: true }));
  });
}
const killButton = (container: HTMLElement) =>
  [...container.querySelectorAll("button")].find(
    (button) => button.textContent === "Kill operation",
  );

describe("production scheduler routing and cancellation UI", () => {
  it("reruns the production operations loader for scheduler-only navigation and rejects late old results", async () => {
    const view = await mount("/operations?scheduler=b");
    try {
      await waitFor(() =>
        expect(view.container.textContent).toContain("b-list"),
      );
      // Fresh matches must still use scheduler identity. Ordinary stale-loader
      // reloads must not hide a missing production scheduler dependency.
      view.router.update({ defaultStaleTime: Number.POSITIVE_INFINITY });
      let finishA:
        | ((value: { operations: OperationState[] }) => void)
        | undefined;
      clients.a.listOperations.mockReturnValueOnce(
        new Promise((resolve) => {
          finishA = resolve;
        }),
      );
      await choose(view.container, "a");
      await waitFor(() =>
        expect(clients.a.listOperations).toHaveBeenCalledTimes(1),
      );
      const signal = clients.a.listOperations.mock.calls[0][1]
        .signal as AbortSignal;
      await choose(view.container, "b");
      await waitFor(() =>
        expect(view.router.state.location.search.scheduler).toBe("b"),
      );
      expect(signal.aborted).toBe(true);
      await act(async () => {
        finishA?.({ operations: [operation("late-a")] });
      });
      await waitFor(() =>
        expect(view.container.textContent).toContain("b-list"),
      );
      expect(view.container.textContent).not.toContain("late-a");
      expect(view.router.state.matches.at(-1)?.loaderData).toMatchObject({
        operations: [{ targetId: "b-list" }],
      });
    } finally {
      await view.close();
    }
  });

  it("uses the production selector handler to leave colliding operation details and scope later cancellation", async () => {
    const view = await mount("/operations/same-id?scheduler=a");
    try {
      await waitFor(() =>
        expect(view.container.textContent).toContain("a-detail"),
      );
      expect(killButton(view.container)?.disabled).toBe(true);
      await choose(view.container, "b");
      await waitFor(() =>
        expect(view.router.state.location.pathname).toBe("/operations"),
      );
      expect(view.router.state.location.search.scheduler).toBe("b");
      expect(clients.b.listOperations).toHaveBeenCalledTimes(1);
      expect(clients.b.getOperation).not.toHaveBeenCalled();
      await waitFor(() =>
        expect(view.container.textContent).toContain("b-list"),
      );
      expect(view.container.textContent).not.toContain("a-detail");

      let grant: ((response: Response) => void) | undefined;
      permission = (id) =>
        id === "b"
          ? new Promise((resolve) => {
              grant = resolve;
            })
          : Promise.resolve(reply({ allowed: false }));
      await act(async () => {
        await view.router.navigate({
          to: "/operations/$operationID",
          params: { operationID: "same-id" },
          search: { scheduler: "b" },
        });
      });
      await waitFor(() =>
        expect(view.container.textContent).toContain("b-detail"),
      );
      expect(
        vi
          .mocked(fetch)
          .mock.calls.some(
            ([url]) =>
              String(url) ===
              "/api/v1/schedulers/b/checkPermissions/killOperation/same-id",
          ),
      ).toBe(true);
      await act(async () => {
        grant?.(reply({ allowed: true }));
      });
      await waitFor(() =>
        expect(killButton(view.container)?.disabled).toBe(false),
      );
      clients.b.getOperation.mockResolvedValue({
        operation: operation("b-refetched-after-kill"),
      });
      const invalidate = vi.spyOn(
        state.cache as QueryClient,
        "invalidateQueries",
      );
      await act(async () => {
        killButton(view.container)?.click();
      });
      await waitFor(() =>
        expect(view.container.textContent).toContain("b-refetched-after-kill"),
      );
      expect(clients.a.killOperations).not.toHaveBeenCalled();
      expect(clients.b.killOperations).toHaveBeenCalledExactlyOnceWith(
        expect.objectContaining({ filter: { operationName: "same-id" } }),
      );
      expect(invalidate).toHaveBeenCalledWith({
        queryKey: ["operationDetails", "b", "same-id"],
        exact: true,
      });
      expect(
        state.cache?.getQueryData(["operationDetails", "a", "same-id"]),
      ).toMatchObject({ operation: { targetId: "a-detail" } });
    } finally {
      await view.close();
    }
  });

  it("can switch to a healthy scheduler after an operations loader failure", async () => {
    clients.b.listOperations.mockRejectedValue(
      new Error("Scheduler unavailable"),
    );
    const view = await mount("/operations?scheduler=b");
    try {
      await waitFor(() =>
        expect(view.container.textContent).toContain(
          "Unable to load operations",
        ),
      );
      expect(
        view.container.querySelectorAll('[aria-label="Scheduler"]'),
      ).toHaveLength(1);
      await choose(view.container, "a");
      await waitFor(() =>
        expect(view.container.textContent).toContain("a-list"),
      );
    } finally {
      await view.close();
    }
  });

  it("reports an invalid invocation filter without a misleading scheduler failure", async () => {
    const filter = encodeURIComponent(JSON.stringify({ "@type": "invalid" }));
    const view = await mount(`/operations?scheduler=b&filter=${filter}`);
    try {
      await waitFor(() =>
        expect(view.container.textContent).toContain("Something went wrong"),
      );
      expect(view.container.textContent).toContain(
        "Either toolInvocationId or correlatedInvocationsId",
      );
      expect(view.container.textContent).not.toContain(
        "Unable to open scheduler",
      );
      expect(
        view.container.querySelectorAll('[aria-label="Scheduler"]'),
      ).toHaveLength(1);
      expect(clients.b.listOperations).not.toHaveBeenCalled();
    } finally {
      await view.close();
    }
  });

  it("preserves invalid scheduler input without injecting a NUL into navigation links", async () => {
    const view = await mount("/operations?scheduler=a&scheduler=b");
    try {
      await waitFor(() =>
        expect(view.container.textContent).toContain(
          "The selected scheduler is not configured",
        ),
      );
      expect(view.router.state.location.search.scheduler).toEqual(["a", "b"]);
      await act(async () => {
        await view.router.navigate({ to: "/" });
      });
      expect(view.router.state.location.href).not.toContain("%00");
      await act(async () => {
        await view.router.navigate({ to: "/operations" });
      });
      await choose(view.container, "a");
      await waitFor(() =>
        expect(view.container.textContent).toContain("a-list"),
      );
    } finally {
      await view.close();
    }
  });

  it.each([
    "",
    "missing",
  ])("rejects explicit scheduler %j without falling back to the default", async (id) => {
    const view = await mount(`/operations?scheduler=${id}`);
    try {
      await waitFor(() =>
        expect(view.container.textContent).toContain(
          "The selected scheduler is not configured",
        ),
      );
      expect(clients.a.listOperations).not.toHaveBeenCalled();
      expect(
        view.container.querySelectorAll('[aria-label="Scheduler"]'),
      ).toHaveLength(1);
      await choose(view.container, "a");
      await waitFor(() =>
        expect(view.container.textContent).toContain("a-list"),
      );
    } finally {
      await view.close();
    }
  });
});

function makeRouter(entry: string) {
  return createPortalRouter(createMemoryHistory({ initialEntries: [entry] }));
}

describe("production scheduler URLs and history", () => {
  it.each([
    "1",
    "2024",
    "true",
    "false",
    "null",
    "001",
    "1e3",
    "999999999999999999999999999999999999999999999999999999999999999",
  ])("preserves the exact scheduler ID %s through loading, navigation and refresh", async (id) => {
    const router = makeRouter(`/scheduler?scheduler=${id}`);
    await router.load();
    expect(router.state.matches.at(-1)?.search.scheduler).toBe(id);
    await router.navigate({ to: "/operations" });
    expect(router.state.matches.at(-1)?.loaderData).toMatchObject({
      operations: [{ targetId: `${id}-list` }],
    });
    expect(clients[id].listOperations).toHaveBeenCalled();
    expect(
      new URLSearchParams(router.state.location.searchStr).get("scheduler"),
    ).toBe(id);
    const refreshed = makeRouter(router.state.location.href);
    await refreshed.load();
    expect(refreshed.state.matches.at(-1)?.search.scheduler).toBe(id);
    const legacyQuoted = makeRouter(
      `/operations?scheduler=${encodeURIComponent(JSON.stringify(id))}`,
    );
    await legacyQuoted.load();
    expect(legacyQuoted.state.matches.at(-1)?.loaderData).toMatchObject({
      operations: [{ targetId: `${id}-list` }],
    });
  });

  it("retains JSON serialization for filters and accepts previously emitted quoted scheduler IDs", async () => {
    const filter = { toolInvocationId: "invocation", nested: [1, true] };
    expect(
      parsePortalSearch(stringifyPortalSearch({ scheduler: "1e3", filter })),
    ).toEqual({ scheduler: "1e3", filter });
    const router = makeRouter("/scheduler?scheduler=%221e3%22");
    await router.load();
    expect(router.state.matches.at(-1)?.search.scheduler).toBe("1e3");
    expect(
      schedulerSearch(parsePortalSearch("?scheduler=a&scheduler=b")).scheduler,
    ).toEqual(["a", "b"]);
  });

  it("preserves selection through links with their own search and outside scheduler pages", async () => {
    const router = makeRouter("/scheduler?scheduler=b");
    await router.load();
    await router.navigate({ to: "/operations", search: { filter: undefined } });
    expect(router.state.location.search.scheduler).toBe("b");
    await router.navigate({ to: "/browser" });
    await router.navigate({ to: "/scheduler" });
    expect(router.state.location.search.scheduler).toBe("b");
  });

  it("keeps independent tab selection and restores selection from history", async () => {
    const first = makeRouter("/scheduler?scheduler=a");
    const second = makeRouter("/scheduler?scheduler=b");
    await Promise.all([first.load(), second.load()]);
    await first.navigate({ to: "/scheduler", search: { scheduler: "b" } });
    first.history.back();
    await first.load();
    expect(first.state.location.search.scheduler).toBe("a");
    expect(second.state.location.search.scheduler).toBe("b");
    const refreshed = makeRouter(second.state.location.href);
    await refreshed.load();
    expect(refreshed.state.location.search.scheduler).toBe("b");
  });
});
