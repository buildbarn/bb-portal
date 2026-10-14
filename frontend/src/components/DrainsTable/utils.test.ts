import { expect, test } from "vitest";
import {
  countMatchingWorkers,
  drainKey,
  modifyDrainsPermissionUrl,
  patternToPropertyList,
  patternToString,
  sizeClassQueueKey,
  workerIdPatternFromFields,
  workerMatchesPattern,
} from "./utils";

test("modifyDrainsPermissionUrl", () => {
  expect(modifyDrainsPermissionUrl("")).toBe(
    "/api/v1/checkPermissions/modifyDrains/",
  );
  expect(modifyDrainsPermissionUrl("foo")).toBe(
    "/api/v1/checkPermissions/modifyDrains/foo",
  );
  expect(modifyDrainsPermissionUrl("foo/bar/baz")).toBe(
    "/api/v1/checkPermissions/modifyDrains/foo/bar/baz",
  );
  expect(modifyDrainsPermissionUrl("my instance/a?b#c")).toBe(
    "/api/v1/checkPermissions/modifyDrains/my%20instance/a%3Fb%23c",
  );
});

test("workerIdPatternFromFields", () => {
  expect(workerIdPatternFromFields([])).toEqual({});
  expect(
    workerIdPatternFromFields([
      { key: "hostname", value: "server123.example.com" },
      { key: "thread", value: "3" },
    ]),
  ).toEqual({ hostname: "server123.example.com", thread: "3" });
});

test("patternToPropertyList", () => {
  expect(patternToPropertyList({})).toEqual([]);
  expect(
    patternToPropertyList({ thread: "3", hostname: "a", datacenter: "x" }),
  ).toEqual([
    { name: "datacenter", value: "x" },
    { name: "hostname", value: "a" },
    { name: "thread", value: "3" },
  ]);
});

test("drainKey", () => {
  expect(drainKey({ thread: "3", hostname: "a" })).toBe(
    drainKey({ hostname: "a", thread: "3" }),
  );
  expect(drainKey({ hostname: "a" })).not.toBe(drainKey({ hostname: "b" }));
});

test("patternToString", () => {
  expect(patternToString({})).toBe("");
  expect(patternToString({ thread: "3", hostname: "a" })).toBe(
    "hostname=a, thread=3",
  );
});

test("workerMatchesPattern", () => {
  const workerId = { hostname: "a", thread: "3", pod: "" };
  expect(workerMatchesPattern(workerId, {})).toBe(true);
  expect(workerMatchesPattern(workerId, { hostname: "a" })).toBe(true);
  expect(workerMatchesPattern(workerId, { hostname: "a", thread: "3" })).toBe(
    true,
  );
  expect(workerMatchesPattern(workerId, { pod: "" })).toBe(true);
  expect(workerMatchesPattern(workerId, { hostname: "b" })).toBe(false);
  expect(workerMatchesPattern(workerId, { hostname: "a", thread: "4" })).toBe(
    false,
  );
  expect(workerMatchesPattern(workerId, { datacenter: "" })).toBe(false);
  expect(workerMatchesPattern(workerId, { constructor: "x" })).toBe(false);
});

test("countMatchingWorkers", () => {
  const workerIds = [
    { hostname: "a", thread: "1" },
    { hostname: "a", thread: "2" },
    { hostname: "b", thread: "1" },
  ];
  expect(countMatchingWorkers(workerIds, {})).toBe(3);
  expect(countMatchingWorkers(workerIds, { hostname: "a" })).toBe(2);
  expect(countMatchingWorkers(workerIds, { thread: "1" })).toBe(2);
  expect(countMatchingWorkers([], { thread: "1" })).toBe(0);
});

test("sizeClassQueueKey", () => {
  const queue = (
    properties: { name: string; value: string }[],
    sizeClass = 0,
  ) => ({
    platformQueueName: { instanceNamePrefix: "foo", platform: { properties } },
    sizeClass,
  });
  const a = { name: "OSFamily", value: "linux" };
  const b = { name: "arch", value: "x86_64" };
  expect(sizeClassQueueKey(queue([a, b]))).toBe(
    sizeClassQueueKey(queue([b, a])),
  );
  expect(sizeClassQueueKey(queue([a]))).not.toBe(
    sizeClassQueueKey(queue([a], 1)),
  );
  expect(sizeClassQueueKey(queue([a]))).not.toBe(sizeClassQueueKey(queue([b])));
});
