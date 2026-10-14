import { useSyncExternalStore } from "react";
import { z } from "zod";
import { sizeClassQueueKey } from "@/components/DrainsTable/utils";
import type { SizeClassQueueNameSearch } from "@/utils/invocationName";

export const DRAIN_PATTERN_KEY_PREFIX = "buildbarn_drain_pattern:";

export type DrainPatternAction =
  | { type: "toggle"; key: string; value: string }
  | { type: "set"; key: string; value: string }
  | { type: "remove"; key: string }
  | { type: "clear" };

const withoutKey = (
  pattern: Record<string, string>,
  key: string,
): Record<string, string> =>
  Object.fromEntries(Object.entries(pattern).filter(([k]) => k !== key));

export const drainPatternReducer = (
  pattern: Record<string, string>,
  action: DrainPatternAction,
): Record<string, string> => {
  switch (action.type) {
    case "toggle":
      return Object.hasOwn(pattern, action.key) &&
        pattern[action.key] === action.value
        ? withoutKey(pattern, action.key)
        : { ...pattern, [action.key]: action.value };
    case "set":
      return { ...pattern, [action.key]: action.value };
    case "remove":
      return withoutKey(pattern, action.key);
    case "clear":
      return {};
  }
};

export const replacedValue = (
  pattern: Record<string, string>,
  action: DrainPatternAction,
): string | undefined =>
  (action.type === "toggle" || action.type === "set") &&
  Object.hasOwn(pattern, action.key) &&
  pattern[action.key] !== action.value
    ? pattern[action.key]
    : undefined;

export const drainPatternStorageKey = (
  sizeClassQueueName: SizeClassQueueNameSearch,
): string => DRAIN_PATTERN_KEY_PREFIX + sizeClassQueueKey(sizeClassQueueName);

const PatternSchema = z.record(z.string(), z.string());

export const parseDrainPattern = (
  raw: string | null,
): Record<string, string> => {
  if (!raw) {
    return {};
  }
  try {
    return PatternSchema.safeParse(JSON.parse(raw)).data ?? {};
  } catch {
    return {};
  }
};

const listeners = new Set<() => void>();

const subscribe = (listener: () => void) => {
  listeners.add(listener);
  return () => {
    listeners.delete(listener);
  };
};

export const getDrainPattern = (
  sizeClassQueueName: SizeClassQueueNameSearch,
): Record<string, string> =>
  parseDrainPattern(
    sessionStorage.getItem(drainPatternStorageKey(sizeClassQueueName)),
  );

export const dispatchDrainPattern = (
  sizeClassQueueName: SizeClassQueueNameSearch,
  action: DrainPatternAction,
) => {
  const key = drainPatternStorageKey(sizeClassQueueName);
  const next = drainPatternReducer(getDrainPattern(sizeClassQueueName), action);
  if (Object.keys(next).length) {
    sessionStorage.setItem(key, JSON.stringify(next));
  } else {
    sessionStorage.removeItem(key);
  }
  for (const listener of listeners) {
    listener();
  }
};

export const useDrainPattern = (
  sizeClassQueueName: SizeClassQueueNameSearch,
): Record<string, string> => {
  const key = drainPatternStorageKey(sizeClassQueueName);
  return parseDrainPattern(
    useSyncExternalStore(
      subscribe,
      () => sessionStorage.getItem(key),
      () => null,
    ),
  );
};
