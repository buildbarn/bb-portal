import { gql } from '@/graphql/__generated__';

export const FIND_BUILD_BY_UUID_QUERY = gql(/* GraphQL */ `
  query FindBuildByUUID($url: String, $uuid: UUID) {
    getBuild(buildURL: $url, buildUUID: $uuid) {
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
        state {
          exitCode {
            name
          }
        }
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


