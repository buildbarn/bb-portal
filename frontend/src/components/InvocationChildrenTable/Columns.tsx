import { Link } from "@tanstack/react-router";
import type { TableColumnsType } from "antd";
import CodeText from "@/components/CodeText";
import type { InvocationChildState } from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import type { Any } from "@/lib/grpc-client/google/protobuf/any";
import { InvocationChildrenFilter } from "@/routes/scheduler.invocation";
import { WorkerListStatus } from "@/routes/scheduler.worker";
import {
  childInvocationNameSearch,
  type InvocationNameSearch,
  invocationIdToString,
} from "@/utils/invocationName";

const emptyId: Any = { typeUrl: "", value: new Uint8Array() };

const getColumns = (
  parent: InvocationNameSearch,
): TableColumnsType<InvocationChildState> => {
  const childLink = (
    record: InvocationChildState,
    childrenFilter: InvocationChildrenFilter,
    children: React.ReactNode,
  ) => (
    <Link
      to="/scheduler/invocation"
      search={{
        invocationName: childInvocationNameSearch(parent, record.id ?? emptyId),
        childrenFilter,
        cursor: undefined,
      }}
    >
      {children}
    </Link>
  );

  const workersLink = (
    record: InvocationChildState,
    workerStatusFilter: WorkerListStatus,
    count: number | undefined,
  ) => (
    <Link
      to="/scheduler/worker"
      search={{
        workerStatusFilter,
        sizeClassQueueName: parent.sizeClassQueueName,
        invocationIds: childInvocationNameSearch(parent, record.id ?? emptyId)
          .ids,
        cursor: undefined,
      }}
    >
      {count ?? 0}
    </Link>
  );

  return [
    {
      key: "invocationId",
      title: "Invocation ID",
      render: (_, record) =>
        childLink(
          record,
          InvocationChildrenFilter.ACTIVE,
          <CodeText style={{ wordBreak: "break-all" }}>
            {invocationIdToString(record.id ?? emptyId)}
          </CodeText>,
        ),
    },
    {
      key: "queuedOperations",
      title: "Queued operations",
      children: [
        {
          key: "queuedOperationsDirect",
          title: "Direct",
          render: (_, record) =>
            childLink(
              record,
              InvocationChildrenFilter.ACTIVE,
              record.state?.queuedOperationsCount?.direct ?? 0,
            ),
        },
        {
          key: "queuedOperationsIndirect",
          title: "Indirect",
          render: (_, record) =>
            record.state?.queuedOperationsCount?.indirect ?? 0,
        },
      ],
    },
    {
      key: "children",
      title: "Children",
      children: [
        {
          key: "queuedChildren",
          title: "Queued",
          render: (_, record) =>
            childLink(
              record,
              InvocationChildrenFilter.QUEUED,
              record.state?.queuedChildrenCount ?? 0,
            ),
        },
        {
          key: "activeChildren",
          title: "Active",
          render: (_, record) =>
            childLink(
              record,
              InvocationChildrenFilter.ACTIVE,
              record.state?.activeChildrenCount ?? 0,
            ),
        },
        {
          key: "allChildren",
          title: "All",
          render: (_, record) =>
            childLink(
              record,
              InvocationChildrenFilter.ALL,
              record.state?.childrenCount ?? 0,
            ),
        },
      ],
    },
    {
      key: "workers",
      title: "Workers",
      children: [
        {
          key: "executingWorkers",
          title: "Executing",
          render: (_, record) =>
            workersLink(
              record,
              WorkerListStatus.EXECUTING,
              record.state?.executingWorkersCount,
            ),
        },
        {
          key: "idleWorkers",
          title: "Idle",
          render: (_, record) => record.state?.idleWorkersCount ?? 0,
        },
        {
          key: "idleSynchronizingWorkers",
          title: "Idle synchronizing",
          render: (_, record) =>
            workersLink(
              record,
              WorkerListStatus.IDLE_SYNCHRONIZING,
              record.state?.idleSynchronizingWorkersCount,
            ),
        },
      ],
    },
  ];
};

export default getColumns;
