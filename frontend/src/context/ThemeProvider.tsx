import { ConfigProvider, type ThemeConfig } from "antd";
import { type FC, type PropsWithChildren, useEffect, useState } from "react";
import dark from "@/theme/dark";
import light from "@/theme/light";
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
      <ConfigProvider theme={configs[theme]}>{children}</ConfigProvider>
    </ThemeContext.Provider>
  );
};
