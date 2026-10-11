import { Link } from "@tanstack/react-router";
import { type TableColumnsType, Typography } from "antd";
import type { ColumnGroupType, ColumnType } from "antd/es/table";
import { InvocationChildrenFilter } from "@/routes/scheduler.invocation";
import { WorkerListStatus } from "@/routes/scheduler.worker";
import {
  rootInvocationNameSearch,
  type SizeClassQueueNameSearch,
} from "@/utils/invocationName";
import { readableDurationFromDates } from "@/utils/time";
import PropertyTagList from "../PropertyTagList";
import type { PlatformQueueTableState } from "./types";

const timeFormat = { precision: 1, smallestUnit: "s" } as const;

const sizeClassQueueNameSearch = (
  record: PlatformQueueTableState,
): SizeClassQueueNameSearch => ({
  platformQueueName: {
    instanceNamePrefix: record.name?.instanceNamePrefix || "",
    platform: {
      properties: record.name?.platform?.properties || [],
    },
  },
  sizeClass: record.sizeClassQueues[0]?.sizeClass,
});

const cellMergingLogic = (value: PlatformQueueTableState) => {
  if (value.isFirstSizeClass) {
    return { rowSpan: value.numberOfSizeClasses };
  }
  return { rowSpan: 0 };
};

const instanceNamePrefixColumn: ColumnType<PlatformQueueTableState> = {
  key: "instanceNamePrefix",
  title: "Instance name prefix",
  onCell: cellMergingLogic,
  render: (_, record) => (
    <Typography.Text>{record.name?.instanceNamePrefix}</Typography.Text>
  ),
};

const platformPropertiesColumn: ColumnType<PlatformQueueTableState> = {
  key: "platformProperties",
  title: "Platform properties",
  onCell: cellMergingLogic,
  render: (_, record) => (
    <PropertyTagList propertyList={record.name?.platform?.properties} />
  ),
};

const sizeClassColumn: ColumnType<PlatformQueueTableState> = {
  key: "sizeClass",
  title: "Size class",
  render: (_, record) => (
    <Typography.Text>{record.sizeClassQueues[0].sizeClass}</Typography.Text>
  ),
};

const queuedOperationsColumn: ColumnType<PlatformQueueTableState> = {
  key: "queuedOperations",
  title: "Queued operations",
  render: (_, record) => (
    <Typography.Text>
      {(record.sizeClassQueues[0].rootInvocation?.queuedOperationsCount
        ?.direct ?? 0) +
        (record.sizeClassQueues[0].rootInvocation?.queuedOperationsCount
          ?.indirect ?? 0)}
    </Typography.Text>
  ),
};

const executingWorkersColumn: ColumnType<PlatformQueueTableState> = {
  key: "executingWorkers",
  title: "Executing",
  render: (_, record) => {
    const sizeClassInfo = record.sizeClassQueues[0];
    let allWorkers = sizeClassInfo.workersCount;
    if (allWorkers === 0) {
      allWorkers = 100;
    }
    const executingWorkers =
      sizeClassInfo.rootInvocation?.executingWorkersCount || 0;
    const percentage = ((executingWorkers / allWorkers) * 100).toFixed(2);
    return (
      <Link
        to="/scheduler/worker"
        search={{
          workerStatusFilter: WorkerListStatus.EXECUTING,
          sizeClassQueueName: sizeClassQueueNameSearch(record),
          cursor: undefined,
        }}
      >
        {executingWorkers} ({percentage}%)
      </Link>
    );
  },
};

const idleWorkersColumn: ColumnType<PlatformQueueTableState> = {
  key: "idleWorkers",
  title: "Idle workers",
  render: (_, record) => (
    <Typography.Text>
      {record.sizeClassQueues[0].rootInvocation?.idleWorkersCount}
    </Typography.Text>
  ),
};

const allWorkersColumn: ColumnType<PlatformQueueTableState> = {
  key: "allWorkers",
  title: "All workers",
  render: (_, record) => (
    <Link
      to="/scheduler/worker"
      search={{
        workerStatusFilter: WorkerListStatus.ALL,
        sizeClassQueueName: sizeClassQueueNameSearch(record),
        cursor: undefined,
      }}
    >
      {record.sizeClassQueues[0].workersCount}
    </Link>
  ),
};

