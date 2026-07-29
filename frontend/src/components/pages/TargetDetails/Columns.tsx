import { SearchOutlined } from "@ant-design/icons";
import { Link } from "@tanstack/react-router";
import type { FilterValue } from "antd/es/table/interface";
import { InvocationTargetAbortReasonTag } from "@/components/InvocationTargetAbortReasonTag";
import { getInvocationTargetAbortReasonFilterOptions } from "@/components/InvocationTargetAbortReasonTag/filter";
import { InvocationTargetTagList } from "@/components/InvocationTargets/InvocationTargetTagList";
import NullBooleanTag from "@/components/NullableBooleanTag";
import { SearchFilterIcon } from "@/components/SearchWidgets";
import type {
  InvocationTargetAbortReason,
  InvocationTargetDetailsFragment,
  InvocationTargetWhereInput,
} from "@/graphql/__generated__/graphql";
import type { TableColumnTypeWithFilter } from "@/types/TableColumnTypeWithFilter";

export const columns: TableColumnTypeWithFilter<
  InvocationTargetDetailsFragment,
  InvocationTargetWhereInput
>[] = [
  {
    title: "Invocation ID",
    key: "name",
    render: (_, record) => (
      <Link
        to="/bazel-invocations/$invocationID"
        params={{ invocationID: record.bazelInvocation.invocationID }}
      >
        {record.bazelInvocation.invocationID}
      </Link>
    ),
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
      return [
        {
          success: value[0] as boolean,
        },
      ];
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
