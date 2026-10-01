import { InfoCircleOutlined } from "@ant-design/icons";
import { Tooltip } from "antd";
import type React from "react";

interface Props {
  text: string;
}

export const PortalTooltip: React.FC<Props> = ({ text }) => {
  return (
    <Tooltip title={text}>
      <InfoCircleOutlined />
    </Tooltip>
  );
};
