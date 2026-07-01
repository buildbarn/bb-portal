import { Link } from "@tanstack/react-router";
import type { TableColumnsType } from "antd";
import {
  type CacheLocation,
  CacheLocationTag,
} from "@/components/CacheLocationTag";
import type { GetTestDetailsQuery } from "@/graphql/__generated__/graphql";
import styles from "@/theme/theme.module.css";
import { readableDurationFromMilliseconds } from "@/utils/time";
import { TestStatusTag } from "../../TestStatusTag";

export type TestDetailsRowType = NonNullable<
  NonNullable<
    NonNullable<GetTestDetailsQuery["findTestSummaries"]["edges"]>[number]
  >["node"]
> & {
  cacheLocation: CacheLocation;
};

export const columns: TableColumnsType<TestDetailsRowType> = [
  {
    key: "status",
    title: "Status",
    render: (_, record) => <TestStatusTag status={record.overallStatus} />,
  },
  {
    title: "Invocation ID",
    dataIndex: "name",
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
