import { blue, generate } from "@ant-design/colors";
import { type ThemeConfig, theme } from "antd";
import { HEADER_OPACITY_HEX } from "@/theme/shared";

const LIGHT_CANVAS_BASE_COLOR = "#bdbdbd";

const lightCanvasPalette = generate(LIGHT_CANVAS_BASE_COLOR);

// Overrides for the shared theme, using nested themes from Ant Design:
// https://ant.design/docs/react/customize-theme#nested-theme
export const light: ThemeConfig = {
  algorithm: theme.defaultAlgorithm,
  components: {
    Card: {
      headerBg: lightCanvasPalette[2],
    },
    Layout: {
      bodyBg: lightCanvasPalette[0],
      footerBg: lightCanvasPalette[0],
      headerBg: `${lightCanvasPalette[1]}${HEADER_OPACITY_HEX}`,
    },
    Menu: {
      activeBarBorderWidth: 0,
      itemBg: lightCanvasPalette[1],
      itemHeight: 32,
      itemHoverBg: lightCanvasPalette[3],
      itemMarginInline: 8,
      itemSelectedBg: lightCanvasPalette[3],
    },
  },
  token: {
    colorLink: blue[6],
    colorBgContainer: lightCanvasPalette[0],
    colorBorder: lightCanvasPalette[1],
  },
};
