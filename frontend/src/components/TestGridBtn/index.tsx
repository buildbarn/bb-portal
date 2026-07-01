import {
  CheckCircleFilled,
  CloseCircleFilled,
  InfoCircleFilled,
  MinusCircleFilled,
  QuestionCircleFilled,
  StopOutlined,
} from "@ant-design/icons";
import type { ButtonProps } from "antd/es/button";
import type React from "react";
import { LinkButton } from "@/components/LinkButton";
import themeStyles from "@/theme/theme.module.css";
import {
  type TestStatusEnum,
  testStatusEnumFromString,
} from "@/types/TestStatus";

const STATUS_CONFIG: Record<TestStatusEnum, ButtonProps> = {
  NO_STATUS: {
    icon: <QuestionCircleFilled />,
    className: themeStyles.colorDisabled,
  },
  PASSED: { icon: <CheckCircleFilled />, className: themeStyles.colorSuccess },
  FLAKY: { icon: <InfoCircleFilled />, className: themeStyles.colorAborted },
  FAILED: {
    icon: <CloseCircleFilled />,
    className: themeStyles.colorFailure,
    color: "danger",
  },
  TIMEOUT: {
    icon: <MinusCircleFilled />,
    className: themeStyles.colorFailure,
    color: "danger",
  },
  INCOMPLETE: { icon: <StopOutlined />, className: themeStyles.colorAborted },
  REMOTE_FAILURE: {
    icon: <CloseCircleFilled />,
    className: themeStyles.colorFailure,
    color: "danger",
  },
  FAILED_TO_BUILD: {
    icon: <QuestionCircleFilled />,
    className: themeStyles.colorFailure,
    color: "danger",
  },
  TOOL_HALTED_BEFORE_TESTING: {
    icon: <QuestionCircleFilled />,
    className: themeStyles.colorDisabled,
  },
};

interface Props {
  status: string | null | undefined;
  invocationId: string;
}

const TestGridBtn: React.FC<Props> = ({ status, invocationId }) => {
  const config = STATUS_CONFIG[testStatusEnumFromString(status)];
  return (
    <LinkButton
      to="/bazel-invocations/$invocationID"
      params={{ invocationID: invocationId }}
      {...config}
    />
  );
};

export default TestGridBtn;
