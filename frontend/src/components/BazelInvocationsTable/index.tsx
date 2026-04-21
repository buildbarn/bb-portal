import { BuildOutlined } from "@ant-design/icons";
import { useQuery } from "@apollo/client/react";
import { Space, Typography } from "antd";
import React from "react";
import {
  buildColumn,
  durationColumn,
  invocationIdColumn,
  startedAtColumn,
  statusColumn,
  userColumn,
} from "@/components/BazelInvocationColumns/Columns";
import {
  BazelInvocationOrderField,
  type BazelInvocationWhereInput,
  OrderDirection,
} from "@/graphql/__generated__/graphql";
import themeStyles from "@/theme/theme.module.css";
import styles from "@/theme/theme.module.css";
import { applyTableFilters } from "@/utils/applyColumnFilters";
import { parseGraphqlEdgeListWithFragment } from "@/utils/parseGraphqlEdgeList";
import { shouldPollInvocation } from "@/utils/shouldPollInvocation";
import { CursorTable, getNewPaginationVariables } from "../CursorTable";
import type { PaginationVariables } from "../CursorTable/types";
import PortalAlert from "../PortalAlert";
import FIND_BAZEL_INVOCATIONS_QUERY, {
  BAZEL_INVOCATION_NODE_FRAGMENT,
} from "./query.graphql";

const BazelInvocationsTable: React.FC = () => {
  const [paginationVariables, setPaginationVariables] =
    React.useState<PaginationVariables>(getNewPaginationVariables());
  const [filterVariables, setFilterVariables] = React.useState<
    BazelInvocationWhereInput[]
  >([]);

  const { loading, data, error } = useQuery(FIND_BAZEL_INVOCATIONS_QUERY, {
    variables: {
      where: {
        and: [...filterVariables, { startedAtNotNil: true }],
      },
      orderBy: {
        direction: OrderDirection.Desc,
        field: BazelInvocationOrderField.StartedAt,
      },
      ...paginationVariables,
    },
    fetchPolicy: "network-only",
  });

  const invocations = parseGraphqlEdgeListWithFragment(
    BAZEL_INVOCATION_NODE_FRAGMENT,
    data?.findBazelInvocations,
  );
  const inProgressInvocations = invocations
    .filter((inv) => shouldPollInvocation(inv))
    .map((inv) => inv.id);

  // Refetch any ongoing invocations periodically. The result of the query is
  // unused, but in the background Apollo updates the result of the original
  // query based on the IDs of the response.
  useQuery(FIND_BAZEL_INVOCATIONS_QUERY, {
    variables: {
      where: {
        idIn: inProgressInvocations,
      },
    },
    skip: inProgressInvocations.length === 0,
    pollInterval: 5000,
  });

  const tableColumns = [
    userColumn,
    invocationIdColumn,
    startedAtColumn,
    durationColumn,
    statusColumn,
    buildColumn,
  ];

  if (error) {
    return (
      <PortalAlert
        type="error"
        message={
          error?.message ||
          "An unknown error occurred while fetching invocations."
        }
        showIcon
        className={styles.alert}
      />
    );
  }

  const emptyText = "No Bazel invocations match the specified search criteria";

  return (
    <CursorTable
      columns={tableColumns}
      dataSource={invocations}
      rowKey={(item) => item.id}
      loading={loading}
      size="small"
      locale={{
        emptyText: (
          <Typography.Text
            disabled
            className={themeStyles.tableEmptyTextTypography}
          >
            <Space>
              <BuildOutlined />
              {emptyText}
            </Space>
          </Typography.Text>
        ),
      }}
      onChange={(_pagination, filters, _sorter, _extra) =>
        applyTableFilters(tableColumns, filters, setFilterVariables)
      }
      pagination={{
        position: "bottom",
        justify: "end",
        size: "middle",
      }}
      pageInfo={{
        startCursor: data?.findBazelInvocations.pageInfo.startCursor || "",
        endCursor: data?.findBazelInvocations.pageInfo.endCursor || "",
        hasNextPage: data?.findBazelInvocations.pageInfo.hasNextPage || false,
        hasPreviousPage:
          data?.findBazelInvocations.pageInfo.hasPreviousPage || false,
      }}
      paginationVariables={paginationVariables}
      setPaginationVariables={setPaginationVariables}
    />
  );
};

export default BazelInvocationsTable;
