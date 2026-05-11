import {
  CheckCircleFilled,
  CloseCircleFilled,
  InfoCircleFilled,
  MinusCircleFilled,
  QuestionCircleFilled,
  StopOutlined,
} from "@ant-design/icons";
import type { ButtonColorType } from "antd/es/button";
import type React from "react";
import { LinkButton } from "@/components/LinkButton";
import type { TestStatusEnum } from "@/components/TestStatusTag";
import themeStyles from "@/theme/theme.module.css";

const STATUS_CONFIG: Record<
  TestStatusEnum,
  {
    icon: React.ReactNode;
    className: string;
    color?: ButtonColorType;
  }
> = {
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
  status: TestStatusEnum;
  invocationId: string;
}

const TestGridBtn: React.FC<Props> = ({ status, invocationId }) => {
  const config = STATUS_CONFIG[status] || STATUS_CONFIG.NO_STATUS;

  return (
    <LinkButton
      to="/bazel-invocations/$invocationID"
      params={{ invocationID: invocationId }}
      icon={config.icon}
      className={config.className}
      color={config.color}
    />
  );
};

export default TestGridBtn;
