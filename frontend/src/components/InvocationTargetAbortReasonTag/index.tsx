import {
  AlertOutlined,
  ClockCircleOutlined,
  CloseCircleFilled,
  CloseOutlined,
  ExclamationCircleOutlined,
  MinusCircleFilled,
  MinusOutlined,
  StopOutlined,
  WarningOutlined,
} from "@ant-design/icons";
import { Tag } from "antd";
import type React from "react";
import { InvocationTargetAbortReason } from "@/graphql/__generated__/graphql";
import themeStyles from "@/theme/theme.module.css";

const invocationTargetAbortReasonLabels: Record<
  InvocationTargetAbortReason,
  string
> = {
  [InvocationTargetAbortReason.AnalysisFailure]: "Analysis Failure",
  [InvocationTargetAbortReason.Incomplete]: "Incomplete",
  [InvocationTargetAbortReason.Internal]: "Internal",
  [InvocationTargetAbortReason.LoadingFailure]: "Loading Failure",
  [InvocationTargetAbortReason.NoAnalyze]: "No Analyze",
  [InvocationTargetAbortReason.NoBuild]: "No Build",
  [InvocationTargetAbortReason.None]: "None",
  [InvocationTargetAbortReason.OutOfMemory]: "Out of Memory",
  [InvocationTargetAbortReason.RemoteEnvironmentFailure]:
    "Remote Environment Failure",
  [InvocationTargetAbortReason.Skipped]: "Skipped",
  [InvocationTargetAbortReason.TimeOut]: "Time Out",
  [InvocationTargetAbortReason.Unknown]: "Unknown",
  [InvocationTargetAbortReason.UserInterrupted]: "User Interrupted",
};

const STATUS_TAGS: Record<InvocationTargetAbortReason, React.ReactNode> = {
  [InvocationTargetAbortReason.AnalysisFailure]: (
    <Tag icon={<CloseOutlined />} color="red" className={themeStyles.tag}>
      <span className={themeStyles.tagContent}>
        {
          invocationTargetAbortReasonLabels[
            InvocationTargetAbortReason.AnalysisFailure
          ]
        }
      </span>
    </Tag>
  ),
  [InvocationTargetAbortReason.Incomplete]: (
    <Tag icon={<CloseCircleFilled />} color="gold" className={themeStyles.tag}>
      <span className={themeStyles.tagContent}>
        {
          invocationTargetAbortReasonLabels[
            InvocationTargetAbortReason.Incomplete
          ]
        }
      </span>
    </Tag>
  ),
  [InvocationTargetAbortReason.Internal]: (
    <Tag
      icon={<ExclamationCircleOutlined />}
      color="cyan"
      className={themeStyles.tag}
    >
      <span className={themeStyles.tagContent}>
        {
          invocationTargetAbortReasonLabels[
            InvocationTargetAbortReason.Internal
          ]
        }
      </span>
    </Tag>
  ),
  [InvocationTargetAbortReason.LoadingFailure]: (
    <Tag icon={<CloseOutlined />} color="red" className={themeStyles.tag}>
      <span className={themeStyles.tagContent}>
        {
          invocationTargetAbortReasonLabels[
            InvocationTargetAbortReason.LoadingFailure
          ]
        }
      </span>
    </Tag>
  ),
  [InvocationTargetAbortReason.NoAnalyze]: (
    <Tag
      icon={<MinusCircleFilled />}
      color="geekblue"
      className={themeStyles.tag}
    >
      <span className={themeStyles.tagContent}>
        {
          invocationTargetAbortReasonLabels[
            InvocationTargetAbortReason.NoAnalyze
          ]
        }
      </span>
    </Tag>
  ),
  [InvocationTargetAbortReason.NoBuild]: (
    <Tag
      icon={<MinusCircleFilled />}
      color="geekblue-inverse"
      className={themeStyles.tag}
    >
      <span className={themeStyles.tagContent}>
        {invocationTargetAbortReasonLabels[InvocationTargetAbortReason.NoBuild]}
      </span>
    </Tag>
  ),
  [InvocationTargetAbortReason.None]: (
    <Tag icon={<MinusOutlined />} color="default" className={themeStyles.tag} />
  ),
  [InvocationTargetAbortReason.OutOfMemory]: (
    <Tag icon={<AlertOutlined />} color="red" className={themeStyles.tag}>
      <span className={themeStyles.tagContent}>
        {
          invocationTargetAbortReasonLabels[
            InvocationTargetAbortReason.OutOfMemory
          ]
        }
      </span>
    </Tag>
  ),
  [InvocationTargetAbortReason.RemoteEnvironmentFailure]: (
    <Tag icon={<WarningOutlined />} color="red" className={themeStyles.tag}>
      <span className={themeStyles.tagContent}>
        {
          invocationTargetAbortReasonLabels[
            InvocationTargetAbortReason.RemoteEnvironmentFailure
          ]
        }
      </span>
    </Tag>
  ),
  [InvocationTargetAbortReason.Skipped]: (
    <Tag
      icon={<MinusCircleFilled />}
      className={themeStyles.tag}
      color="purple"
    >
      <span className={themeStyles.tagContent}>
        {invocationTargetAbortReasonLabels[InvocationTargetAbortReason.Skipped]}
      </span>
    </Tag>
  ),
  [InvocationTargetAbortReason.TimeOut]: (
    <Tag icon={<ClockCircleOutlined />} color="red" className={themeStyles.tag}>
      <span className={themeStyles.tagContent}>
        {invocationTargetAbortReasonLabels[InvocationTargetAbortReason.TimeOut]}
      </span>
    </Tag>
  ),
  [InvocationTargetAbortReason.Unknown]: (
    <Tag icon={<AlertOutlined />} color="red" className={themeStyles.tag}>
      <span className={themeStyles.tagContent}>
        {invocationTargetAbortReasonLabels[InvocationTargetAbortReason.Unknown]}
      </span>
    </Tag>
  ),
  [InvocationTargetAbortReason.UserInterrupted]: (
    <Tag icon={<StopOutlined />} color="blue" className={themeStyles.tag}>
      <span className={themeStyles.tagContent}>
        {
          invocationTargetAbortReasonLabels[
            InvocationTargetAbortReason.UserInterrupted
          ]
        }
      </span>
    </Tag>
  ),
};

interface Props {
  reason: InvocationTargetAbortReason;
}

export const InvocationTargetAbortReasonTag: React.FC<Props> = ({ reason }) => {
  return <>{STATUS_TAGS[reason]}</>;
};
