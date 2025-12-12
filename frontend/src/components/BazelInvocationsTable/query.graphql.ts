import { gql } from "@/graphql/__generated__";

const FIND_BAZEL_INVOCATIONS_QUERY = gql(/* GraphQL */ `
  query FindBazelInvocations(
    $after: Cursor
    $first: Int
    $before: Cursor
    $last: Int
    $orderBy: BazelInvocationOrder
    $where: BazelInvocationWhereInput
  ) {
    findBazelInvocations(after: $after, first: $first, before: $before, last: $last, orderBy: $orderBy, where: $where) {
      pageInfo {
        startCursor
        endCursor
        hasNextPage
        hasPreviousPage
      }
      edges {
        node {
          ...BazelInvocationNode
        }
      }
    }
  }
`);

export const BAZEL_INVOCATION_NODE_FRAGMENT = gql(/* GraphQL */ `
  fragment BazelInvocationNode on BazelInvocation {
    id
    invocationID
    startedAt
    user {
      Email
      LDAP
    }
    authenticatedUser {
      userUUID
      displayName
    }
    endedAt
    exitCodeName
    bepCompleted
    build {
      buildUUID
    }
  }
`);

export default FIND_BAZEL_INVOCATIONS_QUERY;
