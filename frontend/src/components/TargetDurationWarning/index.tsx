import { WarningOutlined } from "@ant-design/icons";
import { Space, Tooltip } from "antd";

const WARNING_TOOLTIP_TEXT =
  "This duration cannot be fully trusted, as it is not measured by Bazel directly. It is based on when the Build Event Protocol messages were received by the server, which may include network delays and other factors.";

interface Props {
  text: string;
}

export const TargetDurationWarning: React.FC<Props> = ({ text }) => {
  return (
    <Tooltip title={WARNING_TOOLTIP_TEXT}>
      <Space direction="horizontal">
        <span>{text}</span>
        <WarningOutlined />
      </Space>
    </Tooltip>
  );
};
