import type React from "react";
import { useMemo } from "react";
import type {
  BuildNodeFragment,
  BuildWhereInput,
  PageInfo,
} from "@/graphql/__generated__/graphql";
import { PageCursorTable } from "../PageCursorTable";
import type { GetPaginationUpdateLinkType } from "../PageCursorTable/types";
import { tableFiltersToGraphqlWhere } from "../PageCursorTable/utils";
import { getColumns } from "./Columns";

interface Props {
  pageSize: number;
  builds: BuildNodeFragment[];
  onFilterChange: (where: BuildWhereInput[]) => void;
  getPaginationUpdateLink: GetPaginationUpdateLinkType;
  pageInfo: PageInfo;
}

const BuildsTable: React.FC<Props> = ({
  pageSize,
  builds,
  onFilterChange,
  getPaginationUpdateLink,
  pageInfo,
}) => {
  const tableColumns = useMemo(getColumns, []);
  return (
    <PageCursorTable<BuildNodeFragment>
      columns={tableColumns}
      size="small"
      rowKey="id"
      onChange={(_pagination, filters, _sorter, _extra) =>
        onFilterChange(tableFiltersToGraphqlWhere(tableColumns, filters))
      }
      dataSource={builds}
      getPaginationUpdateLink={getPaginationUpdateLink}
      pageSize={pageSize}
      pageInfo={pageInfo}
    />
  );
};

export default BuildsTable;
