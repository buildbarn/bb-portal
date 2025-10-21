import { gql } from "@/graphql/__generated__";

export const GET_INVOCATION_TARGETS_FOR_INVOCATION = gql(/* GraphQl */ `
  query GetInvocationTargetsForInvocation(
    $invocationID: String!
    $after: Cursor
    $first: Int
    $before: Cursor
    $last: Int
    $orderBy: InvocationTargetOrder
    $where: InvocationTargetWhereInput
  ){
    bazelInvocation(invocationId: $invocationID) {
      invocationTargets(after: $after, first: $first, before: $before, last: $last, orderBy: $orderBy, where: $where){
        pageInfo {
          startCursor
          endCursor
          hasNextPage
          hasPreviousPage
        }
        edges {
          node {
            id
            success
            abortReason
            durationInMs
            failureMessage
            tags
            target {
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
      numTotal: invocationTargets {
        totalCount
      }
      numSuccessful: invocationTargets(where: { success: true }) {
        totalCount
      }
      numSkipped: invocationTargets(where: {abortReason: SKIPPED}) {
        totalCount
      }
    }
  }
`);
