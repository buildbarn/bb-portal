import { gql } from "@/graphql/__generated__";

export const GET_TARGETS_DATA = gql(/* GraphQl */`
query GetTargetsWithOffset(
  $label: String,
  $offset: Int,
  $limit: Int,
  $sortBy: String,
  $direction: String) {
    getTargetsWithOffset(
      label: $label
      offset: $offset
      limit: $limit
      sortBy: $sortBy
      direction: $direction
    ) {
      total
      result {
        label
        sum
        min
        max
        avg
        count
        passRate
      }
    }
  }
`);

export const FIND_TARGETS = gql(/* GraphQL */ `
  query FindTargets(
     $first: Int!
     $where: TargetPairWhereInput
     $orderBy: TargetPairOrder
     $after: Cursor
   ){
   findTargets (first: $first, where: $where, orderBy: $orderBy, after: $after){
     totalCount
     pageInfo{
       startCursor
       endCursor
       hasNextPage
       hasPreviousPage
     }
     edges {
       node {
         id
         durationInMs
         label
         success
         bazelInvocation {
           invocationID
         }
       }
     }
   }
 }
 `);

export default GET_TARGETS_DATA;
