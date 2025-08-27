import CodeLink from "@/components/CodeLink";
import { operationsStateToActionPageUrl } from "@/components/OperationsGrid/utils";
import type { WorkerState } from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import { readableDurationFromDates } from "@/utils/time";
import { type TableColumnsType, Typography } from "antd";
import type { ColumnType } from "antd/lib/table";
import PropertyTagList from "../PropertyTagList";

const workerIdColumn: ColumnType<WorkerState> = {
  key: "workerId",
  title: "Worker ID",
  render: (_, record) => (
    <PropertyTagList
      propertyList={Object.entries(record.id)
        .sort()
        .map(([property, value]) => ({ name: property, value: value }))}
    />
  ),
};

const workerTimeoutColumn: ColumnType<WorkerState> = {
  key: "workerTimeout",
  title: "Worker timeout",
  render: (_, record) => (
    <Typography.Text>
      {(record.timeout &&
        readableDurationFromDates(new Date(), record.timeout, {
          precision: 1,
          smallestUnit: "s",
        })) ||
        "∞"}
    </Typography.Text>
  ),
};

const operationTimeoutColumn: ColumnType<WorkerState> = {
  key: "operationTimeout",
  title: "Operation timeout",
  onCell: (value, _) => ({
    colSpan: value.currentOperation ? 1 : 4,
    align: "center",
  }),
  render: (_, record) => (
    <Typography.Text>
      {record.currentOperation
        ? (record.currentOperation?.timeout &&
            readableDurationFromDates(
              new Date(),
              record.currentOperation.timeout,
              { precision: 1, smallestUnit: "s" },
            )) ||
          "∞"
        : "Idle"}
    </Typography.Text>
  ),
};

const operationNameColumn: ColumnType<WorkerState> = {
  key: "operationName",
  title: "Operation name",
  onCell: (value, _) => ({ colSpan: value.currentOperation ? 1 : 0 }),
  render: (_, record) => (
    <CodeLink
      text={`${record.currentOperation?.name}`}
      url={
        (record.currentOperation?.name &&
          `/operations/${record.currentOperation?.name}`) ||
        ""
      }
      abbreviate
    />
  ),
};

const actionDigestColumn: ColumnType<WorkerState> = {
  key: "actionDigest",
  title: "Action digest",
  onCell: (value, _) => ({ colSpan: value.currentOperation ? 1 : 0 }),
  render: (_, record) => (
    <CodeLink
      text={`${record.currentOperation?.actionDigest?.hash}`}
      url={
        (record.currentOperation &&
          operationsStateToActionPageUrl(record.currentOperation)) ||
        ""
      }
      abbreviate
    />
  ),
};

const targetIdColumn: ColumnType<WorkerState> = {
  key: "targetId",
  title: "Target ID",
  onCell: (value, _) => ({ colSpan: value.currentOperation ? 1 : 0 }),
  render: (_, record) => (
    <Typography.Text>{record.currentOperation?.targetId}</Typography.Text>
  ),
};

const getColumns = (): TableColumnsType<WorkerState> => {
  return [
    workerIdColumn,
    workerTimeoutColumn,
    operationTimeoutColumn,
    operationNameColumn,
    actionDigestColumn,
    targetIdColumn,
  ];
};

export default getColumns;
