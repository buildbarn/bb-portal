import {
  createMemoryHistory,
  createRootRoute,
  createRoute,
  createRouter,
  Outlet,
  RouterProvider,
} from "@tanstack/react-router";
import { act } from "react";
import { createRoot } from "react-dom/client";
import { afterEach, describe, expect, it, vi } from "vitest";
import { RequestMetadata } from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import { OperationState } from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import { useSchedulerContext } from "@/schedulers/context";
import { schedulerSearch } from "@/schedulers/registry";
import {
  parsePortalSearch,
  stringifyPortalSearch,
} from "@/schedulers/urlSearch";
import { ProtobufTypeUrls } from "@/types/protobufTypeUrls";
import OperationDetails from "./index";

const { useQuery } = vi.hoisted(() => ({ useQuery: vi.fn() }));

vi.mock("@tanstack/react-query", async (importOriginal) => {
  const actual = await importOriginal<typeof import("@tanstack/react-query")>();
  return {
    ...actual,
    useMutation: () => ({ mutate: vi.fn() }),
    useQuery,
    useQueryClient: () => ({ invalidateQueries: vi.fn() }),
  };
});
vi.mock("@/schedulers/context", () => ({
  useScheduler: () => ({
    scheduler: { instanceNamePrefix: "ci/aws" },
    client: { killOperations: vi.fn() },
  }),
  useSchedulerContext: vi.fn(),
}));

afterEach(() => vi.unstubAllGlobals());

describe("operation details invocation links", () => {
  it("omits scheduler selection in single-scheduler mode", async () => {
    Object.assign(globalThis, { IS_REACT_ACT_ENVIRONMENT: true });
    vi.stubGlobal("scrollTo", vi.fn());
    vi.stubGlobal(
      "matchMedia",
      vi.fn((media: string) => ({
        matches: false,
        media,
        addListener: vi.fn(),
        removeListener: vi.fn(),
        addEventListener: vi.fn(),
        removeEventListener: vi.fn(),
        dispatchEvent: vi.fn(),
      })),
    );
    const operation = OperationState.fromPartial({
      name: "op-1",
      invocationName: {
        ids: [
          {
            typeUrl: ProtobufTypeUrls.REQUEST_METADATA,
            value: RequestMetadata.encode(
              RequestMetadata.fromPartial({
                toolInvocationId: "00000000-0000-4000-8000-000000000000",
              }),
            ).finish(),
          },
        ],
      },
    });
    useQuery
      .mockReturnValueOnce({
        data: { operation },
        isLoading: false,
        isError: false,
      })
      .mockReturnValueOnce({ data: true, isLoading: false, isError: false });
    vi.mocked(useSchedulerContext).mockReturnValue({
      schedulerRegistry: { schedulers: [{ instanceNamePrefix: "ci/aws" }] },
    });

    const rootRoute = createRootRoute({
      validateSearch: schedulerSearch,
      component: Outlet,
    });
    const source = createRoute({
      getParentRoute: () => rootRoute,
      path: "source",
      component: () => <OperationDetails operationID="op-1" />,
    });
    const operations = createRoute({
      getParentRoute: () => rootRoute,
      path: "operations",
    });
    const router = createRouter({
      routeTree: rootRoute.addChildren([source, operations]),
      history: createMemoryHistory({ initialEntries: ["/source"] }),
      parseSearch: parsePortalSearch,
      stringifySearch: stringifyPortalSearch,
    });
    const container = document.createElement("div");
    const root = createRoot(container);
    try {
      await act(async () => {
        root.render(<RouterProvider router={router} />);
        await router.load();
      });
      const link = container.querySelector('a[href^="/operations"]');
      expect(link).not.toBeNull();
      const location = new URL(
        link?.getAttribute("href") ?? "",
        "http://portal.test",
      );
      expect(location.searchParams.has("scheduler")).toBe(false);
      expect(parsePortalSearch(location.search).filter).toMatchObject({
        toolInvocationId: "00000000-0000-4000-8000-000000000000",
      });
    } finally {
      await act(async () => root.unmount());
    }
  });
});
