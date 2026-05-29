import { createFileRoute } from "@tanstack/react-router";
import { z } from "zod";
import { apolloClient } from "@/components/ApolloWrapper";
import { DEFAULT_PAGE_SIZE } from "@/components/PageCursorTable";
import {
  type OnTablePaginationChange,
  type TablePaginationVars,
  TablePaginationVarsSchema,
} from "@/components/PageCursorTable/types";
import { UserDetailsPage } from "@/components/pages/UserDetails";
import { getFragmentData, gql } from "@/graphql/__generated__";
import {
  BazelInvocationOrderField,
  type BazelInvocationWhereInput,
  OrderDirection,
} from "@/graphql/__generated__/graphql";
import { BazelInvocationWhereInputSchema } from "@/graphql/__generated__/zod";
import { UserNotFoundError } from "@/main";
import { generatePageTitle } from "@/utils/generatePageTitle";

export const GET_AUTHENTICATED_USER_BY_UUID = gql(/* GraphQL */ `
  query GetAuthenticatedUser(
    $userUUID: UUID!
    $after: Cursor
    $first: Int
    $before: Cursor
    $last: Int
    $bazelInvocationsOrderBy: BazelInvocationOrder
    $bazelInvocationsWhere: BazelInvocationWhereInput
  ) {
    getAuthenticatedUser(userUUID: $userUUID) {
      ...AuthenticatedUserNodeFragment
    }
  }
`);

export const AUTHENTICATED_USER_NODE_FRAGMENT = gql(/* GraphQL */ `
  fragment AuthenticatedUserNodeFragment on AuthenticatedUser {
    id
    displayName
    userInfo
    userUUID
    bazelInvocations(after: $after, first: $first, before: $before, last: $last, orderBy: $bazelInvocationsOrderBy, where: $bazelInvocationsWhere) {
      pageInfo {
        startCursor
        endCursor
        hasNextPage
        hasPreviousPage
      }
      edges {
        node {
          id
          invocationID
          build {
            id
            buildUUID
          }
          endedAt
          startedAt
          exitCodeName
          connectionMetadata {
            id
            connectionLastOpenAt
            timeSinceLastConnectionMillis
          }
        }
      }
    }
  }
`);

const UserDetailsSearchSchema = z.object({
  invocationTable: TablePaginationVarsSchema.extend({
    where: z.array(BazelInvocationWhereInputSchema().partial()).optional(),
  }).optional(),
});

export const Route = createFileRoute("/users/$userUUID")({
  component: RouteComponent,
  validateSearch: (search) => UserDetailsSearchSchema.parse(search),
  loaderDeps: ({ search: { invocationTable } }) => ({ invocationTable }),
  loader: async ({ params, deps }) => {
    // We set the defaults here instead of in the validate search function, as
    // that updates the URL and we don't want to do that on initial load.
    const pageSize = deps.invocationTable?.pageSize ?? DEFAULT_PAGE_SIZE;
    const pagination = deps.invocationTable?.pagination ?? { first: pageSize };
    const where = deps.invocationTable?.where ?? [];

    const { data } = await apolloClient.query({
      query: GET_AUTHENTICATED_USER_BY_UUID,
      variables: {
        userUUID: params.userUUID,
        bazelInvocationsWhere: { and: [...where, { startedAtNotNil: true }] },
        bazelInvocationsOrderBy: {
          field: BazelInvocationOrderField.StartedAt,
          direction: OrderDirection.Desc,
        },
        ...pagination,
      },
      fetchPolicy: "network-only",
    });

    if (!data?.getAuthenticatedUser) {
      throw new UserNotFoundError();
    }

    const user = getFragmentData(
      AUTHENTICATED_USER_NODE_FRAGMENT,
      data?.getAuthenticatedUser,
    );

    return { user, pageSize };
  },
  head: (_ctx) => ({
    meta: [
      {
        title: generatePageTitle([
          "User",
          _ctx.loaderData?.user?.displayName || _ctx.params.userUUID,
        ]),
      },
    ],
  }),
});

function RouteComponent() {
  const navigate = Route.useNavigate();
  const { user, pageSize } = Route.useLoaderData();

  const onPaginationChange: OnTablePaginationChange = (
    vars: TablePaginationVars,
  ) => {
    navigate({
      search: (prev) => ({
        ...prev,
        invocationTable: {
          ...prev.invocationTable,
          pageSize: vars.pageSize,
          pagination: vars.pagination,
        },
      }),
    });
  };

  const onFilterChange = (where: BazelInvocationWhereInput[]) => {
    navigate({
      search: (prev) => ({
        ...prev,
        invocationTable: {
          ...prev.invocationTable,
          where,
        },
      }),
    });
  };

  return (
    <UserDetailsPage
      user={user}
      pageSize={pageSize}
      onPaginationChange={onPaginationChange}
      onFilterChange={onFilterChange}
    />
  );
}
