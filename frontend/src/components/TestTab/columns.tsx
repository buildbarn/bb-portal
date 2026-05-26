import { SearchOutlined } from "@ant-design/icons";
import { Link } from "@tanstack/react-router";
import type { FilterValue } from "antd/es/table/interface";
import { SearchFilterIcon, SearchWidget } from "@/components/SearchWidgets";
import {
  OrderDirection,
  type TestSummaryNodeFragment,
  type TestSummaryOrder,
  TestSummaryOrderField,
  type TestSummaryWhereInput,
} from "@/graphql/__generated__/graphql";
import type { TableColumnTypeWithFilter } from "@/types/TableColumnTypeWithFilter";
import { TEST_STATUS_FILTERS } from "@/types/TestStatus";
import { readableDurationFromMilliseconds } from "@/utils/time";
import styles from "../../theme/theme.module.css";
import { type CacheLocation, CacheLocationTag } from "../CacheLocationTag";
import { TestStatusTag } from "../TestStatusTag";

export type TestTabRowType = TestSummaryNodeFragment & {
  cacheLocation: CacheLocation;
};

export const defaultSorting: TestSummaryOrder = {
  field: TestSummaryOrderField.TotalRunDurationInMs,
  direction: OrderDirection.Desc,
};

export const columns: TableColumnTypeWithFilter<
  TestTabRowType,
  TestSummaryWhereInput
>[] = [
  {
    key: "overallStatus",
    title: "Status",
    render: (_, record) => <TestStatusTag status={record.overallStatus} />,
    filters: TEST_STATUS_FILTERS,
    applyFilter: (value: FilterValue) => {
      if (!value || value.length === 0) {
        return undefined;
      }
      return [
        {
          overallStatusIn: value as TestSummaryWhereInput["overallStatusIn"],
        },
      ];
    },
  },
  {
    title: "Label",
    key: "label",
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
    applyFilter: (value: FilterValue) => {
      if (!value || value.length === 0) {
        return undefined;
      }
      return [
        {
          hasInvocationTargetWith: [
            {
              hasTargetWith: [
                {
                  labelContainsFold: value[0] as string,
                },
              ],
            },
          ],
        },
      ];
    },
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
    key: "totalRunDurationInMs",
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
