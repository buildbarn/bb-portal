import type { OperationState } from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import themeStyles from "@/theme/theme.module.css";
import {
  CheckCircleFilled,
  ClockCircleFilled,
  CloseCircleFilled,
  ExclamationCircleFilled,
  LoadingOutlined,
} from "@ant-design/icons";
import { Tag } from "antd";
import { GrpcErrorCodes } from "../../utils/grpcErrorCodes";

interface Props {
  operation: OperationState;
}

const OperationStatusTag: React.FC<Props> = ({ operation }) => {
  if (operation.queued) {
    return (
      <Tag icon={<ClockCircleFilled />} className={themeStyles.tag}>
        <span className={themeStyles.tagContent}>
          Queued at priority {operation.priority}
        </span>
      </Tag>
    );
  }

  if (operation.executing) {
    return (
      <Tag icon={<LoadingOutlined />} className={themeStyles.tag}>
        <span className={themeStyles.tagContent}>Executing</span>
      </Tag>
    );
  }

  if (operation.completed) {
    if (operation.completed.status?.code) {
      return (
        <Tag
          icon={<CloseCircleFilled />}
          color="red"
          className={themeStyles.tag}
        >
          <span className={themeStyles.tagContent}>
            Failed with status code{" "}
            <code>{GrpcErrorCodes[operation.completed.status?.code]}</code>
          </span>
        </Tag>
      );
    }

    if (operation.completed.result) {
      if (operation.completed.result.exitCode === 0) {
        return (
          <Tag
            icon={<CheckCircleFilled />}
            color="green"
            className={themeStyles.tag}
          >
            <span className={themeStyles.tagContent}>
              Completed with exit code 0
            </span>
          </Tag>
        );
      }

      return (
        <Tag
          icon={<ExclamationCircleFilled />}
          color="yellow"
          className={themeStyles.tag}
        >
          <span className={themeStyles.tagContent}>
            Completed with exit code {operation.completed.result.exitCode}
          </span>
        </Tag>
      );
    }

    return (
      <Tag
        icon={<ExclamationCircleFilled />}
        color="red"
        className={themeStyles.tag}
      >
        <span className={themeStyles.tagContent}>Action result missing</span>
      </Tag>
    );
  }

  return (
    <Tag
      icon={<ExclamationCircleFilled />}
      color="red"
      className={themeStyles.tag}
    >
      <span className={themeStyles.tagContent}>Unknown</span>
    </Tag>
  );
};

export default OperationStatusTag;
