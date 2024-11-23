import { ColumnType } from 'antd/lib/table';
import React from 'react';
import { SearchOutlined } from '@ant-design/icons';
import Link from 'next/link';
import { BuildNodeFragment } from '@/graphql/__generated__/graphql';
import { SearchFilterIcon, SearchWidget, TimeRangeSelector } from '@/components/SearchWidgets';
import { Typography } from 'antd';
import styles from '@/components/BazelInvocationsTable/Columns.module.css'
import dayjs from 'dayjs';

const buildUuidColumn: ColumnType<BuildNodeFragment> = {
  key: 'buildUUID',
  width: 220,
  title: 'Build ID',
  render: (_, record) => <Link href={`/builds/${record.buildUUID}`}>{record.buildUUID}</Link>,
  filterDropdown: filterProps => (
    <SearchWidget placeholder="Provide a build UUID..." {...filterProps} />
  ),
  filterIcon: filtered => <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />,
};

const buildUrlColumn: ColumnType<BuildNodeFragment> = {
  key: 'buildURL',
  width: 220,
  title: 'Build URL',
  render: (_, record) => <Link href={record.buildURL}>{record.buildURL}</Link>,
  filterDropdown: filterProps => (
    <SearchWidget placeholder="Provide a build URL prefix..." {...filterProps} />
  ),
  filterIcon: filtered => <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />,
};

const buildDateColumn: ColumnType<BuildNodeFragment> = {
  key: 'buildDate',
  width: 220,
  title: 'Timestamp',
  render: (_, record) => (
    <Typography.Text code ellipsis className={styles.startedAt}>
      {dayjs(record.timestamp).format('YYYY-MM-DD hh:mm:ss A')}
    </Typography.Text>
  ),
  filterDropdown: filterProps => <TimeRangeSelector {...filterProps} />,
  filterIcon: filtered => <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />,
};

const getColumns = (): ColumnType<BuildNodeFragment>[] => {
  return [
    buildUuidColumn,
    buildUrlColumn,
    buildDateColumn,
  ];
};

export default getColumns;
