import { QueryClient, QueryObserver } from "@tanstack/react-query";
import { afterEach, describe, expect, it, vi } from "vitest";
import { getBuildQueueStateClient } from "@/grpc/buildQueueStateClient";
import type { WorkerListStatus } from "@/routes/scheduler.worker";
import {
  killPermissionOptions,
  operationDetailsOptions,
  platformQueuesOptions,
  schedulerStatisticsOptions,
  workersOptions,
} from "./queries";

vi.mock("@/grpc/buildQueueStateClient", () => ({
  getBuildQueueStateClient: vi.fn(),
}));

afterEach(() => {
  vi.restoreAllMocks();
  vi.unstubAllGlobals();
});

describe("scheduler request isolation", () => {
  it.each([
    "all",
    "executing",
  ])("routes %s worker requests with their queue, cursor and abort signal", async (status) => {
    const listWorkers = vi.fn().mockResolvedValue({ workers: [] });
    vi.mocked(getBuildQueueStateClient).mockReturnValue({
      listWorkers,
    } as unknown as ReturnType<typeof getBuildQueueStateClient>);
    const queue = {
      platformQueueName: {
        instanceNamePrefix: "buildbarn",
        platform: { properties: [] },
      },
      sizeClass: 2,
    };
    const cursor = { worker: "worker-100" };
    const cache = new QueryClient();
    try {
      await cache.fetchQuery(
        workersOptions("b", status as WorkerListStatus, queue, cursor),
      );
      expect(getBuildQueueStateClient).toHaveBeenCalledWith("b");
      expect(listWorkers).toHaveBeenCalledExactlyOnceWith(
        {
          pageSize: 100,
          startAfter: { workerId: cursor },
          filter:
            status === "all"
              ? { all: queue }
              : { executing: { sizeClassQueueName: queue } },
        },
        { signal: expect.any(AbortSignal) },
      );
    } finally {
      cache.clear();
    }
  });

  it("keeps fresh platform queue results separate for each scheduler", async () => {
    vi.mocked(getBuildQueueStateClient).mockImplementation(
      (id) =>
        ({
          listPlatformQueues: async () => ({
            platformQueues: [{ name: { instanceNamePrefix: id } }],
          }),
        }) as unknown as ReturnType<typeof getBuildQueueStateClient>,
    );
    const cache = new QueryClient({
      defaultOptions: { queries: { staleTime: Number.POSITIVE_INFINITY } },
    });
    try {
      const a = await cache.fetchQuery(platformQueuesOptions("a"));
      const b = await cache.fetchQuery(platformQueuesOptions("b"));
      expect(a.platformQueues[0].name?.instanceNamePrefix).toBe("a");
      expect(b.platformQueues[0].name?.instanceNamePrefix).toBe("b");
    } finally {
      cache.clear();
    }
  });

  it("does not publish a late response from the previous scheduler", async () => {
    let finishA:
      | ((result: {
          operations: [];
          paginationInfo: { totalEntries: number; startIndex: number };
        }) => void)
      | undefined;
    let signalA: AbortSignal | undefined;
    vi.mocked(getBuildQueueStateClient).mockImplementation(
      (id) =>
        ({
          listOperations: (
            _request: unknown,
            options: { signal: AbortSignal },
          ) =>
            id === "a"
              ? new Promise((resolve) => {
                  finishA = resolve;
                  signalA = options.signal;
                })
              : Promise.resolve({
                  operations: [],
                  paginationInfo: { totalEntries: 22, startIndex: 0 },
                }),
        }) as unknown as ReturnType<typeof getBuildQueueStateClient>,
    );
    const cache = new QueryClient({
      defaultOptions: { queries: { retry: false } },
    });
    const observer = new QueryObserver(cache, schedulerStatisticsOptions("a"));
    const unsubscribe = observer.subscribe(() => {});
    observer.setOptions(schedulerStatisticsOptions("b"));
    await vi.waitFor(() =>
      expect(
        observer.getCurrentResult().data?.paginationInfo?.totalEntries,
      ).toBe(22),
    );
    expect(signalA?.aborted).toBe(true);
    finishA?.({
      operations: [],
      paginationInfo: { totalEntries: 11, startIndex: 0 },
    });
    await Promise.resolve();
    expect(observer.getCurrentResult().data?.paginationInfo?.totalEntries).toBe(
      22,
    );
    unsubscribe();
    cache.clear();
  });

  it("isolates the same operation ID across schedulers and checks the matching permission endpoint", async () => {
    vi.mocked(getBuildQueueStateClient).mockImplementation(
      (id) =>
        ({
          getOperation: () =>
            Promise.resolve({ operation: { name: "same-id", targetId: id } }),
        }) as unknown as ReturnType<typeof getBuildQueueStateClient>,
    );
    const request = vi
      .fn()
      .mockResolvedValue({ ok: true, json: async () => ({ allowed: true }) });
    vi.stubGlobal("fetch", request);
    const cache = new QueryClient();
    const [a, b] = await Promise.all([
      cache.fetchQuery(operationDetailsOptions("a", "same-id")),
      cache.fetchQuery(operationDetailsOptions("b", "same-id")),
    ]);
    expect(a.operation?.targetId).toBe("a");
    expect(b.operation?.targetId).toBe("b");
    await cache.fetchQuery(killPermissionOptions("b", "same/id"));
    expect(request.mock.calls[0][0]).toBe(
      "/api/v1/schedulers/b/checkPermissions/killOperation/same%2Fid",
    );
    cache.clear();
  });
});
