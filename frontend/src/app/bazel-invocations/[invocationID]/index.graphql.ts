import {gql} from "@/graphql/__generated__";

export const LOAD_FULL_BAZEL_INVOCATION_DETAILS = gql(/* GraphQL */ `
  query LoadFullBazelInvocationDetails($invocationID: String!) {
    bazelInvocation(invocationId: $invocationID) {
      ...FullBazelInvocationDetails
    }
  }
`);

export const BAZEL_INVOCATION_FRAGMENT = gql(/* GraphQL */ `
fragment BazelInvocationInfo on BazelInvocation {
  bazelCommand {
    command
    executable
    id
    buildOptions: options
    residual
  }
  id
  invocationID
  build {
    id
    buildUUID
  }
  relatedFiles {
    name
    url
  }
  startedAt
  endedAt
  state {
    bepCompleted
    buildEndTime
    buildStartTime
    exitCode {
      code
      id
      name
    }
    id
  }
  stepLabel
}
`);
export const PROBLEM_INFO_FRAGMENT = gql(/* GraphQL */ `
 fragment ProblemInfo on Problem {
  id
  label
  __typename
  ... on ActionProblem {
    __typename
    id
    label
    type
    stdout {
      ...BlobReferenceInfo
    }
    stderr {
      ...BlobReferenceInfo
    }
  }
  ... on TestProblem {
    __typename
    id
    label
    status
    results {
      __typename
      id
      run
      shard
      attempt
      status
      actionLogOutput {
        ...BlobReferenceInfo
      }
      undeclaredTestOutputs {
        ...BlobReferenceInfo
      }
    }
  }
  ... on TargetProblem {
    __typename
    id
    label
  }
  ... on ProgressProblem {
    __typename
    id
    output
    label
  }
}
`);

export const BLOB_REFERENCE_INFO_FRAGMENT = gql(/* GraphQL */ `
fragment BlobReferenceInfo on BlobReference {
  availabilityStatus
  name
  sizeInBytes
  downloadURL
}
`)


export const FULL_BAZEL_INVOCATION_DETAILS = gql(/* GraphQL */ `
    fragment FullBazelInvocationDetails on BazelInvocation {
      problems {
        ...ProblemInfo
      }
      ...BazelInvocationInfo
    }
`);




export const GET_ACTION_PROBLEM = gql(/* GraphQL */ `
  query GetActionProblem($id: ID!) {
    node(id: $id) {
      id
      ... on ActionProblem {
        label
        stdout {
          ...BlobReferenceInfo
        }
        stderr {
          ...BlobReferenceInfo
        }
      }
    }
  }
`);

export const TEST_RESULT_FRAGMENT = gql(/* GraphQL */`
fragment TestResultInfo on TestResult {
      actionLogOutput {
  ...BlobReferenceInfo
  }
  attempt
  run
  shard
  status
  undeclaredTestOutputs {
    ...BlobReferenceInfo
  }
}`)