import { useQuery } from "@tanstack/react-query";
import { Table } from "antd";
import { useSearchParams } from "next/navigation";
import type React from "react";
import { useGrpcClients } from "@/context/GrpcClientsContext";
import { RequestMetadata } from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import themeStyles from "@/theme/theme.module.css";
import OperationsInvocationFilter from "../OperationsInvocationFilter";
import PortalAlert from "../PortalAlert";
import getColumns from "./Columns";

const PAGE_SIZE = 1000;

const OperationsTable: React.FC = () => {
  const { buildQueueStateClient } = useGrpcClients();
  const searchParams = useSearchParams();
  const filterInvocationId = searchParams.get("filter_invocation_id");

  const getSerializedFilterInvocationId = (searchParam: string | null) => {
    if (searchParam) {
      const parsedSearchParam = JSON.parse(decodeURIComponent(searchParam));
      return {
        typeUrl: parsedSearchParam["@type"],
        value: RequestMetadata.encode(parsedSearchParam).finish(),
      };
    }
    return undefined;
  };

  const { data, isLoading, isError, error } = useQuery({
    queryKey: ["operationsTable", filterInvocationId],
    queryFn: buildQueueStateClient.listOperations.bind(window, {
      pageSize: PAGE_SIZE,
      filterInvocationId: getSerializedFilterInvocationId(filterInvocationId),
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
      <OperationsInvocationFilter filterInvocationId={filterInvocationId} />
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
