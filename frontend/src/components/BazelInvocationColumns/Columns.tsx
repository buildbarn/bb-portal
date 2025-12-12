import type { ColumnType } from 'antd/lib/table';
import { Typography } from 'antd';
import {
  ClockCircleFilled,
  SearchOutlined,
} from '@ant-design/icons';
import Link from 'next/link';
import dayjs from 'dayjs';
import styles from './Columns.module.css';
import { BazelInvocationNodeFragment, BazelInvocationWhereInput } from '@/graphql/__generated__/graphql';
import { SearchFilterIcon, SearchWidget, TimeRangeSelector } from '@/components/SearchWidgets';
import PortalDuration from "@/components/PortalDuration";
import UserStatusIndicator from '../UserStatusIndicator';
import { InvocationResultTag } from '../InvocationResultTag';
import { FilterValue } from 'antd/es/table/interface';
import { applyInvocationResultTagFilter, invocationResultTagFilters } from '../InvocationResultTag/filters';

type ColumnTypeWithFilter<T> = ColumnType<T> & {
  applyFilter?: (value: FilterValue) => BazelInvocationWhereInput[] | undefined;
};

export const invocationIdColumn: ColumnTypeWithFilter<BazelInvocationNodeFragment> = {
  key: 'invocationID',
  width: 220,
  title: 'Invocation',
  render: (_, record) => <Link href={`/bazel-invocations/${record.invocationID}`}>{record.invocationID}</Link>,
  filterDropdown: filterProps => (
    <SearchWidget placeholder="Provide a Bazel invocation ID..." {...filterProps} />
  ),
  filterIcon: filtered => <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />,
  applyFilter: (value: FilterValue) => {
    if (value.length === 0) {
      return undefined
    }
    return [{ invocationID: value[0] as string }];
  },
};

export const startedAtColumn: ColumnTypeWithFilter<BazelInvocationNodeFragment> = {
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
  applyFilter: (value: FilterValue) => {
    if (value.length !== 2) {
      return undefined;
    }
    const filter: BazelInvocationWhereInput[] = [];
    if (value[0]) {
      filter.push({ startedAtGTE: value[0] });
    }
    if (value[1]) {
      filter.push({ startedAtLTE: value[1] });
    }
    return filter;
  },
};

export const durationColumn: ColumnTypeWithFilter<BazelInvocationNodeFragment> = {
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

export const statusColumn: ColumnTypeWithFilter<BazelInvocationNodeFragment> = {
  key: 'result',
  width: 120,
  title: 'Result',
  render: (_, record) => <InvocationResultTag exitCodeName={record.exitCodeName || undefined} bepCompleted={record.bepCompleted} />,
  filterIcon: filtered => <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />,
  filters: invocationResultTagFilters,
  applyFilter: applyInvocationResultTagFilter,
};

export const buildColumn: ColumnTypeWithFilter<BazelInvocationNodeFragment> = {
  key: 'build',
  width: 220,
  title: 'Build',
  render: (_, record) => record.build && <Link href={`/builds/${record.build.buildUUID}`}>{record.build.buildUUID}</Link>,
  filterDropdown: filterProps => (
    <SearchWidget placeholder="Provide a build UUID..." {...filterProps} />
  ),
  filterIcon: filtered => <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />,
  applyFilter: (value: FilterValue) => {
    if (value.length === 0) {
      return undefined;
    }
    return [{ hasBuildWith: [{ buildUUID: value[0] as string }] }];
  },
};

export const userColumn: ColumnTypeWithFilter<BazelInvocationNodeFragment> = {
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
  applyFilter: (value: FilterValue) => {
    if (value.length === 0) {
      return undefined;
    }
    const user = value[0] as string;
    return [
      {
        or: [
          { hasAuthenticatedUserWith: [{ displayNameContains: user }] },
          {
            and: [{ userLdapContains: user }, { hasAuthenticatedUser: false }],
          },
        ],
      },
    ];
  },
};
