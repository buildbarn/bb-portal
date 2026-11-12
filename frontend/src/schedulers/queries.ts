import { queryOptions } from "@tanstack/react-query";
import { getBuildQueueStateClient } from "@/grpc/buildQueueStateClient";
import type {
  DeepPartial,
  ListWorkersRequest,
  SizeClassQueueName,
  WorkerState,
} from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import type { WorkerListStatus } from "@/routes/scheduler.worker";

export const LIST_WORKERS_PAGE_SIZE = 100;

export const platformQueuesOptions = (schedulerId: string) =>
  queryOptions({
    queryKey: ["listPlatformQueues", schedulerId],
    queryFn: ({ signal }) =>
      getBuildQueueStateClient(schedulerId).listPlatformQueues({}, { signal }),
  });

export const workersOptions = (
  schedulerId: string,
  workerStatusFilter: WorkerListStatus,
  sizeClassQueueName: SizeClassQueueName,
  cursor?: WorkerState["id"],
) =>
  queryOptions({
    queryKey: [
      "listWorkers",
      schedulerId,
      workerStatusFilter,
      sizeClassQueueName,
      cursor,
    ],
    meta: { schedulerId },
    queryFn: ({ signal }) => {
      const request: DeepPartial<ListWorkersRequest> = {
        filter: {},
        pageSize: LIST_WORKERS_PAGE_SIZE,
      };
      if (cursor) request.startAfter = { workerId: cursor };
      switch (workerStatusFilter) {
        case "all":
          request.filter = { all: sizeClassQueueName };
          break;
        case "executing":
          request.filter = { executing: { sizeClassQueueName } };
          break;
      }
      return getBuildQueueStateClient(schedulerId).listWorkers(request, {
        signal,
      });
    },
    // Preserve paging/filtering rows within a scheduler, never across schedulers.
    placeholderData: (previousData, previousQuery) =>
      previousQuery?.meta?.schedulerId === schedulerId
        ? previousData
        : undefined,
  });

export const schedulerStatisticsOptions = (schedulerId: string) =>
  queryOptions({
    queryKey: ["listOperations", schedulerId],
    queryFn: ({ signal }) =>
      getBuildQueueStateClient(schedulerId).listOperations({}, { signal }),
  });

export const operationDetailsOptions = (
  schedulerId: string,
  operationID: string,
) =>
  queryOptions({
    queryKey: ["operationDetails", schedulerId, operationID],
    queryFn: ({ signal }) =>
      getBuildQueueStateClient(schedulerId).getOperation(
        { operationName: operationID },
        { signal },
      ),
    staleTime: Number.POSITIVE_INFINITY,
    refetchOnMount: "always",
  });

export const killPermissionOptions = (
  schedulerId: string,
  operationID: string,
) =>
  queryOptions({
    queryKey: ["killOperationsButtonState", schedulerId, operationID],
    queryFn: async ({ signal }): Promise<boolean> => {
      const response = await fetch(
        `/api/v1/schedulers/${encodeURIComponent(schedulerId)}/checkPermissions/killOperation/${encodeURIComponent(operationID)}`,
        { signal },
      );
      return (await response.json()).allowed;
    },
  });
