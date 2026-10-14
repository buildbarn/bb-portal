import type { PropertyTagListEntry } from "@/components/PropertyTagList/types";
import type { SizeClassQueueNameSearch } from "@/utils/invocationName";

export const modifyDrainsPermissionUrl = (instanceNamePrefix: string): string =>
  `/api/v1/checkPermissions/modifyDrains/${instanceNamePrefix
    .split("/")
    .map(encodeURIComponent)
    .join("/")}`;

export const workerIdPatternFromFields = (
  fields: { key: string; value: string }[],
): Record<string, string> =>
  Object.fromEntries(fields.map(({ key, value }) => [key, value]));

export const patternToPropertyList = (
  pattern: Record<string, string>,
): PropertyTagListEntry[] =>
  Object.entries(pattern)
    .sort(([a], [b]) => (a < b ? -1 : a > b ? 1 : 0))
    .map(([name, value]) => ({ name, value }));

export const drainKey = (pattern: Record<string, string>): string =>
  JSON.stringify(patternToPropertyList(pattern));

export const patternToString = (pattern: Record<string, string>): string =>
  patternToPropertyList(pattern)
    .map(({ name, value }) => `${name}=${value}`)
    .join(", ");

export const workerMatchesPattern = (
  workerId: Record<string, string>,
  pattern: Record<string, string>,
): boolean =>
  Object.entries(pattern).every(
    ([key, value]) => Object.hasOwn(workerId, key) && workerId[key] === value,
  );

export const countMatchingWorkers = (
  workerIds: Record<string, string>[],
  pattern: Record<string, string>,
): number =>
  workerIds.filter((workerId) => workerMatchesPattern(workerId, pattern))
    .length;

export const sizeClassQueueKey = (
  sizeClassQueueName: SizeClassQueueNameSearch,
): string =>
  JSON.stringify([
    sizeClassQueueName.platformQueueName.instanceNamePrefix,
    sizeClassQueueName.platformQueueName.platform.properties
      .map(({ name, value }) => JSON.stringify([name, value]))
      .sort(),
    sizeClassQueueName.sizeClass,
  ]);
