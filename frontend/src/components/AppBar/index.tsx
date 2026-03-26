import { MenuOutlined } from "@ant-design/icons";
import { Button, Divider, Drawer, Layout } from "antd";
import type { ItemType } from "antd/lib/menu/interface";
import type React from "react";
import { useEffect, useState } from "react";
import AppBarButtons from "@/components/AppBar/AppBarButtons";
import AppBarMenu from "@/components/AppBar/AppBarMenu";
import AppBarTitle from "@/components/AppBar/AppBarTitle";
import styles from "@/components/AppBar/index.module.css";
import { SIDER_BAR_MINIMUM_SCREEN_WIDTH } from "@/components/Content";
import FooterBar from "@/components/FooterBar";
import { getItem } from "@/components/Utilities/navigation";
import { env } from "@/utils/env";
import useScreenSize from "@/utils/screen";

const getAppBarMenuItems = (): ItemType[] => {
  const items: (ItemType | undefined)[] = [
    getItem({
      depth: 0,
      href: "/builds",
      title: "Builds",
      requiredFeatures: [env.featureFlags?.bes?.pageBuilds],
    }),
    getItem({
      depth: 0,
      href: "/bazel-invocations",
      title: "Invocations",
      requiredFeatures: [env.featureFlags?.bes?.pageInvocations],
    }),
    getItem({
      depth: 0,
      href: "/trends",
      title: "Trends",
      requiredFeatures: [env.featureFlags?.bes?.pageTrends],
    }),
    getItem({
      depth: 0,
      href: "/tests",
      title: "Tests",
      requiredFeatures: [env.featureFlags?.bes?.pageTests],
    }),
    getItem({
      depth: 0,
      href: "/targets",
      title: "Targets",
      requiredFeatures: [env.featureFlags?.bes?.pageTargets],
    }),
    getItem({
      depth: 0,
      href: "/browser",
      title: "Browser",
      requiredFeatures: [env.featureFlags?.browser],
    }),
    getItem({
      depth: 0,
      href: "/scheduler",
      title: "Scheduler",
      requiredFeatures: [env.featureFlags?.scheduler],
    }),
    getItem({
      depth: 0,
      href: "/operations",
      title: "Operations",
      requiredFeatures: [env.featureFlags?.scheduler],
    }),
  ];
  return items.filter((item): item is ItemType => item !== undefined);
};

const APP_BAR_MENU_ITEMS: ItemType[] = getAppBarMenuItems();

type Props = {
  toggleTheme: () => void;
  prefersDark: boolean;
};

const AppBar: React.FC<Props> = ({ toggleTheme, prefersDark }) => {
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
        footer={<FooterBar className={styles.footerBar} />}
      >
        <AppBarMenu mode="inline" items={APP_BAR_MENU_ITEMS} />
        <Divider orientation="center" type="horizontal" />
        <AppBarButtons toggleTheme={toggleTheme} prefersDark={prefersDark} />
      </Drawer>
    </>
  );
};

export default AppBar;
