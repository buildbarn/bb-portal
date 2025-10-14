import React, { useCallback, useState } from 'react';
import { Space, Table, TableProps, Typography } from 'antd';
import { TablePaginationConfig } from 'antd/lib/table';
import { RocketOutlined } from '@ant-design/icons';
import { useQuery } from '@apollo/client';
import { FilterValue } from 'antd/lib/table/interface';
import getColumns from './Columns';
import {
  BuildNodeFragment,
  BuildOrderField,
  BuildWhereInput,
  FindBuildsQueryVariables,
  OrderDirection,
} from '@/graphql/__generated__/graphql';
import { getFragmentData } from '@/graphql/__generated__';
import FIND_BUILDS_QUERY, {
  BUILD_NODE_FRAGMENT
} from '@/components/BuildsTable/query.graphql';
import themeStyles from '@/theme/theme.module.css';

const PAGE_SIZE = 100;

type Props = {
  height?: number;
};

const BuildsTable: React.FC<Props> = ({ height }) => {
  const [variables, setVariables] = useState<FindBuildsQueryVariables>({
    first: PAGE_SIZE,
    orderBy: {
      direction: OrderDirection.Desc,
      field: BuildOrderField.Timestamp,
    },
  });

  const { loading, data, previousData, error } = useQuery(FIND_BUILDS_QUERY, {
    variables,
    pollInterval: 120000,
    fetchPolicy: 'cache-and-network',
  });

  const onChange: TableProps<BuildNodeFragment>['onChange'] = useCallback(
    (pagination: TablePaginationConfig, filters: Record<string, FilterValue | null>, extra: any) => {
      const wheres: BuildWhereInput[] = [];
      if (filters['buildUUID']?.length) {
        wheres.push({ buildUUID: filters['buildUUID'][0].toString() });
      }
      if (filters['buildURL']?.length) {
        wheres.push({ buildURLHasPrefix: filters['buildURL'][0].toString() });
      }
      if (filters['buildDate']?.length === 2) {
        if (filters['buildDate'][0])  {
          wheres.push({ timestampGTE: filters['buildDate'][0]});
        }
        if (filters['buildDate'][1])  {
          wheres.push({ timestampLTE: filters['buildDate'][1]});
        }
      }

      setVariables({
        first: PAGE_SIZE,
        where: wheres.length ? { and: [...wheres] } : wheres[0],
        orderBy: {
          direction: OrderDirection.Desc,
          field: BuildOrderField.Timestamp,
        },
      });
    },
    [],
  );

  const activeData = loading ? previousData : data;

  let emptyText = 'No builds match the specified search criteria';
  let dataSource: BuildNodeFragment[] = [];
  if (error) {
    emptyText = error.message;
    dataSource = [];
  } else {
    const buildNodes = activeData?.findBuilds.edges?.flatMap(edge => edge?.node) ?? [];
    const buildNodeFragments = buildNodes.map(node => getFragmentData(BUILD_NODE_FRAGMENT, node));
    dataSource = buildNodeFragments.filter((nodeFragment): nodeFragment is BuildNodeFragment => !!nodeFragment);
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
              <RocketOutlined rotate={160} />
              <>Loading...</>
            </Space>
          </Typography.Text>
        ) : (
          <Typography.Text disabled className={themeStyles.tableEmptyTextTypography}>
            <Space>
              <RocketOutlined rotate={160} />
              <>{emptyText}</>
            </Space>
          </Typography.Text>
        ),
      }}
    />
  );
};

export default BuildsTable;
