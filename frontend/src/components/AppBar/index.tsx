'use client';

import AppBarButtons from '@/components/AppBar/AppBarButtons';
import AppBarMenu from '@/components/AppBar/AppBarMenu';
import AppBarTitle from '@/components/AppBar/AppBarTitle';
import styles from '@/components/AppBar/index.module.css';
import { SIDER_BAR_MINIMUM_SCREEN_WIDTH } from '@/components/Content';
import FooterBar from '@/components/FooterBar';
import { getItem } from '@/components/Utilities/navigation';
import { FeatureType, isFeatureEnabled } from '@/utils/isFeatureEnabled';
import useScreenSize from '@/utils/screen';
import { MenuOutlined } from '@ant-design/icons';
import { Button, Divider, Drawer, Layout } from 'antd';
import type { ItemType } from 'antd/lib/menu/interface';
import type React from 'react';
import { createContext, useEffect, useState } from 'react';

export const SetExtraAppBarMenuItemsContext = createContext<
  React.Dispatch<React.SetStateAction<ItemType[]>> | undefined
>(undefined);

const APP_BAR_MENU_ITEMS: ItemType[] = [
  getItem({ depth: 0, href: '/builds', title: 'Builds', requiredFeatures: [FeatureType.BES] }),
  getItem({ depth: 0, href: '/bazel-invocations', title: 'Invocations', requiredFeatures: [FeatureType.BES]}),
  getItem({ depth: 0, href: '/trends', title: 'Trends', requiredFeatures: [FeatureType.BES]}),
  getItem({ depth: 0, href: '/tests', title: 'Tests', requiredFeatures: [FeatureType.BES, FeatureType.BES_PAGE_TESTS] }),
  getItem({ depth: 0, href: '/targets', title: 'Targets', requiredFeatures: [FeatureType.BES, FeatureType.BES_PAGE_TARGETS] }),
  getItem({ depth: 0, href: '/browser', title: 'Browser', requiredFeatures: [FeatureType.BROWSER] }),
  getItem({
    depth: 0,
    href: '/scheduler',
    title: 'Scheduler',
    children: [
      getItem({ depth: 0, href: '/operations', title: 'Operations' }),
    ],
    requiredFeatures: [FeatureType.SCHEDULER],
  }),
].filter(item => item !== undefined);

type Props = {
  toggleTheme: () => void;
  prefersDark: boolean;
  extraMenuItems: ItemType[];
};

const AppBar: React.FC<Props> = ({
  toggleTheme,
  prefersDark,
  extraMenuItems,
}) => {
  const screenSize = useScreenSize();
  const showHeaderMenu = screenSize.width > SIDER_BAR_MINIMUM_SCREEN_WIDTH;
  const [isDrawerOpen, setIsDrawerOpen] = useState(false);

  useEffect(() => {
    if (showHeaderMenu) {
      setIsDrawerOpen(false);
    }
  }, [showHeaderMenu]);

  return (
    <>
      <Layout.Header className={styles.header}>
        <AppBarTitle />
        {showHeaderMenu ? (
          <>
            <AppBarMenu mode="horizontal" items={APP_BAR_MENU_ITEMS} />
            <AppBarButtons
              toggleTheme={toggleTheme}
              prefersDark={prefersDark}
            />
          </>
        ) : (
          <Button
            type="text"
            onClick={() => setIsDrawerOpen(true)}
            className={styles.menuButton}
          >
            <MenuOutlined />
          </Button>
        )}
      </Layout.Header>
      <Drawer
        placement="right"
        closable={true}
        onClose={() => {
          setIsDrawerOpen(false);
        }}
        onClick={() => {
          setIsDrawerOpen(false);
        }}
        open={isDrawerOpen}
        footer={
          <FooterBar
            className={styles.footerBar}
            linkItemClassName={styles.linkItem}
          />
        }
      >
        <AppBarMenu mode="inline" items={APP_BAR_MENU_ITEMS} />
        <Divider orientation="center" type="horizontal" />
        {extraMenuItems.length ? (
          <>
            <AppBarMenu mode="inline" items={extraMenuItems} />
            <Divider orientation="center" type="horizontal" />
          </>
        ) : null}
        <AppBarButtons toggleTheme={toggleTheme} prefersDark={prefersDark} />
      </Drawer>
    </>
  );
};

export default AppBar;
