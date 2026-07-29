import { FloatButton, Layout } from "antd";
import type React from "react";
import AppBar from "@/components/AppBar";
import { Breadcrumbs } from "@/components/Breadcrumbs";
import FooterBar from "@/components/FooterBar";
import styles from "./index.module.css";

interface Props {
  children: React.ReactNode;
  toggleTheme: () => void;
  prefersDark: boolean;
}

export const PageWrapper: React.FC<Props> = ({
  children,
  toggleTheme,
  prefersDark,
}) => {
  return (
    <Layout>
      <AppBar toggleTheme={toggleTheme} prefersDark={prefersDark} />
      <div className={styles.container}>
        <Breadcrumbs />
        <Layout.Content className={styles.content}>{children}</Layout.Content>
        <div className={styles.footer}>
          <FooterBar />
        </div>
      </div>
      <FloatButton.BackTop />
    </Layout>
  );
};
