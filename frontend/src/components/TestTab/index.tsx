import { useQuery } from "@apollo/client";
import type { FilterValue, SorterResult } from "antd/es/table/interface";
import React, { useMemo } from "react";
import {
  type GetTestsForInvocationQuery,
  OrderDirection,
  type TestSummaryOrder,
  TestSummaryOrderField,
  type TestSummaryWhereInput,
} from "@/graphql/__generated__/graphql";
import { parseGraphqlEdgeList } from "@/utils/parseGraphqlEdgeList";
import styles from "../../theme/theme.module.css";
import { CursorTable, getNewPaginationVariables } from "../CursorTable";
import type { PaginationVariables } from "../CursorTable/types";
import PortalAlert from "../PortalAlert";
import { columns, defaultSorting, type TestTabRowType } from "./columns";
import { GET_TESTS_FOR_INVOCATION } from "./graphql";

interface Props {
  invocationId: string;
}

export const TestTab: React.FC<Props> = ({ invocationId }) => {
  const [paginationVariables, setPaginationVariables] =
    React.useState<PaginationVariables>(getNewPaginationVariables());
  const [filterVariables, setFilterVariables] = React.useState<
    Array<TestSummaryWhereInput>
  >([]);
  const [orderBy, setOrderBy] =
    React.useState<TestSummaryOrder>(defaultSorting);

  const { data, error, loading } = useQuery<GetTestsForInvocationQuery>(
    GET_TESTS_FOR_INVOCATION,
    {
      variables: {
        where: {
          and: [
            {
              hasInvocationTargetWith: {
                hasBazelInvocationWith: { invocationID: invocationId },
              },
            },
            ...filterVariables,
          ],
        },
        orderBy,
        ...paginationVariables,
      },
      fetchPolicy: "cache-first",
      notifyOnNetworkStatusChange: true,
    },
  );

  const onFilterChange = (filters: Record<string, FilterValue | null>) => {
    const newFilters: TestSummaryWhereInput[] = [];
    Object.entries(filters).forEach(([key, value]) => {
      if (value && value.length > 0) {
        switch (key) {
          case "overallStatus":
            newFilters.push({
              overallStatusIn:
                value as TestSummaryWhereInput["overallStatusIn"],
            });
            break;
          case "label":
            newFilters.push({
              hasInvocationTargetWith: [
                {
                  hasTargetWith: [
                    {
                      labelContainsFold: value[0] as string,
                    },
                  ],
                },
              ],
            });
            break;
        }
      }
    });
    setFilterVariables(newFilters);
  };

  const onSortChange = (
    sorter: SorterResult<TestTabRowType> | SorterResult<TestTabRowType>[],
  ) => {
    const s = Array.isArray(sorter) ? sorter[0] : sorter;
    if (!s || !s.order) {
      return;
    }
    switch (s.field) {
      case "totalRunDurationInMs":
        setOrderBy({
          field: TestSummaryOrderField.TotalRunDurationInMs,
          direction:
            s.order === "ascend" ? OrderDirection.Asc : OrderDirection.Desc,
        });
        break;
    }
  };

  const parsedData: TestTabRowType[] = useMemo(() => {
    return parseGraphqlEdgeList(data?.findTestSummaries).map((ts) => {
      var cachedLocally: boolean | null = null;
      var cachedRemotely: boolean | null = null;
      if (
        ts.testResults !== null &&
        ts.testResults !== undefined &&
        ts.testResults.length > 0
      ) {
        if (ts.testResults.every((tr) => tr.cachedLocally === true)) {
          cachedLocally = true;
        } else if (ts.testResults.some((tr) => tr.cachedLocally === false)) {
          cachedLocally = false;
        }
        if (ts.testResults.every((tr) => tr.cachedRemotely === true)) {
          cachedRemotely = true;
        } else if (ts.testResults.some((tr) => tr.cachedRemotely === false)) {
          cachedRemotely = false;
        }
      }
      return {
        ...ts,
        cachedLocally: cachedLocally,
        cachedRemotely: cachedRemotely,
      };
    });
  }, [data]);

  if (error) {
    return (
      <PortalAlert
        type="error"
        message={error.message}
        showIcon
        className={styles.alert}
      />
    );
  }

  return (
    <CursorTable
      size="small"
      columns={columns}
      dataSource={parsedData}
      loading={loading}
      rowKey={"id"}
      showSorterTooltip={{ target: "sorter-icon" }}
      pagination={{
        position: "bottom",
        justify: "end",
        size: "middle",
      }}
      onChange={(_pagination, filters, sorter, _extra) => {
        onFilterChange(filters);
        onSortChange(sorter);
      }}
      pageInfo={{
        startCursor: data?.findTestSummaries.pageInfo.startCursor,
        endCursor: data?.findTestSummaries.pageInfo.endCursor,
        hasNextPage: data?.findTestSummaries.pageInfo.hasNextPage,
        hasPreviousPage: data?.findTestSummaries.pageInfo.hasPreviousPage,
      }}
      paginationVariables={paginationVariables}
      setPaginationVariables={setPaginationVariables}
    />
  );
};
