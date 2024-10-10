import { gql } from '@/graphql/__generated__';

export const FIND_TESTS_WITH_CACHE = gql(/* GraphQL */ `
    query FindTestsWithCache(
       $first: Int!
       $where: TestCollectionWhereInput
       $orderBy: TestCollectionOrder
       $after: Cursor
     ){
     findTests (first: $first, where: $where, orderBy: $orderBy, after: $after){
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
           durationMs
           firstSeen
           label
           overallStatus
           cachedLocally
           cachedRemotely
           bazelInvocation {
             invocationID
           }
         }
       }
     }
   }
   `);