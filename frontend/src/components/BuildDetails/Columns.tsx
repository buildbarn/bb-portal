import { FilterOutlined, SearchOutlined } from "@ant-design/icons";
import { Space, type TableColumnsType, Typography } from "antd";
import Link from "next/link";
import { validate as uuidValidate } from "uuid";
import styles from "@/components/AppBar/index.module.css";
import type { GetBuildInvocationFragment } from "@/graphql/__generated__/graphql";
import { InvocationResultTag } from "../InvocationResultTag";
import { invocationResultTagFilters } from "../InvocationResultTag/filters";
import PortalDuration from "../PortalDuration";
import SearchWidget, { SearchFilterIcon } from "../SearchWidgets";

export const columns: TableColumnsType<GetBuildInvocationFragment> = [
  {
    key: "workflow",
    title: "Workflow",
    dataIndex: ["sourceControl", "workflow"],
    filterSearch: true,
    filterDropdown: (filterProps) => (
      <SearchWidget placeholder="Workflow..." {...filterProps} />
    ),
    filterIcon: (filtered) => (
      <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />
    ),
  },
  {
    key: "job",
    title: "Job",
    dataIndex: ["sourceControl", "job"],
    filterSearch: true,
    filterDropdown: (filterProps) => (
      <SearchWidget placeholder="Job..." {...filterProps} />
    ),
    filterIcon: (filtered) => (
      <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />
    ),
  },
  {
    key: "action",
    title: "Action",
    dataIndex: ["sourceControl", "action"],
    filterSearch: true,
    filterDropdown: (filterProps) => (
      <SearchWidget placeholder="Action..." {...filterProps} />
    ),
    filterIcon: (filtered) => (
      <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />
    ),
  },
  {
    key: "invocationID",
    title: "Invocation ID",
    dataIndex: "invocationID",
    filterSearch: true,
    filterDropdown: (filterProps) => (
      <SearchWidget
        placeholder="Invocation ID..."
        {...filterProps}
        dataValidator={uuidValidate}
        validationTooltip="The search string needs to be a valid UUID"
      />
    ),
    filterIcon: (filtered) => (
      <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />
    ),
    render: (_, record) => (
      <Space>
        <span className={styles.copyIcon}>
          <Typography.Text copyable={{ text: record.invocationID ?? "Copy" }} />
        </span>
        <Link href={`/bazel-invocations/${record.invocationID}`}>
          {record.invocationID}
        </Link>
      </Space>
    ),
  },
  {
    key: "duration",
    title: "Duration",
    dataIndex: "startedAt",
    render: (_, record) => (
      <PortalDuration
        key="duration"
        from={record.startedAt || undefined}
        to={
          record.endedAt
            ? record.endedAt
            : record.connectionMetadata?.connectionLastOpenAt || undefined
        }
        includeIcon
        includePopover
        formatConfig={{ smallestUnit: "s" }}
      />
    ),
  },
  {
    key: "status",
    title: "Status",
    dataIndex: "status",
    filterSearch: true,
    render: (_, record) => (
      <InvocationResultTag
        key="result"
        exitCodeName={record.exitCodeName || undefined}
        timeSinceLastConnectionMillis={
          record.connectionMetadata?.timeSinceLastConnectionMillis || undefined
        }
      />
    ),
    filters: invocationResultTagFilters,
    filterIcon: (filtered) => (
      <SearchFilterIcon icon={<FilterOutlined />} filtered={filtered} />
    ),
  },
];