const timeoutColumn: ColumnType<PlatformQueueTableState> = {
  key: "timeout",
  title: "Timeout",
  render: (_, record) => {
    const timeout = record.sizeClassQueues[0].timeout;
    return timeout
      ? readableDurationFromDates(new Date(), timeout, timeFormat)
      : "∞";
  },
};

const invocationLink = (
  record: PlatformQueueTableState,
  childrenFilter: InvocationChildrenFilter,
  count: number | undefined,
) => (
  <Link
    to="/scheduler/invocation"
    search={{
      invocationName: rootInvocationNameSearch(
        sizeClassQueueNameSearch(record),
      ),
      childrenFilter,
      cursor: undefined,
    }}
  >
    {count ?? 0}
  </Link>
);

const rootInvocationColumn: ColumnGroupType<PlatformQueueTableState> = {
  key: "rootInvocation",
  title: "Root invocation",
  children: [
    {
      key: "rootInvocationQueuedOperations",
      title: "Queued operations",
      children: [
        {
          key: "queuedOperationsDirect",
          title: "Direct",
          render: (_, record) =>
            invocationLink(
              record,
              InvocationChildrenFilter.ACTIVE,
              record.sizeClassQueues[0].rootInvocation?.queuedOperationsCount
                ?.direct,
            ),
        },
        {
          key: "queuedOperationsIndirect",
          title: "Indirect",
          render: (_, record) =>
            record.sizeClassQueues[0].rootInvocation?.queuedOperationsCount
              ?.indirect ?? 0,
        },
      ],
    },
    {
      key: "rootInvocationChildren",
      title: "Children",
      children: [
        {
          key: "queuedChildren",
          title: "Queued",
          render: (_, record) =>
            invocationLink(
              record,
              InvocationChildrenFilter.QUEUED,
              record.sizeClassQueues[0].rootInvocation?.queuedChildrenCount,
            ),
        },
        {
          key: "activeChildren",
          title: "Active",
          render: (_, record) =>
            invocationLink(
              record,
              InvocationChildrenFilter.ACTIVE,
              record.sizeClassQueues[0].rootInvocation?.activeChildrenCount,
            ),
        },
        {
          key: "allChildren",
          title: "All",
          render: (_, record) =>
            invocationLink(
              record,
              InvocationChildrenFilter.ALL,
              record.sizeClassQueues[0].rootInvocation?.childrenCount,
            ),
        },
      ],
    },
    {
      key: "rootInvocationWorkers",
      title: "Workers",
      children: [
        executingWorkersColumn,
        {
          key: "idleWorkers",
          title: "Idle",
          render: (_, record) =>
            record.sizeClassQueues[0].rootInvocation?.idleWorkersCount ?? 0,
        },
        {
          key: "idleSynchronizingWorkers",
          title: "Idle synchronizing",
          render: (_, record) => (
            <Link
              to="/scheduler/worker"
              search={{
                workerStatusFilter: WorkerListStatus.IDLE_SYNCHRONIZING,
                sizeClassQueueName: sizeClassQueueNameSearch(record),
                invocationIds: [],
                cursor: undefined,
              }}
            >
              {record.sizeClassQueues[0].rootInvocation
                ?.idleSynchronizingWorkersCount ?? 0}
            </Link>
          ),
        },
      ],
    },
  ],
};

const drainsColumn: ColumnType<PlatformQueueTableState> = {
  key: "drains",
  title: "Drains",
  render: (_, record) => (
    <Link
      to="/scheduler/drain"
      search={{ sizeClassQueueName: sizeClassQueueNameSearch(record) }}
    >
      {record.sizeClassQueues[0].drainsCount ?? 0}
    </Link>
  ),
};

const getColumns = (
  operatorMode: boolean,
): TableColumnsType<PlatformQueueTableState> => {
  if (operatorMode) {
    return [
      instanceNamePrefixColumn,
      platformPropertiesColumn,
      sizeClassColumn,
      timeoutColumn,
      rootInvocationColumn,
      allWorkersColumn,
      drainsColumn,
    ];
  }
  return [
    instanceNamePrefixColumn,
    platformPropertiesColumn,
    sizeClassColumn,
    queuedOperationsColumn,
    executingWorkersColumn,
    idleWorkersColumn,
    allWorkersColumn,
  ];
};

export default getColumns;
