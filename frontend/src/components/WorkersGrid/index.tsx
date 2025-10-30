import { useGrpcClients } from "@/context/GrpcClientsContext";
import type {
  BuildQueueStateClient,
  DeepPartial,
  ListWorkersRequest,
  SizeClassQueueName,
  WorkerState,
} from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import { ListWorkerFilterType } from "@/types/ListWorkerFilterType";
import { useQuery } from "@tanstack/react-query";
import { Row, Space, Typography } from "antd";
import type React from "react";
import PortalAlert from "../PortalAlert";
import WorkersInfo from "../WorkersInfo";
import WorkersTable from "../WorkersTable";

const LIST_WORKERS_PAGE_SIZE = 100;

interface Props {
  listWorkerFilterType: ListWorkerFilterType;
  sizeClassQueueName: SizeClassQueueName;
  paginationCursor?: WorkerState["id"];
}

const fetchWorkers = async (
  client: BuildQueueStateClient,
  listWorkerFilterType: ListWorkerFilterType,
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

  switch (listWorkerFilterType) {
    case ListWorkerFilterType.ALL:
      requestBody.filter = {
        all: sizeClassQueueName,
      };
      break;
    case ListWorkerFilterType.EXECUTING:
      requestBody.filter = {
        executing: {
          sizeClassQueueName: sizeClassQueueName,
        },
      };
      break;
  }

  return client.listWorkers(requestBody);
};

const WorkersGrid: React.FC<Props> = ({
  listWorkerFilterType,
  sizeClassQueueName,
  paginationCursor,
}) => {
  const { buildQueueStateClient } = useGrpcClients();

  const { data, isError, isLoading, error } = useQuery({
    queryKey: [
      "listWorkers",
      listWorkerFilterType,
      sizeClassQueueName,
      paginationCursor,
    ],
    queryFn: () =>
      fetchWorkers(
        buildQueueStateClient,
        listWorkerFilterType,
        sizeClassQueueName,
        paginationCursor,
      ),
    placeholderData: (previousData, _) => previousData,
  });

  if (isError) {
    return (
      <PortalAlert
        className="error"
        message={
          <>
            <Typography.Text>
              There was a problem communicating with the backend server:
            </Typography.Text>
            <pre>{String(error)}</pre>
          </>
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
          listWorkerFilterType={listWorkerFilterType}
          data={data?.workers.map((value: WorkerState) => {
            // Size class queue is not included in
            // workerState, so we set it manually.
            return {
              ...value,
              currentOperation: value.currentOperation && {
                ...value.currentOperation,
                invocationName: {
                  ...value.currentOperation?.invocationName,
                  sizeClassQueueName: sizeClassQueueName,
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
