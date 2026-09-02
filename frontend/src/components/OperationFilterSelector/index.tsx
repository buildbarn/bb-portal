import { Radio, type RadioChangeEvent } from "antd";
import type React from "react";
import { ExecutionStage_Value } from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";

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

const operationStatusMap: Record<
  OperationStatus,
  ExecutionStage_Value | undefined
> = {
  [OperationStatus.ALL]: undefined,
  [OperationStatus.EXECUTING]: ExecutionStage_Value.EXECUTING,
  [OperationStatus.QUEUED]: ExecutionStage_Value.QUEUED,
  [OperationStatus.COMPLETED]: ExecutionStage_Value.COMPLETED,
};

export const getExecutionStageFromOperationStatus = (
  operationState: OperationStatus,
): ExecutionStage_Value | undefined => {
  return operationStatusMap[operationState];
};

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
