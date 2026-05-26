import type { SorterResult } from "antd/es/table/interface";
import type React from "react";
import type { TestSummaryWhereInput } from "@/graphql/__generated__/graphql";
import { PageCursorTable } from "../PageCursorTable";
import type {
  GetPaginationUpdateLinkType,
  PageInfo,
} from "../PageCursorTable/types";
import { tableFiltersToGraphqlWhere } from "../PageCursorTable/utils";
import { columns, type TestTabRowType } from "./columns";

interface Props {
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
  testSummaryData,
  pageSize,
  onFilterChange,
  onSortChange,
  getPaginationUpdateLink,
  pageInfo,
}) => {
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
