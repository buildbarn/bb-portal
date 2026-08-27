import { MoonOutlined, SunOutlined } from "@ant-design/icons";
import { theme as antTheme, Space, Switch, Tooltip } from "antd";
import type { FC } from "react";
import { useTheme } from "@/context/ThemeContext";

const { useToken } = antTheme;

export const ThemeSwitch: FC = () => {
  const { theme, setTheme } = useTheme();
  const { token } = useToken();

  const iconStyle = (active: boolean) => ({
    color: active ? token.colorText : token.colorTextDisabled,
  });

  return (
    <Space size="small">
      <MoonOutlined style={iconStyle(theme === "dark")} />
      <Tooltip placement="bottom" title="Switch themes">
        <Switch
          onChange={(v) => setTheme(v ? "light" : "dark")}
          value={theme === "light"}
        />
      </Tooltip>
      <SunOutlined style={iconStyle(theme === "light")} />
    </Space>
  );
};
