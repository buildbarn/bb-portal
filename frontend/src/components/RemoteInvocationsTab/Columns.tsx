import { Link } from "@tanstack/react-router";
import type { ColumnsType } from "antd/es/table";
import type { CompletedActionDetailsFragment } from "@/graphql/__generated__/graphql";

export const columns: ColumnsType<CompletedActionDetailsFragment> = [
  {
    title: "Action ID",
    dataIndex: "actionID",
    key: "actionId",
    render: (_, record) => (
      <Link
        to="/browser/$"
        params={{
          _splat: `${record.historicalExecuteResponseDigest.rev2InstanceName}/blobs/${record.historicalExecuteResponseDigest.digestFunction}/historical_execute_response/${record.historicalExecuteResponseDigest.hash}-${record.historicalExecuteResponseDigest.sizeBytes}`,
        }}
      >
        {record.actionID}
      </Link>
    ),
  },
  {
    title: "Target Label",
    dataIndex: "targetID",
    key: "targetID",
  },
  {
    title: "Instance Name",
    dataIndex: ["historicalExecuteResponseDigest", "rev2InstanceName"],
    key: "rev2InstanceName",
  },
  {
    title: "Mnemonic",
    dataIndex: "mnemonic",
    key: "mnemonic",
  },
  {
    title: "Tool Name",
    dataIndex: "toolName",
    key: "toolName",
  },
];
