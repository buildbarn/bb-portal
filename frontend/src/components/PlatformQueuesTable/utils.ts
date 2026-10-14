import { sizeClassQueueKey } from "@/components/DrainsTable/utils";
import type { PlatformQueueState } from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import type { SizeClassQueueNameSearch } from "@/utils/invocationName";

export interface DrainTarget {
  sizeClassQueueName: SizeClassQueueNameSearch;
  workersCount: number;
  drainsCount: number;
}

export type QueueSelector =
  | { property: { name: string; value: string } }
  | { instanceNamePrefix: string };

const platformQueueMatches = (
  queue: PlatformQueueState,
  selector: QueueSelector,
): boolean => {
  if ("instanceNamePrefix" in selector) {
    return (
      (queue.name?.instanceNamePrefix ?? "") === selector.instanceNamePrefix
    );
  }
  const { name, value } = selector.property;
  return (queue.name?.platform?.properties ?? []).some(
    (property) => property.name === name && property.value === value,
  );
};

export const selectDrainTargets = (
  queues: PlatformQueueState[],
  selector: QueueSelector,
): DrainTarget[] =>
  queues
    .filter((queue) => platformQueueMatches(queue, selector))
    .flatMap((queue) =>
      queue.sizeClassQueues.map((sizeClassQueue) => ({
        sizeClassQueueName: {
          platformQueueName: {
            instanceNamePrefix: queue.name?.instanceNamePrefix ?? "",
            platform: { properties: queue.name?.platform?.properties ?? [] },
          },
          sizeClass: sizeClassQueue.sizeClass,
        },
        workersCount: sizeClassQueue.workersCount,
        drainsCount: sizeClassQueue.drainsCount,
      })),
    );

export const describeSelector = (selector: QueueSelector): string =>
  "instanceNamePrefix" in selector
    ? `under prefix "${selector.instanceNamePrefix}"`
    : `with ${selector.property.name}=${selector.property.value}`;

export const findDrainTarget = (
  queues: PlatformQueueState[],
  sizeClassQueueName: SizeClassQueueNameSearch,
): DrainTarget | undefined => {
  const key = sizeClassQueueKey(sizeClassQueueName);
  return selectDrainTargets(queues, {
    instanceNamePrefix: sizeClassQueueName.platformQueueName.instanceNamePrefix,
  }).find((target) => sizeClassQueueKey(target.sizeClassQueueName) === key);
};
