import { createFileRoute, linkOptions } from "@tanstack/react-router";
import { z } from "zod";
import { apolloClient } from "@/components/ApolloWrapper";
import { InvocationTargetsTab } from "@/components/InvocationTargets/InvocationTargetsTab";
import { DEFAULT_PAGE_SIZE } from "@/components/PageCursorTable";
import {
  type TablePaginationVars,
  TablePaginationVarsSchema,
} from "@/components/PageCursorTable/types";
import { getFragmentData, gql } from "@/graphql/__generated__";
import type { InvocationTargetWhereInput } from "@/graphql/__generated__/graphql";
import { InvocationTargetWhereInputSchema } from "@/graphql/__generated__/zod";
import { NotFoundError } from "@/main";
import { env } from "@/utils/env";
import { requireFeature } from "@/utils/featureGuard";
import { generatePageTitle } from "@/utils/generatePageTitle";
import { parseGraphqlEdgeListWithFragment } from "@/utils/parseGraphqlEdgeList";

const GET_BAZEL_INVOCATION_TARGETS = gql(/* GraphQL */ `
  query GetBazelInvocationTargets(
    $invocationID: UUID!
    $after: Cursor
    $first: Int
    $before: Cursor
    $last: Int
    $where: InvocationTargetWhereInput
  ) {
    getBazelInvocation(invocationID: $invocationID) {
      id
      metrics {
        id
        targetMetrics {
          ...BazelInvocationTargetMetrics
        }
      }
      invocationTargets(after: $after, first: $first, before: $before, last: $last, where: $where){
        pageInfo {
          startCursor
          endCursor
          hasNextPage
          hasPreviousPage
        }
        edges {
          node {
            ...BazelInvocationTargets
          }
        }
      }
      ...BazelInvocationTargetCounts
    }
  }
`);

const BAZEL_INVOCATION_TARGET_METRICS_FRAGMENT = gql(/* GraphQL */ `
  fragment BazelInvocationTargetMetrics on TargetMetrics {
    id
    targetsLoaded
    targetsConfigured
    targetsConfiguredNotIncludingAspects
  }
`);
const BAZEL_INVOCATION_TARGETS_FRAGMENT = gql(/* GraphQL */ `
  fragment BazelInvocationTargets on InvocationTarget {
    id
    success
    abortReason
    failureMessage
    tags
    target {
      id
      label
      aspect
      targetKind
      instanceName {
        name
      }
    }
  }
`);
const BAZEL_INVOCATION_TARGET_COUNTS_FRAGMENT = gql(/* GraphQL */ `
  fragment BazelInvocationTargetCounts on BazelInvocation {
    numTotal: invocationTargets {
      totalCount
    }
    numSuccessful: invocationTargets(where: { success: true }) {
      totalCount
    }
    numSkipped: invocationTargets(where: {abortReason: SKIPPED}) {
      totalCount
    }
  }
`);

const InvocationTargetsSearchSchema = z.object({
  invocationTable: TablePaginationVarsSchema.extend({
    where: z.array(InvocationTargetWhereInputSchema().partial()).optional(),
  }).optional(),
});

export const Route = createFileRoute(
  "/bazel-invocations/$invocationID/targets",
)({
  component: RouteComponent,
  validateSearch: (search) => InvocationTargetsSearchSchema.parse(search),
  beforeLoad: requireFeature(env.featureFlags?.bes?.pageTargets),
  loaderDeps: ({ search: { invocationTable } }) => ({ invocationTable }),
  loader: async ({ params, deps }) => {
    // We set the defaults here instead of in the validate search function, as
    // that updates the URL and we don't want to do that on initial load.
    const pageSize = deps.invocationTable?.pageSize ?? DEFAULT_PAGE_SIZE;
    const pagination = deps.invocationTable?.pagination ?? {
      first: pageSize,
    };
    const where: InvocationTargetWhereInput[] =
      deps.invocationTable?.where ?? [];

    const { data, error } = await apolloClient.query({
      errorPolicy: "all",
      query: GET_BAZEL_INVOCATION_TARGETS,
      variables: {
        invocationID: params.invocationID,
        where: { and: [...where] },
        ...pagination,
      },
      fetchPolicy: "network-only",
    });

    if (!data?.getBazelInvocation) {
      throw new NotFoundError("invocation", error?.message);
    }

    const targetMetrics = getFragmentData(
      BAZEL_INVOCATION_TARGET_METRICS_FRAGMENT,
      data.getBazelInvocation.metrics?.targetMetrics,
    );
    const targetCounts = getFragmentData(
      BAZEL_INVOCATION_TARGET_COUNTS_FRAGMENT,
      data.getBazelInvocation,
    );
    const invocationTargets = parseGraphqlEdgeListWithFragment(
      BAZEL_INVOCATION_TARGETS_FRAGMENT,
      data.getBazelInvocation.invocationTargets,
    );
    return {
      targetMetrics,
      targetCounts,
      invocationTargets,
      pageSize,
      pageInfo: data.getBazelInvocation.invocationTargets.pageInfo,
    };
  },
  head: (_ctx) => ({
    meta: [
      {
        title: generatePageTitle([
          "Invocation",
          "Targets",
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
      invocationTable: {
        ...prev.invocationTable,
        ...newPagination,
      },
    }),
  });

function RouteComponent() {
  const navigate = Route.useNavigate();
  const { targetMetrics, targetCounts, invocationTargets, pageSize, pageInfo } =
    Route.useLoaderData();

  const onFilterChange = (where: InvocationTargetWhereInput[]) => {
    navigate({
      from: Route.id,
      to: ".",
      search: (prev): typeof prev => ({
        ...prev,
        invocationTable: {
          ...prev.invocationTable,
          where,
          pagination: undefined,
        },
      }),
    });
  };
  return (
    <InvocationTargetsTab
      targetMetrics={targetMetrics}
      targetCounts={targetCounts}
      invocationTargets={invocationTargets}
      onFilterChange={onFilterChange}
      getPaginationUpdateLink={getPaginationUpdateLink}
      pageSize={pageSize}
      pageInfo={pageInfo}
    />
  );
}
