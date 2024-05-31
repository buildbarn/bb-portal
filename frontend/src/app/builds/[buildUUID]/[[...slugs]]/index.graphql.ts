import { gql } from '@/graphql/__generated__';

export const FIND_BUILD_BY_UUID_QUERY = gql(/* GraphQL */ `
  query FindBuildByUUID($url: String, $uuid: UUID) {
    getBuild(buildURL: $url, buildUUID: $uuid) {
      id
      buildURL
      buildUUID
      invocations {
        ...FullBazelInvocationDetails
      }
      env {
        key
        value
      }
    }
  }
`);


