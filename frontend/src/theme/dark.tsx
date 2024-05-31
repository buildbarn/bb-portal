import { theme, ThemeConfig } from 'antd';
import { generate } from '@ant-design/colors';
import shared, { HEADER_OPACITY_HEX } from '@/theme/shared';

const DARK_CANVAS_BASE_COLOR = '#001d66';

const darkCanvasPalette = generate(DARK_CANVAS_BASE_COLOR, { theme: 'dark' });

const dark: ThemeConfig = {
  algorithm: theme.darkAlgorithm,
  components: {
    Card: {
      headerBg: darkCanvasPalette[3],
    },
    Layout: {
      bodyBg: darkCanvasPalette[0],
      footerBg: darkCanvasPalette[0],
      headerBg: `${darkCanvasPalette[1]}${HEADER_OPACITY_HEX}`,
      headerPadding: '0 32px',
      siderBg: darkCanvasPalette[1],
    },
    Menu: {
      activeBarBorderWidth: 0,
      itemBg: darkCanvasPalette[1],
      itemHeight: 32,
      itemHoverBg: darkCanvasPalette[1],
      itemMarginInline: 8,
      itemSelectedBg: darkCanvasPalette[3],
    },
    ...shared.components,
  },
  token: {
    colorBgContainer: darkCanvasPalette[2],
    colorBorder: darkCanvasPalette[1],
    ...shared.token,
  },
};

export default dark;
