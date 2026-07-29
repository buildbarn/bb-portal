import { createFileRoute } from "@tanstack/react-router";
import { ActionsTab } from "@/components/ActionsTab";
import { apolloClient } from "@/components/ApolloWrapper";
import { InvocationDataNotFoundAlert } from "@/components/pages/InvocationDataNotFoundAlert";
import { getFragmentData, gql } from "@/graphql/__generated__";
import { NotFoundError } from "@/main";
import { generatePageTitle } from "@/utils/generatePageTitle";

const GET_BAZEL_INVOCATION_ACTIONS = gql(/* GraphQL */ `
  query GetBazelInvocationActions($invocationID: UUID!) {
    getBazelInvocation(invocationID: $invocationID) {
      id
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
    configuration {
      id
      configurationID
      mnemonic
      platformName
      cpu
      makeVariables
    }
    stdout {
      ...FileDetails
    }
    stderr {
      ...FileDetails
    }
  }
`);

export const Route = createFileRoute(
  "/bazel-invocations/$invocationID/actions",
)({
  component: RouteComponent,
  loader: async ({ params }) => {
    const { data, error } = await apolloClient.query({
      errorPolicy: "all",
      query: GET_BAZEL_INVOCATION_ACTIONS,
      variables: { invocationID: params.invocationID },
      fetchPolicy: "network-only",
    });

    if (!data?.getBazelInvocation) {
      throw new NotFoundError("invocation", error?.message);
    }

    if (!data.getBazelInvocation.actions) {
      return { actions: undefined };
    }

    const actions = getFragmentData(
      BAZEL_INVOCATION_ACTIONS_FRAGMENT,
      data.getBazelInvocation?.actions,
    );

    return { actions };
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
  const { actions } = Route.useLoaderData();

  if (actions === undefined || actions.length === 0) {
    return <InvocationDataNotFoundAlert type="actions" />;
  }

  // TODO (isakstenstrom): Maybe we should fetch the logs here instead?
  return <ActionsTab actions={actions} />;
}
