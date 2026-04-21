import type { TableColumnType } from "antd";
import type { FilterValue } from "antd/es/table/interface";

export type TableColumnTypeWithFilter<T, F> = TableColumnType<T> & {
  applyFilter?: (value: FilterValue) => F[] | undefined;
};
