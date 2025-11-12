import { gql } from "@/graphql/__generated__";

const GET_AUTHENTICATED_USER_BY_UUID = gql(/* GraphQL */ `
  query GetAuthenticatedUser(
    $userUUID: UUID
    $bazelInvocationsOrderBy: BazelInvocationOrder
  ) {
    getAuthenticatedUser(userUUID: $userUUID) {
      ...AuthenticatedUserNodeFragment
    }
  }
`);

export const AUTHENTICATED_USER_NODE_FRAGMENT = gql(/* GraphQL */ `
  fragment AuthenticatedUserNodeFragment on AuthenticatedUser {
    displayName
    userInfo
    bazelInvocations(orderBy: $bazelInvocationsOrderBy) {
      edges {
        node {
          invocationID
          build {
            id
            buildUUID
          }
          endedAt
          startedAt
          state {
            bepCompleted
            exitCode {
              code
              name
            }
          }
        }
      }
    }
  }
`);

export default GET_AUTHENTICATED_USER_BY_UUID;
