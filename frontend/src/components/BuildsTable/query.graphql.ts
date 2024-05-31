import { gql } from "@/graphql/__generated__";

const FIND_BUILDS_QUERY = gql(/* GraphQL */ `
  query FindBuilds(
    $first: Int!
    $where: BuildWhereInput
  ) {
    findBuilds(first: $first, where: $where) {
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
  }
`);

export default FIND_BUILDS_QUERY;
