import { useQuery } from "@apollo/client/react";
import React, { useMemo } from "react";
import {
  type BuildNodeFragment,
  BuildOrderField,
  type BuildWhereInput,
  OrderDirection,
} from "@/graphql/__generated__/graphql";
import { applyTableFilters } from "@/utils/applyColumnFilters";
import { parseGraphqlEdgeListWithFragment } from "@/utils/parseGraphqlEdgeList";
import { CursorTable, getNewPaginationVariables } from "../CursorTable";
import type { PaginationVariables } from "../CursorTable/types";
import PortalAlert from "../PortalAlert";
import { getColumns } from "./Columns";
import { BUILD_NODE_FRAGMENT, FIND_BUILDS_QUERY } from "./query.graphql";

const BuildsTable: React.FC = () => {
  const [paginationVariables, setPaginationVariables] =
    React.useState<PaginationVariables>(getNewPaginationVariables());
  const [filterVariables, setFilterVariables] = React.useState<
    BuildWhereInput[]
  >([]);

  const { data, loading, error } = useQuery(FIND_BUILDS_QUERY, {
    variables: {
      ...paginationVariables,
      where: { and: filterVariables },
      orderBy: {
        direction: OrderDirection.Desc,
        field: BuildOrderField.Timestamp,
      },
    },
  });

  const tableColumns = useMemo(getColumns, []);

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
      columns={tableColumns}
      loading={loading}
      size="small"
      rowKey="id"
      onChange={(_pagination, filters, _sorter, _extra) =>
        applyTableFilters(tableColumns, filters, setFilterVariables)
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
