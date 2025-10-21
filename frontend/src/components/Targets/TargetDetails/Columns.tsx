import { SearchOutlined } from "@ant-design/icons";
import type { TableColumnsType } from "antd";
import Link from "next/link";
import { InvocationTargetAbortReasonTag } from "@/components/InvocationTargetAbortReasonTag";
import { getInvocationTargetAbortReasonFilterOptions } from "@/components/InvocationTargetAbortReasonTag/filter";
import { InvocationTargetTagList } from "@/components/InvocationTargets/InvocationTargetTagList";
import NullBooleanTag from "@/components/NullableBooleanTag";
import { SearchFilterIcon } from "@/components/SearchWidgets";
import { TargetDurationWarning } from "@/components/TargetDurationWarning";
import type { GetTargetDetailsQuery } from "@/graphql/__generated__/graphql";
import styles from "@/theme/theme.module.css";
import { readableDurationFromMilliseconds } from "@/utils/time";

export type InvocationTargetRowType = NonNullable<
  NonNullable<
    NonNullable<
      NonNullable<
        GetTargetDetailsQuery["getTarget"]
      >["invocationTargets"]["edges"]
    >[number]
  >["node"]
>;

export const columns: TableColumnsType<InvocationTargetRowType> = [
  {
    title: "Invocation ID",
    dataIndex: "name",
    render: (_, record) => (
      <Link href={`/bazel-invocations/${record.bazelInvocation.invocationID}`}>
        {record.bazelInvocation.invocationID}
      </Link>
    ),
  },

  {
    title: <TargetDurationWarning text="Duration" />,
    dataIndex: "duration",
    align: "right",
    render: (_, record) =>
      record.durationInMs !== undefined &&
      record.durationInMs !== null && (
        <span className={styles.numberFormat}>
          {readableDurationFromMilliseconds(record.durationInMs, {
            smallestUnit: "ms",
          })}
        </span>
      ),
  },
  {
    title: "Overall Success",
    dataIndex: "success",
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
  },
  {
    title: "Abort Reason",
    dataIndex: "abort-reason",
    filters: getInvocationTargetAbortReasonFilterOptions(),
    render: (_, record) => (
      <InvocationTargetAbortReasonTag reason={record.abortReason} />
    ),
    filterIcon: (filtered) => (
      <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />
    ),
  },
  {
    title: "Tags",
    dataIndex: "tags",
    render: (_, record) => <InvocationTargetTagList tags={record.tags} />,
  },
];
