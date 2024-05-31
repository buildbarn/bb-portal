import { ColumnType } from 'antd/lib/table';
import React from 'react';
import { SearchOutlined } from '@ant-design/icons';
import Link from 'next/link';
import { BuildNodeFragment } from '@/graphql/__generated__/graphql';
import { SearchFilterIcon, SearchWidget } from '@/components/SearchWidgets';

const buildUuidColumn: ColumnType<BuildNodeFragment> = {
  key: 'buildUUID',
  width: 220,
  title: 'Build',
  render: (_, record) => <Link href={`/builds/${record.buildUUID}`}>{record.buildUUID}</Link>,
  filterDropdown: filterProps => (
    <SearchWidget placeholder="Provide a build UUID..." {...filterProps} />
  ),
  filterIcon: filtered => <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />,
};

const buildUrlColumn: ColumnType<BuildNodeFragment> = {
  key: 'buildURL',
  width: 220,
  title: 'URL',
  render: (_, record) => <Link href={record.buildURL}>{record.buildURL}</Link>,
  filterDropdown: filterProps => (
    <SearchWidget placeholder="Provide a build URL prefix..." {...filterProps} />
  ),
  filterIcon: filtered => <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />,
};

const getColumns = (): ColumnType<BuildNodeFragment>[] => {
  return [
    buildUuidColumn,
    buildUrlColumn,
  ];
};

export default getColumns;
