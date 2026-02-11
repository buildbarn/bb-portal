import { gql } from "@/graphql/__generated__";

export const FIND_BUILDS_QUERY = gql(/* GraphQL */ `
  query FindBuilds(
    $after: Cursor
    $first: Int
    $before: Cursor
    $last: Int
    $orderBy: BuildOrder
    $where: BuildWhereInput
  ) {
    findBuilds(after: $after, first: $first, before: $before, last: $last, orderBy: $orderBy, where: $where) {
      pageInfo {
        startCursor
        endCursor
        hasNextPage
        hasPreviousPage
      }
      edges {
        node {
          ...BuildNode
        }
      }
    }
  }
`);

export const BUILD_NODE_FRAGMENT = gql(/* GraphQL */ `
  fragment BuildNode on Build {
    id
    buildUUID
    buildURL
    timestamp
  }
`);
