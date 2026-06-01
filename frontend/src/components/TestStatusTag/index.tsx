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
import {
  type TestStatusEnum,
  testStatusEnumFromString,
} from "@/types/TestStatus";

interface Props {
  status: string | null | undefined;
}

const TEST_RESULT_TAGS: Record<TestStatusEnum, TagVariables> = {
  NO_STATUS: {
    icon: <QuestionCircleFilled />,
    color: "default",
    text: "No Status",
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
    text: "Tool halted before testing",
  },
  UNKNOWN: {
    icon: <QuestionCircleFilled />,
    color: "default",
    text: "Status Unknown",
  },
};

export const TestStatusTag: React.FC<Props> = ({ status }) => {
  const statusEnum = testStatusEnumFromString(status);
  return <ResultTag tagVars={TEST_RESULT_TAGS[statusEnum]} />;
};
