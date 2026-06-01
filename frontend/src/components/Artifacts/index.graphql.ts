import { gql } from "@/graphql/__generated__";

export const ARTIFACT_GRAPH_QUERY = gql(/* GraphQL */ `
  query ArtifactGraph($id: UUID!) {
    getBazelInvocation(invocationID: $id) {
      id
      artifactGraph {
        namedSets {
          id
          childSetIds
          files {
            name
            uri
            digest
            sizeBytes
            downloadUrl
          }
        }
        targets {
          label
          aspect
          outputGroups {
            name
            incomplete
            rootSetIds
          }
        }
      }
    }
  }
`);
