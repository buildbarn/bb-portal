import { FilterOutlined, SearchOutlined } from "@ant-design/icons";
import { Space, Typography } from "antd";
import type { FilterValue } from "antd/es/table/interface";
import { validate as uuidValidate } from "uuid";
import appbarStyles from "@/components/AppBar/index.module.css";
import { CodeLink } from "@/components/CodeLink";
import type { CommandLineData } from "@/components/CommandLine";
import CommandLinePreview from "@/components/CommandLinePreview";
import { InvocationResultTag } from "@/components/InvocationResultTag";
import {
  applyInvocationResultTagFilter,
  invocationResultTagFilters,
} from "@/components/InvocationResultTag/filters";
import { OptionalLinkWrapper } from "@/components/OptionalLinkWrapper";
import PortalDuration from "@/components/PortalDuration";
import SearchWidget, { SearchFilterIcon } from "@/components/SearchWidgets";
import type {
  BazelInvocationWhereInput,
  GetBuildInvocationFragment,
} from "@/graphql/__generated__/graphql";
import type { TableColumnTypeWithFilter } from "@/types/TableColumnTypeWithFilter";
import { env } from "@/utils/env";
import { parseGraphqlEdgeList } from "@/utils/parseGraphqlEdgeList";
import buildDetailsStyles from "./index.module.css";

export const getColumns = (): TableColumnTypeWithFilter<
  GetBuildInvocationFragment,
  BazelInvocationWhereInput
>[] => {
  const columns: TableColumnTypeWithFilter<
    GetBuildInvocationFragment,
    BazelInvocationWhereInput
  >[] = [];

  const additionalColumns = env.additionalBuildInvocationColumns;
  for (const column of additionalColumns) {
    columns.push({
      key: column.valueKey,
      title: column.title,
      filterSearch: true,
      render: (_, record) => {
        const tags = parseGraphqlEdgeList(record.tags);
        const valueTag = tags.find((tag) => tag.key === column.valueKey);
        const urlTag = tags.find((tag) => tag.key === column.urlKey);
        return (
          <OptionalLinkWrapper url={urlTag?.value}>
            {valueTag?.value || ""}
          </OptionalLinkWrapper>
        );
      },
      filterDropdown: (filterProps) => (
        <SearchWidget placeholder={`${column.title}...`} {...filterProps} />
      ),
      filterIcon: (filtered) => (
        <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />
      ),
      applyFilter: (value: FilterValue) => {
        if (value.length === 0) {
          return undefined;
        }
        return [
          {
            hasTagsWith: [
              { key: column.valueKey, valueContainsFold: value[0] as string },
            ],
          },
        ];
      },
    });
  }

  columns.push(
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
      title: "Invocation",
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
            <Typography.Text
              copyable={{ text: record.invocationID ?? "Copy" }}
            />
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
      applyFilter: (value: FilterValue) => {
        if (value.length === 0) {
          return undefined;
        }
        return [{ invocationID: value[0] as string }];
      },
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
            record.connectionMetadata?.timeSinceLastConnectionMillis ||
            undefined
          }
        />
      ),
      filters: invocationResultTagFilters,
      applyFilter: applyInvocationResultTagFilter,
      filterIcon: (filtered) => (
        <SearchFilterIcon icon={<FilterOutlined />} filtered={filtered} />
      ),
    },
  );

  return columns;
};
