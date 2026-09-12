import { blue, generate } from "@ant-design/colors";
import { type ThemeConfig, theme } from "antd";
import shared, { HEADER_OPACITY_HEX } from "@/theme/shared";

const LIGHT_CANVAS_BASE_COLOR = "#bdbdbd";

const lightCanvasPalette = generate(LIGHT_CANVAS_BASE_COLOR);

const light: ThemeConfig = {
  ...shared.token,
  algorithm: theme.defaultAlgorithm,
  components: {
    ...shared.components,
    Card: {
      ...shared.components?.Card,
      headerBg: lightCanvasPalette[2],
    },
    Layout: {
      ...shared.components?.Layout,
      bodyBg: lightCanvasPalette[0],
      footerBg: lightCanvasPalette[0],
      headerBg: `${lightCanvasPalette[1]}${HEADER_OPACITY_HEX}`,
    },
    Menu: {
      ...shared.components?.Menu,
      activeBarBorderWidth: 0,
      itemBg: lightCanvasPalette[1],
      itemHeight: 32,
      itemHoverBg: lightCanvasPalette[3],
      itemMarginInline: 8,
      itemSelectedBg: lightCanvasPalette[3],
    },
  },
  token: {
    ...shared.token,
    colorLink: blue[6],
    colorBgContainer: lightCanvasPalette[0],
    colorBorder: lightCanvasPalette[1],
  },
};

export default light;
