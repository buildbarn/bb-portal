import { gql } from "@/graphql/__generated__";

export const GET_INVOCATION_TARGETS_FOR_TARGET = gql(/* GraphQl */ `
  query GetInvocationTargetsForTarget(
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
      invocationTargets(after: $after, first: $first, before: $before, last: $last, orderBy: $orderBy, where: $where) {
        edges {
          node {
            id
            success
            bazelInvocation {
              invocationID
            }
          }
        }
      }
    }
  }
`);
