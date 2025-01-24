import { gql } from "@/graphql/__generated__";

export const GET_TEST_GRID_DATA = gql(/* GraphQl */`
query GetTestsWithOffset(
  $label: String,
  $offset: Int,
  $limit: Int,
  $sortBy: String,
  $direction: String) {
    getTestsWithOffset(
      label: $label
      offset: $offset
      limit: $limit
      sortBy: $sortBy
      direction: $direction
    ) {
      total
      result {
        label
        sum
        min
        max
        avg
        count
        passRate
      }
    }
  }
`);


export const GET_AVERAGE_PASS_PERCENTAGE_FOR_LABEL = gql(/* GraphQL */ `

  query GetAveragePassPercentageForLabel(
    $label: String!
  ) {
    getAveragePassPercentageForLabel(label:$label)
  }

`);

export const GET_TEST_DURATION_AGGREGATION = gql(/* GraphQL */ `
  query GetTestDurationAggregation(
    $label: String
  ) {
    getTestDurationAggregation(label:$label) {
      label
      count
      sum
      min
      max
    }
  }
`);

export const FIND_TESTS = gql(/* GraphQL */ `
 query FindTests(
    $first: Int!
    $where: TestCollectionWhereInput
    $orderBy: TestCollectionOrder
    $after: Cursor
  ){
  findTests (first: $first, where: $where, orderBy: $orderBy, after: $after){
    totalCount
    pageInfo{
      startCursor
      endCursor
      hasNextPage
      hasPreviousPage
    }
    edges {
      node {
        id
        durationMs
        firstSeen
        label
        overallStatus
        bazelInvocation {
          invocationID
        }
      }
    }
  }
}
`);

export default FIND_TESTS;
