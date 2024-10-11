import { gql } from "@/graphql/__generated__";

const FIND_BUILD_DURATIONS = gql(/* GraphQL */ `
  query FindBuildTimes(
    $first: Int!
  	$where: BazelInvocationWhereInput
  ) {
    findBazelInvocations(first: $first, where: $where ) {
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
