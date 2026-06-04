import { createFileRoute, linkOptions } from "@tanstack/react-router";
import z from "zod";
import { apolloClient } from "@/components/ApolloWrapper";
import { DEFAULT_PAGE_SIZE } from "@/components/PageCursorTable";
import {
  type TablePaginationVars,
  TablePaginationVarsSchema,
} from "@/components/PageCursorTable/types";
import { TargetDetailsPage } from "@/components/pages/TargetDetails";
import { getFragmentData, gql } from "@/graphql/__generated__";
import type { InvocationTargetWhereInput } from "@/graphql/__generated__/graphql";
import { InvocationTargetWhereInputSchema } from "@/graphql/__generated__/zod";
import { NotFoundError } from "@/main";
import { env } from "@/utils/env";
import { requireFeature } from "@/utils/featureGuard";
import { generatePageTitle } from "@/utils/generatePageTitle";
import { parseGraphqlEdgeListWithFragment } from "@/utils/parseGraphqlEdgeList";

const GET_TARGET_DETAILS = gql(/* GraphQl */ `
  query GetTargetDetails(
    $id: ID!
    $after: Cursor
    $first: Int
    $before: Cursor
    $last: Int
    $where: InvocationTargetWhereInput
  ) {
    getTarget(id: $id ) {
      ...TargetDetails
      invocationTargets(after: $after, first: $first, before: $before, last: $last, where: $where) {
        pageInfo {
          startCursor
          endCursor
          hasNextPage
          hasPreviousPage
        }
        totalCount
        edges {
          node {
            ...InvocationTargetDetails
          }
        }
      }
    }
  }
`);

const TARGET_DETAILS_FRAGMENT = gql(/* GraphQL */ `
  fragment TargetDetails on Target {
    id
    aspect
    instanceName {
      name
    }
    label
    targetKind
  }
`);

const INVOCATION_TARGET_DETAILS_FRAGMENT = gql(/* GraphQL */ `
  fragment InvocationTargetDetails on InvocationTarget {
    id
    success
    abortReason
    failureMessage
    tags
    bazelInvocation {
      invocationID
    }
  }
`);

const TargetsSearchSchema = z.object({
  targetsTable: TablePaginationVarsSchema.extend({
    where: z.array(InvocationTargetWhereInputSchema().partial()).optional(),
  }).optional(),
});

export const Route = createFileRoute("/targets/$targetID/")({
  component: RouteComponent,
  validateSearch: (search) => TargetsSearchSchema.parse(search),
  beforeLoad: requireFeature(env.featureFlags?.bes?.pageTargets),
  loaderDeps: ({ search: { targetsTable } }) => ({ targetsTable }),
  loader: async ({ params, deps }) => {
    // We set the defaults here instead of in the validate search function, as
    // that updates the URL and we don't want to do that on initial load.
    const pageSize = deps.targetsTable?.pageSize ?? DEFAULT_PAGE_SIZE;
    const pagination = deps.targetsTable?.pagination ?? {
      first: pageSize,
    };
    const where: InvocationTargetWhereInput[] = deps.targetsTable?.where ?? [];

    const { data, error } = await apolloClient.query({
      errorPolicy: "all",
      query: GET_TARGET_DETAILS,
      variables: {
        id: params.targetID,
        where: { and: [...where] },
        ...pagination,
      },
      fetchPolicy: "cache-first",
    });
    if (!data?.getTarget) {
      throw new NotFoundError("target", error?.message);
    }
    const targetData = getFragmentData(TARGET_DETAILS_FRAGMENT, data.getTarget);
    const invocationTargets = parseGraphqlEdgeListWithFragment(
      INVOCATION_TARGET_DETAILS_FRAGMENT,
      data.getTarget.invocationTargets,
    );
    const pageInfo = data.getTarget.invocationTargets.pageInfo;

    return {
      targetData,
      pageSize,
      pageInfo,
      invocationTargets,
    };
  },
  head: (_ctx) => {
    const label = _ctx.loaderData?.targetData.label;
    if (label === undefined) {
      return { meta: [{ title: generatePageTitle(["Target", "Not Found"]) }] };
    }
    return { meta: [{ title: generatePageTitle(["Target", label]) }] };
  },
});

const getPaginationUpdateLink = (newPagination: TablePaginationVars) =>
  linkOptions({
    from: Route.id,
    to: ".",
    search: (prev): typeof prev => ({
      ...prev,
      targetsTable: {
        ...prev.targetsTable,
        ...newPagination,
      },
    }),
  });

function RouteComponent() {
  const navigate = Route.useNavigate();
  const { targetData, pageSize, pageInfo, invocationTargets } =
    Route.useLoaderData();
  const onFilterChange = (where: InvocationTargetWhereInput[]) => {
    navigate({
      search: (prev): typeof prev => ({
        ...prev,
        targetsTable: {
          ...prev.targetsTable,
          where,
          pagination: undefined,
        },
      }),
    });
  };
  return (
    <TargetDetailsPage
      targetData={targetData}
      pageSize={pageSize}
      pageInfo={pageInfo}
      getPaginationUpdateLink={getPaginationUpdateLink}
      onFilterChange={onFilterChange}
      invocationTargets={invocationTargets}
    />
  );
}
