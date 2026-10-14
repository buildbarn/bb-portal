import { queryOptions, useQuery } from "@tanstack/react-query";
import type { MenuProps } from "antd";
import type React from "react";
import { useBulkDrain, useDrainActions } from "@/components/DrainActions";
import { useCanModifyDrainsByPrefix } from "@/components/DrainsTable/hooks";
import type { DrainMenus } from "@/components/WorkersInfo";
import { buildQueueStateClient } from "@/grpc/buildQueueStateClient";
import type { SizeClassQueueNameSearch } from "@/utils/invocationName";
import type { PlatformQueueTableState } from "./types";
import {
  type DrainTarget,
  describeSelector,
  findDrainTarget,
  type QueueSelector,
  selectDrainTargets,
} from "./utils";

export const platformQueuesQuery = queryOptions({
  queryKey: ["listPlatformQueues"],
  queryFn: async (): Promise<PlatformQueueTableState[]> => {
    const queues = await buildQueueStateClient.listPlatformQueues({});
    // Convert PlatoformQueueState to PlatformQueueTableState, making it
    // suitable for the table component. This is done by flattening the
    // sizeClassQueues array into a single element array, and adding
    // additional properties to track the number of size classes and if
    // the current size class is the first size class in the queue.
    return queues.platformQueues.flatMap((queue) => {
      return queue.sizeClassQueues.map((sizeClassQueue, index) => {
        return {
          ...queue,
          sizeClassQueues: [sizeClassQueue],
          numberOfSizeClasses: queue.sizeClassQueues.length,
          isFirstSizeClass: index === 0,
        };
      });
    });
  },
});

export interface BulkDrain {
  rows: PlatformQueueTableState[];
  isAllowed: (instanceNamePrefix: string) => boolean | undefined;
  drainTargets: (targets: DrainTarget[], description: string) => void;
}

export const usePlatformQueuesBulkDrain = (
  rows: PlatformQueueTableState[] | undefined,
  enabled: boolean,
): { bulk: BulkDrain | undefined; contextHolder: React.ReactNode } => {
  const allowedByPrefix = useCanModifyDrainsByPrefix(
    (rows ?? []).map((row) => row.name?.instanceNamePrefix ?? ""),
    enabled,
  );
  const isAllowed = (instanceNamePrefix: string) =>
    allowedByPrefix.get(instanceNamePrefix);
  const { drainTargets, contextHolder } = useBulkDrain(isAllowed);
  return {
    bulk: enabled && rows ? { rows, isAllowed, drainTargets } : undefined,
    contextHolder,
  };
};

export const usePlatformQueues = (enabled = true) =>
  useQuery({ ...platformQueuesQuery, enabled });

export const bulkDrainMenu = (
  bulk: BulkDrain | undefined,
  selector: QueueSelector,
): MenuProps["items"] | undefined => {
  if (!bulk) {
    return undefined;
  }
  const targets = selectDrainTargets(bulk.rows, selector);
  if (
    !targets.some(({ sizeClassQueueName }) =>
      bulk.isAllowed(sizeClassQueueName.platformQueueName.instanceNamePrefix),
    )
  ) {
    return undefined;
  }
  const description = describeSelector(selector);
  return [
    {
      key: "drainAll",
      danger: true,
      label: `Drain all workers in every queue ${description}`,
      onClick: () => bulk.drainTargets(targets, description),
    },
  ];
};

export const useQueueDrainMenus = (
  sizeClassQueueName: SizeClassQueueNameSearch,
  enabled: boolean,
): { drainMenus: DrainMenus | undefined; contextHolder: React.ReactNode } => {
  const { data } = usePlatformQueues(enabled);
  const { bulk, contextHolder: bulkContextHolder } = usePlatformQueuesBulkDrain(
    data,
    enabled,
  );
  const { actions, contextHolder } = useDrainActions(
    sizeClassQueueName,
    enabled,
  );
  const workersCount =
    data && findDrainTarget(data, sizeClassQueueName)?.workersCount;
  const queueItems: NonNullable<MenuProps["items"]> = actions
    ? [
        {
          key: "drainQueue",
          danger: true,
          label: "Drain all workers in this size class queue",
          onClick: () => actions.drainQueue(workersCount),
        },
      ]
    : [];
  const menu = (selector: QueueSelector): MenuProps["items"] => {
    const bulkItems = bulkDrainMenu(bulk, selector) ?? [];
    const items = [
      ...queueItems,
      ...(queueItems.length && bulkItems.length
        ? [{ type: "divider" as const }]
        : []),
      ...bulkItems,
    ];
    return items.length ? items : undefined;
  };
  return {
    drainMenus: enabled
      ? {
          instanceNamePrefix: menu({
            instanceNamePrefix:
              sizeClassQueueName.platformQueueName.instanceNamePrefix,
          }),
          platformProperty: ({ name, value }) =>
            menu({ property: { name, value } }),
        }
      : undefined,
    contextHolder: (
      <>
        {contextHolder}
        {bulkContextHolder}
      </>
    ),
  };
};
