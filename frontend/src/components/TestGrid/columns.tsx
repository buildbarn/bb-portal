import { SearchOutlined } from "@ant-design/icons";
import type { TableColumnsType } from "antd/lib";
import Link from "next/link";
import { SearchFilterIcon, SearchWidget } from "@/components/SearchWidgets";
import type { GetTestsQuery } from "@/graphql/__generated__/graphql";
import { generateLinkToTestPage } from "@/utils/urlGenerator";

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
      <Link
        href={generateLinkToTestPage(
          record.instanceName.name,
          record.label,
          record.aspect,
          record.targetKind,
        )}
      >
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
