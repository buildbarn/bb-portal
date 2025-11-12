import type { ColumnType } from 'antd/lib/table';
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
import UserStatusIndicator from '../UserStatusIndicator';

export const invocationIdColumn: ColumnType<BazelInvocationNodeFragment> = {
  key: 'invocationID',
  width: 220,
  title: 'Invocation',
  render: (_, record) => <Link href={`/bazel-invocations/${record.invocationID}`}>{record.invocationID}</Link>,
  filterDropdown: filterProps => (
    <SearchWidget placeholder="Provide a Bazel invocation ID..." {...filterProps} />
  ),
  filterIcon: filtered => <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />,
};

export const startedAtColumn: ColumnType<BazelInvocationNodeFragment> = {
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

export const durationColumn: ColumnType<BazelInvocationNodeFragment> = {
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

export const statusColumn: ColumnType<BazelInvocationNodeFragment> = {
  key: 'result',
  width: 120,
  title: 'Result',
  render: (_, record) => <BuildStepResultTag result={record.state.exitCode?.name as BuildStepResultEnum} />,
  filterIcon: filtered => <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />,
  onFilter:  (value, record) => record.state.exitCode?.name == value,
  filters:[
    {
      text: "Succeeded",
      value: "SUCCESS",
    },
    {
      text: "Unstable",
      value: "UNSTABLE",
    },
    {
      text: "Parsing Failed",
      value: "PARSING_FAILURE",
    },
    {
      text: "Build Failed",
      value: "BUILD_FAILURE",
    },
    {
      text: "Tests Failed",
      value: "TESTS_FAILED",
    },
    {
      text: "Not Built",
      value: "NOT_BUILT",
    },
    {
      text: "Aborted",
      value: "ABORTED",
    },
    {
      text: "Interrupted",
      value: "INTERRUPTED",
    },
    {
      text: "Status Unknown",
      value: "UNKNOWN",
    },
  ]
};

export const buildColumn: ColumnType<BazelInvocationNodeFragment> = {
  key: 'build',
  width: 220,
  title: 'Build',
  render: (_, record) => record.build && <Link href={`/builds/${record.build.buildUUID}`}>{record.build.buildUUID}</Link>,
  filterDropdown: filterProps => (
    <SearchWidget placeholder="Provide a build UUID..." {...filterProps} />
  ),
  filterIcon: filtered => <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />,
};

export const userColumn: ColumnType<BazelInvocationNodeFragment> = {
  key: 'user',
  width: 120,
  title: "User",
  render: (_, record) => {
      return <UserStatusIndicator
          authenticatedUser={record.authenticatedUser}
          user={record.user}
        />
  },
  filterDropdown: filterProps => (
    <SearchWidget placeholder="Provide a username..." {...filterProps} />
  ),
  filterIcon: filtered => <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />,
}
