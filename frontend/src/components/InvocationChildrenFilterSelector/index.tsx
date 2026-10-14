import { Radio, type RadioChangeEvent } from "antd";
import type React from "react";
import { InvocationChildrenFilter } from "@/routes/scheduler.invocation";

interface Props {
  childrenFilter: InvocationChildrenFilter;
  setChildrenFilter: (value: InvocationChildrenFilter) => void;
}

const InvocationChildrenFilterSelector: React.FC<Props> = ({
  childrenFilter,
  setChildrenFilter,
}) => {
  return (
    <Radio.Group
      buttonStyle="solid"
      value={childrenFilter}
      onChange={(e: RadioChangeEvent) => setChildrenFilter(e.target.value)}
    >
      <Radio.Button value={InvocationChildrenFilter.QUEUED}>
        Queued
      </Radio.Button>
      <Radio.Button value={InvocationChildrenFilter.ACTIVE}>
        Active
      </Radio.Button>
      <Radio.Button value={InvocationChildrenFilter.ALL}>All</Radio.Button>
    </Radio.Group>
  );
};

export default InvocationChildrenFilterSelector;
