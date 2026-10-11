import { ToolOutlined } from "@ant-design/icons";
import { theme as antTheme, Space, Switch, Tooltip, Typography } from "antd";
import type { FC } from "react";
import { useOperatorMode } from "@/utils/operatorMode";

const { useToken } = antTheme;

export const OperatorModeSwitch: FC = () => {
  const [operatorMode, setOperatorMode] = useOperatorMode();
  const { token } = useToken();

  const iconStyle = (active: boolean) => ({
    color: active ? token.colorText : token.colorTextDisabled,
  });

  return (
    <Space size="small">
      <ToolOutlined style={iconStyle(operatorMode)} />
      <Tooltip placement="bottom" title="Show detailed scheduler state">
        <Switch value={operatorMode} onChange={setOperatorMode} />
      </Tooltip>
      <Typography.Text>Operator</Typography.Text>
    </Space>
  );
};
