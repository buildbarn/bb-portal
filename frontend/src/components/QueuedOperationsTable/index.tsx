import { useQuery } from "@tanstack/react-query";
import { Flex, Row, Space, Table } from "antd";
import type React from "react";
import PortalAlert from "@/components/PortalAlert";
import WorkersTablePageSelector from "@/components/WorkersTablePageSelector";
import { buildQueueStateClient } from "@/grpc/buildQueueStateClient";
import { Route } from "@/routes/scheduler.invocation";
import themeStyles from "@/theme/theme.module.css";
import { isNotFoundGrpcError } from "@/utils/grpcStatus";
import {
  type InvocationNameSearch,
  invocationNameFromSearch,
} from "@/utils/invocationName";
import getColumns from "./Columns";
import {
  cursorFromOperation,
  fetchQueuedOperationsPage,
  type QueuedOperationsCursor,
} from "./utils";

const LIST_QUEUED_OPERATIONS_PAGE_SIZE = 100;

interface Props {
  invocationName: InvocationNameSearch;
  cursor: QueuedOperationsCursor | undefined;
}

const QueuedOperationsTable: React.FC<Props> = ({ invocationName, cursor }) => {
  const navigate = Route.useNavigate();
  const decodedInvocationName = invocationNameFromSearch(invocationName);

  const { data, isLoading, isPlaceholderData, isError, error } = useQuery({
    queryKey: ["listQueuedOperations", invocationName, cursor],
    queryFn: () =>
      fetchQueuedOperationsPage(
        (pageSize) =>
          buildQueueStateClient.listQueuedOperations({
            invocationName: invocationNameFromSearch(invocationName),
            pageSize,
            startAfter: cursor && {
              priority: cursor.priority,
              expectedDuration: cursor.expectedDuration,
              queuedTimestamp: new Date(cursor.queuedTimestamp),
            },
          }),
        cursor,
        LIST_QUEUED_OPERATIONS_PAGE_SIZE,
      ),
    placeholderData: (previousData, _) => previousData,
  });

  if (isError) {
    if (isNotFoundGrpcError(error)) {
      return null;
    }
    return (
      <PortalAlert
        showIcon
        type="error"
        title="Error fetching queued operations"
        description={
          error.message ||
          "Unknown error occurred while fetching data from the server."
        }
      />
    );
  }

  const goToFirstPage = () => {
    navigate({
      search: (prev) => ({ ...prev, cursor: undefined }),
    });
  };

  const goToNextPage = () => {
    const last = data?.operations.at(-1);
    if (!last) return;
    navigate({
      search: (prev) => ({ ...prev, cursor: cursorFromOperation(last) }),
    });
  };

  return (
    <Space orientation="vertical" style={{ width: "100%" }}>
      {data && (
        <Row>
          <Flex style={{ width: "100%" }} justify="flex-end" wrap>
            <WorkersTablePageSelector
              noun="queued operations"
              paginationInfo={data}
              goToFirstPage={goToFirstPage}
              goToNextPage={goToNextPage}
              pageSize={data.operations.length}
            />
          </Flex>
        </Row>
      )}
      <Row>
        <Table
          columns={getColumns(decodedInvocationName)}
          loading={isLoading || isPlaceholderData}
          bordered={true}
          style={{ width: "100%" }}
          dataSource={data?.operations}
          size="small"
          rowClassName={() => themeStyles.compactTable}
          pagination={false}
          rowKey="key"
          locale={{
            emptyText:
              "No operations queued directly in this invocation (that you have access to).",
          }}
        />
      </Row>
    </Space>
  );
};

export default QueuedOperationsTable;
