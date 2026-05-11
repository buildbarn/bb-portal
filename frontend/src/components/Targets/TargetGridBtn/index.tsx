import {
  CheckCircleFilled,
  CloseCircleFilled,
  QuestionCircleFilled,
} from "@ant-design/icons";
import { Link } from "@tanstack/react-router";
import { Button } from "antd";
import type React from "react";
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
  const resultTag = (
    <Link
      to="/bazel-invocations/$invocationID"
      params={{ invocationID: invocationId }}
    >
      <Button
        icon={getIconForStatus(status)}
        className={getClassForStatus(status)}
      />
    </Link>
  );
  return <>{resultTag}</>;
};

export default TargetGridBtn;
