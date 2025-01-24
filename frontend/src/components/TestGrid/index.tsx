import React, { useCallback, useState } from 'react';
import { TableColumnsType } from "antd/lib"
import { Space, Row, Statistic, Table, TableProps, TablePaginationConfig, Pagination } from 'antd';
import { TestStatusEnum } from '../TestStatusTag';
import type { StatisticProps } from "antd/lib";
import CountUp from 'react-countup';
import { SearchFilterIcon, SearchWidget } from '@/components/SearchWidgets';
import { SearchOutlined } from '@ant-design/icons';
import { useQuery } from '@apollo/client';
import { GET_TEST_GRID_DATA } from '@/app/tests/index.graphql';
import { FilterValue } from 'antd/es/table/interface';
import { uniqueId } from 'lodash';
import { GetTestsWithOffsetQueryVariables, GetUniqueTestLabelsQueryVariables } from '@/graphql/__generated__/graphql';
import TestGridRow from '../TestGridRow';
import PortalAlert from '../PortalAlert';
import Link from 'next/link';
import styles from "../../theme/theme.module.css"
import { millisecondsToTime } from '../Utilities/time';
import { GET_TEST_LABELS } from './graphql';
interface Props {}

const formatter: StatisticProps['formatter'] = (value) => (
  <CountUp end={value as number} separator="," />
);
export interface TestStatusType {
  label: string
  invocationId: string,
  status: TestStatusEnum
}

interface TestGridRowDataType {
  key: React.Key;
  label: string;
}

const PAGE_SIZE = 20
const columns: TableColumnsType<TestGridRowDataType> = [
  {
    title: "Label",
    dataIndex: "label",
    filterSearch: true,
    render: (_, record) =>

      <Link href={"tests/" + btoa(encodeURIComponent(record.label))}>{record.label}</Link>,
    filterDropdown: filterProps => (
      <SearchWidget placeholder="Target Pattern..." {...filterProps} />
    ),
    filterIcon: filtered => <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />,
    onFilter: (value, record) => (record.label.includes(value.toString()) ? true : false)
  },
]

const TestGrid: React.FC<Props> = () => {

  const [variables, setVariables] = useState<GetUniqueTestLabelsQueryVariables>({})

  const { loading: labelLoading, data: labelData, previousData: labelPreviousData, error: labelError } = useQuery(GET_TEST_LABELS, {
    variables: variables,
    pollInterval: 300000
  });

  const data = labelLoading ? labelPreviousData : labelData;
  var result: TestGridRowDataType[] = []

  if (labelError) {
    <PortalAlert className="error" message="There was a problem communicating w/the backend server." />
  } else {
    data?.getUniqueTestLabels?.map(dataRow => {
      var row: TestGridRowDataType = {
        key: "test-grid-row-data-" + uniqueId(),
        label: dataRow ?? "",
      }
      result.push(row)
    })
  }
  const onChange: TableProps<TestGridRowDataType>['onChange'] = useCallback(
    (pagination: TablePaginationConfig,
      filters: Record<string, FilterValue | null>, extra: any) => {
      var vars: GetTestsWithOffsetQueryVariables = {}
      if (filters['label']?.length) {
        var label = filters['label']?.[0]?.toString() ?? ""
        vars.label = label
      } else {
        vars.label = ""
      }
      vars.offset = ((pagination.current ?? 1) - 1) * PAGE_SIZE;
      setVariables(vars)
    },
    [variables],
  );
  return (
    <Space direction="vertical" size="middle" style={{ display: 'flex' }}>
      <Row>
        <Table<TestGridRowDataType>
          columns={columns}
          loading={labelLoading}
          rowKey="key"
          onChange={onChange}
          expandable={{
            indentSize: 100,
            expandedRowRender: (record) => (
              <TestGridRow rowLabel={record.label} first={20} reverseOrder={true} />
            ),
            rowExpandable: (_) => true,
          }}
          pagination = {false}
          dataSource={result} />
      </Row>
    </Space>
  );
};

export default TestGrid;