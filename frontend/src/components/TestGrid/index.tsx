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
import { GetTestsWithOffsetQueryVariables } from '@/graphql/__generated__/graphql';
import TestGridRow from '../TestGridRow';
import PortalAlert from '../PortalAlert';
import Link from 'next/link';
import styles from "../../theme/theme.module.css"
import { millisecondsToTime } from '../Utilities/time';
interface Props {
  //labelData: GetTestsWithOffsetQuery | undefined
}

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
  average_duration: number;
  min_duration: number;
  max_duration: number;
  total_count: number;
  status: TestStatusType[];
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
  {
    title: "Average Duration",
    dataIndex: "average_duration",
    render: (_, record) => <span className={styles.numberFormat}>{millisecondsToTime(record.average_duration)}</span>
  },
  {
    title: "Min Duration",
    dataIndex: "min_duration",
    render: (_, record) => <span className={styles.numberFormat}>{millisecondsToTime(record.min_duration)}</span>
  },
  {
    title: "Max Duration",
    dataIndex: "max_duration",
    render: (_, record) => <span className={styles.numberFormat}>{millisecondsToTime(record.max_duration)}</span>
  },
  {
    title: "# Runs",
    dataIndex: "total_count",
    align: "right",
    render: (_, record) => <span className={styles.numberFormat}>{record.total_count}</span>,
  },
]

const TestGrid: React.FC<Props> = () => {

  const [variables, setVariables] = useState<GetTestsWithOffsetQueryVariables>({ limit:  PAGE_SIZE})

  const { loading: labelLoading, data: labelData, previousData: labelPreviousData, error: labelError } = useQuery(GET_TEST_GRID_DATA, {
    variables: variables,
    pollInterval: 300000
  });

  const data = labelLoading ? labelPreviousData : labelData;
  var result: TestGridRowDataType[] = []

  if (labelError) {
    <PortalAlert className="error" message="There was a problem communicating w/the backend server." />
  } else {
    data?.getTestsWithOffset?.result?.map(dataRow => {
      var row: TestGridRowDataType = {
        key: "test-grid-row-data-" + uniqueId(),
        label: dataRow?.label ?? "",
        status: [],
        average_duration: dataRow?.avg ?? 0,
        min_duration: dataRow?.min ?? 0,
        max_duration: dataRow?.max ?? 0,
        total_count: dataRow?.count ?? 0,
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