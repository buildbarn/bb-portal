import { ClockCircleFilled, SearchOutlined } from "@ant-design/icons";
import { Link } from "@tanstack/react-router";
import { Space, Tag, Tooltip } from "antd";
import type { FilterValue } from "antd/es/table/interface";
import dayjs from "dayjs";
import { validate as uuidValidate } from "uuid";
import CodeText from "@/components/CodeText";
import { InvocationResultTag } from "@/components/InvocationResultTag";
import {
  applyInvocationResultTagFilter,
  invocationResultTagFilters,
} from "@/components/InvocationResultTag/filters";
import PortalDuration from "@/components/PortalDuration";
import {
  FilterPickerDropdown,
  SearchFilterIcon,
  SearchWidget,
  TimeRangeSelector,
} from "@/components/SearchWidgets";
import UserStatusIndicator from "@/components/UserStatusIndicator";
import type {
  BazelInvocationNodeFragment,
  BazelInvocationWhereInput,
} from "@/graphql/__generated__/graphql";
import type { TableColumnTypeWithFilter } from "@/types/TableColumnTypeWithFilter";

export const invocationIdColumn: TableColumnTypeWithFilter<
  BazelInvocationNodeFragment,
  BazelInvocationWhereInput
> = {
  key: "invocationID",
  width: 220,
  title: "Invocation",
  render: (_, record) => (
    <Link
      to={`/bazel-invocations/$invocationID`}
      params={{ invocationID: record.invocationID }}
    >
      {record.invocationID}
    </Link>
  ),
  filterDropdown: (filterProps) => (
    <SearchWidget
      placeholder="Provide a Bazel invocation ID..."
      dataValidator={uuidValidate}
      validationTooltip="The search string needs to be a valid UUID"
      {...filterProps}
    />
  ),
  filterIcon: (filtered) => (
    <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />
  ),
  applyFilter: (value: FilterValue) => {
    if (value.length === 0) {
      return undefined;
    }
    return [{ invocationID: value[0] as string }];
  },
};

export const startedAtColumn: TableColumnTypeWithFilter<
  BazelInvocationNodeFragment,
  BazelInvocationWhereInput
