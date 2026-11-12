import { useQuery } from "@tanstack/react-query";
import { Row, Space } from "antd";
import type React from "react";
import PortalAlert from "@/components/PortalAlert";
import WorkersInfo from "@/components/WorkersInfo";
import WorkersTable from "@/components/WorkersTable";
import type { WorkerState } from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import type { WorkerSearchParams } from "@/routes/scheduler.worker";
import { useScheduler } from "@/schedulers/context";
import { LIST_WORKERS_PAGE_SIZE, workersOptions } from "@/schedulers/queries";

const WorkersGrid: React.FC<WorkerSearchParams> = ({
  workerStatusFilter,
  sizeClassQueueName,
  cursor,
}) => {
  const { scheduler } = useScheduler();
  const { data, isError, isLoading, error } = useQuery(
    workersOptions(
      scheduler.id,
      workerStatusFilter,
      sizeClassQueueName,
      cursor,
    ),
  );

  if (isError) {
    return (
      <PortalAlert
        showIcon
        type="error"
        title="Error fetching workers"
        description={
          error.message ||
          "Unknown error occurred while fetching data from the server."
        }
      />
    );
  }

  return (
    <Space orientation="vertical" size="middle" style={{ display: "flex" }}>
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
