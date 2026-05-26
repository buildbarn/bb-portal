import { SearchOutlined } from "@ant-design/icons";
import { Link } from "@tanstack/react-router";
import type { FilterValue } from "antd/es/table/interface";
import { SearchFilterIcon, SearchWidget } from "@/components/SearchWidgets";
import type {
  TargetWhereInput,
  TestListRowFragment,
} from "@/graphql/__generated__/graphql";
import type { TableColumnTypeWithFilter } from "@/types/TableColumnTypeWithFilter";

export const columns: TableColumnTypeWithFilter<
  TestListRowFragment,
  TargetWhereInput
>[] = [
  {
    key: "target",
    title: "Target",
    render: (_, record) => (
      <Link to="/targets/$targetID/tests" params={{ targetID: record.id }}>
        {record.label}
      </Link>
    ),
    filterDropdown: (filterProps) => (
      <SearchWidget placeholder="Target Pattern..." {...filterProps} />
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
          labelContainsFold: value[0] as string,
        },
      ];
    },
  },
  {
    key: "instanceName",
    title: "Instance Name",
    filterDropdown: (filterProps) => (
      <SearchWidget placeholder="Instance Name Pattern..." {...filterProps} />
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
          hasInstanceNameWith: [
            {
              nameContainsFold: value[0] as string,
            },
          ],
        },
      ];
    },
  },
];
