import { useQuery } from "@apollo/client/react";
import { createFileRoute, linkOptions } from "@tanstack/react-router";
import z from "zod";
import { apolloClient } from "@/components/ApolloWrapper";
import { DEFAULT_PAGE_SIZE } from "@/components/PageCursorTable";
import {
  type TablePaginationVars,
  TablePaginationVarsSchema,
} from "@/components/PageCursorTable/types";
import { BazelInvocationsPage } from "@/components/pages/BazelInvocations";
import { gql } from "@/graphql/__generated__/gql";
import {
  BazelInvocationOrderField,
  type BazelInvocationWhereInput,
  OrderDirection,
} from "@/graphql/__generated__/graphql";
import { BazelInvocationWhereInputSchema } from "@/graphql/__generated__/zod";
import { NotFoundError } from "@/main";
import { generatePageTitle } from "@/utils/generatePageTitle";
import { parseGraphqlEdgeListWithFragment } from "@/utils/parseGraphqlEdgeList";
import { shouldPollInvocation } from "@/utils/shouldPollInvocation";

const FIND_BAZEL_INVOCATIONS_QUERY = gql(/* GraphQL */ `
  query FindBazelInvocations(
    $after: Cursor
    $first: Int
    $before: Cursor
    $last: Int
    $orderBy: BazelInvocationOrder
    $where: BazelInvocationWhereInput
  ) {
    findBazelInvocations(after: $after, first: $first, before: $before, last: $last, orderBy: $orderBy, where: $where) {
      pageInfo {
        startCursor
        endCursor
        hasNextPage
        hasPreviousPage
      }
      edges {
        node {
          ...BazelInvocationNode
        }
      }
    }
  }
`);

const BAZEL_INVOCATION_NODE_FRAGMENT = gql(/* GraphQL */ `
  fragment BazelInvocationNode on BazelInvocation {
    id
    invocationID
    startedAt
    username
    authenticatedUser {
      userUUID
      displayName
    }
    endedAt
    exitCodeName
    connectionMetadata {
      connectionLastOpenAt
      timeSinceLastConnectionMillis
    }
    build {
      buildUUID
    }
    instanceName {
      id
      name
    }
  }
`);

const UserDetailsSearchSchema = z.object({
  invocationTable: TablePaginationVarsSchema.extend({
    where: z.array(BazelInvocationWhereInputSchema().partial()).optional(),
  }).optional(),
});

export const Route = createFileRoute("/bazel-invocations/")({
  component: RouteComponent,
  validateSearch: (search) => UserDetailsSearchSchema.parse(search),
  loaderDeps: ({ search: { invocationTable } }) => ({ invocationTable }),
  loader: async ({ deps }) => {
    // We set the defaults here instead of in the validate search function, as
    // that updates the URL and we don't want to do that on initial load.
    const pageSize = deps.invocationTable?.pageSize ?? DEFAULT_PAGE_SIZE;
    const pagination = deps.invocationTable?.pagination ?? { first: pageSize };
    const where = deps.invocationTable?.where ?? [];

    const variables = {
      where: {
        and: [...where, { startedAtNotNil: true }],
      },
      orderBy: {
        direction: OrderDirection.Desc,
        field: BazelInvocationOrderField.StartedAt,
      },
      ...pagination,
    };
    const { data, error } = await apolloClient.query({
      errorPolicy: "all",
      query: FIND_BAZEL_INVOCATIONS_QUERY,
      variables,
      fetchPolicy: "network-only",
    });
    if (!data) {
      throw new NotFoundError("invocations", error?.message);
    }

    return { data, variables, pageSize };
  },
  head: (_ctx) => ({ meta: [{ title: generatePageTitle(["Invocations"]) }] }),
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
  const { data: initialData, variables, pageSize } = Route.useLoaderData();
  // Subscribe to the loaded page so normalized updates from polling are rendered.
  const { data = initialData } = useQuery(FIND_BAZEL_INVOCATIONS_QUERY, {
    variables,
    fetchPolicy: "cache-only",
    errorPolicy: "all",
  });
  const invocations = parseGraphqlEdgeListWithFragment(
    BAZEL_INVOCATION_NODE_FRAGMENT,
    data.findBazelInvocations,
  );
  const pageInfo = data.findBazelInvocations.pageInfo;

  const inProgressInvocations = invocations
    .filter((inv) => shouldPollInvocation(inv))
    .map((inv) => inv.id);

  // Poll only ongoing invocations; the page subscription receives their updates.
  useQuery(FIND_BAZEL_INVOCATIONS_QUERY, {
    variables: {
      where: {
        idIn: inProgressInvocations,
      },
    },
    skip: inProgressInvocations.length === 0,
    pollInterval: 5000,
  });

  const onFilterChange = (where: BazelInvocationWhereInput[]) => {
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
    <BazelInvocationsPage
      invocations={invocations}
      pageSize={pageSize}
      onFilterChange={onFilterChange}
      getPaginationUpdateLink={getPaginationUpdateLink}
      pageInfo={pageInfo}
    />
  );
}
