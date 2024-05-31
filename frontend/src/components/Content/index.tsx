'use client';

import React, { Key } from 'react';
import { FloatButton, Layout, MenuProps } from 'antd';
import { ItemType } from 'antd/lib/menu/hooks/useItems';
import styles from '@/components/Content/index.module.css';
import SiderBar from '@/components/SiderBar';
import Breadcrumbs from '@/components/Breadcrumbs';
import FooterBar from '@/components/FooterBar';
import useScreenSize from '@/utils/screen';

export const SIDER_BAR_MINIMUM_SCREEN_WIDTH = 932;

interface Props {
  breadcrumbSegmentTitles?: string[];
  sidebarMenuKey?: Key;
  sidebarMenuItems?: ItemType[];
  sidebarMenuDefaultSelectedKeys?: Key[];
  sidebarMenuDefaultOpenKeys?: Key[];
  sidebarMenuOpenKeys?: Key[];
  sidebarMenuOnOpenChange?: MenuProps['onOpenChange'];
  sidebarMenuExpandedWidth?: number;
  content: React.ReactNode;
}

const Content: React.FC<Props> = ({
  breadcrumbSegmentTitles,
  sidebarMenuKey,
  sidebarMenuItems,
  sidebarMenuDefaultSelectedKeys,
  sidebarMenuDefaultOpenKeys,
  sidebarMenuOpenKeys,
  sidebarMenuOnOpenChange,
  sidebarMenuExpandedWidth,
  content,
}) => {
  const screenSize = useScreenSize();
  const showSiderBar = sidebarMenuItems?.length && screenSize.width > SIDER_BAR_MINIMUM_SCREEN_WIDTH;
  return (
    <Layout>
      {showSiderBar ? (
        <SiderBar
          key={sidebarMenuKey}
          menuKey={sidebarMenuKey}
          items={sidebarMenuItems}
          defaultSelectedKeys={sidebarMenuDefaultSelectedKeys}
          defaultOpenKeys={sidebarMenuDefaultOpenKeys}
          openKeys={sidebarMenuOpenKeys}
          onOpenChange={sidebarMenuOnOpenChange}
          expandedWidth={sidebarMenuExpandedWidth}
        />
      ) : null}
      <Layout className={showSiderBar ? styles.contentWithSiderBar : undefined}>
        <div className={styles.container}>
          <Breadcrumbs segmentTitles={breadcrumbSegmentTitles} />
          <Layout.Content className={styles.content}>{content}</Layout.Content>
          <div className={styles.footer}>
            <FooterBar />
          </div>
        </div>
        <FloatButton.BackTop />
      </Layout>
    </Layout>
  );
};

export default Content;
