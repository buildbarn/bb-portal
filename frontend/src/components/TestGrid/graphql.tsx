import { gql } from "@/graphql/__generated__";

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
          id
          label
          aspect
          targetKind
          instanceName {
            name
          }
        }
      }
    }
  }
`);
