import { GithubOutlined } from "@ant-design/icons";
import { Space } from "antd";
import type { FC } from "react";
import AppBarButton from "@/components/AppBar/AppBarButton";
import styles from "@/components/AppBar/index.module.css";
import { ThemeSwitch } from "@/components/ThemeSwitch";

const AppBarButtons: FC = () => {
  return (
    <Space size="large" className={styles.buttonContainer}>
      <AppBarButton
        icon={<GithubOutlined />}
        title="Github"
        href="https://github.com/buildbarn/bb-portal"
      />
      <ThemeSwitch />
    </Space>
  );
};

export default AppBarButtons;
