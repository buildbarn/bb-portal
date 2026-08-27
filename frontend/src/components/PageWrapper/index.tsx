import { FloatButton, Grid, Layout, theme } from "antd";
import type React from "react";
import type { PropsWithChildren } from "react";
import AppBar from "@/components/AppBar";
import { Breadcrumbs } from "@/components/Breadcrumbs";
import FooterBar from "@/components/FooterBar";
import styles from "./index.module.css";

const { useToken } = theme;
const { useBreakpoint } = Grid;

export const PageWrapper: React.FC<PropsWithChildren> = ({ children }) => {
  const { token } = useToken();
  const bp = useBreakpoint();

  return (
    <Layout className={styles.layout}>
      <AppBar />
      <div
        style={{ marginTop: !bp.xl ? token.Layout?.headerHeight : 0 }}
        className={styles.container}
      >
        <Layout.Content className={styles.content}>
          <Breadcrumbs
            style={{ marginTop: token.marginSM, marginBottom: token.marginSM }}
          />
          {children}
        </Layout.Content>
        <div className={styles.footer}>
          <FooterBar />
        </div>
      </div>
      <FloatButton.BackTop />
    </Layout>
  );
};
