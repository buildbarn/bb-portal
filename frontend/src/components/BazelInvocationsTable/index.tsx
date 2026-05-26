import { BuildOutlined } from "@ant-design/icons";
import { Space, Typography } from "antd";
import type React from "react";
import {
  buildColumn,
  durationColumn,
  invocationIdColumn,
  startedAtColumn,
  statusColumn,
  userColumn,
} from "@/components/BazelInvocationColumns/Columns";
import type {
  BazelInvocationNodeFragment,
  PageInfo,
} from "@/graphql/__generated__/graphql";
import themeStyles from "@/theme/theme.module.css";
import { PageCursorTable } from "../PageCursorTable";
import type {
  GetPaginationUpdateLinkType,
  OnBazelInvocationFilterChange,
} from "../PageCursorTable/types";
import { tableFiltersToGraphqlWhere } from "../PageCursorTable/utils";

interface Props {
  pageSize: number;
  invocations: BazelInvocationNodeFragment[];
  onFilterChange: OnBazelInvocationFilterChange;
  getPaginationUpdateLink: GetPaginationUpdateLinkType;
  pageInfo: PageInfo;
}

const BazelInvocationsTable: React.FC<Props> = ({
  pageSize,
  invocations,
  onFilterChange,
  getPaginationUpdateLink,
  pageInfo,
}) => {
  const tableColumns = [
    userColumn,
    invocationIdColumn,
    startedAtColumn,
    durationColumn,
    statusColumn,
    buildColumn,
  ];

  const emptyText = "No Bazel invocations match the specified search criteria";

  return (
    <PageCursorTable
      columns={tableColumns}
      dataSource={invocations}
      rowKey={(item) => item.id}
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
        onFilterChange(tableFiltersToGraphqlWhere(tableColumns, filters))
      }
      pageInfo={pageInfo}
      pageSize={pageSize}
      getPaginationUpdateLink={getPaginationUpdateLink}
    />
  );
};

export default BazelInvocationsTable;
