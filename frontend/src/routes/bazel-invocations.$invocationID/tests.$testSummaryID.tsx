import { createFileRoute } from "@tanstack/react-router";
import { apolloClient } from "@/components/ApolloWrapper";
import { TestSummaryPage } from "@/components/pages/TestSummary";
import { getFragmentData, gql } from "@/graphql/__generated__";
import type { InvocationTestResultDetailsFragment } from "@/graphql/__generated__/graphql";
import { NotFoundError } from "@/main";
import { env } from "@/utils/env";
import { requireFeature } from "@/utils/featureGuard";
import { generatePageTitle } from "@/utils/generatePageTitle";

const GET_BAZEL_INVOCATION_TAGS = gql(/* GraphQL */ `
  query GetTestSummary($invocationID: UUID!, $testSummaryID: String!) {
    getTestSummary(invocationID: $invocationID, testSummaryID: $testSummaryID) {
      ...InvocationTestSummaryDetails
    }
  }
`);

const INVOCATION_TEST_SUMMARY_DETAILS_FRAGMENT = gql(/* GraphQL */ `
  fragment InvocationTestSummaryDetails on TestSummary {
    id
    overallStatus
    totalRunCount
    runCount
    attemptCount
    shardCount
    totalNumCached
    firstStartTime
    lastStopTime
    totalRunDurationInMs
    testResults {
      ...InvocationTestResultDetails
    }
    invocationTarget {
      id
      target {
        id
        targetKind
        label
        aspect
        instanceName {
          id
          name
        }
      }
    }
  }
`);

export const INVOCATION_TEST_RESULT_DETAILS_FRAGMENT = gql(/* GraphQL */ `
  fragment InvocationTestResultDetails on TestResult {
    id
    run
    shard
    attempt
    status
    testAttemptDurationInMs
    cachedLocally
    cachedRemotely
    exitCode
    strategy
    testActionOutput {
      ...FileDetails
    }
  }
`);

export const Route = createFileRoute(
  "/bazel-invocations/$invocationID/tests/$testSummaryID",
)({
  component: RouteComponent,
  beforeLoad: requireFeature(env.featureFlags?.bes?.pageTargets),
  loader: async ({ params }) => {
    const { data, error } = await apolloClient.query({
      errorPolicy: "all",
      query: GET_BAZEL_INVOCATION_TAGS,
      variables: {
        invocationID: params.invocationID,
        testSummaryID: params.testSummaryID,
      },
      fetchPolicy: "network-only",
    });

    if (!data?.getTestSummary) {
      throw new NotFoundError("test summary", error?.message);
    }

    const testSummary = getFragmentData(
      INVOCATION_TEST_SUMMARY_DETAILS_FRAGMENT,
      data.getTestSummary,
    );
    const testResults: InvocationTestResultDetailsFragment[] =
      testSummary.testResults?.map((tr) =>
        getFragmentData(INVOCATION_TEST_RESULT_DETAILS_FRAGMENT, tr),
      ) ?? [];

    return { testSummary, testResults };
  },
  head: (_ctx) => ({
    meta: [
      {
        title: generatePageTitle([
          "Invocation",
          "Tests",
          _ctx.params.testSummaryID,
          _ctx.params.invocationID,
        ]),
      },
    ],
  }),
});

function RouteComponent() {
  const { testSummary, testResults } = Route.useLoaderData();
  return (
    <TestSummaryPage testSummary={testSummary} testResults={testResults} />
  );
}
