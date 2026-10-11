import { type TableColumnsType, Tooltip } from "antd";
import { CodeLink } from "@/components/CodeLink";
import { operationsStateToBrowserSplat } from "@/components/OperationStateDisplay/utils";
import type {
  InvocationName,
  OperationState,
} from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import {
  readableDurationFromDates,
  readableDurationFromProtobufDuration,
} from "@/utils/time";
import { withInvocationName } from "./utils";

const timeFormat = { precision: 1, smallestUnit: "s" } as const;

const censored = (
  <Tooltip title="This operation belongs to an instance name you are not allowed to view.">
    -
  </Tooltip>
);

const getColumns = (
  invocationName: InvocationName,
): TableColumnsType<OperationState> => [
  {
    key: "priority",
    title: "Priority",
    render: (_, record) => record.priority,
  },
  {
    key: "expectedDuration",
    title: "Expected duration",
    render: (_, record) =>
      record.expectedDuration &&
      readableDurationFromProtobufDuration(record.expectedDuration),
  },
  {
    key: "age",
    title: "Age",
    render: (_, record) =>
      record.queuedTimestamp &&
      readableDurationFromDates(record.queuedTimestamp, new Date(), timeFormat),
  },
  {
    key: "timeout",
    title: "Timeout",
    render: (_, record) =>
      record.timeout
        ? readableDurationFromDates(new Date(), record.timeout, timeFormat)
        : "∞",
  },
  {
    key: "operationName",
    title: "Operation name",
    render: (_, record) =>
      record.name ? (
        <CodeLink
          truncate
          text={record.name}
          link={{
            to: "/operations/$operationID",
            params: { operationID: record.name },
          }}
        />
      ) : (
        censored
      ),
  },
  {
    key: "actionDigest",
    title: "Action digest",
    render: (_, record) => {
      if (!record.actionDigest) {
        return censored;
      }
      return (
        <CodeLink
          truncate
          text={`${record.actionDigest.hash}-${record.actionDigest.sizeBytes}`}
          link={{
            to: "/browser/$",
            params: {
              _splat: operationsStateToBrowserSplat(
                withInvocationName(record, invocationName),
              ),
            },
          }}
        />
      );
    },
  },
  {
    key: "targetId",
    title: "Target ID",
    render: (_, record) => record.targetId || "-",
  },
];

export default getColumns;
