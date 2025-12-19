import { SearchOutlined } from "@ant-design/icons";
import type { TableColumnsType } from "antd/lib";
import { SearchFilterIcon, SearchWidget } from "@/components/SearchWidgets";
import {
  type GetTestsForInvocationQuery,
  OrderDirection,
  type TestSummaryOrder,
  TestSummaryOrderField,
} from "@/graphql/__generated__/graphql";
import { readableDurationFromMilliseconds } from "@/utils/time";
import { generateLinkToTargetsPage } from "@/utils/urlGenerator";
import styles from "../../theme/theme.module.css";
import Link from "../Link";
import NullBooleanTag from "../NullableBooleanTag";
import TestStatusTag, { type TestStatusEnum } from "../TestStatusTag";

export type TestTabRowType = NonNullable<
  NonNullable<
    NonNullable<
      GetTestsForInvocationQuery["findTestSummaries"]["edges"]
    >[number]
  >["node"]
> & {
  cachedLocally: boolean | null;
  cachedRemotely: boolean | null;
};

export const defaultSorting: TestSummaryOrder = {
  field: TestSummaryOrderField.TotalRunDurationInMs,
  direction: OrderDirection.Desc,
};

export const columns: TableColumnsType<TestTabRowType> = [
  {
    title: "Status",
    dataIndex: "overallStatus",
    render: (x) => (
      <TestStatusTag
        displayText={true}
        key="status"
        status={x as TestStatusEnum}
      />
    ),
    filters: [
      {
        text: "No Status",
        value: "NO_STATUS",
      },
      {
        text: "Passed",
        value: "PASSED",
      },
      {
        text: "Flaky",
        value: "FLAKY",
      },
      {
        text: "Timeout",
        value: "TIMEOUT",
      },
      {
        text: "Failed",
        value: "FAILED",
      },
      {
        text: "Incomplete",
        value: "INCOMPLETE",
      },
      {
        text: "Remote Failure",
        value: "REMOTE_FAILURE",
      },
      {
        text: "Failed to Build",
        value: "FAILED_TO_BUILD",
      },
      {
        text: "Tool Halted Before Testing",
        value: "TOOL_HALTED_BEFORE_TESTING",
      },
    ],
  },
  {
    title: "Label",
    dataIndex: "label",
    render: (_, record) => (
      <Link
        href={generateLinkToTargetsPage(
          record.invocationTarget.target.instanceName.name,
          record.invocationTarget.target.label,
          record.invocationTarget.target.aspect,
          record.invocationTarget.target.targetKind,
        )}
      >
        {record.invocationTarget.target.label}
      </Link>
    ),
    filterSearch: true,
    filterDropdown: (filterProps) => (
      <SearchWidget placeholder="Target Pattern..." {...filterProps} />
    ),
    filterIcon: (filtered) => (
      <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />
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
    dataIndex: "totalRunDurationInMs",
    render: (_, record) => (
      <span className={styles.numberFormat}>
        {readableDurationFromMilliseconds(record.totalRunDurationInMs, {
          smallestUnit: "ms",
        })}
      </span>
    ),
    align: "right",
    sortDirections: ["ascend", "descend", "ascend"],
    defaultSortOrder:
      defaultSorting.direction === OrderDirection.Asc ? "ascend" : "descend",
    // Using a dummy sorter function as sorting is handled server-side
    sorter: (_a, _b) => 0,
  },
];
