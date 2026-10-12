import { gql } from "@/graphql/__generated__";

const FIND_BUILD_DURATIONS = gql(/* GraphQL */ `
  query FindBuildTimes(
    $first: Int!
  	$where: BazelInvocationWhereInput
    $orderBy: BazelInvocationOrder
  ) {
    findBazelInvocations(first: $first, where: $where, orderBy: $orderBy) {
      pageInfo{
        startCursor
        endCursor
        hasNextPage
        hasPreviousPage

      }
      totalCount
      edges {
        node {
          invocationID
          startedAt
          endedAt
        }
      }
    }
  }
`);

export default FIND_BUILD_DURATIONS;
