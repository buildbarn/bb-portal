import { SearchOutlined } from "@ant-design/icons";
import { Link } from "@tanstack/react-router";
import type { TableColumnsType } from "antd/lib";
import { SearchFilterIcon, SearchWidget } from "@/components/SearchWidgets";
import {
  type GetTestsForInvocationQuery,
  OrderDirection,
  type TestSummaryOrder,
  TestSummaryOrderField,
} from "@/graphql/__generated__/graphql";
import { TEST_STATUS_FILTERS } from "@/types/TestStatus";
import { readableDurationFromMilliseconds } from "@/utils/time";
import styles from "../../theme/theme.module.css";
import { type CacheLocation, CacheLocationTag } from "../CacheLocationTag";
import { TestStatusTag } from "../TestStatusTag";

export type TestTabRowType = NonNullable<
  NonNullable<
    NonNullable<
      GetTestsForInvocationQuery["findTestSummaries"]["edges"]
    >[number]
  >["node"]
> & {
  cacheLocation: CacheLocation;
};

export const defaultSorting: TestSummaryOrder = {
  field: TestSummaryOrderField.TotalRunDurationInMs,
  direction: OrderDirection.Desc,
};

export const columns: TableColumnsType<TestTabRowType> = [
  {
    key: "overallStatus",
    title: "Status",
    render: (_, record) => <TestStatusTag status={record.overallStatus} />,
    filters: TEST_STATUS_FILTERS,
  },
  {
    title: "Label",
    dataIndex: "label",
    render: (_, record) => (
      <Link
        to="/targets/$targetID/tests"
        params={{ targetID: record.invocationTarget.target.id }}
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
    key: "cacheLocation",
    title: "Cache Location",
    render: (_, record) => (
      <CacheLocationTag cacheLocation={record.cacheLocation} />
    ),
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
