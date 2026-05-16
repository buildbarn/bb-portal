import { gql } from "@/graphql/__generated__";

export const ARTIFACT_GRAPH_QUERY = gql(/* GraphQL */ `
  query ArtifactGraph($id: UUID!) {
    getBazelInvocation(invocationID: $id) {
      id
      artifactGraph {
        payload
        uncompressedSize
      }
    }
  }
`);
