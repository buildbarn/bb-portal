import { gql } from "@/graphql/__generated__";

export const GET_TARGET_DETAILS = gql(/* GraphQl */ `
  query GetTargetDetails(
    $instanceName: String!
    $label: String!
    $aspect: String!
    $targetKind: String!

    $after: Cursor
    $first: Int
    $before: Cursor
    $last: Int
    $orderBy: InvocationTargetOrder
    $where: InvocationTargetWhereInput
  ){
    getTarget (instanceName: $instanceName, label: $label, aspect: $aspect, targetKind: $targetKind){
      invocationTargetsTotalDurationMillis
      invocationTargets(after: $after, first: $first, before: $before, last: $last, orderBy: $orderBy, where: $where) {
        pageInfo {
          startCursor
          endCursor
          hasNextPage
          hasPreviousPage
        }
        totalCount
        edges {
          node {
            id
            success
            durationInMs
            abortReason
            failureMessage
            tags
            bazelInvocation {
              invocationID
            }
          }
        }
      }
    }
  }
`);
