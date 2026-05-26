import { createFileRoute, linkOptions } from "@tanstack/react-router";
import z from "zod";
import { apolloClient } from "@/components/ApolloWrapper";
import { DEFAULT_PAGE_SIZE } from "@/components/PageCursorTable";
import {
  type TablePaginationVars,
  TablePaginationVarsSchema,
} from "@/components/PageCursorTable/types";
import { TestsPage } from "@/components/pages/Tests";
import { gql } from "@/graphql/__generated__/gql";
import type {
  TargetWhereInput,
  TestListRowFragment,
} from "@/graphql/__generated__/graphql";
import { TargetWhereInputSchema } from "@/graphql/__generated__/zod";
import { NotFoundError } from "@/main";
import { generatePageTitle } from "@/utils/generatePageTitle";
import { parseGraphqlEdgeListWithFragment } from "@/utils/parseGraphqlEdgeList";

export const GET_TESTS = gql(/* GraphQl */ `
  query GetTests(
    $after: Cursor
    $first: Int
    $before: Cursor
    $last: Int
    $where: TargetWhereInput
  ) {
    findTargets(
      after: $after
      first: $first
      before: $before
      last: $last
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
          ...TestListRow
        }
      }
    }
  }
`);

const TESTS_FRAGMENT = gql(/* GraphQL */ `
  fragment TestListRow on Target {
    id
    label
    aspect
    targetKind
    instanceName {
      name
    }
  }
`);

const TestSearchSchema = z.object({
  testTable: TablePaginationVarsSchema.extend({
    where: z.array(TargetWhereInputSchema().partial()).optional(),
  }).optional(),
});

export const Route = createFileRoute("/tests/")({
  component: RouteComponent,
  validateSearch: (search) => TestSearchSchema.parse(search),
  loaderDeps: ({ search: { testTable } }) => ({ testTable }),
  loader: async ({ deps }) => {
    // We set the defaults here instead of in the validate search function, as
    // that updates the URL and we don't want to do that on initial load.
    const pageSize = deps.testTable?.pageSize ?? DEFAULT_PAGE_SIZE;
    const pagination = deps.testTable?.pagination ?? { first: pageSize };
    const where = deps.testTable?.where ?? [];

    const { data, error } = await apolloClient.query({
      errorPolicy: "all",
      query: GET_TESTS,
      variables: {
        where: {
          and: [{ hasTestTarget: true }, ...where],
        },
        ...pagination,
      },
    });
    if (!data) {
      // TODO: Should this perhaps give a "failed to connect to server" error instead?
      throw new NotFoundError("tests", error?.message);
    }
    const testsData: TestListRowFragment[] = parseGraphqlEdgeListWithFragment(
      TESTS_FRAGMENT,
      data?.findTargets,
    );

    return { testsData, pageSize, pageInfo: data.findTargets.pageInfo };
  },
  head: (_ctx) => ({ meta: [{ title: generatePageTitle(["Tests"]) }] }),
});

const getPaginationUpdateLink = (newPagination: TablePaginationVars) =>
  linkOptions({
    from: Route.id,
    to: ".",
    search: (prev): typeof prev => ({
      ...prev,
      testTable: {
        ...prev.testTable,
        ...newPagination,
      },
    }),
  });

function RouteComponent() {
  const navigate = Route.useNavigate();
  const { testsData, pageSize, pageInfo } = Route.useLoaderData();
  const onFilterChange = (where: TargetWhereInput[]) => {
    navigate({
      search: (prev): typeof prev => ({
        ...prev,
        testTable: {
          ...prev.testTable,
          where,
          pagination: undefined,
        },
      }),
    });
  };
  return (
    <TestsPage
      testsData={testsData}
      pageSize={pageSize}
      pageInfo={pageInfo}
      getPaginationUpdateLink={getPaginationUpdateLink}
      onFilterChange={onFilterChange}
    />
  );
}
