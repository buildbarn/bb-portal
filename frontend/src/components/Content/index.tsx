import { FloatButton, Layout } from "antd";
import type React from "react";
import { Breadcrumbs } from "@/components/Breadcrumbs";
import styles from "@/components/Content/index.module.css";
import FooterBar from "@/components/FooterBar";

interface Props {
  content: React.ReactNode;
}

const Content: React.FC<Props> = ({ content }) => {
  return (
    <Layout>
      <div className={styles.container}>
        <Breadcrumbs />
        <Layout.Content className={styles.content}>{content}</Layout.Content>
        <div className={styles.footer}>
          <FooterBar />
        </div>
      </div>
      <FloatButton.BackTop />
    </Layout>
  );
};

export default Content;
