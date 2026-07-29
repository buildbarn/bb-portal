import { Link } from "@tanstack/react-router";
import type { TableColumnsType } from "antd";
import {
  type CacheLocation,
  CacheLocationTag,
} from "@/components/CacheLocationTag";
import type { TestSummaryRowFragment } from "@/graphql/__generated__/graphql";
import styles from "@/theme/theme.module.css";
import { readableDurationFromMilliseconds } from "@/utils/time";
import { TestStatusTag } from "../../TestStatusTag";

export type TestDetailsRowType = TestSummaryRowFragment & {
  cacheLocation: CacheLocation;
};

export const columns: TableColumnsType<TestDetailsRowType> = [
  {
    key: "status",
    title: "Status",
    render: (_, record) => <TestStatusTag status={record.overallStatus} />,
  },
  {
    key: "invocationID",
    title: "Invocation ID",
    render: (_, record) => (
      <Link
        to="/bazel-invocations/$invocationID"
        params={{
          invocationID: record.invocationTarget.bazelInvocation.invocationID,
        }}
      >
        {record.invocationTarget.bazelInvocation.invocationID}
      </Link>
    ),
  },
  {
    key: "cacheLocation",
    title: "Cache Location",
    render: (_, record) => (
      <CacheLocationTag cacheLocation={record.cacheLocation} />
    ),
  },
  {
    key: "duration",
    title: "Duration",
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