> = {
  key: "startedAt",
  width: 165,
  title: "Start Time",
  render: (_, record) => (
    <CodeText>
      {dayjs(record.startedAt).format("YYYY-MM-DD hh:mm:ss A")}
    </CodeText>
  ),
  filterDropdown: (filterProps) => <TimeRangeSelector {...filterProps} />,
  filterIcon: (filtered) => (
    <SearchFilterIcon icon={<ClockCircleFilled />} filtered={filtered} />
  ),
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

export const durationColumn: TableColumnTypeWithFilter<
  BazelInvocationNodeFragment,
  BazelInvocationWhereInput
> = {
  key: "duration",
  width: 100,
  title: "Duration",
  render: (_, record) => (
    <PortalDuration
      from={record.startedAt || undefined}
      to={
        record.endedAt
          ? record.endedAt
          : record.connectionMetadata?.connectionLastOpenAt || undefined
      }
      formatConfig={{ smallestUnit: "s" }}
    />
  ),
};

export const statusColumn: TableColumnTypeWithFilter<
  BazelInvocationNodeFragment,
  BazelInvocationWhereInput
> = {
  key: "result",
  width: 120,
  title: "Result",
  render: (_, record) => (
    <InvocationResultTag
      exitCodeName={record.exitCodeName || undefined}
      timeSinceLastConnectionMillis={
        record.connectionMetadata?.timeSinceLastConnectionMillis || undefined
      }
    />
  ),
  filterIcon: (filtered) => (
    <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />
  ),
  filterDropdown: (filterProps) => <FilterPickerDropdown {...filterProps} />,
  filters: invocationResultTagFilters,
  applyFilter: applyInvocationResultTagFilter,
};

export const buildColumn: TableColumnTypeWithFilter<
  BazelInvocationNodeFragment,
  BazelInvocationWhereInput
> = {
  key: "build",
  width: 220,
  title: "Build",
  render: (_, record) =>
    record.build && (
      <Link
        to={`/builds/$buildUUID`}
        params={{ buildUUID: record.build.buildUUID }}
      >
        {record.build.buildUUID}
      </Link>
    ),
  filterDropdown: (filterProps) => (
    <SearchWidget
      placeholder="Provide a build UUID..."
      {...filterProps}
      dataValidator={uuidValidate}
      validationTooltip="The search string needs to be a valid UUID"
    />
  ),
  filterIcon: (filtered) => (
    <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />
  ),
  applyFilter: (value: FilterValue) => {
    if (value.length === 0) {
      return undefined;
    }
    return [{ hasBuildWith: [{ buildUUID: value[0] as string }] }];
  },
};

const RATIO_COLORS: Record<string, string> = {
  remote: "#52C41A",
  "remote cache hit": "#1890FF",
  local: "#FA8C16",
};

const RATIO_ABBR: Record<string, string> = {
  remote: "R",
  "remote cache hit": "C",
  local: "L",
};

export const executionRatioColumn: TableColumnTypeWithFilter<
  BazelInvocationNodeFragment,
  BazelInvocationWhereInput
> = {
  key: "executionRatio",
  width: 160,
  title: "Execution Ratio",
  render: (_, record) => {
    const runnerCounts = record.metrics?.actionSummary?.runnerCount;
    if (!runnerCounts?.length) {
      return <span style={{ color: "#8C8C8C" }}>—</span>;
    }
    const SHOWN_KEYS = ["remote", "remote cache hit", "local"] as const;
    const counts = SHOWN_KEYS.map((key) => ({
      key,
      count: runnerCounts.find((r) => r.name === key)?.actionsExecuted ?? 0,
    }));
    const otherRunners = runnerCounts.filter(
      (r) => r.name && r.name !== "total" && !(SHOWN_KEYS as readonly string[]).includes(r.name),
    );
    const otherCount = otherRunners.reduce((sum, r) => sum + (r.actionsExecuted ?? 0), 0);
    const allTotal = counts.reduce((sum, { count }) => sum + count, 0) + otherCount;
    if (allTotal === 0) {
      return <span style={{ color: "#8C8C8C" }}>—</span>;
    }
    const segments = counts.flatMap(({ key, count }) => {
      if (count === 0) return [];
      const pct = ((count / allTotal) * 100).toFixed(0);
      return [{ key, pct }];
    });
    if (segments.length === 0) {
      return <span style={{ color: "#8C8C8C" }}>—</span>;
    }
    const otherPct = otherCount > 0 ? ((otherCount / allTotal) * 100).toFixed(0) : null;
    const otherTitle = otherRunners
      .map((r) => `${r.name}: ${r.actionsExecuted ?? 0}`)
      .join(", ");
    return (
      <Space size={2} wrap>
        {segments.map((s) => (
          <Tag
            key={s.key}
            color={RATIO_COLORS[s.key]}
            style={{ margin: 0, fontSize: 11 }}
          >
            {RATIO_ABBR[s.key]} {s.pct}%
          </Tag>
        ))}
        {otherPct && (
          <Tooltip title={otherTitle}>
            <Tag style={{ margin: 0, fontSize: 11, cursor: "default" }}>
              O {otherPct}%
            </Tag>
          </Tooltip>
        )}
      </Space>
    );
  },
};

export const userColumn: TableColumnTypeWithFilter<
  BazelInvocationNodeFragment,
  BazelInvocationWhereInput
> = {
  key: "user",
  width: 120,
  title: "User",
  render: (_, record) => {
    return (
      <UserStatusIndicator
        authenticatedUser={record.authenticatedUser}
        username={record.username || undefined}
        showIcon
      />
    );
  },
  filterDropdown: (filterProps) => (
    <SearchWidget placeholder="Provide a username..." {...filterProps} />
  ),
  filterIcon: (filtered) => (
    <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />
  ),
  applyFilter: (value: FilterValue) => {
    if (value.length === 0) {
      return undefined;
    }
    const username = value[0] as string;
    return [
      {
        or: [
          { hasAuthenticatedUserWith: [{ displayNameContainsFold: username }] },
          {
            and: [
              { usernameContainsFold: username },
              { hasAuthenticatedUser: false },
            ],
          },
        ],
      },
    ];
  },
};
