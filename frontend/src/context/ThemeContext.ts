import { createContext, useContext } from "react";

export const themes = ["light", "dark"] as const;
export type Theme = (typeof themes)[number];

interface ThemeContextState {
  setTheme: (theme: Theme) => void;
  theme: Theme;
}

export const ThemeContext = createContext<ThemeContextState>({
  setTheme: () => {},
  theme: "light",
});

export const useTheme = () => useContext(ThemeContext);
