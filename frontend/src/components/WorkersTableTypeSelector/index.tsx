import { ListWorkerFilterType } from "@/types/ListWorkerFilterType";
import { Radio, type RadioChangeEvent } from "antd";
import type React from "react";

interface Props {
  listWorkerFilterType: ListWorkerFilterType;
  setListWorkerFilterType: (listWorkerFilterType: ListWorkerFilterType) => void;
}

const WorkersTableTypeSelector: React.FC<Props> = ({
  listWorkerFilterType,
  setListWorkerFilterType,
}) => {
  return (
    <Radio.Group
      buttonStyle="solid"
      defaultValue={listWorkerFilterType}
      onChange={(e: RadioChangeEvent) =>
        setListWorkerFilterType(e.target.value)
      }
    >
      <Radio.Button value={ListWorkerFilterType.EXECUTING}>
        Executing
      </Radio.Button>
      <Radio.Button value={ListWorkerFilterType.ALL}>All</Radio.Button>
    </Radio.Group>
  );
};

export default WorkersTableTypeSelector;
