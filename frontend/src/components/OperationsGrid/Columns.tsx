import CodeLink from "@/components/CodeLink";
import type { OperationState } from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import { readableDurationFromDates } from "@/utils/time";
import { type TableColumnsType, Typography } from "antd";
import type { ColumnType } from "antd/lib/table";
import {
  historicalExecuteResponseDigestFromUrl,
  historicalExecuteResponseUrlFromOperation,
} from "../OperationStateDisplay/utils";
import OperationStatusTag from "../OperationStatusTag";
import { operationsStateToActionPageUrl } from "./utils";

const operationNameColumn: ColumnType<OperationState> = {
  title: "Operation name",
  dataIndex: "name",
  render: (value: string) => (
    <CodeLink url={`/operations/${value}`} text={value} abbreviate />
  ),
};

const timeoutColumn: ColumnType<OperationState> = {
  title: "Timeout",
  dataIndex: "timeout",
  render: (value: Date | undefined) => (
    <Typography.Text>
      {value
        ? readableDurationFromDates(new Date(), value, {
            precision: 1,
            smallestUnit: "s",
          })
        : "∞"}
    </Typography.Text>
  ),
};

const actionDigestColumn: ColumnType<OperationState> = {
  title: "Action digest / Historical execute response digest",
  key: "actionDigest",
  render: (record: OperationState) => {
    const historical_execute_response_url =
      historicalExecuteResponseUrlFromOperation(record);
    const historical_execute_response_digest =
      historicalExecuteResponseDigestFromUrl(historical_execute_response_url);

    if (historical_execute_response_digest && historical_execute_response_url) {
      return (
        <CodeLink
          text={historical_execute_response_digest}
          url={historical_execute_response_url}
          abbreviate
        />
      );
    }

    return (
      <CodeLink
        text={`${record.actionDigest?.hash}`}
        url={`${operationsStateToActionPageUrl(record)}`}
        abbreviate
      />
    );
  },
};

const targetIdColumn: ColumnType<OperationState> = {
  title: "Target ID",
  dataIndex: "targetId",
};

const statusColumn: ColumnType<OperationState> = {
  title: "Status",
  key: "status",
  render: (record: OperationState) => <OperationStatusTag operation={record} />,
};

const getColumns = (): TableColumnsType<OperationState> => {
  return [
    timeoutColumn,
    operationNameColumn,
    actionDigestColumn,
    targetIdColumn,
    statusColumn,
  ];
};

export default getColumns;
