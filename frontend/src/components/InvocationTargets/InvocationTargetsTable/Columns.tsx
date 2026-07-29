import { SearchOutlined } from "@ant-design/icons";
import { Link } from "@tanstack/react-router";
import type { FilterValue } from "antd/es/table/interface";
import { getInvocationTargetAbortReasonFilterOptions } from "@/components/InvocationTargetAbortReasonTag/filter";
import NullBooleanTag from "@/components/NullableBooleanTag";
import SearchWidget, { SearchFilterIcon } from "@/components/SearchWidgets";
import type {
  BazelInvocationTargetsFragment,
  InvocationTargetAbortReason,
  InvocationTargetWhereInput,
} from "@/graphql/__generated__/graphql";
import type { TableColumnTypeWithFilter } from "@/types/TableColumnTypeWithFilter";
import { InvocationTargetAbortReasonTag } from "../../InvocationTargetAbortReasonTag";
import { InvocationTargetTagList } from "../InvocationTargetTagList";

export const columns: TableColumnTypeWithFilter<
  BazelInvocationTargetsFragment,
  InvocationTargetWhereInput
>[] = [
  {
    title: "Target kind",
    key: "target-kind",
    filterSearch: true,
    render: (_, record) => (
      <span>
        {record.target.aspect === ""
          ? record.target.targetKind
          : `${record.target.targetKind} (${record.target.aspect})`}
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
      return [
        {
          hasTargetWith: [{ targetKindContainsFold: value[0] as string }],
        },
      ];
    },
  },
  {
    title: "Label",
    key: "label",
    filterSearch: true,
    render: (_, record) => (
      <Link to="/targets/$targetID" params={{ targetID: record.target.id }}>
        {record.target.label}
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
          hasTargetWith: [{ labelContainsFold: value[0] as string }],
        },
      ];
    },
  },
  {
    title: "Overall Success",
    key: "success",
    render: (_, record) => (
      <NullBooleanTag key="success" status={record.success} />
    ),
    filters: [
      {
        text: "Yes",
        value: true,
      },
      {
        text: "No",
        value: false,
      },
    ],
    filterIcon: (filtered) => (
      <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />
    ),
    applyFilter: (value: FilterValue) => {
      if (value.length === 0) {
        return undefined;
      }
      return [{ success: value[0] as boolean }];
    },
  },
  {
    title: "Abort Reason",
    key: "abort-reason",
    filters: getInvocationTargetAbortReasonFilterOptions(),
    render: (_, record) => (
      <InvocationTargetAbortReasonTag reason={record.abortReason} />
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
          abortReasonIn: value as InvocationTargetAbortReason[],
        },
      ];
    },
  },
  {
    title: "Tags",
    key: "tags",
    render: (_, record) => <InvocationTargetTagList tags={record.tags} />,
  },
];
