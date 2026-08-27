import { MenuOutlined } from "@ant-design/icons";
import { Button, Divider, Drawer, Grid, Layout } from "antd";
import { type FC, useState } from "react";
import AppBarButtons from "@/components/AppBar/AppBarButtons";
import { AppBarMenu } from "@/components/AppBar/AppBarMenu";
import AppBarTitle from "@/components/AppBar/AppBarTitle";
import styles from "@/components/AppBar/index.module.css";
import FooterBar from "@/components/FooterBar";

const { useBreakpoint } = Grid;

const AppBar: FC = () => {
  const bp = useBreakpoint();
  const [isDrawerOpen, setIsDrawerOpen] = useState(false);

  return (
    <>
      <Layout.Header
        style={{
          inset: 0,
          display: "grid",
          gridTemplateColumns: "max-content 1fr max-content",
          alignItems: "center",
          position: bp.xl ? "static" : "fixed",
          zIndex: 3,
        }}
      >
        <AppBarTitle />
        {bp.xl ? (
          <>
            <AppBarMenu mode="horizontal" />
            <AppBarButtons />
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
      {bp.xl ? null : (
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
          <Divider orientation="horizontal" />
          <AppBarButtons />
        </Drawer>
      )}
    </>
  );
};

export default AppBar;
