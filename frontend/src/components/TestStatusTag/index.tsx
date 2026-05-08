import {
  CheckCircleFilled,
  CloseCircleFilled,
  InfoCircleFilled,
  MinusCircleFilled,
  QuestionCircleFilled,
  StopOutlined,
} from "@ant-design/icons";
import type React from "react";
import ResultTag, { type TagVariables } from "@/components/ResultTag";

export const ALL_STATUS_VALUES = [
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

export type StatusTuple = typeof ALL_STATUS_VALUES;
export type TestStatusEnum = StatusTuple[number];

interface Props {
  status: TestStatusEnum;
  displayText: boolean;
}

const TEST_RESULT_TAGS: { [key in TestStatusEnum]: TagVariables } = {
  NO_STATUS: {
    icon: <QuestionCircleFilled />,
    text: "No Status",
    color: "default",
  },
  PASSED: {
    icon: <CheckCircleFilled />,
    color: "green",
    text: "Passed",
  },
  FLAKY: {
    icon: <InfoCircleFilled />,
    color: "orange",
    text: "Flaky",
  },
  FAILED: {
    icon: <CloseCircleFilled />,
    color: "red",
    text: "Failed",
  },
  TIMEOUT: {
    icon: <MinusCircleFilled />,
    color: "red",
    text: "Timeout",
  },
  INCOMPLETE: {
    icon: <StopOutlined />,
    color: "blue",
    text: "Incomplete",
  },
  REMOTE_FAILURE: {
    icon: <CloseCircleFilled />,
    color: "red",
    text: "Remote Failure",
  },
  FAILED_TO_BUILD: {
    icon: <QuestionCircleFilled />,
    color: "red",
    text: "Failed to Build",
  },
  TOOL_HALTED_BEFORE_TESTING: {
    icon: <QuestionCircleFilled />,
    color: "blue",
    text: "Status Unknown",
  },
};

const TestStatusTag: React.FC<Props> = ({ status, displayText }) => {
  const tagVars = TEST_RESULT_TAGS[status];

  if (displayText) {
    return <ResultTag tagVars={tagVars} />;
  }

  return <ResultTag tagVars={{ ...tagVars, text: undefined }} />;
};

export default TestStatusTag;
