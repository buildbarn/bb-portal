import { SearchOutlined } from "@ant-design/icons";
import { Link } from "@tanstack/react-router";
import type { TableColumnsType } from "antd/lib";
import { SearchFilterIcon, SearchWidget } from "@/components/SearchWidgets";
import type { GetTargetsListQuery } from "@/graphql/__generated__/graphql";

export type TargetGridRowType = NonNullable<
  NonNullable<
    NonNullable<
      NonNullable<GetTargetsListQuery["findTargets"]>["edges"]
    >[number]
  >["node"]
>;

export const columns: TableColumnsType<TargetGridRowType> = [
  {
    title: "Target kind",
    dataIndex: "target-kind",
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
  },
  {
    title: "Label",
    dataIndex: "label",
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
  },
];
