const ALL_STATUS_VALUES = [
  "NO_STATUS",
  "PASSED",
  "FLAKY",
  "TIMEOUT",
  "FAILED",
  "INCOMPLETE",
  "REMOTE_FAILURE",
  "FAILED_TO_BUILD",
  "TOOL_HALTED_BEFORE_TESTING",
] as const;

type StatusTuple = typeof ALL_STATUS_VALUES;
export type TestStatusEnum = StatusTuple[number];

const STATUS_SET = new Set<string>(ALL_STATUS_VALUES);

export const testStatusEnumFromString = (
  status: string | null | undefined,
): TestStatusEnum => {
  const isValidStatus = status && STATUS_SET.has(status);
  return isValidStatus ? (status as TestStatusEnum) : "NO_STATUS";
};

export const TEST_STATUS_FILTERS = [
  {
    text: "No Status",
    value: "NO_STATUS",
  },
  {
    text: "Passed",
    value: "PASSED",
  },
  {
    text: "Flaky",
    value: "FLAKY",
  },
  {
    text: "Timeout",
    value: "TIMEOUT",
  },
  {
    text: "Failed",
    value: "FAILED",
  },
  {
    text: "Incomplete",
    value: "INCOMPLETE",
  },
  {
    text: "Remote Failure",
    value: "REMOTE_FAILURE",
  },
  {
    text: "Failed to Build",
    value: "FAILED_TO_BUILD",
  },
  {
    text: "Tool Halted Before Testing",
    value: "TOOL_HALTED_BEFORE_TESTING",
  },
];
