import { DownloadOutlined } from "@ant-design/icons";
import {
  Button,
  Dropdown,
  type DropdownProps,
  type MenuProps,
  Popover,
  Space,
} from "antd";
import type { MenuItemType } from "antd/es/menu/interface";
import type React from "react";
import { useState } from "react";

export interface DownloadOpts {
  url?: string;
  fileName?: string;
  buttonLabel?: string;
  popoverTitle?: string;
  enabled?: boolean;
  items?: MenuItemType[];
  preventDefaultForKeys?: string[];
}

interface Props extends DownloadOpts {
  renderIfDisabled?: boolean;
}

const DownloadButton: React.FC<Props> = ({
  url,
  fileName,
  buttonLabel,
  enabled,
  popoverTitle,
  items,
  preventDefaultForKeys,
  renderIfDisabled = true,
}) => {
  const [open, setOpen] = useState(false);

  const handleMenuClick: MenuProps["onClick"] = (e) => {
    // Closing is delegated to item.
    if (preventDefaultForKeys?.includes(e.key)) {
      e.domEvent.preventDefault();
    }
  };
  const onItemClick = () => {
    setTimeout(() => {
      setOpen(false);
    }, 750);
  };
  const handleOpenChange: DropdownProps["onOpenChange"] = (
    nextOpen: boolean,
  ) => {
    setOpen(nextOpen);
  };
  if (items) {
    for (const item of items) {
      item.onClick = onItemClick;
    }
  }

  if (!enabled && !renderIfDisabled) {
    return null;
  }

  if (!items || items.length === 0) {
    return (
      <Button icon={<DownloadOutlined />} disabled={!enabled}>
        <a href={url ?? "#"} download={fileName} target="_self">
          {buttonLabel}
        </a>
      </Button>
    );
  }

  // eslint-disable-next-line prefer-const
  const button = (
    <Dropdown.Button
      disabled={!enabled}
      menu={{ items: items, onClick: handleMenuClick }}
      onOpenChange={handleOpenChange}
      open={open}
    >
      <Space>
        <DownloadOutlined />
        <a href={url ?? "#"} download={fileName} target="_self">
          {buttonLabel}
        </a>
      </Space>
    </Dropdown.Button>
  );

  if (popoverTitle) {
    return <Popover content={popoverTitle}>{button}</Popover>;
  }
  return button;
};

export default DownloadButton;
