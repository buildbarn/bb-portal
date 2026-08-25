import { type TableColumnsType, Tooltip } from "antd";
import { CodeLink } from "@/components/CodeLink";
import type { WorkerState } from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import { readableDurationFromDates } from "@/utils/time";
import { operationsStateToBrowserSplat } from "../OperationStateDisplay/utils";
import PropertyTagList from "../PropertyTagList";

export const columns: TableColumnsType<WorkerState> = [
  {
    key: "workerId",
    title: "Worker ID",
    render: (_, record) => (
      <PropertyTagList
        propertyList={Object.entries(record.id)
          .sort()
          .map(([property, value]) => ({ name: property, value: value }))}
      />
    ),
  },
  {
    key: "workerTimeout",
    title: "Worker timeout",
    render: (_, record) =>
      (record.timeout &&
        readableDurationFromDates(new Date(), record.timeout, {
          precision: 1,
          smallestUnit: "s",
        })) ||
      "∞",
  },
  {
    key: "operationTimeout",
    title: "Operation timeout",
    onCell: (value, _) => ({
      colSpan: value.currentOperation ? 1 : 4,
      align: "center",
    }),
    render: (_, record) => {
      if (!record.currentOperation) {
        // These values span the rest of the table
        if (record.timeout) {
          return (
            <Tooltip title="This worker is executing an action which you are not allowed to view.">
              Executing
            </Tooltip>
          );
        }
        return "Idle";
      }
      if (!record.currentOperation.timeout) {
        return "∞";
      }
      return readableDurationFromDates(
        new Date(),
        record.currentOperation.timeout,
        { precision: 1, smallestUnit: "s" },
      );
    },
  },
  {
    key: "operationName",
    title: "Operation name",
    onCell: (value, _) => ({ colSpan: value.currentOperation ? 1 : 0 }),
    render: (_, record) => {
      const operationID = record.currentOperation?.name;
      if (operationID) {
        return (
          <CodeLink
            truncate
            text={operationID}
            link={{
              to: "/operations/$operationID",
              params: { operationID },
            }}
          />
        );
      } else {
        return "-";
      }
    },
  },
  {
    key: "actionDigest",
    title: "Action digest",
    onCell: (value, _) => ({ colSpan: value.currentOperation ? 1 : 0 }),
    render: (_, record) => {
      if (record.currentOperation?.actionDigest) {
        return (
          <CodeLink
            truncate
            text={`${record.currentOperation.actionDigest.hash}-${record.currentOperation.actionDigest.sizeBytes}`}
            link={{
              to: "/browser/$",
              params: {
                _splat: operationsStateToBrowserSplat(record.currentOperation),
              },
            }}
          />
        );
      } else {
        return "-";
      }
    },
  },
  {
    key: "targetId",
    title: "Target ID",
    onCell: (value, _) => ({ colSpan: value.currentOperation ? 1 : 0 }),
    render: (_, record) => record.currentOperation?.targetId,
  },
];
