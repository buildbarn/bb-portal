import { InfoCircleOutlined, SearchOutlined } from "@ant-design/icons";
import { Link } from "@tanstack/react-router";
import { Popover } from "antd";
import type { FilterValue } from "antd/es/table/interface";
import dayjs from "dayjs";
import { validate as uuidValidate } from "uuid";
import {
  SearchFilterIcon,
  SearchWidget,
  TimeRangeSelector,
} from "@/components/SearchWidgets";
import type {
  BuildNodeFragment,
  BuildWhereInput,
} from "@/graphql/__generated__/graphql";
import type { TableColumnTypeWithFilter } from "@/types/TableColumnTypeWithFilter";
import { env } from "@/utils/env";
import { parseGraphqlEdgeList } from "@/utils/parseGraphqlEdgeList";
import CodeText from "../CodeText";
import { OptionalLinkWrapper } from "../OptionalLinkWrapper";

export const getColumns = (): TableColumnTypeWithFilter<
  BuildNodeFragment,
  BuildWhereInput
>[] => {
  const columns: TableColumnTypeWithFilter<
    BuildNodeFragment,
    BuildWhereInput
  >[] = [];

  const additionalColumns = env.additionalBuildColumns;
  for (const column of additionalColumns) {
    columns.push({
      key: column.valueKey,
      title: column.title,
      filterSearch: true,
      render: (_, record) => {
        const tags = parseGraphqlEdgeList(record.tags);
        const valueTags = tags.filter((tag) => tag.key === column.valueKey);
        const urlTags = tags.filter((tag) => tag.key === column.urlKey);
        const singleUrl = urlTags.length === 1 ? urlTags[0].value : undefined;

        return (
          <>
            <OptionalLinkWrapper url={singleUrl}>
              {valueTags.map((tag) => tag.value).join(", ")}
            </OptionalLinkWrapper>
            {urlTags.length > 1 && (
              <Popover
                title="This field has multiple urls:"
                content={urlTags.map((tag) => (
                  <a key={tag.id} href={tag.value}>
                    {tag.value}
                  </a>
                ))}
              >
                <InfoCircleOutlined style={{ marginLeft: "7px" }} />
              </Popover>
            )}
          </>
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
      key: "buildUUID",
      title: "Build ID",
      render: (_, record) => (
        <Link
          to={`/builds/$buildUUID`}
          params={{ buildUUID: record.buildUUID }}
        >
          {record.buildUUID}
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
        const buildUUID = value[0] as string;
        if (!uuidValidate(buildUUID)) {
          return undefined;
        }
        return [{ buildUUID: buildUUID as string }];
      },
    },
    {
      key: "buildDate",
      title: "Timestamp",
      render: (_, record) => (
        <CodeText>
          {dayjs(record.timestamp).format("YYYY-MM-DD hh:mm:ss A")}
        </CodeText>
      ),
      filterDropdown: (filterProps) => <TimeRangeSelector {...filterProps} />,
      filterIcon: (filtered) => (
        <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />
      ),
      applyFilter: (value: FilterValue) => {
        if (value.length !== 2) {
          return undefined;
        }
        const newFilters = [];
        if (value[0]) {
          newFilters.push({ timestampGTE: value[0] });
        }
        if (value[1]) {
          newFilters.push({ timestampLTE: value[1] });
        }
        return newFilters;
      },
    },
  );

  return columns;
};
