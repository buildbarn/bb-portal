import { useGrpcClients } from "@/context/GrpcClientsContext";
import { RequestMetadata } from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import themeStyles from "@/theme/theme.module.css";
import useScreenSize from "@/utils/screen";
import { CodeOutlined } from "@ant-design/icons";
import { useQuery } from "@tanstack/react-query";
import { Space, Table, Typography } from "antd";
import { useSearchParams } from "next/navigation";
import type React from "react";
import OperationsInvocationFilter from "../OperationsInvocationFilter";
import getColumns from "./Columns";

const PAGE_SIZE = 1000;

const OperationsTable: React.FC = () => {
  const screenSize = useScreenSize();
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

  const { data, isLoading, isError } = useQuery({
    queryKey: ["operationsTable", filterInvocationId],
    queryFn: buildQueueStateClient.listOperations.bind(window, {
      pageSize: PAGE_SIZE,
      filterInvocationId: getSerializedFilterInvocationId(filterInvocationId),
    }),
    staleTime: Number.POSITIVE_INFINITY,
  });

  let emptyText = "No active operations";
  if (isError) emptyText = "There was a problem fetching the operations";

  return (
    <>
      <OperationsInvocationFilter filterInvocationId={filterInvocationId} />
      <Table
        dataSource={data?.operations}
        columns={getColumns()}
        pagination={{ pageSize: PAGE_SIZE }}
        size="small"
        rowClassName={() => themeStyles.compactTable}
        rowKey={(item) => item.name}
        locale={{
          emptyText: isLoading ? (
            <Typography.Text
              disabled
              className={themeStyles.tableEmptyTextTypography}
            >
              <Space>
                <CodeOutlined />
                Loading...
              </Space>
            </Typography.Text>
          ) : (
            <Typography.Text
              disabled
              className={themeStyles.tableEmptyTextTypography}
            >
              <Space>
                <CodeOutlined />
                {emptyText}
              </Space>
            </Typography.Text>
          ),
        }}
      />
    </>
  );
};

export default OperationsTable;
