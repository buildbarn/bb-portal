import { createFileRoute, linkOptions } from "@tanstack/react-router";
import { z } from "zod";
import { apolloClient } from "@/components/ApolloWrapper";
import { DEFAULT_PAGE_SIZE } from "@/components/PageCursorTable";
import {
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
import { NotFoundError } from "@/main";
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

    const { data, error } = await apolloClient.query({
      errorPolicy: "all",
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
      throw new NotFoundError("user", error?.message);
    }

    const user = getFragmentData(
      AUTHENTICATED_USER_NODE_FRAGMENT,
      data.getAuthenticatedUser,
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
  const { user, pageSize } = Route.useLoaderData();

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
    <UserDetailsPage
      user={user}
      pageSize={pageSize}
      onFilterChange={onFilterChange}
      getPaginationUpdateLink={getPaginationUpdateLink}
    />
  );
}
