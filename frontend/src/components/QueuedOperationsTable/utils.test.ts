import { expect, test } from "vitest";
import {
  type OperationState,
  OperationState as OperationStateMessage,
} from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import {
  cursorFromOperation,
  dropRepeatedQueuedOperations,
  fetchQueuedOperationsPage,
  withInvocationName,
} from "./utils";

const operation = (
  name: string,
  queuedTimestamp: Date,
  priority = 0,
): OperationState =>
  OperationStateMessage.fromPartial({
    name,
    priority,
    expectedDuration: { seconds: "5", nanos: 100 },
    queuedTimestamp,
  });

const t1 = new Date("2026-10-11T12:00:00.123Z");
const t2 = new Date("2026-10-11T12:00:00.124Z");

test("cursorFromOperation", () => {
  expect(cursorFromOperation(operation("a", t1, 3))).toEqual({
    priority: 3,
    expectedDuration: { seconds: "5", nanos: 100 },
    queuedTimestamp: "2026-10-11T12:00:00.123Z",
    operationName: "a",
  });
});

test("no cursor", () => {
  const operations = [operation("a", t1), operation("b", t1)];
  expect(dropRepeatedQueuedOperations(operations, undefined)).toEqual({
    operations,
    dropped: 0,
  });
});

test("cursor row at index 0", () => {
  const operations = [operation("a", t1), operation("b", t2)];
  const cursor = cursorFromOperation(operations[0]);
  expect(dropRepeatedQueuedOperations(operations, cursor)).toEqual({
    operations: [operations[1]],
    dropped: 1,
  });
});

test("cursor row at index 2 among ties", () => {
  const operations = [
    operation("a", t1),
    operation("b", t1),
    operation("c", t1),
    operation("d", t1),
    operation("e", t2),
  ];
  const cursor = cursorFromOperation(operations[2]);
  expect(dropRepeatedQueuedOperations(operations, cursor)).toEqual({
    operations: operations.slice(3),
    dropped: 3,
  });
});

test("cursor row absent", () => {
  const operations = [operation("b", t1), operation("c", t1)];
  const cursor = cursorFromOperation(operation("a", t1));
  expect(dropRepeatedQueuedOperations(operations, cursor)).toEqual({
    operations,
    dropped: 0,
  });
});

test("different tuple stops the scan", () => {
  const operations = [operation("b", t1), operation("a", t1, 1)];
  const cursor = cursorFromOperation(operation("a", t1));
  expect(dropRepeatedQueuedOperations(operations, cursor)).toEqual({
    operations,
    dropped: 0,
  });
});

const page = (operations: OperationState[], startIndex: number) => ({
  queuedOperations: operations,
  paginationInfo: { startIndex, totalEntries: 300 },
});

test("page grows past repeated operations", async () => {
  const group = Array.from({ length: 150 }, (_, i) => operation(`op${i}`, t1));
  const requested: number[] = [];
  const result = await fetchQueuedOperationsPage(
    async (pageSize) => {
      requested.push(pageSize);
      return page(group.slice(0, pageSize), 0);
    },
    cursorFromOperation(group[99]),
    100,
  );
  expect(requested).toEqual([100, 200]);
  expect(result.startIndex).toBe(100);
  expect(result.totalEntries).toBe(300);
  expect(result.operations.map((o) => o.key)).toEqual(
    group.slice(100).map((o) => o.name),
  );
});

test("page growth stops when the backend returns less", async () => {
  const group = Array.from({ length: 50 }, (_, i) => operation(`op${i}`, t1));
  const requested: number[] = [];
  const result = await fetchQueuedOperationsPage(
    async (pageSize) => {
      requested.push(pageSize);
      return page(group, 0);
    },
    cursorFromOperation(group[49]),
    50,
  );
  expect(requested).toEqual([50, 100]);
  expect(result.operations).toEqual([]);
  expect(result.startIndex).toBe(50);
});

test("censored operations are keyed by index", async () => {
  const result = await fetchQueuedOperationsPage(
    async () => page([operation("", t1), operation("", t2)], 10),
    undefined,
    100,
  );
  expect(result.operations.map((o) => o.key)).toEqual(["10", "11"]);
});

test("withInvocationName", () => {
  const invocationName = { sizeClassQueueName: undefined, ids: [] };
  const original = operation("a", t1);
  const result = withInvocationName(original, invocationName);
  expect(result.invocationName).toBe(invocationName);
  expect(original.invocationName).toBeUndefined();
});
