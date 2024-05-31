'use client';

import React, { Key, useState } from 'react';
import { Layout, Menu, MenuProps } from 'antd';
import { usePathname, useSearchParams } from 'next/navigation';
import { ItemType } from 'antd/lib/menu/hooks/useItems';
import styles from '@/components/SiderBar/index.module.css';

export const SIDEBAR_MENU_INLINE_INDENT = 24;

const DEFAULT_SIDER_WIDTH = 280;

const SIDER_EXPANDED_WIDTH_LOCAL_STORAGE_KEY = 'sider-width';

interface Props {
  menuKey?: Key;
  items?: ItemType[];
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
  // This should be amended to put the search parameters in a known order
  // Logic should then also be added to getItem() to sort them in the same order
  // This is to enable the right menu items to be highlighted whenever they include search parameters
  const searchParams = useSearchParams().toString();
  const searchParamsKeySuffix = searchParams.length ? `?${searchParams}` : '';
  const currentKey = `${usePathname()}${searchParamsKeySuffix}`;

  const [siderWidth, setSiderWidth] = useState<number>(() => {
    const cachedExpandedState = window.localStorage.getItem(SIDER_EXPANDED_WIDTH_LOCAL_STORAGE_KEY);
    if (cachedExpandedState) {
      return parseInt(cachedExpandedState);
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
        selectedKeys={[currentKey]}
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
