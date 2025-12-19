import type { TableColumnsType } from "antd";
import Link from "next/link";
import type { GetTestDetailsQuery } from "@/graphql/__generated__/graphql";
import styles from "@/theme/theme.module.css";
import { readableDurationFromMilliseconds } from "@/utils/time";
import NullBooleanTag from "../NullableBooleanTag";
import TestStatusTag, { type TestStatusEnum } from "../TestStatusTag";

export type TestDetailsRowType = NonNullable<
  NonNullable<
    NonNullable<GetTestDetailsQuery["findTestSummaries"]["edges"]>[number]
  >["node"]
> & {
  cachedLocally: boolean | null;
  cachedRemotely: boolean | null;
};

export const columns: TableColumnsType<TestDetailsRowType> = [
  {
    title: "Status",
    dataIndex: "status",
    render: (_, record) => (
      <TestStatusTag
        displayText={true}
        key="status"
        status={record.overallStatus as TestStatusEnum}
      />
    ),
  },
  {
    title: "Invocation ID",
    dataIndex: "name",
    render: (_, record) => (
      <Link
        href={
          "/bazel-invocations/" +
          record.invocationTarget.bazelInvocation.invocationID
        }
      >
        {record.invocationTarget.bazelInvocation.invocationID}
      </Link>
    ),
  },
  {
    title: "Cached Locally",
    dataIndex: "cachedLocally",
    render: (x) => <NullBooleanTag key="local" status={x as boolean | null} />,
  },
  {
    title: "Cached Remotely",
    dataIndex: "cachedRemotely",
    render: (x) => <NullBooleanTag key="remote" status={x as boolean | null} />,
  },
  {
    title: "Duration",
    dataIndex: "duration",
    render: (_, record) => (
      <span className={styles.numberFormat}>
        {readableDurationFromMilliseconds(record.totalRunDurationInMs, {
          smallestUnit: "ms",
        })}
      </span>
    ),
    align: "right",
  },
];
