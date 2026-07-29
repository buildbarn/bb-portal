import { Space, Statistic } from "antd";
import type React from "react";

interface MultiStatisticValue {
  key: React.Key;
  value: number | string;
}

interface Props {
  title: string;
  values: MultiStatisticValue[];
}

export const MultiStatistic: React.FC<Props> = ({ title, values }) => {
  return (
    <Space.Compact direction="vertical">
      {values.map((elem, index) => (
        <Statistic
          key={elem.key}
          title={index === 0 ? title : undefined}
          value={elem.value}
        />
      ))}
    </Space.Compact>
  );
};
