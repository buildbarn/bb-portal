import { useQuery } from "@apollo/client";
import type { FilterValue } from "antd/es/table/interface";
import React from "react";
import { validate as uuidValidate } from "uuid";
import {
  type BuildNodeFragment,
  BuildOrderField,
  type BuildWhereInput,
  OrderDirection,
} from "@/graphql/__generated__/graphql";
import { parseGraphqlEdgeListWithFragment } from "@/utils/parseGraphqlEdgeList";
import { CursorTable, getNewPaginationVariables } from "../CursorTable";
import type { PaginationVariables } from "../CursorTable/types";
import PortalAlert from "../PortalAlert";
import { columns } from "./Columns";
import { BUILD_NODE_FRAGMENT, FIND_BUILDS_QUERY } from "./query.graphql";

const BuildsTable: React.FC = () => {
  const [paginationVariables, setPaginationVariables] =
    React.useState<PaginationVariables>(getNewPaginationVariables());

  const [filterVariables, setFilterVariables] = React.useState<BuildWhereInput>(
    {},
  );

  const { data, loading, error } = useQuery(FIND_BUILDS_QUERY, {
    variables: {
      ...paginationVariables,
      where: filterVariables,
      orderBy: {
        direction: OrderDirection.Desc,
        field: BuildOrderField.Timestamp,
      },
    },
  });

  const onFilterChange = (filters: Record<string, FilterValue | null>) => {
    const newFilters: BuildWhereInput[] = [];
    Object.entries(filters).forEach(([key, value]) => {
      if (value && value.length > 0) {
        switch (key) {
          case "buildUUID": {
            const buildUUID = value[0] as string;
            if (uuidValidate(buildUUID)) {
              newFilters.push({ buildUUID: buildUUID as string });
            }
            break;
          }
          case "buildURL":
            newFilters.push({ buildURLContainsFold: value[0] as string });
            break;
          case "buildDate":
            if (value.length === 2) {
              if (value[0]) {
                newFilters.push({ timestampGTE: value[0] });
              }
              if (value[1]) {
                newFilters.push({ timestampLTE: value[1] });
              }
            }
            break;
        }
      }
    });
    setFilterVariables({ and: newFilters });
  };

  if (error) {
    return (
      <PortalAlert
        className="error"
        message="There was a problem communicating with the backend server."
      />
    );
  }

  const rowData = parseGraphqlEdgeListWithFragment(
    BUILD_NODE_FRAGMENT,
    data?.findBuilds,
  );

  return (
    <CursorTable<BuildNodeFragment>
      columns={columns}
      loading={loading}
      size="small"
      rowKey="id"
      onChange={(_pagination, filters, _sorter, _extra) =>
        onFilterChange(filters)
      }
      dataSource={rowData}
      pagination={{
        position: "bottom",
        justify: "end",
        size: "middle",
      }}
      pageInfo={{
        startCursor: data?.findBuilds.pageInfo.startCursor || "",
        endCursor: data?.findBuilds.pageInfo.endCursor || "",
        hasNextPage: data?.findBuilds.pageInfo.hasNextPage || false,
        hasPreviousPage: data?.findBuilds.pageInfo.hasPreviousPage || false,
      }}
      paginationVariables={paginationVariables}
      setPaginationVariables={setPaginationVariables}
    />
  );
};

export default BuildsTable;
