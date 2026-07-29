import {
  CheckCircleFilled,
  CloseCircleFilled,
  QuestionCircleFilled,
} from "@ant-design/icons";
import type React from "react";
import { LinkButton } from "@/components/LinkButton";
import themeStyles from "@/theme/theme.module.css";

interface Props {
  status: boolean | null;
  invocationId: string;
}

function getIconForStatus(status: boolean | null) {
  if (status == null) {
    return <QuestionCircleFilled />;
  }
  if (status === true) {
    return <CheckCircleFilled />;
  }
  return <CloseCircleFilled />;
}

function getClassForStatus(status: boolean | null) {
  if (status == null) {
    return themeStyles.colorDisabled;
  }
  if (status === true) {
    return themeStyles.colorSuccess;
  }
  return themeStyles.colorFailure;
}

const TargetGridBtn: React.FC<Props> = ({ status, invocationId }) => {
  return (
    <LinkButton
      to="/bazel-invocations/$invocationID"
      params={{ invocationID: invocationId }}
      icon={getIconForStatus(status)}
      className={getClassForStatus(status)}
    />
  );
};

export default TargetGridBtn;
