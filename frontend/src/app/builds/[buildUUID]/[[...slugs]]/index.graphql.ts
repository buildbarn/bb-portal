import { gql } from "@/graphql/__generated__";

export const FIND_BUILD_BY_UUID_QUERY = gql(/* GraphQL */ `
  query FindBuildByUUID($uuid: UUID) {
    getBuild(buildUUID: $uuid) {
      id
      buildURL
      buildUUID
      timestamp
      invocations {
        id
        invocationID
        userLdap
        endedAt
        startedAt
        exitCodeName
        bepCompleted
        sourceControl{
          job
          action
          workflow
          runnerName
        }
      }
    }
  }
`);
