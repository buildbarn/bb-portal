import { MenuOutlined } from "@ant-design/icons";
import { Button, Divider, Drawer, Layout } from "antd";
import type React from "react";
import { useEffect, useState } from "react";
import AppBarButtons from "@/components/AppBar/AppBarButtons";
import { AppBarMenu } from "@/components/AppBar/AppBarMenu";
import AppBarTitle from "@/components/AppBar/AppBarTitle";
import styles from "@/components/AppBar/index.module.css";
import FooterBar from "@/components/FooterBar";
import useScreenSize from "@/utils/screen";

const SIDE_BAR_MINIMUM_SCREEN_WIDTH = 932;

type Props = {
  toggleTheme: () => void;
  prefersDark: boolean;
};

const AppBar: React.FC<Props> = ({ toggleTheme, prefersDark }) => {
  const screenSize = useScreenSize();
  const showHeaderMenu = screenSize.width > SIDE_BAR_MINIMUM_SCREEN_WIDTH;
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
            <AppBarMenu mode="horizontal" />
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
      {!showHeaderMenu && (
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
          <AppBarMenu mode="inline" />
          <Divider orientation="center" type="horizontal" />
          <AppBarButtons toggleTheme={toggleTheme} prefersDark={prefersDark} />
        </Drawer>
      )}
    </>
  );
};

export default AppBar;
