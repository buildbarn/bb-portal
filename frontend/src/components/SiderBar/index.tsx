
import type React from 'react';
import { type Key, useState } from 'react'
import { Layout, Menu, type MenuProps } from 'antd';
import type { ItemType } from 'antd/es/menu/interface';
import styles from '@/components/SiderBar/index.module.css';
import { useLocation } from '@tanstack/react-router';
import { getClosestKey } from '../Utilities/navigation';

export const SIDEBAR_MENU_INLINE_INDENT = 24;

const DEFAULT_SIDER_WIDTH = 280;

const SIDER_EXPANDED_WIDTH_LOCAL_STORAGE_KEY = 'sider-width';

interface Props {
  menuKey?: Key;
  items: ItemType[];
  defaultSelectedKeys?: Key[];
  defaultOpenKeys?: Key[];
  openKeys?: Key[];
  onOpenChange?: MenuProps['onOpenChange'];
  expandedWidth?: number;
}

const SiderBar: React.FC<Props> = ({
  menuKey,
  items,
  defaultSelectedKeys,
  defaultOpenKeys,
  openKeys,
  onOpenChange,
  expandedWidth,
}) => {
  const { pathname } = useLocation();
  const closestKeyToPathname = getClosestKey(pathname, items);
  const currentKeys = closestKeyToPathname
    ? [closestKeyToPathname.toString()]
    : [];

  const [siderWidth, setSiderWidth] = useState<number>(() => {
    const cachedExpandedState = window.localStorage.getItem(SIDER_EXPANDED_WIDTH_LOCAL_STORAGE_KEY);
    if (cachedExpandedState) {
      return parseInt(cachedExpandedState, 10);
    }
    return DEFAULT_SIDER_WIDTH;
  });

  return (
    <Layout.Sider width={siderWidth}>
      <Menu
        key={menuKey}
        mode="inline"
        defaultSelectedKeys={defaultSelectedKeys?.map(key => key.toString())}
        defaultOpenKeys={defaultOpenKeys?.map(key => key.toString())}
        selectedKeys={currentKeys}
        openKeys={openKeys?.map(key => key.toString())}
        onOpenChange={onOpenChange}
        onMouseMove={() => {
          const updatedSiderWidth =
            expandedWidth && expandedWidth > DEFAULT_SIDER_WIDTH ? expandedWidth : DEFAULT_SIDER_WIDTH;
          setSiderWidth(updatedSiderWidth);
          localStorage.setItem(SIDER_EXPANDED_WIDTH_LOCAL_STORAGE_KEY, updatedSiderWidth.toString());
        }}
        onMouseLeave={() => {
          setSiderWidth(DEFAULT_SIDER_WIDTH);
          localStorage.setItem(SIDER_EXPANDED_WIDTH_LOCAL_STORAGE_KEY, DEFAULT_SIDER_WIDTH.toString());
        }}
        className={styles.menu}
        style={{ width: `${siderWidth}px` }}
        inlineIndent={SIDEBAR_MENU_INLINE_INDENT}
        items={items}
      />
    </Layout.Sider>
  );
};

export default SiderBar;
