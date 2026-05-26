import { createFileRoute, linkOptions } from "@tanstack/react-router";
import type { SorterResult } from "antd/es/table/interface";
import z from "zod";
import { apolloClient } from "@/components/ApolloWrapper";
import { cacheLocationFromTestResults } from "@/components/CacheLocationTag";
import { DEFAULT_PAGE_SIZE } from "@/components/PageCursorTable";
import {
  type TablePaginationVars,
  TablePaginationVarsSchema,
} from "@/components/PageCursorTable/types";
import { TestTab } from "@/components/TestTab";
import {
  defaultSorting,
  type TestTabRowType,
} from "@/components/TestTab/columns";
import { gql } from "@/graphql/__generated__";
import {
  OrderDirection,
  type TestSummaryOrder,
  TestSummaryOrderField,
  type TestSummaryWhereInput,
} from "@/graphql/__generated__/graphql";
import {
  TestSummaryOrderSchema,
  TestSummaryWhereInputSchema,
} from "@/graphql/__generated__/zod";
import { NotFoundError } from "@/main";
import { env } from "@/utils/env";
import { requireFeature } from "@/utils/featureGuard";
import { generatePageTitle } from "@/utils/generatePageTitle";
import { parseGraphqlEdgeListWithFragment } from "@/utils/parseGraphqlEdgeList";

export const GET_TESTS_FOR_INVOCATION = gql(/* GraphQl */ `
  query GetTestsForInvocation(
    $after: Cursor
    $first: Int
    $before: Cursor
    $last: Int
    $orderBy: TestSummaryOrder
    $where: TestSummaryWhereInput
  ) {
    findTestSummaries(
      after: $after
      first: $first
      before: $before
      last: $last
      orderBy: $orderBy
      where: $where
    ) {
      pageInfo {
        startCursor
        endCursor
        hasNextPage
        hasPreviousPage
      }
      edges {
        node {
          ...TestSummaryNode
        }
      }
    }
  }
`);

const TEST_SUMMARY_FRAGMENT = gql(/* GraphQL */ `
  fragment TestSummaryNode on TestSummary {
    id
    overallStatus
    totalRunDurationInMs
    testResults {
      id
      cachedLocally
      cachedRemotely
    }
    invocationTarget {
      id
      target {
        id
        instanceName {
          id
          name
        }
        label
        aspect
        targetKind
      }
    }
  }
`);

const TestSummarySearchSchema = z.object({
  testSummaryTable: TablePaginationVarsSchema.extend({
    where: z.array(TestSummaryWhereInputSchema().partial()).optional(),
    orderBy: TestSummaryOrderSchema().optional(),
  }).optional(),
});

export const Route = createFileRoute("/bazel-invocations/$invocationID/tests")({
  component: RouteComponent,
  validateSearch: (search) => TestSummarySearchSchema.parse(search),
  beforeLoad: requireFeature(env.featureFlags?.bes?.pageTargets),
  // TODO: Add backend integration test for this
  loaderDeps: ({ search: { testSummaryTable } }) => ({ testSummaryTable }),
  loader: async ({ params, deps }) => {
    // We set the defaults here instead of in the validate search function, as
    // that updates the URL and we don't want to do that on initial load.
    const pageSize = deps.testSummaryTable?.pageSize ?? DEFAULT_PAGE_SIZE;
    const pagination = deps.testSummaryTable?.pagination ?? {
      first: pageSize,
    };
    const where: TestSummaryWhereInput[] = deps.testSummaryTable?.where ?? [];

    const orderBy: TestSummaryOrder =
      deps.testSummaryTable?.orderBy ?? defaultSorting;

    const { data, error } = await apolloClient.query({
      errorPolicy: "all",
      query: GET_TESTS_FOR_INVOCATION,
      variables: {
        where: {
          and: [
            {
              hasInvocationTargetWith: [
                {
                  hasBazelInvocationWith: [
                    { invocationID: params.invocationID },
                  ],
                },
              ],
            },
            ...where,
          ],
        },
        orderBy,
        ...pagination,
      },
      fetchPolicy: "cache-first",
    });
    if (!data) {
      throw new NotFoundError("tests", error?.message);
    }

    const testSummaryData: TestTabRowType[] = parseGraphqlEdgeListWithFragment(
      TEST_SUMMARY_FRAGMENT,
      data?.findTestSummaries,
    ).map((ts) => ({
      ...ts,
      cacheLocation: cacheLocationFromTestResults(ts.testResults),
    }));
    return {
      testSummaryData,
      pageSize,
      pageInfo: data.findTestSummaries.pageInfo,
    };
  },
  head: (_ctx) => ({
    meta: [
      {
        title: generatePageTitle([
          "Invocation",
          "Tests",
          _ctx.params.invocationID,
        ]),
      },
    ],
  }),
});

const getPaginationUpdateLink = (newPagination: TablePaginationVars) =>
  linkOptions({
    from: Route.id,
    to: ".",
    search: (prev): typeof prev => ({
      ...prev,
      testSummaryTable: {
        ...prev.testSummaryTable,
        ...newPagination,
      },
    }),
  });

function RouteComponent() {
  const { testSummaryData, pageSize, pageInfo } = Route.useLoaderData();
  const navigate = Route.useNavigate();

  const onFilterChange = (where: TestSummaryWhereInput[]) => {
    navigate({
      from: Route.id,
      to: ".",
      search: (prev): typeof prev => ({
        ...prev,
        testSummaryTable: {
          ...prev.testSummaryTable,
          where,
          pagination: undefined,
        },
      }),
    });
  };
  const onSortChange = (
    sorter: SorterResult<TestTabRowType> | SorterResult<TestTabRowType>[],
  ) => {
    const s = Array.isArray(sorter) ? sorter[0] : sorter;
    if (!s || !s.order) {
      return;
    }
    switch (s.columnKey) {
      case "totalRunDurationInMs":
        navigate({
          from: Route.id,
          to: ".",
          search: (prev): typeof prev => ({
            ...prev,
            testSummaryTable: {
              ...prev.testSummaryTable,
              pagination: undefined,
              orderBy: {
                field: TestSummaryOrderField.TotalRunDurationInMs,
                direction:
                  s.order === "ascend"
                    ? OrderDirection.Asc
                    : OrderDirection.Desc,
              },
            },
          }),
        });
        break;
    }
  };
  return (
    <TestTab
      getPaginationUpdateLink={getPaginationUpdateLink}
      onFilterChange={onFilterChange}
      onSortChange={onSortChange}
      pageInfo={pageInfo}
      pageSize={pageSize}
      testSummaryData={testSummaryData}
    />
  );
}
