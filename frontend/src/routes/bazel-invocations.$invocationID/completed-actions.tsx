import { createFileRoute, linkOptions } from "@tanstack/react-router";
import { z } from "zod";
import { apolloClient } from "@/components/ApolloWrapper";
import { DEFAULT_PAGE_SIZE } from "@/components/PageCursorTable";
import type { TablePaginationVars } from "@/components/PageCursorTable/types";
import { TablePaginationVarsSchema } from "@/components/PageCursorTable/types";
import { CompletedAction } from "@/components/RemoteInvocationsTab";
import { gql } from "@/graphql/__generated__";
import type { CompletedActionWhereInput } from "@/graphql/__generated__/graphql";
import { CompletedActionWhereInputSchema } from "@/graphql/__generated__/zod";
import { NotFoundError } from "@/main";
import { generatePageTitle } from "@/utils/generatePageTitle";
import { parseGraphqlEdgeListWithFragment } from "@/utils/parseGraphqlEdgeList";

const GET_COMPLETED_ACTIONS_LIST = gql(/* GraphQL */ `
  query GetCompletedActionsList(
    $after: Cursor
    $first: Int
    $before: Cursor
    $last: Int
    $where: CompletedActionWhereInput
  ) {
    findCompletedActions(
      after: $after
      first: $first
      before: $before
      last: $last
      where: $where
    ) {
      pageInfo {
        startCursor
        endCursor
        hasNextPage
        hasPreviousPage
      }
      edges {
        node {
          ...CompletedActionDetails
        }
      }
    }
  }
`);

export const COMPLETED_ACTION_DETAILS_FRAGMENT = gql(`
  fragment CompletedActionDetails on CompletedAction {
    id
    targetID
    invocationID
    actionID
    toolName
    mnemonic
    actionDigest {
      hash
      sizeBytes
      digestFunction
      rev2InstanceName
    }
    historicalExecuteResponseDigest {
      hash
      sizeBytes
      digestFunction
      rev2InstanceName
    }
  }
`);

const CompletedActionSearchSchema = z.object({
  completedActionTable: TablePaginationVarsSchema.extend({
    where: z.array(CompletedActionWhereInputSchema().partial()).optional(),
  }).optional(),
});

export const Route = createFileRoute(
  "/bazel-invocations/$invocationID/completed-actions",
)({
  component: RouteComponent,
  validateSearch: (search) => CompletedActionSearchSchema.parse(search),
  loaderDeps: ({ search: { completedActionTable } }) => ({
    completedActionTable,
  }),
  loader: async ({ params, deps }) => {
    const pageSize = deps.completedActionTable?.pageSize ?? DEFAULT_PAGE_SIZE;
    const pagination = deps.completedActionTable?.pagination ?? {
      first: pageSize,
    };

    const where: CompletedActionWhereInput[] =
      deps.completedActionTable?.where ?? [];

    const { data, error } = await apolloClient.query({
      errorPolicy: "all",
      query: GET_COMPLETED_ACTIONS_LIST,
      variables: {
        ...pagination,
        where: {
          and: [{ invocationID: params.invocationID }, ...where],
        },
      },
      fetchPolicy: "cache-first",
    });

    if (!data?.findCompletedActions) {
      throw new NotFoundError("completed actions", error?.message);
    }
    const actions = parseGraphqlEdgeListWithFragment(
      COMPLETED_ACTION_DETAILS_FRAGMENT,
      data.findCompletedActions,
    );

    return { actions, pageSize, pageInfo: data.findCompletedActions.pageInfo };
  },

  head: (_ctx) => ({
    meta: [
      {
        title: generatePageTitle([
          "Invocation",
          "Completed Actions",
          _ctx.params.invocationID,
        ]),
      },
    ],
  }),
});

const getPaginationUpdateLink = (newPagination: TablePaginationVars) =>
  linkOptions({
    from: Route.id,
    to: ".",
    search: (prev): typeof prev => ({
      ...prev,
      completedActionTable: {
        ...prev.completedActionTable,
        ...newPagination,
      },
    }),
  });

function RouteComponent() {
  const { actions, pageSize, pageInfo } = Route.useLoaderData();
  return (
    <CompletedAction
      actions={actions}
      getPaginationUpdateLink={getPaginationUpdateLink}
      pageSize={pageSize}
      pageInfo={pageInfo}
    />
  );
}
