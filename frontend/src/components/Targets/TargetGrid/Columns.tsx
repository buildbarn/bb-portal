import { SearchOutlined } from "@ant-design/icons";
import { Link } from "@tanstack/react-router";
import type { FilterValue } from "antd/es/table/interface";
import { SearchFilterIcon, SearchWidget } from "@/components/SearchWidgets";
import type {
  TargetListDetailsFragment,
  TargetWhereInput,
} from "@/graphql/__generated__/graphql";
import type { TableColumnTypeWithFilter } from "@/types/TableColumnTypeWithFilter";

export const columns: TableColumnTypeWithFilter<
  TargetListDetailsFragment,
  TargetWhereInput
>[] = [
  {
    title: "Target kind",
    key: "target-kind",
    filterSearch: true,
    render: (_, record) => (
      <span>
        {record.aspect === ""
          ? record.targetKind
          : `${record.targetKind} (${record.aspect})`}
      </span>
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
      return [{ targetKindContainsFold: value[0] as string }];
    },
  },
  {
    title: "Label",
    key: "label",
    filterSearch: true,
    render: (_, record) => (
      <Link to="/targets/$targetID" params={{ targetID: record.id }}>
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
      return [{ labelContainsFold: value[0] as string }];
    },
  },
];
