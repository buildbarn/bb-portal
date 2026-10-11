import type { PropertyTagListEntry } from "@/components/PropertyTagList/types";

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
