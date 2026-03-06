
import styles from "@/components/FooterBar/index.module.css";
import {
  DiscordOutlined,
  GithubOutlined,
  SlackOutlined
} from "@ant-design/icons";
import { Layout, Space } from "antd";
import { env } from "@/utils/env";
import { Link } from '@tanstack/react-router';
import React from "react";
import { PortalFrontendConfiguration_FooterElement } from "@/lib/grpc-client/portal/frontend/frontend";

const FooterLink: React.FC<PortalFrontendConfiguration_FooterElement> = ({ text, href, icon }) => {
  let iconElement: React.ReactElement | undefined = undefined;
  if (icon?.url) {
    iconElement = <img src={icon.url} width={20} height={20} alt="Footer icon" />
  } else if (icon?.slack) {
    iconElement = <SlackOutlined />;
  } else if (icon?.github) {
    iconElement = <GithubOutlined />;
  } else if (icon?.discord) {
    iconElement = <DiscordOutlined />;
  }

  if (!href) {
    return (
      <Space>
        {iconElement}
        {text}
      </Space>
    );
  }

  return (
    <a href={href} target="_blank">
      <Space>
        {iconElement}
        {text}
      </Space>
    </a>
  );
};

interface Props {
  className?: string;
}

const FooterBar: React.FC<Props> = ({ className }) => {
  return (
    <Layout.Footer className={`${className} ${styles.footerBar}`}>
      <Space size="large">
        {env.footerContent.map((element: PortalFrontendConfiguration_FooterElement, index: number) => (
          <FooterLink key={index} {...element} />
        ))}
      </Space>
    </Layout.Footer>
  );
};

export default FooterBar;
