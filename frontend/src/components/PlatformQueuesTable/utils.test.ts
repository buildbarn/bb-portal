import { expect, test } from "vitest";
import type {
  PlatformQueueState,
  SizeClassQueueState,
} from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import { describeSelector, findDrainTarget, selectDrainTargets } from "./utils";

const sizeClassQueue = (
  sizeClass: number,
  workersCount: number,
  drainsCount = 0,
): SizeClassQueueState => ({
  sizeClass,
  timeout: undefined,
  workersCount,
  drainsCount,
  rootInvocation: undefined,
});

const platformQueue = (
  instanceNamePrefix: string,
  properties: [string, string][],
  sizeClassQueues: SizeClassQueueState[],
): PlatformQueueState => ({
  name: {
    instanceNamePrefix,
    platform: {
      properties: properties.map(([name, value]) => ({ name, value })),
    },
  },
  sizeClassQueues,
});

const queues = [
  platformQueue(
    "a",
    [
      ["OSFamily", "linux"],
      ["arch", "x86_64"],
    ],
    [sizeClassQueue(1, 10), sizeClassQueue(8, 2, 1)],
  ),
  platformQueue(
    "b",
    [
      ["arch", "x86_64"],
      ["OSFamily", "linux"],
    ],
    [sizeClassQueue(1, 5)],
  ),
  platformQueue(
    "a",
    [
      ["OSFamily", "linux"],
      ["arch", "aarch64"],
    ],
    [sizeClassQueue(1, 3)],
  ),
  platformQueue("", [["OSFamily", "windows"]], [sizeClassQueue(0, 1)]),
];

const summary = (selector: Parameters<typeof selectDrainTargets>[1]) =>
  selectDrainTargets(queues, selector).map((target) => [
    target.sizeClassQueueName.platformQueueName.instanceNamePrefix,
    target.sizeClassQueueName.sizeClass,
    target.workersCount,
    target.drainsCount,
  ]);

test("selects by exact property name and value regardless of order", () => {
  expect(summary({ property: { name: "arch", value: "x86_64" } })).toEqual([
    ["a", 1, 10, 0],
    ["a", 8, 2, 1],
    ["b", 1, 5, 0],
  ]);
  expect(summary({ property: { name: "OSFamily", value: "linux" } })).toEqual([
    ["a", 1, 10, 0],
    ["a", 8, 2, 1],
    ["b", 1, 5, 0],
    ["a", 1, 3, 0],
  ]);
  expect(summary({ property: { name: "arch", value: "x86" } })).toEqual([]);
  expect(summary({ property: { name: "Arch", value: "x86_64" } })).toEqual([]);
  expect(summary({ property: { name: "x86_64", value: "arch" } })).toEqual([]);
});

test("selects by exact instance name prefix", () => {
  expect(summary({ instanceNamePrefix: "a" })).toEqual([
    ["a", 1, 10, 0],
    ["a", 8, 2, 1],
    ["a", 1, 3, 0],
  ]);
  expect(summary({ instanceNamePrefix: "" })).toEqual([["", 0, 1, 0]]);
  expect(summary({ instanceNamePrefix: "a/b" })).toEqual([]);
});

test("targets carry the full size class queue name", () => {
  expect(
    selectDrainTargets(queues, { instanceNamePrefix: "b" })[0]
      .sizeClassQueueName,
  ).toEqual({
    platformQueueName: {
      instanceNamePrefix: "b",
      platform: {
        properties: [
          { name: "arch", value: "x86_64" },
          { name: "OSFamily", value: "linux" },
        ],
      },
    },
    sizeClass: 1,
  });
});

test("selects from flattened table rows", () => {
  const rows = queues.flatMap((queue) =>
    queue.sizeClassQueues.map((q) => ({ ...queue, sizeClassQueues: [q] })),
  );
  expect(selectDrainTargets(rows, { instanceNamePrefix: "a" })).toEqual(
    selectDrainTargets(queues, { instanceNamePrefix: "a" }),
  );
});

test("describeSelector", () => {
  expect(
    describeSelector({ property: { name: "arch", value: "x86_64" } }),
  ).toBe("with arch=x86_64");
  expect(describeSelector({ instanceNamePrefix: "a" })).toBe(
    'under prefix "a"',
  );
});

test("findDrainTarget", () => {
  const sizeClassQueueName = (
    instanceNamePrefix: string,
    properties: [string, string][],
    sizeClass: number,
  ) => ({
    platformQueueName: {
      instanceNamePrefix,
      platform: {
        properties: properties.map(([name, value]) => ({ name, value })),
      },
    },
    sizeClass,
  });
  expect(
    findDrainTarget(
      queues,
      sizeClassQueueName(
        "a",
        [
          ["arch", "x86_64"],
          ["OSFamily", "linux"],
        ],
        8,
      ),
    ),
  ).toMatchObject({ workersCount: 2, drainsCount: 1 });
  expect(
    findDrainTarget(queues, sizeClassQueueName("b", [["arch", "x86_64"]], 1)),
  ).toBeUndefined();
  expect(
    findDrainTarget(
      queues,
      sizeClassQueueName(
        "a",
        [
          ["OSFamily", "linux"],
          ["arch", "x86_64"],
        ],
        4,
      ),
    ),
  ).toBeUndefined();
});
