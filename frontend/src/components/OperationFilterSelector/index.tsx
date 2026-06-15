import { Radio, type RadioChangeEvent } from "antd";
import type React from "react";

interface Props {
  value: OperationStatus;
  onChange: (value: OperationStatus) => void;
}

export enum OperationStatus {
  ALL = "all",
  EXECUTING = "executing",
  QUEUED = "queued",
  COMPLETED = "completed",
}

export const OperationFilterSelector: React.FC<Props> = ({
  value,
  onChange,
}) => {
  return (
    <Radio.Group
      buttonStyle="solid"
      defaultValue={value}
      onChange={(e: RadioChangeEvent) => onChange(e.target.value)}
    >
      <Radio.Button value={OperationStatus.ALL}>All</Radio.Button>
      <Radio.Button value={OperationStatus.EXECUTING}>Executing</Radio.Button>
      <Radio.Button value={OperationStatus.QUEUED}>Queued</Radio.Button>
      <Radio.Button value={OperationStatus.COMPLETED}>Completed</Radio.Button>
    </Radio.Group>
  );
};
