import { FilterOutlined, SearchOutlined } from "@ant-design/icons";
import { Space, type TableColumnsType, Typography } from "antd";
import { validate as uuidValidate } from "uuid";
import appbarStyles from "@/components/AppBar/index.module.css";
import type { GetBuildInvocationFragment } from "@/graphql/__generated__/graphql";
import { CodeLink } from "@/components/CodeLink";
import type { CommandLineData } from "@/components/CommandLine";
import CommandLinePreview from "@/components/CommandLinePreview";
import { InvocationResultTag } from "@/components/InvocationResultTag";
import { invocationResultTagFilters } from "@/components/InvocationResultTag/filters";
import PortalDuration from "@/components/PortalDuration";
import SearchWidget, { SearchFilterIcon } from "@/components/SearchWidgets";
import buildDetailsStyles from "./index.module.css";
import { Link } from "@tanstack/react-router";

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
    key: "command",
    title: "Command",
    filterSearch: false,
    className: buildDetailsStyles.commandColumnCell,
    render: (_, record) => (
      <div className={buildDetailsStyles.commandWrapper}>
        <CommandLinePreview
          command={record.originalCommandLine as CommandLineData}
          copyable
        />
      </div>
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
        <span className={appbarStyles.copyIcon}>
          <Typography.Text copyable={{ text: record.invocationID ?? "Copy" }} />
        </span>
          <CodeLink
            text={record.invocationID}
            link={{
              to: "/bazel-invocations/$invocationID",
              params: { invocationID: record.invocationID },
            }}
          />
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
