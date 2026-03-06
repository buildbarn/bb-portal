import { useQuery } from "@tanstack/react-query";
import { Row, Space } from "antd";
import type React from "react";
import { useGrpcClients } from "@/context/GrpcClientsContext";
import type {
  BuildQueueStateClient,
  DeepPartial,
  ListWorkersRequest,
  ListWorkersResponse,
  SizeClassQueueName,
  WorkerState,
} from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import PortalAlert from "@/components/PortalAlert";
import WorkersInfo from "@/components/WorkersInfo";
import WorkersTable from "@/components/WorkersTable";
import { WorkerListStatus, WorkerSearchParams } from "@/routes/scheduler.worker";

const LIST_WORKERS_PAGE_SIZE = 100;

const fetchWorkers = async (
  client: BuildQueueStateClient,
  workerStatusFilter: WorkerListStatus,
  sizeClassQueueName: SizeClassQueueName,
  paginationCursor?: WorkerState["id"],
) => {
  const requestBody: DeepPartial<ListWorkersRequest> = {
    filter: {},
    pageSize: LIST_WORKERS_PAGE_SIZE,
  };

  if (paginationCursor) {
    requestBody.startAfter = { workerId: paginationCursor };
  }

  switch (workerStatusFilter) {
    case WorkerListStatus.ALL:
      requestBody.filter = {
        all: sizeClassQueueName,
      };
      break;
    case WorkerListStatus.EXECUTING:
      requestBody.filter = {
        executing: {
          sizeClassQueueName: sizeClassQueueName,
        },
      };
      break;
  }

  return client.listWorkers(requestBody);
};

const WorkersGrid: React.FC<WorkerSearchParams> = ({
  workerStatusFilter, 
  sizeClassQueueName, 
  cursor
}) => {
  const { buildQueueStateClient } = useGrpcClients();

  const { data, isError, isLoading, error } = useQuery<ListWorkersResponse>({
    queryKey: [
      "listWorkers",
      workerStatusFilter,
      sizeClassQueueName,
      cursor,
    ],
    queryFn: () =>
      fetchWorkers(
        buildQueueStateClient,
        workerStatusFilter,
        sizeClassQueueName,
        cursor,
      ),
    placeholderData: (previousData, _) => previousData,
  });

  if (isError) {
    return (
      <PortalAlert
        showIcon
        type="error"
        message="Error fetching workers"
        description={
          error.message ||
          "Unknown error occurred while fetching data from the server."
        }
      />
    );
  }

  return (
    <Space direction="vertical" size="middle" style={{ display: "flex" }}>
      <Row>
        <WorkersInfo sizeClassQueueName={sizeClassQueueName} />
      </Row>
      <Row>
        <WorkersTable
          workerStatusFilter={workerStatusFilter}
          data={data?.workers.map((value: WorkerState) => {
            // Size class queue is not included in
            // workerState, so we set it manually.
            return {
              ...value,
              currentOperation: value.currentOperation && {
                ...value.currentOperation,
                invocationName: {
                  ...value.currentOperation.invocationName,
                  sizeClassQueueName: sizeClassQueueName,
                  ids: value.currentOperation.invocationName?.ids ?? [],
                },
              },
            };
          })}
          paginationInfo={data?.paginationInfo}
          isLoading={isLoading}
          pageSize={LIST_WORKERS_PAGE_SIZE}
        />
      </Row>
    </Space>
  );
};

export default WorkersGrid;
