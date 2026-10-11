import { z } from "zod";
import type {
  InvocationName,
  ListQueuedOperationsResponse,
  OperationState,
} from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";

export const QueuedOperationsCursorSchema = z.object({
  priority: z.number().int(),
  expectedDuration: z.object({
    seconds: z.string(),
    nanos: z.number().int(),
  }),
  queuedTimestamp: z.string(),
  operationName: z.string(),
});

export type QueuedOperationsCursor = z.infer<
  typeof QueuedOperationsCursorSchema
>;

export const cursorFromOperation = (
  operation: OperationState,
): QueuedOperationsCursor => ({
  priority: operation.priority,
  expectedDuration: {
    seconds: operation.expectedDuration?.seconds ?? "0",
    nanos: operation.expectedDuration?.nanos ?? 0,
  },
  queuedTimestamp: (operation.queuedTimestamp ?? new Date(0)).toISOString(),
  operationName: operation.name,
});

const hasCursorPosition = (
  operation: OperationState,
  cursor: QueuedOperationsCursor,
): boolean =>
  operation.priority === cursor.priority &&
  (operation.expectedDuration?.seconds ?? "0") ===
    cursor.expectedDuration.seconds &&
  (operation.expectedDuration?.nanos ?? 0) === cursor.expectedDuration.nanos &&
  (operation.queuedTimestamp?.getTime() ?? 0) ===
    new Date(cursor.queuedTimestamp).getTime();

// Timestamps are truncated to milliseconds by the client, so the page
// following a cursor may start with operations already shown.
export const dropRepeatedQueuedOperations = (
  operations: OperationState[],
  cursor: QueuedOperationsCursor | undefined,
): { operations: OperationState[]; dropped: number } => {
  if (cursor) {
    for (let i = 0; i < operations.length; i++) {
      if (!hasCursorPosition(operations[i], cursor)) {
        break;
      }
      if (operations[i].name === cursor.operationName) {
        return { operations: operations.slice(i + 1), dropped: i + 1 };
      }
    }
  }
  return { operations, dropped: 0 };
};

export interface QueuedOperationsPage {
  operations: (OperationState & { key: string })[];
  startIndex: number;
  totalEntries: number;
}

// A page may consist only of operations already shown when more of them
// share the cursor's position than fit in it, so grow the page until it
// makes progress or the backend returns less than requested.
export const fetchQueuedOperationsPage = async (
  list: (pageSize: number) => Promise<ListQueuedOperationsResponse>,
  cursor: QueuedOperationsCursor | undefined,
  pageSize: number,
): Promise<QueuedOperationsPage> => {
  for (let requestSize = pageSize; ; requestSize *= 2) {
    const response = await list(requestSize);
    const { operations, dropped } = dropRepeatedQueuedOperations(
      response.queuedOperations,
      cursor,
    );
    if (
      operations.length > 0 ||
      response.queuedOperations.length < requestSize
    ) {
      const startIndex = (response.paginationInfo?.startIndex ?? 0) + dropped;
      return {
        operations: operations.slice(0, pageSize).map((operation, i) => ({
          ...operation,
          key: operation.name || `${startIndex + i}`,
        })),
        startIndex,
        totalEntries: response.paginationInfo?.totalEntries ?? 0,
      };
    }
  }
};

export const withInvocationName = (
  operation: OperationState,
  invocationName: InvocationName,
): OperationState => ({ ...operation, invocationName });
