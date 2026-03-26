import { type TableColumnsType, Typography } from "antd";
import type { ColumnType } from "antd/lib/table";
import { CodeLink } from "@/components/CodeLink";
import type { OperationState } from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import { BrowserPageType } from "@/types/BrowserPageType";
import { readableDurationFromDates } from "@/utils/time";
import { generateBrowserSplat } from "@/utils/urlGenerator";
import { historicalExecuteResponseDigestFromOperation } from "../OperationStateDisplay/utils";
import OperationStatusTag from "../OperationStatusTag";
import { operationsStateToBrowserSplat } from "./utils";

const operationNameColumn: ColumnType<OperationState> = {
  key: "name",
  title: "Operation name",
  render: (_, record) => (
    <CodeLink
      text={record.name}
      link={{
        to: "/operations/$operationID",
        params: { operationID: record.name },
      }}
    />
  ),
};

const timeoutColumn: ColumnType<OperationState> = {
  key: "timeout",
  title: "Timeout",
  render: (_, record) => (
    <Typography.Text>
      {record.timeout
        ? readableDurationFromDates(new Date(), record.timeout, {
            precision: 1,
            smallestUnit: "s",
          })
        : "∞"}
    </Typography.Text>
  ),
};

const actionDigestColumn: ColumnType<OperationState> = {
  key: "actionDigest",
  title: "Action digest / Historical execute response digest",
  render: (_, record: OperationState) => {
    const historicalExecuteResponseBrowserPageParams =
      historicalExecuteResponseDigestFromOperation(record);

    if (historicalExecuteResponseBrowserPageParams) {
      return (
        <CodeLink
          text={historicalExecuteResponseBrowserPageParams.digest.hash}
          link={{
            to: "/browser/$",
            params: {
              _splat: generateBrowserSplat(
                historicalExecuteResponseBrowserPageParams.instanceName,
                historicalExecuteResponseBrowserPageParams.digestFunction,
                historicalExecuteResponseBrowserPageParams.digest,
                BrowserPageType.HistoricalExecuteResponse,
              ),
            },
          }}
        />
      );
    }

    return (
      <CodeLink
        text={`${record.actionDigest?.hash}-${record.actionDigest?.sizeBytes}`}
        link={{
          to: "/browser/$",
          params: { _splat: operationsStateToBrowserSplat(record) },
        }}
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
