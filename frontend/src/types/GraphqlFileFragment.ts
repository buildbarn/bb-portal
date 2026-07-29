import { gql } from "@/graphql/__generated__";

export const FILE_DETAILS_FRAGMENT = gql(/* GraphQL */ `
  fragment FileDetails on File {
    id
    filePath {
      id
      path
    }
    digest {
      id
      rev2InstanceName
      digestFunction
      hash
      sizeBytes
    }
  }
`);
