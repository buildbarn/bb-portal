import { createFileRoute, linkOptions } from "@tanstack/react-router";
import z from "zod";
import { ActionsTab } from "@/components/ActionsTab";
import { apolloClient } from "@/components/ApolloWrapper";
import { DEFAULT_PAGE_SIZE } from "@/components/PageCursorTable";
import {
  type TablePaginationVars,
  TablePaginationVarsSchema,
} from "@/components/PageCursorTable/types";
import { InvocationDataNotFoundAlert } from "@/components/pages/InvocationDataNotFoundAlert";
import { gql } from "@/graphql/__generated__";
import type { ActionWhereInput } from "@/graphql/__generated__/graphql";
import { ActionWhereInputSchema } from "@/graphql/__generated__/zod";
import { NotFoundError } from "@/main";
import { generatePageTitle } from "@/utils/generatePageTitle";
import { parseGraphqlEdgeListWithFragment } from "@/utils/parseGraphqlEdgeList";

export const GET_BAZEL_INVOCATION_ACTIONS = gql(/* GraphQL */ `
  query GetBazelInvocationActions(
    $invocationID: UUID!
    $after: Cursor
    $first: Int
    $before: Cursor
    $last: Int
    $where: ActionWhereInput
  ) {
    getBazelInvocation(invocationID: $invocationID) {
      id
      actionTimingMetrics {
        totalExpectedTimeInMs
        timeSavedByCacheHitsInMs
        totalActions
        timedActions
        cacheHitActions
        timedCacheHitActions
      }
      configurations {
        mnemonic
      }
      metrics {
        actionSummary {
          actionData {
            mnemonic
          }
        }
        timingMetrics {
          executionPhaseTimeInMs
        }
      }
      actions(
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
            ...BazelInvocationAction
          }
        }
      }
    }
  }
`);

export const BAZEL_INVOCATION_ACTION_FRAGMENT = gql(/* GraphQL */ `
  fragment BazelInvocationAction on Action {
    id
    label
    type
    runner
    cacheHit
    success
    exitCode
    commandLine
    startTime
    endTime
    failureCode
    failureMessage
    primaryOutput
    primaryOutputURI
    actionDigest {
      rev2InstanceName
      digestFunction
      hash
      sizeBytes
    }
    stdoutURI
    stderrURI
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

const ActionSearchSchema = z.object({
  actionTable: TablePaginationVarsSchema.extend({
    where: z.array(ActionWhereInputSchema().partial()).optional(),
  }).optional(),
});

export const Route = createFileRoute(
  "/bazel-invocations/$invocationID/actions",
)({
  component: RouteComponent,
  validateSearch: (search) => ActionSearchSchema.parse(search),
  loaderDeps: ({ search: { actionTable } }) => ({ actionTable }),
  loader: async ({ params, deps }) => {
    const pageSize = deps.actionTable?.pageSize ?? DEFAULT_PAGE_SIZE;
    const pagination = deps.actionTable?.pagination ?? { first: pageSize };
    const where: ActionWhereInput[] = deps.actionTable?.where ?? [];

    const { data, error } = await apolloClient.query({
      errorPolicy: "all",
      query: GET_BAZEL_INVOCATION_ACTIONS,
      variables: {
        invocationID: params.invocationID,
        where: where.length > 0 ? { and: where } : undefined,
        ...pagination,
      },
      fetchPolicy: "cache-first",
    });

    if (!data?.getBazelInvocation) {
      throw new NotFoundError("actions", error?.message);
    }

    const actions = parseGraphqlEdgeListWithFragment(
      BAZEL_INVOCATION_ACTION_FRAGMENT,
      data.getBazelInvocation.actions,
    );

    return {
      actions,
      actionTimingMetrics: data.getBazelInvocation.actionTimingMetrics,
      executionPhaseTimeInMs:
        data.getBazelInvocation.metrics?.timingMetrics?.executionPhaseTimeInMs,
      actionMnemonics: Array.from(
        new Set(
          [
            ...(data.getBazelInvocation.metrics?.actionSummary?.actionData?.map(
              ({ mnemonic }) => mnemonic,
            ) ?? []),
            ...actions.map(({ type }) => type),
          ].filter((mnemonic): mnemonic is string => Boolean(mnemonic)),
        ),
      ).sort((left, right) => left.localeCompare(right)),
      configurationMnemonics: Array.from(
        new Set(
          [
            ...(data.getBazelInvocation.configurations ?? []),
            ...actions.map(({ configuration }) => configuration),
          ]
            .map((configuration) => configuration?.mnemonic)
            .filter((mnemonic): mnemonic is string => Boolean(mnemonic)),
        ),
      ).sort((left, right) => left.localeCompare(right)),
      pageSize,
      pageInfo: data.getBazelInvocation.actions.pageInfo,
    };
  },
  head: (_ctx) => ({
    meta: [
      {
        title: generatePageTitle([
          "Invocation",
          "Actions",
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
      actionTable: {
        ...prev.actionTable,
        ...newPagination,
      },
    }),
  });

function RouteComponent() {
  const {
    actions,
    actionTimingMetrics,
    actionMnemonics,
    configurationMnemonics,
    executionPhaseTimeInMs,
    pageSize,
    pageInfo,
  } = Route.useLoaderData();
  const { actionTable } = Route.useSearch();
  const navigate = Route.useNavigate();

  if (actions.length === 0 && !actionTable?.where?.length) {
    return <InvocationDataNotFoundAlert type="actions" />;
  }

  const onFilterChange = (where: ActionWhereInput[]) => {
    navigate({
      from: Route.id,
      to: ".",
      search: (prev): typeof prev => ({
        ...prev,
        actionTable: {
          ...prev.actionTable,
          where,
          pagination: undefined,
        },
      }),
    });
  };

  return (
    <ActionsTab
      actions={actions}
      actionTimingMetrics={actionTimingMetrics}
      actionMnemonics={actionMnemonics}
      configurationMnemonics={configurationMnemonics}
      executionPhaseTimeInMs={executionPhaseTimeInMs}
      getPaginationUpdateLink={getPaginationUpdateLink}
      onFilterChange={onFilterChange}
      pageInfo={pageInfo}
      pageSize={pageSize}
    />
  );
}
