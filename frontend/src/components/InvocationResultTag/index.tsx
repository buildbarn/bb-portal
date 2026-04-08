import {
  CheckCircleFilled,
  CloseCircleFilled,
  ExclamationCircleFilled,
  InfoCircleFilled,
  LoadingOutlined,
  QuestionCircleFilled,
  StopFilled,
} from "@ant-design/icons";
import { Tag } from "antd";
import type React from "react";
import themeStyles from "@/theme/theme.module.css";
import { getInvocationResultTagEnum, type InvocationResult } from "./enum";

type TagVariables = {
  icon: React.ReactNode;
  color: string;
  text: string;
};

export const INVOCATION_RESULT_TAGS: {
  [key in InvocationResult]: TagVariables;
} = {
  SUCCESS: {
    icon: <CheckCircleFilled />,
    color: "green",
    text: "Succeeded",
  },
  UNSTABLE: {
    icon: <InfoCircleFilled />,
    color: "orange",
    text: "Unstable",
  },
  PARSING_FAILURE: {
    icon: <CloseCircleFilled />,
    color: "red",
    text: "Parsing Failed",
  },
  BUILD_FAILURE: {
    icon: <CloseCircleFilled />,
    color: "red",
    text: "Build Failed",
  },
  TESTS_FAILED: {
    icon: <CloseCircleFilled />,
    color: "red",
    text: "Tests Failed",
  },
  REMOTE_ERROR: {
    icon: <ExclamationCircleFilled />,
    color: "red",
    text: "Remote error",
  },
  NOT_BUILT: {
    icon: <StopFilled />,
    color: "purple",
    text: "Not Built",
  },
  ABORTED: {
    icon: <ExclamationCircleFilled />,
    color: "cyan",
    text: "Aborted",
  },
  INTERRUPTED: {
    icon: <ExclamationCircleFilled />,
    color: "cyan",
    text: "Interrupted",
  },
  // TODO: Make this show the code
  UNKNOWN_EXIT_CODE: {
    icon: <QuestionCircleFilled />,
    color: "default",
    text: "Status Unknown",
  },
  IN_PROGRESS: {
    icon: <LoadingOutlined />,
    color: "blue",
    text: "In Progress",
  },
  DISCONNECTED: {
    icon: <ExclamationCircleFilled />,
    color: "cyan",
    text: "Disconnected",
  },
};

interface Props {
  exitCodeName: string | undefined;
  timeSinceLastConnectionMillis: number | undefined;
}

export const InvocationResultTag: React.FC<Props> = ({
  exitCodeName,
  timeSinceLastConnectionMillis,
}) => {
  const tagEnum = getInvocationResultTagEnum(
    exitCodeName,
    timeSinceLastConnectionMillis,
  );
  const tagVars = INVOCATION_RESULT_TAGS[tagEnum];
  return (
    <Tag icon={tagVars.icon} color={tagVars.color} className={themeStyles.tag}>
      <span className={themeStyles.tagContent}>{tagVars.text}</span>
    </Tag>
  );
};
