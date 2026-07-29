import { gql } from "@/graphql/__generated__";

export const GET_INVOCATION_TARGETS_FOR_TARGET = gql(/* GraphQl */ `
  query GetInvocationTargetsForTarget(
    $targetId: ID!

    $after: Cursor
    $first: Int
    $before: Cursor
    $last: Int
    $where: InvocationTargetWhereInput
  ){
    getTarget (id: $targetId){
      invocationTargets(after: $after, first: $first, before: $before, last: $last, where: $where) {
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
