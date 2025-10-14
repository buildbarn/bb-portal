import React, { useCallback, useState } from 'react';
import { Space, Table, TableProps, Typography } from 'antd';
import { TablePaginationConfig } from 'antd/lib/table';
import { BuildOutlined } from '@ant-design/icons';
import { useQuery } from '@apollo/client';
import { FilterValue } from 'antd/lib/table/interface';
import getColumns from './Columns';
import {
  BazelInvocationNodeFragment,
  BazelInvocationOrderField,
  BazelInvocationWhereInput,
  FindBazelInvocationsQueryVariables,
  OrderDirection,
} from '@/graphql/__generated__/graphql';
import { getFragmentData } from '@/graphql/__generated__';
import FIND_BAZEL_INVOCATIONS_QUERY, {
  BAZEL_INVOCATION_NODE_FRAGMENT
} from '@/components/BazelInvocationsTable/query.graphql';
import themeStyles from '@/theme/theme.module.css';

const PAGE_SIZE = 100;

type Props = {
  height?: number;
};

const BazelInvocationsTable: React.FC<Props> = ({ height }) => {
  const [variables, setVariables] = useState<FindBazelInvocationsQueryVariables>({
    first: PAGE_SIZE,
    where: {startedAtNotNil: true},
    orderBy: {
      direction: OrderDirection.Desc,
      field: BazelInvocationOrderField.StartedAt,
    },
  });

  const { loading, data, previousData, error } = useQuery(FIND_BAZEL_INVOCATIONS_QUERY, {
    variables,
    pollInterval: 120000,
    fetchPolicy: "network-only",
  });

  const onChange: TableProps<BazelInvocationNodeFragment>['onChange'] = useCallback(
    (pagination: TablePaginationConfig, filters: Record<string, FilterValue | null>, extra: any) => {
      const wheres: BazelInvocationWhereInput[] = [];
      if (filters['invocationID']?.length) {
        wheres.push({ invocationID: filters['invocationID'][0].toString() });
      }
      if (filters['startedAt']?.length === 2) {
        if (filters['startedAt'][0]) {
          wheres.push({ startedAtGTE: filters['startedAt'][0] });
        }
        if (filters['startedAt'][1]) {
          wheres.push({ startedAtLTE: filters['startedAt'][1] });
        }
      }
      if (filters['build']?.length) {
        wheres.push({ hasBuildWith: [{ buildUUID: filters['build'][0].toString() }] });
      }
      if (filters["user"]?.length){
        const userFilterValue = filters["user"][0].toString()
        wheres.push({
            or: [
              {
                hasAuthenticatedUserWith: [
                  { displayNameContains: userFilterValue },
                ],
              },
              {
                and: [
                  { userLdapContains: userFilterValue },
                  { hasAuthenticatedUser: false },
                ],
              },
            ],
          });
      }
      wheres.push({startedAtNotNil: true})

      setVariables({
        first: PAGE_SIZE,
        where: wheres.length ? { and: [...wheres] } : undefined,
        orderBy: {
          direction: OrderDirection.Desc,
          field: BazelInvocationOrderField.StartedAt,
        },
      });
    },
    [],
  );

  const activeData = loading ? previousData : data;

  let emptyText = 'No Bazel invocations match the specified search criteria';
  let dataSource: BazelInvocationNodeFragment[] = [];
  if (error) {
    emptyText = error.message;
    dataSource = [];
  } else {
    const bazelInvocationNodes = activeData?.findBazelInvocations.edges?.flatMap(edge => edge?.node) ?? [];
    const bazelInvocationNodeFragments = bazelInvocationNodes.map(node => getFragmentData(BAZEL_INVOCATION_NODE_FRAGMENT, node));
    dataSource = bazelInvocationNodeFragments.filter((nodeFragment): nodeFragment is BazelInvocationNodeFragment => !!nodeFragment);
  }

  return (
    <Table
      columns={getColumns()}
      virtual
      scroll={{ y: height ? height : 320, scrollToFirstRowOnChange: true }}
      dataSource={dataSource}
      pagination={false}
      rowKey={item => item.id}
      onChange={onChange}
      locale={{
        emptyText: loading ? (
          <Typography.Text disabled className={themeStyles.tableEmptyTextTypography}>
            <Space>
              <BuildOutlined />
              <>Loading...</>
            </Space>
          </Typography.Text>
        ) : (
          <Typography.Text disabled className={themeStyles.tableEmptyTextTypography}>
            <Space>
              <BuildOutlined />
              <>{emptyText}</>
            </Space>
          </Typography.Text>
        ),
      }}
    />
  );
};

export default BazelInvocationsTable;
