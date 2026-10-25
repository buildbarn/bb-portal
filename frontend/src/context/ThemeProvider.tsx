import { ConfigProvider, type ThemeConfig } from "antd";
import { type FC, type PropsWithChildren, useEffect, useState } from "react";
import { dark } from "@/theme/dark";
import { light } from "@/theme/light";
import { shared } from "@/theme/shared";
import { type Theme, ThemeContext } from "./ThemeContext";

const configs: Record<Theme, ThemeConfig> = { light, dark };

export const ThemeProvider: FC<PropsWithChildren> = ({ children }) => {
  const [theme, setTheme] = useState<Theme>(
    document.documentElement.getAttribute("data-theme") as Theme,
  );

  useEffect(() => {
    if (window.matchMedia(`(prefers-color-scheme: ${theme})`).matches) {
      localStorage.removeItem("theme");
    } else {
      localStorage.setItem("theme", theme);
    }
    document.documentElement.setAttribute("data-theme", theme);
  }, [theme]);

  return (
    <ThemeContext.Provider
      value={{
        theme,
        setTheme,
      }}
    >
      {/* This seems to be the best way to have a common shared theme for both light and dark. https://ant.design/docs/react/customize-theme#nested-theme */}
      <ConfigProvider theme={shared}>
        <ConfigProvider theme={configs[theme]}>{children}</ConfigProvider>
      </ConfigProvider>
    </ThemeContext.Provider>
  );
};
