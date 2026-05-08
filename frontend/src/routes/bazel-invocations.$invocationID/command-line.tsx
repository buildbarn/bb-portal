import { createFileRoute, notFound } from "@tanstack/react-router";
import { apolloClient } from "@/components/ApolloWrapper";
import CommandLineDisplay from "@/components/CommandLine";
import { getFragmentData, gql } from "@/graphql/__generated__";
import { commandLineDataToString } from "@/utils/commandLineDataToString";
import { generatePageTitle } from "@/utils/generatePageTitle";

const GET_BAZEL_INVOCATION_COMMANDLINE = gql(/* GraphQL */ `
  query GetBazelInvocationCommandline($invocationID: UUID!) {
    getBazelInvocation(invocationID: $invocationID) {
      ...BazelInvocationCommandline
    }
  }
`);

const BAZEL_INVOCATION_COMMANDLINE_FRAGMENT = gql(/* GraphQL */ `
  fragment BazelInvocationCommandline on BazelInvocation {
    id
    invocationID
    canonicalCommandLine
    originalCommandLine
    optionsParsed
    environmentVariables
  }
`);

export const Route = createFileRoute(
  "/bazel-invocations/$invocationID/command-line",
)({
  component: RouteComponent,
  loader: async ({ params }) => {
    const { data } = await apolloClient.query({
      query: GET_BAZEL_INVOCATION_COMMANDLINE,
      variables: { invocationID: params.invocationID },
      fetchPolicy: "network-only",
    });

    if (!data?.getBazelInvocation) {
      throw notFound();
    }

    return {
      invocation: getFragmentData(
        BAZEL_INVOCATION_COMMANDLINE_FRAGMENT,
        data.getBazelInvocation,
      ),
    };
  },
  head: (_ctx) => ({
    meta: [
      {
        title: generatePageTitle([
          "Invocation",
          "Commandline",
          _ctx.params.invocationID,
        ]),
      },
    ],
  }),
});

function RouteComponent() {
  const { invocation } = Route.useLoaderData();
  const command = commandLineDataToString(invocation.originalCommandLine);
  return (
    <CommandLineDisplay
      parsedOptions={invocation.optionsParsed}
      rawCommand={command}
      canonicalCommandLine={invocation.canonicalCommandLine}
      environmentVariables={invocation.environmentVariables}
    />
  );
}
