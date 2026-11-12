import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import {
  createMemoryHistory,
  createRootRoute,
  createRoute,
  createRouter,
  Outlet,
  RouterProvider,
  useSearch,
} from "@tanstack/react-router";
import { act, type ReactNode, useEffect } from "react";
import { createRoot } from "react-dom/client";
import { afterEach, describe, expect, it, vi } from "vitest";
import { SchedulerLayout } from "@/components/SchedulerLayout";
import WorkersGrid from "@/components/WorkersGrid";
import { getBuildQueueStateClient } from "@/grpc/buildQueueStateClient";
import type {
  ListWorkersResponse,
  WorkerState,
} from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import type { WorkerListStatus } from "@/routes/scheduler.worker";
import { resolveScheduler, schedulerSearch } from "./registry";

vi.mock("@/routes/scheduler.worker", () => ({
  WorkerListStatus: { ALL: "all", EXECUTING: "executing" },
}));

const mounts = vi.hoisted(() => ({ count: 0 }));
vi.mock("@/grpc/buildQueueStateClient", () => ({
  getBuildQueueStateClient: vi.fn(),
}));
vi.mock("@/components/WorkersInfo", () => ({ default: () => null }));
vi.mock("@/components/WorkersTable", () => ({
  default: ({ data }: { data?: WorkerState[] }) => {
    useEffect(() => {
      mounts.count++;
    }, []);
    return (
      <div data-testid="workers">
        {data?.map((worker) => worker.id.name).join(",")}
      </div>
    );
  },
}));
vi.mock("antd", () => ({
  Row: ({ children }: { children: ReactNode }) => <div>{children}</div>,
  Space: ({ children }: { children: ReactNode }) => <div>{children}</div>,
  Select: () => <select aria-label="Scheduler" />,
  Alert: () => null,
}));

const queue = {
  platformQueueName: {
    instanceNamePrefix: "instance",
    platform: { properties: [] },
  },
  sizeClass: 1,
};
const status = "all" as WorkerListStatus;
function response(name: string): ListWorkersResponse {
  return {
    workers: [
      {
        id: { name },
        drained: false,
        timeout: undefined,
        currentOperation: undefined,
      },
    ],
    paginationInfo: { totalEntries: 2, startIndex: 0 },
  };
}
function fakeBackends() {
  let finishPage: ((value: ListWorkersResponse) => void) | undefined;
  vi.mocked(getBuildQueueStateClient).mockImplementation(
    (id) =>
      ({
        listWorkers: (request: {
          startAfter?: unknown;
          filter?: { executing?: unknown };
        }) => {
          if (id === "b" || request.startAfter || request.filter?.executing)
            return new Promise<ListWorkersResponse>((resolve) => {
              finishPage = resolve;
            });
          return Promise.resolve(response("a-first"));
        },
      }) as unknown as ReturnType<typeof getBuildQueueStateClient>,
  );
  return (value: ListWorkersResponse) => finishPage?.(value);
}

afterEach(() => {
  vi.restoreAllMocks();
  vi.unstubAllGlobals();
  mounts.count = 0;
});

describe("worker pagination", () => {
  it.each([
    true,
    false,
  ])("retains page/filter rows and isolates scheduler switches (keyed layout: %s)", async (keyedLayout) => {
    Object.assign(globalThis, { IS_REACT_ACT_ENVIRONMENT: true });
    vi.stubGlobal("scrollTo", vi.fn());
    const finish = fakeBackends();
    const cache = new QueryClient();
    const registry = {
      schedulers: [
        { id: "a", label: "A" },
        { id: "b", label: "B" },
      ],
      defaultSchedulerId: "a",
    };
    const rootRoute = createRootRoute({
      validateSearch: (search: Record<string, unknown>) => ({
        ...schedulerSearch(search),
        workerStatusFilter: search.workerStatusFilter as
          | WorkerListStatus
          | undefined,
        cursor: search.cursor as Record<string, string> | undefined,
      }),
      component: () => (
        <QueryClientProvider client={cache}>
          <Outlet />
        </QueryClientProvider>
      ),
    });
    const schedulerRoute = createRoute({
      getParentRoute: () => rootRoute,
      path: "scheduler",
      beforeLoad: ({ search }) => ({
        schedulerRegistry: registry,
        scheduler: resolveScheduler(registry, search.scheduler),
      }),
      component: keyedLayout ? SchedulerLayout : Outlet,
    });
    const workerRoute = createRoute({
      getParentRoute: () => schedulerRoute,
      path: "worker",
      component: () => {
        const { cursor, workerStatusFilter } = useSearch({ strict: false }) as {
          cursor?: Record<string, string>;
          workerStatusFilter?: WorkerListStatus;
        };
        return (
          <WorkersGrid
            workerStatusFilter={workerStatusFilter ?? status}
            sizeClassQueueName={queue}
            cursor={cursor}
          />
        );
      },
    });
    const router = createRouter({
      routeTree: rootRoute.addChildren([
        schedulerRoute.addChildren([workerRoute]),
      ]),
      history: createMemoryHistory({
        initialEntries: ["/scheduler/worker?scheduler=a"],
      }),
    });
    const container = document.createElement("div");
    const root = createRoot(container);
    try {
      await act(async () => {
        root.render(<RouterProvider router={router} />);
        await router.load();
      });
      await act(async () => {
        await vi.waitFor(() =>
          expect(container.textContent).toContain("a-first"),
        );
      });
      const originalMounts = mounts.count;
      await act(async () => {
        await router.navigate({
          to: "/scheduler/worker",
          search: {
            scheduler: "a",
            cursor: { name: "cursor" },
            workerStatusFilter: status,
            sizeClassQueueName: queue,
          },
        });
      });
      expect(container.textContent).toContain("a-first");
      expect(mounts.count).toBe(originalMounts);
      await act(async () => {
        finish(response("a-next"));
      });
      await act(async () => {
        await vi.waitFor(() =>
          expect(container.textContent).toContain("a-next"),
        );
      });
      await act(async () => {
        await router.navigate({
          to: "/scheduler/worker",
          search: {
            scheduler: "a",
            workerStatusFilter: "executing" as WorkerListStatus,
            sizeClassQueueName: queue,
          },
        });
      });
      expect(container.textContent).toContain("a-next");
      expect(mounts.count).toBe(originalMounts);
      await act(async () => {
        finish(response("a-filter"));
      });
      await act(async () => {
        await vi.waitFor(() =>
          expect(container.textContent).toContain("a-filter"),
        );
      });
      await act(async () => {
        await router.navigate({
          to: "/scheduler/worker",
          search: {
            scheduler: "b",
            workerStatusFilter: status,
            sizeClassQueueName: queue,
          },
        });
      });
      expect(container.textContent).not.toContain("a-filter");
      expect(mounts.count).toBe(originalMounts + (keyedLayout ? 1 : 0));
      await act(async () => {
        finish(response("b-first"));
      });
      await act(async () => {
        await vi.waitFor(() =>
          expect(container.textContent).toContain("b-first"),
        );
      });
    } finally {
      await act(async () => root.unmount());
      cache.clear();
    }
  });
});
