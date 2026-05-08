import { createFileRoute, notFound } from "@tanstack/react-router";
import { ActionsTab } from "@/components/ActionsTab";
import { apolloClient } from "@/components/ApolloWrapper";
import { getFragmentData, gql } from "@/graphql/__generated__";
import { generatePageTitle } from "@/utils/generatePageTitle";

const GET_BAZEL_INVOCATION_ACTIONS = gql(/* GraphQL */ `
  query GetBazelInvocationActions($invocationID: UUID!) {
    getBazelInvocation(invocationID: $invocationID) {
      id
      instanceName {
        id
        name
      }
      actions {
        ...BazelInvocationActions
      }
    }
  }
`);

const BAZEL_INVOCATION_ACTIONS_FRAGMENT = gql(/* GraphQL */ `
  fragment BazelInvocationActions on Action {
    id
    label
    type
    success
    exitCode
    commandLine
    startTime
    endTime
    failureCode
    failureMessage
    stdoutHash
    stdoutSizeBytes
    stdoutHashFunction
    stderrHash
    stderrSizeBytes
    stderrHashFunction
    configuration {
      id
      configurationID
      mnemonic
      platformName
      cpu
      makeVariables
    }
  }
`);

export const Route = createFileRoute(
  "/bazel-invocations/$invocationID/actions",
)({
  component: RouteComponent,
  loader: async ({ params }) => {
    const { data } = await apolloClient.query({
      query: GET_BAZEL_INVOCATION_ACTIONS,
      variables: { invocationID: params.invocationID },
      fetchPolicy: "network-only",
    });

    if (!data?.getBazelInvocation) {
      throw notFound();
    }

    if (!data.getBazelInvocation.actions) {
      throw notFound();
    }

    const actions = getFragmentData(
      BAZEL_INVOCATION_ACTIONS_FRAGMENT,
      data.getBazelInvocation?.actions,
    );
    return { instanceName: data.getBazelInvocation.instanceName.name, actions };
  },
  head: (_ctx) => ({
    meta: [
      {
        title: generatePageTitle([
          "Invocation",
          "Failed actions",
          _ctx.params.invocationID,
        ]),
      },
    ],
  }),
});

function RouteComponent() {
  const { instanceName, actions } = Route.useLoaderData();
  // TODO (isakstenstrom): Maybe we should fetch the logs here instead?
  return <ActionsTab instanceName={instanceName} actions={actions} />;
}
