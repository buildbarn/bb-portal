import { gql } from "@/graphql/__generated__";

export const CHECK_IF_INVOCATION_EXISTS = gql(/* GraphQl */ `
  query CheckIfInvocationExists(
    $invocationID: UUID!
  ){
    getBazelInvocation(invocationID: $invocationID){
      id
    }
  }
`);
