import { ColumnType } from 'antd/lib/table';
import React from 'react';
import { Typography } from 'antd';
import {
  ClockCircleFilled,
  SearchOutlined,
} from '@ant-design/icons';
import Link from 'next/link';
import dayjs from 'dayjs';
import styles from './Columns.module.css';
import { BazelInvocationNodeFragment } from '@/graphql/__generated__/graphql';
import { SearchFilterIcon, SearchWidget, TimeRangeSelector } from '@/components/SearchWidgets';
import PortalDuration from "@/components/PortalDuration";
import BuildStepResultTag, { BuildStepResultEnum } from "@/components/BuildStepResultTag";

const invocationIdColumn: ColumnType<BazelInvocationNodeFragment> = {
  key: 'invocationID',
  width: 220,
  title: 'Invocation',
  render: (_, record) => <Link href={`/bazel-invocations/${record.invocationID}`}>{record.invocationID}</Link>,
  filterDropdown: filterProps => (
    <SearchWidget placeholder="Provide a Bazel invocation ID..." {...filterProps} />
  ),
  filterIcon: filtered => <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />,
};

const startedAtColumn: ColumnType<BazelInvocationNodeFragment> = {
  key: 'startedAt',
  width: 165,
  title: 'Start Time',
  render: (_, record) => (
    <Typography.Text code ellipsis className={styles.startedAt}>
      {dayjs(record.startedAt).format('YYYY-MM-DD hh:mm:ss A')}
    </Typography.Text>
  ),
  filterDropdown: filterProps => <TimeRangeSelector {...filterProps} />,
  filterIcon: filtered => <SearchFilterIcon icon={<ClockCircleFilled />} filtered={filtered} />,
};

const durationColumn: ColumnType<BazelInvocationNodeFragment> = {
  key: 'duration',
  width: 100,
  title: 'Duration',
  render: (_, record) => (
    <PortalDuration
      from={record.startedAt}
      to={record.endedAt}
      includePopover
    />
  ),
};

const statusColumn: ColumnType<BazelInvocationNodeFragment> = {
  key: 'result',
  width: 120,
  title: 'Result',
  render: (_, record) => <BuildStepResultTag result={record.state.exitCode?.name as BuildStepResultEnum} />,
};

const buildColumn: ColumnType<BazelInvocationNodeFragment> = {
  key: 'build',
  width: 220,
  title: 'Build',
  render: (_, record) => record.build && <Link href={`/builds/${record.build.buildUUID}`}>{record.build.buildUUID}</Link>,
  filterDropdown: filterProps => (
    <SearchWidget placeholder="Provide a build UUID..." {...filterProps} />
  ),
  filterIcon: filtered => <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />,
};

const userColumn: ColumnType<BazelInvocationNodeFragment> = {
  key: 'user',
  width: 120,
  title: "User",
  render: (_, record) => <Link href={`mailto:${record.user?.Email}`}>{record.user?.LDAP}</Link>,

  filterIcon: filtered => <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />,
}

const getColumns = (): ColumnType<BazelInvocationNodeFragment>[] => {
  return [
    userColumn,
    invocationIdColumn,
    startedAtColumn,
    durationColumn,
    statusColumn,
    buildColumn,
  ];
};

export default getColumns;
