import { createFileRoute, linkOptions } from "@tanstack/react-router";
import z from "zod";
import { apolloClient } from "@/components/ApolloWrapper";
import { cacheLocationFromTestResults } from "@/components/CacheLocationTag";
import { DEFAULT_PAGE_SIZE } from "@/components/PageCursorTable";
import {
  type TablePaginationVars,
  TablePaginationVarsSchema,
} from "@/components/PageCursorTable/types";
import { TestDetailsPage } from "@/components/pages/TestDetails";
import { getFragmentData, gql } from "@/graphql/__generated__";
import { NotFoundError } from "@/main";
import { env } from "@/utils/env";
import { requireFeature } from "@/utils/featureGuard";
import { generatePageTitle } from "@/utils/generatePageTitle";
import { parseGraphqlEdgeListWithFragment } from "@/utils/parseGraphqlEdgeList";

export const GET_TEST_DETAILS = gql(/* GraphQL */ `
  query GetTestDetails(
    $targetID: ID!
    $after: Cursor
    $first: Int
    $before: Cursor
    $last: Int
    $orderBy: TestSummaryOrder
    $where: TestSummaryWhereInput
  ){
    getTarget(id: $targetID ) {
      ...TestSummaryTargetDetails
    }
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
          ...TestSummaryRow
        }
      }
    }
  }
`);

const TEST_SUMMARY_TARGET_DETAILS_FRAGMENT = gql(/* GraphQL */ `
  fragment TestSummaryTargetDetails on Target {
    id
    aspect
    instanceName {
      id
      name
    }
    label
    targetKind
  } 
`);

const TEST_SUMMARY_ROW_FRAGMENT = gql(/* GraphQL */ `
  fragment TestSummaryRow on TestSummary {
    id
    overallStatus
    runCount
    attemptCount
    shardCount
    firstStartTime
    totalRunDurationInMs
    testResults {
      id
      cachedLocally
      cachedRemotely
    }
    invocationTarget {
      id
      bazelInvocation {
        id
        invocationID
      }
    }
  } 
`);

const TargetTestSearchSchema = z.object({
  targetTestTable: TablePaginationVarsSchema.optional(),
});

export const Route = createFileRoute("/targets/$targetID/tests")({
  component: RouteComponent,
  validateSearch: (search) => TargetTestSearchSchema.parse(search),
  beforeLoad: requireFeature(env.featureFlags?.bes?.pageTests),
  loaderDeps: ({ search: { targetTestTable } }) => ({ targetTestTable }),
  loader: async ({ params, deps }) => {
    // We set the defaults here instead of in the validate search function, as
    // that updates the URL and we don't want to do that on initial load.
    const pageSize = deps.targetTestTable?.pageSize ?? DEFAULT_PAGE_SIZE;
    const pagination = deps.targetTestTable?.pagination ?? {
      first: pageSize,
    };
    const { data, error } = await apolloClient.query({
      errorPolicy: "all",
      query: GET_TEST_DETAILS,
      variables: {
        targetID: params.targetID,
        where: {
          hasInvocationTargetWith: [
            {
              hasTargetWith: [
                {
                  id: params.targetID,
                },
              ],
            },
          ],
        },
        ...pagination,
      },
      fetchPolicy: "network-only",
    });
    if (!data?.getTarget || !data?.findTestSummaries) {
      throw new NotFoundError("test", error?.message);
    }
    const target = getFragmentData(
      TEST_SUMMARY_TARGET_DETAILS_FRAGMENT,
      data.getTarget,
    );
    const testSummaries = parseGraphqlEdgeListWithFragment(
      TEST_SUMMARY_ROW_FRAGMENT,
      data.findTestSummaries,
    ).map((ts) => ({
      ...ts,
      cacheLocation: cacheLocationFromTestResults(ts.testResults),
    }));

    return {
      target,
      testSummaries,
      pageSize,
      pageInfo: data.findTestSummaries.pageInfo,
    };
  },
  head: (_ctx) => {
    const label = _ctx.loaderData?.target.label;
    if (label === undefined) {
      return { meta: [{ title: generatePageTitle(["Test", "Not Found"]) }] };
    }
    return { meta: [{ title: generatePageTitle(["Test", label]) }] };
  },
});

const getPaginationUpdateLink = (newPagination: TablePaginationVars) =>
  linkOptions({
    from: Route.id,
    to: ".",
    search: (prev): typeof prev => ({
      ...prev,
      targetTestTable: {
        ...prev.targetTestTable,
        ...newPagination,
      },
    }),
  });

function RouteComponent() {
  const { target, testSummaries, pageSize, pageInfo } = Route.useLoaderData();

  return (
    <TestDetailsPage
      target={target}
      getPaginationUpdateLink={getPaginationUpdateLink}
      pageInfo={pageInfo}
      pageSize={pageSize}
      testSummaries={testSummaries}
    />
  );
}
