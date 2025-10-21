import { SearchOutlined } from "@ant-design/icons";
import type { TableColumnsType } from "antd/lib";
import Link from "next/link";
import { SearchFilterIcon, SearchWidget } from "@/components/SearchWidgets";
import type { GetTargetsListQuery } from "@/graphql/__generated__/graphql";
import { generateLinkToTargetsPage } from "@/utils/urlGenerator";

export type TargetGridRowType = NonNullable<
  NonNullable<
    NonNullable<
      NonNullable<GetTargetsListQuery["findTargets"]>["edges"]
    >[number]
  >["node"]
>;

const getTargetPageLink = (record: TargetGridRowType) => {
  return generateLinkToTargetsPage(
    record.instanceName.name,
    record.label,
    record.aspect,
    record.targetKind,
  );
};

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
      <Link href={getTargetPageLink(record)}>{record.label}</Link>
    ),
    filterDropdown: (filterProps) => (
      <SearchWidget placeholder="Target Pattern..." {...filterProps} />
    ),
    filterIcon: (filtered) => (
      <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />
    ),
  },
];
