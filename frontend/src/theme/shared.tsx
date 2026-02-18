// The hexadecimal string to be used as the A of RGBA values in configuring Application Bar backgrounds
export const HEADER_OPACITY_HEX = 'BB';

const shared = {
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
  },
  token: {
    fontFamily: '-apple-system, system-ui, BlinkMacSystemFont;',
  },
};

export default shared;
