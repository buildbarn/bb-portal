import { ApolloLink, Observable } from "@apollo/client";
import { ApolloProvider } from "@apollo/client/react";
import { act, type ComponentType } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, expect, test, vi } from "vitest";
import { apolloClient } from "@/components/ApolloWrapper";
import { Route as InvocationsRoute } from "@/routes/bazel-invocations.index";
import { Route as BuildRoute } from "@/routes/builds.$buildUUID";

const pages = vi.hoisted(() => {
  type PageProps = {
    invocations: { id: string; exitCodeName?: string | null }[];
    pageInfo: unknown;
    pageSize: number;
  };
  return {
    invocations: vi.fn<(_props: PageProps) => null>(() => null),
    build: vi.fn<(_props: PageProps) => null>(() => null),
  };
});

vi.mock("@/components/pages/BazelInvocations", () => ({
  BazelInvocationsPage: pages.invocations,
}));
vi.mock("@/components/pages/BuildDetails", () => ({
  BuildDetailsPage: pages.build,
}));
vi.mock("@/components/PageCursorTable", () => ({ DEFAULT_PAGE_SIZE: 10 }));
vi.mock("@/utils/generatePageTitle", () => ({ generatePageTitle: vi.fn() }));
vi.mock("@/main", () => ({ NotFoundError: class extends Error {} }));

const pageInfo = {
  __typename: "PageInfo",
  startCursor: "page-start",
  endCursor: "page-end",
  hasNextPage: true,
  hasPreviousPage: true,
};
const completed = new Set<string>();
const requests: Record<string, unknown>[] = [];
let root: Root;

function invocation(id: string) {
  return {
    __typename: "BazelInvocation",
    id,
    invocationID: `invocation-${id}`,
    startedAt: "2026-01-01T00:00:00Z",
    endedAt: completed.has(id) ? "2026-01-01T00:01:00Z" : null,
    exitCodeName: completed.has(id) ? "SUCCESS" : null,
    username: "test",
    authenticatedUser: null,
    build: null,
    instanceName: { __typename: "InstanceName", id: "instance", name: "test" },
    connectionMetadata: {
      __typename: "ConnectionMetadata",
      connectionLastOpenAt: "2026-01-01T00:00:00Z",
      timeSinceLastConnectionMillis: 1000,
    },
    tags: { __typename: "InvocationTagConnection", edges: [] },
    profile: null,
    originalCommandLine: null,
  };
}

beforeEach(async () => {
  Object.assign(globalThis, { IS_REACT_ACT_ENVIRONMENT: true });
  vi.useFakeTimers();
  completed.clear();
  completed.add("C");
  requests.length = 0;
  pages.invocations.mockClear();
  pages.build.mockClear();
  await apolloClient.clearStore();
  apolloClient.setLink(
    new ApolloLink((operation) => {
      requests.push(operation.variables);
      const ids = operation.variables.where?.idIn as string[] | undefined;
      const connection = {
        __typename: "BazelInvocationConnection",
        // Polling must not replace the loaded page's cursors or membership.
        pageInfo: ids ? { ...pageInfo, endCursor: "poll-end" } : pageInfo,
        edges: (ids ?? ["A", "B", "C"]).map((id) => ({
          __typename: "BazelInvocationEdge",
          node: invocation(id),
        })),
      };
      const data =
        operation.operationName === "FindBuildByUUID"
          ? {
              getBuild: {
                __typename: "Build",
                id: "build",
                buildUUID: "build-uuid",
                timestamp: "2026-01-01T00:00:00Z",
                tags: { __typename: "BuildTagConnection", edges: [] },
                invocations: connection,
              },
            }
          : { findBazelInvocations: connection };
      return new Observable((observer) => {
        observer.next({ data });
        observer.complete();
      });
    }),
  );
  root = createRoot(document.createElement("div"));
});

afterEach(async () => {
  await act(async () => root.unmount());
  await apolloClient.clearStore();
  vi.restoreAllMocks();
  vi.useRealTimers();
});

async function checkPolling(
  Component: ComponentType,
  page: typeof pages.build,
) {
  await act(async () => {
    root.render(
      <ApolloProvider client={apolloClient}>
        <Component />
      </ApolloProvider>,
    );
  });
  await act(async () => {
    await vi.dynamicImportSettled();
    await vi.advanceTimersByTimeAsync(20);
  });
  expect(requests[0]).toMatchObject({
    first: 3,
    after: "previous-page",
    where: { and: expect.arrayContaining([{ username: "test" }]) },
  });
  expect(requests.at(-1)?.where).toEqual({ idIn: ["A", "B"] });

  completed.add("A");
  await act(async () => {
    await vi.advanceTimersByTimeAsync(5020);
  });
  expect(page.mock.lastCall?.[0].invocations).toMatchObject([
    { id: "A", exitCodeName: "SUCCESS" },
    { id: "B", exitCodeName: null },
    { id: "C", exitCodeName: "SUCCESS" },
  ]);
  expect(requests.at(-1)?.where).toEqual({ idIn: ["B"] });
  expect(page.mock.lastCall?.[0].pageInfo).toEqual(pageInfo);
  expect(page.mock.lastCall?.[0].pageSize).toBe(3);

  completed.add("B");
  await act(async () => {
    await vi.advanceTimersByTimeAsync(5020);
  });
  expect(page.mock.lastCall?.[0].invocations[1]).toMatchObject({
    id: "B",
    exitCodeName: "SUCCESS",
  });
  const requestCount = requests.length;
  await act(async () => {
    await vi.advanceTimersByTimeAsync(10000);
  });
  expect(requests).toHaveLength(requestCount);
}

const table = {
  pageSize: 3,
  pagination: { first: 3, after: "previous-page" },
  where: [{ username: "test" }],
};

test("invocation polling updates the loaded page and stops polling completed IDs", async () => {
  const loader = InvocationsRoute.options.loader;
  if (typeof loader !== "function") throw new Error("Missing route loader");
  const loaded = await loader({ deps: { invocationTable: table } } as never);
  vi.spyOn(InvocationsRoute, "useLoaderData").mockReturnValue(loaded);
  vi.spyOn(InvocationsRoute, "useNavigate").mockReturnValue(vi.fn());
  await checkPolling(
    InvocationsRoute.options.component as ComponentType,
    pages.invocations,
  );
});

test("build polling updates the loaded page and preserves its pagination", async () => {
  const loader = BuildRoute.options.loader;
  if (typeof loader !== "function") throw new Error("Missing route loader");
  const loaded = await loader({
    deps: { buildTable: table },
    params: { buildUUID: "build-uuid" },
  } as never);
  vi.spyOn(BuildRoute, "useLoaderData").mockReturnValue(loaded);
  vi.spyOn(BuildRoute, "useNavigate").mockReturnValue(vi.fn());
  vi.spyOn(BuildRoute, "useParams").mockReturnValue({
    buildUUID: "build-uuid",
  });
  await checkPolling(
    BuildRoute.options.component as ComponentType,
    pages.build,
  );
});
