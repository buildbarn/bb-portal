import { createFileRoute, linkOptions } from "@tanstack/react-router";
import { z } from "zod";
import { apolloClient } from "@/components/ApolloWrapper";
import { DEFAULT_PAGE_SIZE } from "@/components/PageCursorTable";
import {
  type TablePaginationVars,
  TablePaginationVarsSchema,
} from "@/components/PageCursorTable/types";
import { TargetsPage } from "@/components/pages/Targets";
import { gql } from "@/graphql/__generated__";
import type { TargetWhereInput } from "@/graphql/__generated__/graphql";
import { TargetWhereInputSchema } from "@/graphql/__generated__/zod";
import { NotFoundError } from "@/main";
import { env } from "@/utils/env";
import { requireFeature } from "@/utils/featureGuard";
import { generatePageTitle } from "@/utils/generatePageTitle";
import { parseGraphqlEdgeListWithFragment } from "@/utils/parseGraphqlEdgeList";

const GET_TARGETS_LIST = gql(/* GraphQl */ `
  query GetTargetsList(
    $after: Cursor
    $first: Int
    $before: Cursor
    $last: Int
    $where: TargetWhereInput
  ){
    findTargets (after: $after, first: $first, before: $before, last: $last, where: $where){
      pageInfo {
        startCursor
        endCursor
        hasNextPage
        hasPreviousPage
      }
      edges {
        node {
          ...TargetListDetails
        }
      }
    }
  }
`);

const TARGET_LIST_DETAILS_FRAGMENT = gql(/* GraphQL */ `
  fragment TargetListDetails on Target {
    id
    label
    aspect
    targetKind
    instanceName {
      id
      name
    }
  } 
`);

const TargetsSearchSchema = z.object({
  targetsTable: TablePaginationVarsSchema.extend({
    where: z.array(TargetWhereInputSchema().partial()).optional(),
  }).optional(),
});

export const Route = createFileRoute("/targets/")({
  component: RouteComponent,
  validateSearch: (search) => TargetsSearchSchema.parse(search),
  beforeLoad: requireFeature(env.featureFlags?.bes?.pageTargets),
  loaderDeps: ({ search: { targetsTable } }) => ({ targetsTable }),
  loader: async ({ deps }) => {
    // We set the defaults here instead of in the validate search function, as
    // that updates the URL and we don't want to do that on initial load.
    const pageSize = deps.targetsTable?.pageSize ?? DEFAULT_PAGE_SIZE;
    const pagination = deps.targetsTable?.pagination ?? {
      first: pageSize,
    };
    const where: TargetWhereInput[] = deps.targetsTable?.where ?? [];
    const { data, error } = await apolloClient.query({
      errorPolicy: "all",
      query: GET_TARGETS_LIST,
      variables: {
        ...pagination,
        where: { and: [...where] },
      },
      fetchPolicy: "cache-first",
    });
    if (!data?.findTargets) {
      throw new NotFoundError("targets", error?.message);
    }
    const targets = parseGraphqlEdgeListWithFragment(
      TARGET_LIST_DETAILS_FRAGMENT,
      data.findTargets,
    );
    return { targets, pageSize, pageInfo: data.findTargets.pageInfo };
  },
  head: (_ctx) => ({ meta: [{ title: generatePageTitle(["Targets"]) }] }),
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
  const { targets, pageSize, pageInfo } = Route.useLoaderData();

  const onFilterChange = (where: TargetWhereInput[]) => {
    navigate({
      from: Route.id,
      to: ".",
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
    <TargetsPage
      targets={targets}
      onFilterChange={onFilterChange}
      getPaginationUpdateLink={getPaginationUpdateLink}
      pageSize={pageSize}
      pageInfo={pageInfo}
    />
  );
}
