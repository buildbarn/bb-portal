import { SearchOutlined } from "@ant-design/icons";
import type { TableColumnsType } from "antd/lib";
import { Link } from '@tanstack/react-router';
import { SearchFilterIcon, SearchWidget } from "@/components/SearchWidgets";
import type { GetTestsQuery } from "@/graphql/__generated__/graphql";

export type TestGridRowDataType = NonNullable<
  NonNullable<
    NonNullable<NonNullable<GetTestsQuery["findTargets"]>["edges"]>[number]
  >["node"]
>;

export const columns: TableColumnsType<TestGridRowDataType> = [
  {
    title: "Target",
    dataIndex: "target",
    render: (_, record) => (
      <Link to="/targets/$targetID/tests" params={{targetID: record.id}}>
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
  {
    key: "instanceName",
    title: "Instance Name",
    dataIndex: ["instanceName", "name"],
    filterDropdown: (filterProps) => (
      <SearchWidget placeholder="Instance Name Pattern..." {...filterProps} />
    ),
    filterIcon: (filtered) => (
      <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />
    ),
  },
];
