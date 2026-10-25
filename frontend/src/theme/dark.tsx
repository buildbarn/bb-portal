import { generate } from "@ant-design/colors";
import { type ThemeConfig, theme } from "antd";
import { HEADER_OPACITY_HEX } from "@/theme/shared";

const DARK_CANVAS_BASE_COLOR = "#001d66";

const darkCanvasPalette = generate(DARK_CANVAS_BASE_COLOR, { theme: "dark" });

// Overrides for the shared theme, using nested themes from Ant Design:
// https://ant.design/docs/react/customize-theme#nested-theme
export const dark: ThemeConfig = {
  algorithm: theme.darkAlgorithm,
  components: {
    Card: {
      headerBg: darkCanvasPalette[3],
    },
    Layout: {
      bodyBg: darkCanvasPalette[0],
      footerBg: darkCanvasPalette[0],
      headerBg: `${darkCanvasPalette[1]}${HEADER_OPACITY_HEX}`,
    },
    Menu: {
      activeBarBorderWidth: 0,
      itemBg: darkCanvasPalette[1],
      itemHeight: 32,
      itemHoverBg: darkCanvasPalette[1],
      itemMarginInline: 8,
      itemSelectedBg: darkCanvasPalette[3],
    },
  },
  token: {
    colorBgContainer: darkCanvasPalette[2],
    colorBorder: darkCanvasPalette[1],
  },
};
