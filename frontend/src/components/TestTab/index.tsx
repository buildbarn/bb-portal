import type { SorterResult } from "antd/es/table/interface";
import type React from "react";
import { useMemo } from "react";
import type { TestSummaryWhereInput } from "@/graphql/__generated__/graphql";
import { PageCursorTable } from "../PageCursorTable";
import type {
  GetPaginationUpdateLinkType,
  PageInfo,
} from "../PageCursorTable/types";
import { tableFiltersToGraphqlWhere } from "../PageCursorTable/utils";
import { getColumns, type TestTabRowType } from "./columns";

interface Props {
  invocationId: string;
  testSummaryData: TestTabRowType[];
  pageSize: number;
  onFilterChange: (where: TestSummaryWhereInput[]) => void;
  onSortChange: (
    sorter: SorterResult<TestTabRowType> | SorterResult<TestTabRowType>[],
  ) => void;
  getPaginationUpdateLink: GetPaginationUpdateLinkType;
  pageInfo: PageInfo;
}

export const TestTab: React.FC<Props> = ({
  invocationId,
  testSummaryData,
  pageSize,
  onFilterChange,
  onSortChange,
  getPaginationUpdateLink,
  pageInfo,
}) => {
  const columns = useMemo(() => getColumns(invocationId), [invocationId]);

  return (
    <PageCursorTable
      size="small"
      columns={columns}
      dataSource={testSummaryData}
      rowKey={"id"}
      showSorterTooltip={{ target: "sorter-icon" }}
      onChange={(_pagination, filters, sorter, _extra) => {
        onSortChange(sorter);
        onFilterChange(tableFiltersToGraphqlWhere(columns, filters));
      }}
      getPaginationUpdateLink={getPaginationUpdateLink}
      pageInfo={pageInfo}
      pageSize={pageSize}
    />
  );
};
