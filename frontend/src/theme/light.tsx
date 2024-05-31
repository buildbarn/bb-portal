import { theme, ThemeConfig } from 'antd';
import { blue, generate } from '@ant-design/colors';
import shared, { HEADER_OPACITY_HEX } from '@/theme/shared';

const LIGHT_CANVAS_BASE_COLOR = '#bdbdbd';

const lightCanvasPalette = generate(LIGHT_CANVAS_BASE_COLOR);

const light: ThemeConfig = {
  algorithm: theme.defaultAlgorithm,
  components: {
    Card: {
      headerBg: lightCanvasPalette[2],
    },
    Layout: {
      bodyBg: lightCanvasPalette[0],
      footerBg: lightCanvasPalette[0],
      headerBg: `${lightCanvasPalette[1]}${HEADER_OPACITY_HEX}`,
      headerPadding: '0 32px',
      siderBg: lightCanvasPalette[1],
    },
    Menu: {
      activeBarBorderWidth: 0,
      itemBg: lightCanvasPalette[1],
      itemHeight: 32,
      itemHoverBg: lightCanvasPalette[3],
      itemMarginInline: 8,
      itemSelectedBg: lightCanvasPalette[3],
    },
    ...shared.components,
  },
  token: {
    colorLink: blue[6],
    colorBgContainer: lightCanvasPalette[0],
    colorBorder: lightCanvasPalette[1],
    ...shared.token,
  },
};

export default light;
