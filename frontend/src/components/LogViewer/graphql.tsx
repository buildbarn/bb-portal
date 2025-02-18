import { gql } from "@/graphql/__generated__";

export const GET_BUILD_LOGS = gql(/* GraphQl */ `
query GetBuildLogs ($invocationId: String!){
  bazelInvocation(invocationId: $invocationId){
    invocationID
    buildLogs    
  }
}
`);
