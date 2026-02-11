import React, { useState } from 'react';
import { Button, Dropdown, DropdownProps, MenuProps, Popover, Space } from 'antd';
import { DownloadOutlined } from '@ant-design/icons';
import { MenuItemType } from 'antd/es/menu/interface';
import Link from '@/components/Link';

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

  const handleMenuClick: MenuProps['onClick'] = e => {
    // Closing is delegated to item.
    if (preventDefaultForKeys && preventDefaultForKeys.includes(e.key)) {
      e.domEvent.preventDefault();
    }
  };
  const onItemClick = () => {
    setTimeout(() => {
      setOpen(false);
    }, 750);
  };
  const handleOpenChange: DropdownProps['onOpenChange'] = (nextOpen: boolean) => {
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
        <Link href={url ?? '#'} download={fileName} target="_self">
          {buttonLabel}
        </Link>
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
        <Link href={url ?? '#'} download={fileName} target="_self">
          {buttonLabel}
        </Link>
      </Space>
    </Dropdown.Button>
  );

  if (popoverTitle) {
    return <Popover content={popoverTitle}>{button}</Popover>;
  }
  return button;
};

export default DownloadButton;
