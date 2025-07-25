import type { SizeClassQueueName } from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import { ListWorkerFilterType } from "@/types/ListWorkerFilterType";
import { type TableColumnsType, Typography } from "antd";
import type { ColumnType } from "antd/lib/table";
import Link from "next/link";
import PropertyTagList from "../PropertyTagList";
import type { PlatformQueueTableState } from "./types";

const getWorkerPageUrlObject = (
  record: PlatformQueueTableState,
  listWorkerFilterType: ListWorkerFilterType,
) => {
  const sizeClassQueueName: SizeClassQueueName = {
    platformQueueName: record.name,
    sizeClass: record.sizeClassQueues[0].sizeClass,
  };

  return {
    pathname: "/scheduler/workers",
    query: {
      listWorkerFilterType: listWorkerFilterType,
      sizeClassQueueName: JSON.stringify(sizeClassQueueName),
    },
  };
};

const cellMergingLogic = (value: PlatformQueueTableState, index?: number) => {
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
    let allWorkers = record.sizeClassQueues[0].workersCount;
    if (allWorkers === 0) {
      allWorkers = 100;
    }
    const executingWorkers =
      record.sizeClassQueues[0].rootInvocation?.executingWorkersCount || 0;
    const percentage = ((executingWorkers / allWorkers) * 100).toFixed(2);
    return (
      <Link
        href={getWorkerPageUrlObject(record, ListWorkerFilterType.EXECUTING)}
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
    <Link href={getWorkerPageUrlObject(record, ListWorkerFilterType.ALL)}>
      {record.sizeClassQueues[0].workersCount}
    </Link>
  ),
};

const getColumns = (): TableColumnsType<PlatformQueueTableState> => {
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
