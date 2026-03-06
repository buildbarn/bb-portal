import { WorkerListStatus } from "@/routes/scheduler.worker";
import { Radio, type RadioChangeEvent } from "antd";
import type React from "react";

interface Props {
  workerStatusFilter: WorkerListStatus;
  setWorkerStatusFilter: (value: WorkerListStatus) => void;
}

const WorkersTableTypeSelector: React.FC<Props> = ({
  workerStatusFilter,
  setWorkerStatusFilter,
}) => {
  return (
    <Radio.Group
      buttonStyle="solid"
      defaultValue={workerStatusFilter}
      onChange={(e: RadioChangeEvent) =>
        setWorkerStatusFilter(e.target.value)
      }
    >
      <Radio.Button value={WorkerListStatus.EXECUTING}>
        Executing
      </Radio.Button>
      <Radio.Button value={WorkerListStatus.ALL}>All</Radio.Button>
    </Radio.Group>
  );
};

export default WorkersTableTypeSelector;
