import { SearchOutlined } from "@ant-design/icons";
import { type TableColumnsType, Typography } from "antd";
import dayjs from "dayjs";
import Link from "next/link";
import { validate as uuidValidate } from "uuid";
import styles from "@/components/BazelInvocationColumns/Columns.module.css";
import {
  SearchFilterIcon,
  SearchWidget,
  TimeRangeSelector,
} from "@/components/SearchWidgets";
import type { BuildNodeFragment } from "@/graphql/__generated__/graphql";

export const columns: TableColumnsType<BuildNodeFragment> = [
  {
    key: "repo",
    width: 200,
    title: "Repository",
    filterDropdown: (filterProps) => (
      <SearchWidget placeholder="Repository..." {...filterProps} />
    ),
    filterIcon: (filtered) => (
      <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />
    ),
  },
  {
    key: "buildUUID",
    width: 220,
    title: "Build ID",
    render: (_, record) => (
      <Link href={`/builds/${record.buildUUID}`}>{record.buildUUID}</Link>
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
  },
  {
    key: "buildURL",
    width: 220,
    title: "Build URL",
    render: (_, record) => (
      <Link href={record.buildURL}>{record.buildURL}</Link>
    ),
    filterDropdown: (filterProps) => (
      <SearchWidget
        placeholder="Provide a build URL prefix..."
        {...filterProps}
      />
    ),
    filterIcon: (filtered) => (
      <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />
    ),
  },
  {
    key: "buildDate",
    width: 220,
    title: "Timestamp",
    render: (_, record) => (
      <Typography.Text code ellipsis className={styles.startedAt}>
        {dayjs(record.timestamp).format("YYYY-MM-DD hh:mm:ss A")}
      </Typography.Text>
    ),
    filterDropdown: (filterProps) => <TimeRangeSelector {...filterProps} />,
    filterIcon: (filtered) => (
      <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />
    ),
  },
];
