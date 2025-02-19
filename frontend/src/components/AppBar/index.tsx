'use client';

import AppBarButtons from '@/components/AppBar/AppBarButtons';
import AppBarMenu from '@/components/AppBar/AppBarMenu';
import AppBarTitle from '@/components/AppBar/AppBarTitle';
import styles from '@/components/AppBar/index.module.css';
import { SIDER_BAR_MINIMUM_SCREEN_WIDTH } from '@/components/Content';
import FooterBar from '@/components/FooterBar';
import { getItem } from '@/components/Utilities/navigation';
import useScreenSize from '@/utils/screen';
import { MenuOutlined } from '@ant-design/icons';
import { Button, Divider, Drawer, Input, Layout } from 'antd';
import type { ItemType } from 'antd/lib/menu/hooks/useItems';
import type React from 'react';
import { createContext, useEffect, useState } from 'react';

export const SetExtraAppBarMenuItemsContext = createContext<
  React.Dispatch<React.SetStateAction<ItemType[]>> | undefined
>(undefined);

const APP_BAR_MENU_ITEMS: ItemType[] = [
  getItem({ depth: 0, href: '/builds', title: 'Builds' }),
  getItem({ depth: 0, href: '/bazel-invocations', title: 'Invocations' }),
  getItem({ depth: 0, href: '/trends', title: 'Trends' }),
  getItem({ depth: 0, href: '/tests', title: 'Tests' }),
  getItem({ depth: 0, href: '/targets', title: 'Targets' }),
  getItem({
    depth: 0,
    href: '/scheduler',
    title: 'Scheduler',
    children: [getItem({ depth: 0, href: '/operations', title: 'Operations' })],
  }),
];

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
  const showSiderBar = screenSize.width > SIDER_BAR_MINIMUM_SCREEN_WIDTH;

  const [isDrawerOpen, setIsDrawerOpen] = useState(false);
  const showDrawer = () => {
    setIsDrawerOpen(true);
  };
  useEffect(() => {
    setIsDrawerOpen(false);
  }, []);
  useEffect(() => {
    if (showSiderBar) {
      setIsDrawerOpen(false);
    }
  }, [showSiderBar]);

  return (
    <div>
      <Layout.Header className={styles.header}>
        <div className={styles.appBar}>
          <AppBarTitle />
          {showSiderBar ? (
            <>
              <div className={styles.appBarMenuContainer}>
                <AppBarMenu
                  mode="horizontal"
                  items={APP_BAR_MENU_ITEMS}
                  className={styles.appBarMenu}
                />
              </div>
              <div className={styles.buttons}>
                <AppBarButtons
                  toggleTheme={toggleTheme}
                  prefersDark={prefersDark}
                />
              </div>
            </>
          ) : (
            <div className={styles.menuButton}>
              <Button type="text" onClick={showDrawer}>
                <MenuOutlined />
              </Button>
            </div>
          )}
        </div>
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
    </div>
  );
};

export default AppBar;
