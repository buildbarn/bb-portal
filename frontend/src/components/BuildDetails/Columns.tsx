import { FilterOutlined, SearchOutlined } from "@ant-design/icons";
import { Space, type TableColumnsType, Typography } from "antd";
import Link from "next/link";
import type { FindBuildFromUuidFragment } from "@/app/builds/[buildUUID]/[[...slugs]]/types";
import styles from "@/components/AppBar/index.module.css";
import type { FindBuildByUuidQuery } from "@/graphql/__generated__/graphql";
import { InvocationResultTag } from "../InvocationResultTag";
import { invocationResultTagFilters } from "../InvocationResultTag/filters";
import PortalDuration from "../PortalDuration";
import SearchWidget, { SearchFilterIcon } from "../SearchWidgets";

export const getColumns = (
  data: FindBuildByUuidQuery,
): TableColumnsType<FindBuildFromUuidFragment> => {
  const workflow_filters: string[] = Array.from(
    new Set(
      data?.getBuild?.invocations?.map(
        (x) => x.sourceControl?.workflow ?? "",
      ) ?? [],
    ),
  );

  const job_filters: string[] = Array.from(
    new Set(
      data?.getBuild?.invocations?.map((x) => x.sourceControl?.job ?? "") ?? [],
    ),
  );

  const action_filters: string[] = Array.from(
    new Set(
      data?.getBuild?.invocations?.map((x) => x.sourceControl?.action ?? "") ??
        [],
    ),
  );

  return [
    {
      title: "Workflow",
      dataIndex: ["sourceControl","workflow"],
      filterSearch: true,
      filterIcon: (filtered) => (
        <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />
      ),
      onFilter: (value, record) =>
        record.sourceControl?.workflow?.includes(value.toString())
          ? true
          : false,
      sorter: (a, b) =>
        (a.sourceControl?.workflow ?? "").localeCompare(
          b.sourceControl?.workflow ?? "",
        ),
      filters: workflow_filters.map((x) => ({ text: x, value: x })),
    },
    {
      title: "Job",
      dataIndex: ["sourceControl","job"],
      filterSearch: true,
      filterIcon: (filtered) => (
        <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />
      ),
      onFilter: (value, record) =>
        record.sourceControl?.job?.includes(value.toString()) ? true : false,
      sorter: (a, b) =>
        (a.sourceControl?.job ?? "").localeCompare(b.sourceControl?.job ?? ""),
      filters: job_filters.map((x) => ({ text: x, value: x })),
    },
    {
      title: "Action",
      dataIndex: ["sourceControl", "action"],
      filterSearch: true,
      filterIcon: (filtered) => (
        <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />
      ),
      onFilter: (value, record) =>
        record.sourceControl?.action?.includes(value.toString()) ? true : false,
      sorter: (a, b) =>
        (a.sourceControl?.action ?? "").localeCompare(
          b.sourceControl?.action ?? "",
        ),
      filters: action_filters.map((x) => ({ text: x, value: x })),
    },
    {
      title: "Invocation ID",
      dataIndex: "invocationID",
      filterSearch: true,
      filterDropdown: (filterProps) => (
        <SearchWidget placeholder="Target Pattern..." {...filterProps} />
      ),
      filterIcon: (filtered) => (
        <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />
      ),
      onFilter: (value, record) =>
        record.invocationID.includes(value.toString()) ? true : false,
      render: (_, record) => (
        <Space>
          <span className={styles.copyIcon}>
            <Typography.Text
              copyable={{ text: record.invocationID ?? "Copy" }}
            />
          </span>
          <Link href={"/bazel-invocations/" + record.invocationID}>
            {record.invocationID}
          </Link>
        </Space>
      ),
    },
    {
      title: "Duration",
      dataIndex: "startedAt",
      render: (_, record) => (
        <PortalDuration
          key="duration"
          from={record.startedAt}
          to={record.endedAt}
          includeIcon
          includePopover
        />
      ),
    },
    {
      title: "Status",
      dataIndex: "status",
      filterSearch: true,
      render: (_, record) => (
        <InvocationResultTag
          key="result"
          exitCodeName={record.exitCodeName || undefined}
          bepCompleted={record.bepCompleted}
        />
      ),
      filterIcon: (filtered) => (
        <SearchFilterIcon icon={<FilterOutlined />} filtered={filtered} />
      ),
      onFilter: (value, record) => record.exitCodeName == value,
      filters: invocationResultTagFilters,
    },
  ];
};
