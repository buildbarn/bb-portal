import type { ThemeConfig } from "antd";

// The hexadecimal string to be used as the A of RGBA values in configuring Application Bar backgrounds
export const HEADER_OPACITY_HEX = "BB";

const shared: ThemeConfig = {
  components: {
    Alert: {
      withDescriptionIconSize: 24,
    },
    Divider: {
      verticalMarginInline: 16,
    },
    Form: {
      itemMarginBottom: 12,
    },
    Typography: {
      titleMarginBottom: 0,
    },
    Popover: {
      titleMinWidth: 0,
    },
    Table: {
      cellPaddingBlockSM: 2,
    },
    Layout: {
      headerHeight: 64,
      headerPadding: "0 32px",
    },
  },
  token: {
    fontFamily:
      "-apple-system, system-ui, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;",
    fontFamilyCode:
      "'SF Mono', 'Cascadia Code', Consolas, 'Roboto Mono', 'Liberation Mono', monospace;",
  },
};

export default shared;
