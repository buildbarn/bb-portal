import { gql } from "@/graphql/__generated__";

const FIND_BAZEL_INVOCATIONS_QUERY = gql(/* GraphQL */ `
  query FindBazelInvocations(
    $first: Int!
    $orderBy: BazelInvocationOrder
    $where: BazelInvocationWhereInput
  ) {
    findBazelInvocations(first: $first, orderBy: $orderBy, where: $where) {
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
    endedAt
    state {
      bepCompleted
      exitCode {
        code
        name
      }
    }
    build {
      buildUUID
    }
  }
`);

export default FIND_BAZEL_INVOCATIONS_QUERY;
