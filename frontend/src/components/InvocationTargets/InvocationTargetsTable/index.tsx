import { useQuery } from "@apollo/client";
import { Space } from "antd";
import type { FilterValue } from "antd/es/table/interface";
import React from "react";
import {
  CursorTable,
  getNewPaginationVariables,
} from "@/components/CursorTable";
import type { PaginationVariables } from "@/components/CursorTable/types";
import PortalAlert from "@/components/PortalAlert";
import type {
  GetInvocationTargetsForInvocationQuery,
  InvocationTargetAbortReason,
  InvocationTargetWhereInput,
  TargetMetrics,
} from "@/graphql/__generated__/graphql";
import styles from "@/theme/theme.module.css";
import { parseGraphqlEdgeList } from "@/utils/parseGraphqlEdgeList";
import { InvocationTargetsMetrics } from "../InvocationTargetsMetrics";
import { columns, type InvocationTargetsTableRowType } from "./Columns";
import { GET_INVOCATION_TARGETS_FOR_INVOCATION } from "./graphql";

interface Props {
  invocationId: string;
  targetMetrics: TargetMetrics | undefined;
}

export const InvocationTargetsTable: React.FC<Props> = ({
  invocationId,
  targetMetrics,
}) => {
  const [paginationVariables, setPaginationVariables] =
    React.useState<PaginationVariables>(getNewPaginationVariables());
  const [filterVariables, setFilterVariables] =
    React.useState<InvocationTargetWhereInput>({});

  const { data, error, loading } =
    useQuery<GetInvocationTargetsForInvocationQuery>(
      GET_INVOCATION_TARGETS_FOR_INVOCATION,
      {
        variables: {
          invocationID: invocationId,
          where: filterVariables,
          ...paginationVariables,
        },
        fetchPolicy: "cache-first",
        notifyOnNetworkStatusChange: true,
      },
    );

  const onFilterChange = (filters: Record<string, FilterValue | null>) => {
    const newFilters: InvocationTargetWhereInput[] = [];
    Object.entries(filters).forEach(([key, value]) => {
      if (value && value.length > 0) {
        switch (key) {
          case "target-kind":
            newFilters.push({
              hasTargetWith: [{ targetKindContainsFold: value[0] as string }],
            });
            break;
          case "label":
            newFilters.push({
              hasTargetWith: [{ labelContainsFold: value[0] as string }],
            });
            break;
          case "aspect":
            newFilters.push({
              hasTargetWith: [{ aspectContainsFold: value[0] as string }],
            });
            break;
          case "success":
            if (value.length === 1) {
              newFilters.push({ success: value[0] as boolean });
            }
            break;
          case "abort-reason":
            newFilters.push({
              abortReasonIn: value as InvocationTargetAbortReason[],
            });
            break;
        }
      }
    });
    setFilterVariables({ and: newFilters });
  };

  if (error) {
    return (
      <PortalAlert
        type="error"
        message={
          error?.message ||
          "An unknown error occurred while fetching invocation targets."
        }
        showIcon
        className={styles.alert}
      />
    );
  }

  const invocationTargets = parseGraphqlEdgeList<InvocationTargetsTableRowType>(
    data?.bazelInvocation.invocationTargets,
  );

  return (
    <>
      <InvocationTargetsMetrics
        targetMetrics={targetMetrics}
        invocationTargetsCount={data?.bazelInvocation.numTotal.totalCount}
        invocationTargetsBuiltSuccessfully={
          data?.bazelInvocation.numSuccessful.totalCount
        }
        invocationTargetsSkipped={data?.bazelInvocation.numSkipped.totalCount}
      />
      <CursorTable<InvocationTargetsTableRowType>
        rowKey={"id"}
        size="small"
        columns={columns}
        loading={loading}
        dataSource={invocationTargets}
        expandable={{
          rowExpandable: (record) =>
            record.failureMessage !== null &&
            record.failureMessage !== undefined &&
            record.failureMessage !== "",
          expandedRowRender: (record) => (
            <Space direction="vertical" style={{ paddingLeft: "48px" }}>
              <strong>Failure Message:</strong>
              <pre>{record.failureMessage}</pre>
            </Space>
          ),
        }}
        showSorterTooltip={{ target: "sorter-icon" }}
        onChange={(_pagination, filters, _sorter, _extra) =>
          onFilterChange(filters)
        }
        pagination={{
          position: "bottom",
          justify: "end",
          size: "middle",
        }}
        pageInfo={{
          startCursor:
            data?.bazelInvocation.invocationTargets.pageInfo.startCursor,
          endCursor: data?.bazelInvocation.invocationTargets.pageInfo.endCursor,
          hasNextPage:
            data?.bazelInvocation.invocationTargets.pageInfo.hasNextPage,
          hasPreviousPage:
            data?.bazelInvocation.invocationTargets.pageInfo.hasPreviousPage,
        }}
        paginationVariables={paginationVariables}
        setPaginationVariables={setPaginationVariables}
      />
    </>
  );
};
