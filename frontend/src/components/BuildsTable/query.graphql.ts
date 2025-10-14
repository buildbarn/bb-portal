import { gql } from "@/graphql/__generated__";

const FIND_BUILDS_QUERY = gql(/* GraphQL */ `
  query FindBuilds(
    $first: Int!
    $orderBy: BuildOrder
    $where: BuildWhereInput
  ) {
    findBuilds(first: $first, orderBy: $orderBy, where: $where) {
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

export default FIND_BUILDS_QUERY;
