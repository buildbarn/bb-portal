import type { FilterValue } from "antd/es/table/interface";
import type { TableColumnTypeWithFilter } from "@/types/TableColumnTypeWithFilter";

export const tableFiltersToGraphqlWhere = <ColumnType, WhereInput>(
  columns: TableColumnTypeWithFilter<ColumnType, WhereInput>[],
  filters: Record<string, FilterValue | null>,
): WhereInput[] => {
  return columns.flatMap((column) => {
    if (!column.applyFilter) {
      return [];
    }
    const filterValue = filters[`${column.key}`];
    if (filterValue === null) {
      return [];
    }
    const filter = column.applyFilter(filterValue);
    if (filter === undefined) {
      return [];
    }
    return filter;
  });
};
