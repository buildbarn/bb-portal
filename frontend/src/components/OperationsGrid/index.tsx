import { useQuery } from "@tanstack/react-query";
import { Table } from "antd";
import type React from "react";
import { useGrpcClients } from "@/context/GrpcClientsContext";
import { RequestMetadata } from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import type { OperationsFilterParams } from "@/routes/operations.index";
import themeStyles from "@/theme/theme.module.css";
import OperationsInvocationFilter from "../OperationsInvocationFilter";
import PortalAlert from "../PortalAlert";
import getColumns from "./Columns";

const PAGE_SIZE = 1000;

interface Props {
  filter: OperationsFilterParams;
}

const OperationsTable: React.FC<Props> = ({ filter }) => {
  const { buildQueueStateClient } = useGrpcClients();

  const { data, isLoading, isError, error } = useQuery({
    queryKey: ["operationsTable", filter],
    queryFn: buildQueueStateClient.listOperations.bind(window, {
      pageSize: PAGE_SIZE,
      filterInvocationId: filter
        ? {
            typeUrl: filter?.["@type"],
            value: RequestMetadata.encode(
              RequestMetadata.fromPartial(filter),
            ).finish(),
          }
        : undefined,
    }),
    staleTime: Number.POSITIVE_INFINITY,
    refetchOnMount: "always",
  });

  if (isError) {
    return (
      <PortalAlert
        showIcon
        type="error"
        message="Error fetching operation"
        description={
          error.message ||
          "Unknown error occurred while fetching data from the server."
        }
      />
    );
  }

  return (
    <>
      <OperationsInvocationFilter filter={filter} />
      <Table
        loading={isLoading}
        dataSource={data?.operations}
        columns={getColumns()}
        pagination={{ pageSize: PAGE_SIZE }}
        size="small"
        rowClassName={() => themeStyles.compactTable}
        rowKey={(item) => item.name}
        locale={{
          emptyText: "No active operations found (that you have access to).",
        }}
      />
    </>
  );
};

export default OperationsTable;
