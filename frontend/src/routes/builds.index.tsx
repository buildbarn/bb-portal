import { createFileRoute, linkOptions } from "@tanstack/react-router";
import z from "zod";
import { apolloClient } from "@/components/ApolloWrapper";
import { DEFAULT_PAGE_SIZE } from "@/components/PageCursorTable";
import {
  type TablePaginationVars,
  TablePaginationVarsSchema,
} from "@/components/PageCursorTable/types";
import { BuildsPage } from "@/components/pages/Builds";
import { gql } from "@/graphql/__generated__";
import {
  BuildOrderField,
  type BuildWhereInput,
  OrderDirection,
} from "@/graphql/__generated__/graphql";
import { BuildWhereInputSchema } from "@/graphql/__generated__/zod";
import { NotFoundError } from "@/main";
import { generatePageTitle } from "@/utils/generatePageTitle";
import { parseGraphqlEdgeListWithFragment } from "@/utils/parseGraphqlEdgeList";

export const FIND_BUILDS_QUERY = gql(/* GraphQL */ `
  query FindBuilds(
    $after: Cursor
    $first: Int
    $before: Cursor
    $last: Int
    $orderBy: BuildOrder
    $where: BuildWhereInput
  ) {
    findBuilds(after: $after, first: $first, before: $before, last: $last, orderBy: $orderBy, where: $where) {
      pageInfo {
        startCursor
        endCursor
        hasNextPage
        hasPreviousPage
      }
      edges {
        node {
          ...BuildNode
        }
      }
    }
  }
`);

export const BUILD_NODE_FRAGMENT = gql(/* GraphQL */ `
  fragment BuildNode on Build {
    id
    buildUUID
    timestamp
    tags {
      edges {
        node {
          id
          key
          value
        }
      }
    }
  }
`);

const BuildSearchSchema = z.object({
  buildTable: TablePaginationVarsSchema.extend({
    where: z.array(BuildWhereInputSchema().partial()).optional(),
  }).optional(),
});

export const Route = createFileRoute("/builds/")({
  component: RouteComponent,
  validateSearch: (search) => BuildSearchSchema.parse(search),
  loaderDeps: ({ search: { buildTable } }) => ({ buildTable }),
  loader: async ({ deps }) => {
    // We set the defaults here instead of in the validate search function, as
    // that updates the URL and we don't want to do that on initial load.
    const pageSize = deps.buildTable?.pageSize ?? DEFAULT_PAGE_SIZE;
    const pagination = deps.buildTable?.pagination ?? {
      first: pageSize,
    };
    const where: BuildWhereInput[] = deps.buildTable?.where ?? [];

    const { data, error } = await apolloClient.query({
      errorPolicy: "all",
      query: FIND_BUILDS_QUERY,
      variables: {
        ...pagination,
        where: { and: where },
        orderBy: {
          direction: OrderDirection.Desc,
          field: BuildOrderField.Timestamp,
        },
      },
    });

    if (!data?.findBuilds) {
      throw new NotFoundError("builds", error?.message);
    }
    const builds = parseGraphqlEdgeListWithFragment(
      BUILD_NODE_FRAGMENT,
      data?.findBuilds,
    );
    return {
      builds,
      pageSize,
      pageInfo: data.findBuilds.pageInfo,
    };
  },
  head: (_ctx) => ({ meta: [{ title: generatePageTitle(["Builds"]) }] }),
});

const getPaginationUpdateLink = (newPagination: TablePaginationVars) =>
  linkOptions({
    from: Route.id,
    to: ".",
    search: (prev): typeof prev => ({
      ...prev,
      buildTable: {
        ...prev.buildTable,
        ...newPagination,
      },
    }),
  });

function RouteComponent() {
  const navigate = Route.useNavigate();
  const { builds, pageSize, pageInfo } = Route.useLoaderData();
  const onFilterChange = (where: BuildWhereInput[]) => {
    navigate({
      search: (prev): typeof prev => ({
        ...prev,
        buildTable: {
          ...prev.buildTable,
          where,
          pagination: undefined,
        },
      }),
    });
  };
  return (
    <BuildsPage
      getPaginationUpdateLink={getPaginationUpdateLink}
      onFilterChange={onFilterChange}
      pageSize={pageSize}
      pageInfo={pageInfo}
      builds={builds}
    />
  );
}
