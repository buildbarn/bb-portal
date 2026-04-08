import { gql } from "@/graphql/__generated__";

export const GET_BUILD_BY_UUID_QUERY = gql(/* GraphQL */ `
  query FindBuildByUUID(
    $buildUUID: UUID!
    $after: Cursor
    $first: Int
    $before: Cursor
    $last: Int
    $orderBy: BazelInvocationOrder
    $where: BazelInvocationWhereInput
  ) {
    getBuild(buildUUID: $buildUUID) {
      id
      buildURL
      buildUUID
      timestamp
      invocations(after: $after, first: $first, before: $before, last: $last, orderBy: $orderBy, where: $where) {
        pageInfo {
          startCursor
          endCursor
          hasNextPage
          hasPreviousPage
        }
        edges {
          node {
            ...GetBuildInvocation
          }
        }
      }
    }
  }
`);

export const GET_BUILD_INVOCATION_FRAGMENT = gql(/* GraphQL */ `
  fragment GetBuildInvocation on BazelInvocation {
    id
    invocationID
    userLdap
    endedAt
    startedAt
    exitCodeName
    sourceControl{
      job
      action
      workflow
      runnerName
    }
    connectionMetadata {
      connectionLastOpenAt
      timeSinceLastConnectionMillis
    }
  }
`);
