import { Space } from "antd";
import type React from "react";
import { PageCursorTable } from "@/components/PageCursorTable";
import type { GetPaginationUpdateLinkType } from "@/components/PageCursorTable/types";
import { tableFiltersToGraphqlWhere } from "@/components/PageCursorTable/utils";
import type {
  BazelInvocationTargetCountsFragment,
  BazelInvocationTargetMetricsFragment,
  BazelInvocationTargetsFragment,
  InvocationTargetWhereInput,
  PageInfo,
} from "@/graphql/__generated__/graphql";
import { InvocationTargetsMetrics } from "../InvocationTargetsMetrics";
import { columns } from "./Columns";

interface Props {
  targetMetrics: BazelInvocationTargetMetricsFragment | null | undefined;
  invocationTargets: BazelInvocationTargetsFragment[];
  targetCounts: BazelInvocationTargetCountsFragment;
  onFilterChange: (where: InvocationTargetWhereInput[]) => void;
  getPaginationUpdateLink: GetPaginationUpdateLinkType;
  pageSize: number;
  pageInfo: PageInfo;
}

export const InvocationTargetsTable: React.FC<Props> = ({
  targetMetrics,
  invocationTargets,
  targetCounts,
  onFilterChange,
  getPaginationUpdateLink,
  pageSize,
  pageInfo,
}) => {
  return (
    <>
      <InvocationTargetsMetrics
        targetMetrics={targetMetrics}
        targetCounts={targetCounts}
      />
      <PageCursorTable<BazelInvocationTargetsFragment>
        rowKey={"id"}
        size="small"
        columns={columns}
        dataSource={invocationTargets}
        expandable={{
          rowExpandable: (record) =>
            record.failureMessage !== null &&
            record.failureMessage !== undefined &&
            record.failureMessage !== "",
          expandedRowRender: (record) => (
            <Space direction="vertical" style={{ paddingLeft: "48px" }}>
              <strong>Failure Message:</strong>
              <pre
                style={{ whiteSpace: "pre-wrap", overflowWrap: "break-word" }}
              >
                {record.failureMessage}
              </pre>
            </Space>
          ),
        }}
        showSorterTooltip={{ target: "sorter-icon" }}
        onChange={(_pagination, filters, _sorter, _extra) =>
          onFilterChange(tableFiltersToGraphqlWhere(columns, filters))
        }
        pageInfo={pageInfo}
        getPaginationUpdateLink={getPaginationUpdateLink}
        pageSize={pageSize}
      />
    </>
  );
};
