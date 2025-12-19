import { gql } from "@/graphql/__generated__";

export const GET_TEST_DETAILS = gql(/* GraphQL */ `
    query GetTestDetails(
    $after: Cursor
    $first: Int
    $before: Cursor
    $last: Int
    $orderBy: TestSummaryOrder
    $where: TestSummaryWhereInput
  ){
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
          id
          overallStatus
          runCount
          attemptCount
          shardCount
          firstStartTime
          totalRunDurationInMs
          testResults {
            cachedLocally
            cachedRemotely
          }
          invocationTarget {
            bazelInvocation {
              invocationID
            }
          }
        }
      }
    }
  }
`);
