import { afterEach, expect, test } from "vitest";
import {
  DRAIN_PATTERN_KEY_PREFIX,
  dispatchDrainPattern,
  drainPatternReducer,
  drainPatternStorageKey,
  getDrainPattern,
  parseDrainPattern,
  replacedValue,
} from "./drainPattern";
import type { SizeClassQueueNameSearch } from "./invocationName";

const queue = (sizeClass: number): SizeClassQueueNameSearch => ({
  platformQueueName: {
    instanceNamePrefix: "foo",
    platform: {
      properties: [
        { name: "OSFamily", value: "linux" },
        { name: "container-image", value: "ubuntu" },
      ],
    },
  },
  sizeClass,
});

afterEach(() => {
  sessionStorage.clear();
});

test("toggle adds, replaces and removes", () => {
  const toggle = (key: string, value: string) =>
    ({ type: "toggle", key, value }) as const;
  expect(drainPatternReducer({}, toggle("hostname", "a"))).toEqual({
    hostname: "a",
  });
  expect(
    drainPatternReducer(
      { hostname: "a", thread: "1" },
      toggle("hostname", "b"),
    ),
  ).toEqual({ hostname: "b", thread: "1" });
  expect(
    drainPatternReducer(
      { hostname: "a", thread: "1" },
      toggle("hostname", "a"),
    ),
  ).toEqual({ thread: "1" });
  expect(drainPatternReducer({ pod: "" }, toggle("pod", ""))).toEqual({});
});

test("set, remove and clear", () => {
  expect(
    drainPatternReducer(
      { hostname: "a" },
      { type: "set", key: "hostname", value: "a" },
    ),
  ).toEqual({ hostname: "a" });
  expect(
    drainPatternReducer(
      { hostname: "a" },
      { type: "set", key: "thread", value: "1" },
    ),
  ).toEqual({ hostname: "a", thread: "1" });
  expect(
    drainPatternReducer(
      { hostname: "a", thread: "1" },
      { type: "remove", key: "hostname" },
    ),
  ).toEqual({ thread: "1" });
  expect(
    drainPatternReducer({ hostname: "a" }, { type: "remove", key: "pod" }),
  ).toEqual({ hostname: "a" });
  expect(drainPatternReducer({ hostname: "a" }, { type: "clear" })).toEqual({});
});

test("replacedValue", () => {
  const pattern = { hostname: "a" };
  expect(
    replacedValue(pattern, { type: "toggle", key: "hostname", value: "b" }),
  ).toBe("a");
  expect(
    replacedValue(pattern, { type: "set", key: "hostname", value: "b" }),
  ).toBe("a");
  expect(
    replacedValue(pattern, { type: "toggle", key: "hostname", value: "a" }),
  ).toBeUndefined();
  expect(
    replacedValue(pattern, { type: "set", key: "thread", value: "1" }),
  ).toBeUndefined();
  expect(
    replacedValue(pattern, { type: "remove", key: "hostname" }),
  ).toBeUndefined();
});

test("parseDrainPattern rejects malformed state", () => {
  expect(parseDrainPattern(null)).toEqual({});
  expect(parseDrainPattern("not json")).toEqual({});
  expect(parseDrainPattern('{"hostname":1}')).toEqual({});
  expect(parseDrainPattern('{"hostname":"a"}')).toEqual({ hostname: "a" });
});

test("patterns persist per size class queue", () => {
  expect(drainPatternStorageKey(queue(1))).not.toBe(
    drainPatternStorageKey(queue(2)),
  );
  expect(drainPatternStorageKey(queue(1))).toMatch(
    new RegExp(`^${DRAIN_PATTERN_KEY_PREFIX}`),
  );
  dispatchDrainPattern(queue(1), {
    type: "toggle",
    key: "hostname",
    value: "a",
  });
  expect(getDrainPattern(queue(1))).toEqual({ hostname: "a" });
  expect(getDrainPattern(queue(2))).toEqual({});
  dispatchDrainPattern(queue(1), { type: "clear" });
  expect(sessionStorage.getItem(drainPatternStorageKey(queue(1)))).toBeNull();
});
