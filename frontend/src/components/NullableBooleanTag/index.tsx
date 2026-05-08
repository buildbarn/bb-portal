import {
  CheckCircleFilled,
  CloseCircleFilled,
  QuestionCircleFilled,
} from "@ant-design/icons";
import type React from "react";
import ResultTag from "@/components/ResultTag";

interface Props {
  status: boolean | null;
  hideText?: boolean;
}

const NullBooleanTag: React.FC<Props> = ({ status, hideText }) => {
  let text: string;
  let color: string;
  let icon: React.ReactNode;

  switch (status) {
    case true:
      color = "green";
      text = "Yes";
      icon = <CheckCircleFilled />;
      break;
    case false:
      color = "red";
      text = "No";
      icon = <CloseCircleFilled />;
      break;
    default:
      color = "orange";
      text = "?";
      icon = <QuestionCircleFilled />;
      break;
  }

  return (
    <ResultTag
      tagVars={{
        color: color,
        icon: icon,
        text: hideText ? undefined : text,
      }}
    />
  );
};

export default NullBooleanTag;
