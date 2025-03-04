import dayjs from "@/lib/dayjs";
import type { OperationState } from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import { type TableColumnsType, Typography } from "antd";
import type { ColumnType } from "antd/lib/table";
import Link from "next/link";
import OperationStatusTag from "../OperationStatusTag";
import { operationsStateToActionPageUrl } from "./utils";

const operationNameColumn: ColumnType<OperationState> = {
  title: "Operation name",
  dataIndex: "name",
  render: (value: string) => <Link href={`operations/${value}`}>{value}</Link>,
};

const timeoutColumn: ColumnType<OperationState> = {
  title: "Timeout",
  dataIndex: "timeout",
  render: (value: Date | undefined) => (
    <Typography.Text>
      {value ? `${dayjs(value).diff(undefined, "seconds")}s` : "∞"}
    </Typography.Text>
  ),
};

const actionDigestColumn: ColumnType<OperationState> = {
  title: "Action digest",
  key: "actionDigest",
  render: (record: OperationState) => (
    <Link
      href={operationsStateToActionPageUrl(record) || ""}
    >{`${record.actionDigest?.hash}-${record.actionDigest?.sizeBytes}`}</Link>
  ),
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
