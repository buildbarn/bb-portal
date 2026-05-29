import { CalendarFilled } from "@ant-design/icons";
import { Collapse, Descriptions, Space, Typography } from "antd";
import type React from "react";
import {
  buildColumn,
  durationColumn,
  invocationIdColumn,
  startedAtColumn,
  statusColumn,
} from "@/components/BazelInvocationColumns/Columns";
import { PageCursorTable } from "@/components/PageCursorTable";
import type { OnTablePaginationChange } from "@/components/PageCursorTable/types";
import { tableFiltersToGraphqlWhere } from "@/components/PageCursorTable/utils";
import PortalCard from "@/components/PortalCard";
import type {
  AuthenticatedUserNodeFragmentFragment,
  BazelInvocationNodeFragment,
  BazelInvocationWhereInput,
} from "@/graphql/__generated__/graphql";
import themeStyles from "@/theme/theme.module.css";
import { parseGraphqlEdgeList } from "@/utils/parseGraphqlEdgeList";

interface Props {
  pageSize: number | undefined;
  user: AuthenticatedUserNodeFragmentFragment;
  onPaginationChange: OnTablePaginationChange;
  onFilterChange: (where: BazelInvocationWhereInput[]) => void;
}

export const UserDetailsPage: React.FC<Props> = ({
  pageSize,
  user,
  onPaginationChange,
  onFilterChange,
}) => {
  const invocations = parseGraphqlEdgeList(user.bazelInvocations);
  const userInfo = user.userInfo || {};

  const tableColumns = [
    invocationIdColumn,
    startedAtColumn,
    durationColumn,
    statusColumn,
    buildColumn,
  ];

  return (
    <PortalCard
      icon={<CalendarFilled />}
      titleBits={[
        <span key="title">User {user.displayName || user.userUUID}</span>,
      ]}
    >
      <Space direction="vertical" className={themeStyles.space}>
        {Object.keys(userInfo).length > 0 && (
          <Collapse
            bordered={false}
            items={[
              {
                key: 0,
                label: (
                  <Typography.Text strong>User information</Typography.Text>
                ),
                children: (
                  <Descriptions column={1} bordered size="small">
                    {Object.keys(userInfo).map((value) => {
                      return (
                        <Descriptions.Item key={value} label={value}>
                          {userInfo[value]}
                        </Descriptions.Item>
                      );
                    })}
                  </Descriptions>
                ),
              },
            ]}
          />
        )}

        <PageCursorTable
          dataSource={invocations as BazelInvocationNodeFragment[]}
          columns={tableColumns}
          size="small"
          pageInfo={user.bazelInvocations.pageInfo}
          pageSize={pageSize}
          onPaginationChange={onPaginationChange}
          onChange={(_pagination, filters, _sorter, _extra) => {
            onFilterChange(tableFiltersToGraphqlWhere(tableColumns, filters));
          }}
        />
      </Space>
    </PortalCard>
  );
};
