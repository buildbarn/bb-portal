import { gql } from "@/graphql/__generated__";

export const GET_TESTS_FOR_TARGET = gql(/* GraphQl */ `
  query GetTestsForTarget(
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
