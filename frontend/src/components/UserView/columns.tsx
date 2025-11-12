import {
  buildColumn,
  durationColumn,
  invocationIdColumn,
  startedAtColumn,
  statusColumn,
} from "@/components/BazelInvocationColumns/Columns";

const getColumns = () => {
  return [
    invocationIdColumn,
    startedAtColumn,
    durationColumn,
    statusColumn,
    buildColumn,
  ];
};

export default getColumns;
