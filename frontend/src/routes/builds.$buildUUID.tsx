import { useQuery } from "@apollo/client/react";
import { createFileRoute, linkOptions } from "@tanstack/react-router";
import z from "zod";
import { apolloClient } from "@/components/ApolloWrapper";
import { DEFAULT_PAGE_SIZE } from "@/components/PageCursorTable";
import {
  type TablePaginationVars,
  TablePaginationVarsSchema,
} from "@/components/PageCursorTable/types";
import { BuildDetailsPage } from "@/components/pages/BuildDetails";
import { getFragmentData, gql } from "@/graphql/__generated__";
import type {
  BazelInvocationWhereInput,
  BuildWhereInput,
} from "@/graphql/__generated__/graphql";
import { BazelInvocationWhereInputSchema } from "@/graphql/__generated__/zod";
import { NotFoundError } from "@/main";
import { generatePageTitle } from "@/utils/generatePageTitle";
import { parseGraphqlEdgeListWithFragment } from "@/utils/parseGraphqlEdgeList";
import { shouldPollInvocation } from "@/utils/shouldPollInvocation";

export const GET_BUILD_BY_UUID_QUERY = gql(/* GraphQL */ `
  query FindBuildByUUID(
    $buildUUID: UUID!
    $after: Cursor
    $first: Int
    $before: Cursor
    $last: Int
    $orderBy: BazelInvocationOrder
    $where: BazelInvocationWhereInput
  ) {
    getBuild(buildUUID: $buildUUID) {
      ...GetBuild
      invocations(after: $after, first: $first, before: $before, last: $last, orderBy: $orderBy, where: $where) {
        pageInfo {
          startCursor
          endCursor
          hasNextPage
          hasPreviousPage
        }
        edges {
          node {
            ...GetBuildInvocation
          }
        }
      }
    }
  }
`);

export const GET_BUILD_FRAGMENT = gql(/* GraphQL */ `
  fragment GetBuild on Build {
    id
    buildUUID
    timestamp
    tags(orderBy: { field: KEY, direction: ASC }) {
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

export const GET_BUILD_INVOCATION_FRAGMENT = gql(/* GraphQL */ `
  fragment GetBuildInvocation on BazelInvocation {
    id
    invocationID
    username
    endedAt
    startedAt
    exitCodeName
    tags {
      edges {
        node {
          id
          key
          value
        }
      }
    }
    connectionMetadata {
      connectionLastOpenAt
      timeSinceLastConnectionMillis
    }
    originalCommandLine
  }
`);

const BuildSearchSchema = z.object({
  buildTable: TablePaginationVarsSchema.extend({
    where: z.array(BazelInvocationWhereInputSchema().partial()).optional(),
  }).optional(),
});

export const Route = createFileRoute("/builds/$buildUUID")({
  component: RouteComponent,
  validateSearch: (search) => BuildSearchSchema.parse(search),
  loaderDeps: ({ search: { buildTable } }) => ({ buildTable }),
  loader: async ({ params, deps }) => {
    // We set the defaults here instead of in the validate search function, as
    // that updates the URL and we don't want to do that on initial load.
    const pageSize = deps.buildTable?.pageSize ?? DEFAULT_PAGE_SIZE;
    const pagination = deps.buildTable?.pagination ?? {
      first: pageSize,
    };
    const where: BuildWhereInput[] = deps.buildTable?.where ?? [];

    const { data, error } = await apolloClient.query({
      errorPolicy: "all",
      query: GET_BUILD_BY_UUID_QUERY,
      variables: {
        ...pagination,
        where: { and: where },
        buildUUID: params.buildUUID,
      },
    });

    if (!data?.getBuild) {
      throw new NotFoundError("build", error?.message);
    }
    const invocations = parseGraphqlEdgeListWithFragment(
      GET_BUILD_INVOCATION_FRAGMENT,
      data.getBuild.invocations,
    );
    const build = getFragmentData(GET_BUILD_FRAGMENT, data.getBuild);

    return {
      build,
      invocations,
      pageInfo: data.getBuild.invocations.pageInfo,
      pageSize,
    };
  },
  head: (_ctx) => ({
    meta: [{ title: generatePageTitle(["Build", _ctx.params.buildUUID]) }],
  }),
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
  const { buildUUID } = Route.useParams();
  const { build, invocations, pageInfo, pageSize } = Route.useLoaderData();
  const inProgressInvocations = invocations
    .filter((inv) => shouldPollInvocation(inv))
    .map((inv) => inv.id);
  const onFilterChange = (where: BazelInvocationWhereInput[]) => {
    navigate({
      from: Route.id,
      to: ".",
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
  // Refetch any ongoing invocations periodically. The result of the query is
  // unused, but in the background Apollo updates the result of the original
  // query based on the IDs of the response.
  useQuery(GET_BUILD_BY_UUID_QUERY, {
    variables: {
      where: {
        idIn: inProgressInvocations,
      },
      buildUUID: buildUUID,
    },
    skip: inProgressInvocations.length === 0,
    pollInterval: 5000,
  });
  return (
    <BuildDetailsPage
      build={build}
      invocations={invocations}
      pageInfo={pageInfo}
      pageSize={pageSize}
      getPaginationUpdateLink={getPaginationUpdateLink}
      onFilterChange={onFilterChange}
    />
  );
}
