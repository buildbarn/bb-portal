import { Radio, type RadioChangeEvent } from "antd";
import type React from "react";
import { WorkerListStatus } from "@/routes/scheduler.worker";
import { useOperatorMode } from "@/utils/operatorMode";

interface Props {
  workerStatusFilter: WorkerListStatus;
  setWorkerStatusFilter: (value: WorkerListStatus) => void;
}

const WorkersTableTypeSelector: React.FC<Props> = ({
  workerStatusFilter,
  setWorkerStatusFilter,
}) => {
  const [operatorMode] = useOperatorMode();
  return (
    <Radio.Group
      buttonStyle="solid"
      value={workerStatusFilter}
      onChange={(e: RadioChangeEvent) => setWorkerStatusFilter(e.target.value)}
    >
      <Radio.Button value={WorkerListStatus.EXECUTING}>Executing</Radio.Button>
      {(operatorMode ||
        workerStatusFilter === WorkerListStatus.IDLE_SYNCHRONIZING) && (
        <Radio.Button value={WorkerListStatus.IDLE_SYNCHRONIZING}>
          Idle synchronizing
        </Radio.Button>
      )}
      <Radio.Button value={WorkerListStatus.ALL}>All</Radio.Button>
    </Radio.Group>
  );
};

export default WorkersTableTypeSelector;
