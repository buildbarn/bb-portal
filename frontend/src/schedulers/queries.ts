import { queryOptions } from "@tanstack/react-query";
import {
  getBuildQueueStateClient,
  SCHEDULER_SELECTION_HEADER,
} from "@/grpc/buildQueueStateClient";
import type {
  DeepPartial,
  ListWorkersRequest,
  SizeClassQueueName,
  WorkerState,
} from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import type { WorkerListStatus } from "@/routes/scheduler.worker";

export const LIST_WORKERS_PAGE_SIZE = 100;

// Every key includes the scheduler's instance name prefix: operation IDs and
// worker IDs are only unique within one scheduler.

export const platformQueuesOptions = (instanceNamePrefix: string) =>
  queryOptions({
    queryKey: ["listPlatformQueues", instanceNamePrefix],
    queryFn: ({ signal }) =>
      getBuildQueueStateClient(instanceNamePrefix).listPlatformQueues(
        {},
        { signal },
      ),
  });

export const workersOptions = (
  instanceNamePrefix: string,
  workerStatusFilter: WorkerListStatus,
  sizeClassQueueName: SizeClassQueueName,
  cursor?: WorkerState["id"],
) =>
  queryOptions({
    queryKey: [
      "listWorkers",
      instanceNamePrefix,
      workerStatusFilter,
      sizeClassQueueName,
      cursor,
    ],
    meta: { instanceNamePrefix },
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
      return getBuildQueueStateClient(instanceNamePrefix).listWorkers(request, {
        signal,
      });
    },
    // Preserve paging/filtering rows within a scheduler, never across schedulers.
    placeholderData: (previousData, previousQuery) =>
      previousQuery?.meta?.instanceNamePrefix === instanceNamePrefix
        ? previousData
        : undefined,
  });

export const schedulerStatisticsOptions = (instanceNamePrefix: string) =>
  queryOptions({
    queryKey: ["listOperations", instanceNamePrefix],
    queryFn: ({ signal }) =>
      getBuildQueueStateClient(instanceNamePrefix).listOperations(
        {},
        { signal },
      ),
  });

export const operationDetailsOptions = (
  instanceNamePrefix: string,
  operationID: string,
) =>
  queryOptions({
    queryKey: ["operationDetails", instanceNamePrefix, operationID],
    queryFn: ({ signal }) =>
      getBuildQueueStateClient(instanceNamePrefix).getOperation(
        { operationName: operationID },
        { signal },
      ),
    staleTime: Number.POSITIVE_INFINITY,
    refetchOnMount: "always",
  });

export const killPermissionOptions = (
  instanceNamePrefix: string,
  operationID: string,
) =>
  queryOptions({
    queryKey: ["killOperationsButtonState", instanceNamePrefix, operationID],
    queryFn: async ({ signal }): Promise<boolean> => {
      const response = await fetch(
        `/api/v1/checkPermissions/killOperation/${encodeURIComponent(operationID)}`,
        {
          signal,
          headers: { [SCHEDULER_SELECTION_HEADER]: instanceNamePrefix },
        },
      );
      return (await response.json()).allowed;
    },
  });
