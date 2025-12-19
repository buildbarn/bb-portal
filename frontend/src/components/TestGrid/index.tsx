import { useQuery } from "@apollo/client";
import type { FilterValue } from "antd/es/table/interface";
import React from "react";
import type {
  GetTestsQuery,
  TargetWhereInput,
} from "@/graphql/__generated__/graphql";
import { parseGraphqlEdgeList } from "@/utils/parseGraphqlEdgeList";
import styles from "../../theme/theme.module.css";
import { CursorTable, getNewPaginationVariables } from "../CursorTable";
import type { PaginationVariables } from "../CursorTable/types";
import PortalAlert from "../PortalAlert";
import TestGridRow from "../TestGridRow";
import type { TestStatusEnum } from "../TestStatusTag";
import { columns, type TestGridRowDataType } from "./columns";
import { GET_TESTS } from "./graphql";

export interface TestStatusType {
  label: string;
  invocationId: string;
  status: TestStatusEnum;
}

const TestGrid: React.FC = () => {
  const [paginationVariables, setPaginationVariables] =
    React.useState<PaginationVariables>(getNewPaginationVariables());
  const [filterVariables, setFilterVariables] = React.useState<
    Array<TargetWhereInput>
  >([]);
  const { loading, data, error } = useQuery<GetTestsQuery>(GET_TESTS, {
    variables: {
      where: {
        and: [
          { hasInvocationTargetsWith: { hasTestSummary: true } },
          ...filterVariables,
        ],
      },
      ...paginationVariables,
    },
    pollInterval: 300000,
  });

  const onFilterChange = (filters: Record<string, FilterValue | null>) => {
    const newFilters: TargetWhereInput[] = [];
    Object.entries(filters).forEach(([key, value]) => {
      if (value && value.length > 0) {
        switch (key) {
          case "instanceName":
            newFilters.push({
              hasInstanceNameWith: [
                {
                  nameContainsFold: value[0] as string,
                },
              ],
            });
            break;
          case "target":
            newFilters.push({
              labelContainsFold: value[0] as string,
            });
            break;
        }
      }
    });
    setFilterVariables(newFilters);
  };

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

  const parsedData: TestGridRowDataType[] = parseGraphqlEdgeList(
    data?.findTargets,
  );

  return (
    <CursorTable<TestGridRowDataType>
      columns={columns}
      loading={loading}
      dataSource={parsedData}
      size="small"
      rowKey="id"
      expandable={{
        expandedRowRender: (record) => (
          <TestGridRow
            instanceName={record.instanceName.name}
            label={record.label}
            aspect={record.aspect}
            targetKind={record.targetKind}
            numberOfElements={40}
            direction={"newToOld"}
          />
        ),
        rowExpandable: (_) => true,
      }}
      pagination={{
        position: "bottom",
        justify: "end",
        size: "middle",
      }}
      onChange={(_pagination, filters, _sorter, _extra) => {
        onFilterChange(filters);
      }}
      pageInfo={{
        startCursor: data?.findTargets.pageInfo.startCursor,
        endCursor: data?.findTargets.pageInfo.endCursor,
        hasNextPage: data?.findTargets.pageInfo.hasNextPage,
        hasPreviousPage: data?.findTargets.pageInfo.hasPreviousPage,
      }}
      paginationVariables={paginationVariables}
      setPaginationVariables={setPaginationVariables}
    />
  );
};

export default TestGrid;
