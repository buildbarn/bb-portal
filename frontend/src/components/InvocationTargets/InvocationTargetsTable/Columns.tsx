import { SearchOutlined } from "@ant-design/icons";
import type { TableColumnsType } from "antd";
import Link from "next/link";
import { getInvocationTargetAbortReasonFilterOptions } from "@/components/InvocationTargetAbortReasonTag/filter";
import NullBooleanTag from "@/components/NullableBooleanTag";
import SearchWidget, { SearchFilterIcon } from "@/components/SearchWidgets";
import { TargetDurationWarning } from "@/components/TargetDurationWarning";
import type { GetInvocationTargetsForInvocationQuery } from "@/graphql/__generated__/graphql";
import styles from "@/theme/theme.module.css";
import { readableDurationFromMilliseconds } from "@/utils/time";
import { generateLinkToTargetsPage } from "@/utils/urlGenerator";
import { InvocationTargetAbortReasonTag } from "../../InvocationTargetAbortReasonTag";
import { InvocationTargetTagList } from "../InvocationTargetTagList";

export type InvocationTargetsTableRowType = NonNullable<
  NonNullable<
    NonNullable<
      GetInvocationTargetsForInvocationQuery["bazelInvocation"]["invocationTargets"]["edges"]
    >[number]
  >["node"]
>;

const getTargetPageLink = (record: InvocationTargetsTableRowType) => {
  return generateLinkToTargetsPage(
    record.target.instanceName.name,
    record.target.label,
    record.target.aspect,
    record.target.targetKind,
  );
};

export const columns: TableColumnsType<InvocationTargetsTableRowType> = [
  {
    title: "Target kind",
    dataIndex: "target-kind",
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
  },
  {
    title: "Label",
    dataIndex: "label",
    filterSearch: true,
    render: (_, record) => (
      <Link href={getTargetPageLink(record)}>{record.target.label}</Link>
    ),
    filterDropdown: (filterProps) => (
      <SearchWidget placeholder="Target Pattern..." {...filterProps} />
    ),
    filterIcon: (filtered) => (
      <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />
    ),
  },
  {
    title: <TargetDurationWarning text="Duration" />,
    dataIndex: "duration",
    align: "right",
    render: (_, record) =>
      record.durationInMs && (
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
