"use client";

import styles from "@/components/FooterBar/index.module.css";
import {
  DisconnectOutlined,
  GithubOutlined,
  SlackOutlined,
} from "@ant-design/icons";
import { Layout, Space } from "antd";
import { env } from "next-runtime-env";
import Link from "next/link";
import React from "react";

interface FooterLinkProps {
  text: string;
  href?: string;
  icon?: string;
}

const FooterLink: React.FC<FooterLinkProps> = ({ text, href, icon }) => {
  let iconElement = undefined;
  switch (icon) {
    case "slack":
      iconElement = <SlackOutlined />;
      break;
    case "github":
      iconElement = <GithubOutlined />;
      break;
    case "discord":
      iconElement = <DisconnectOutlined />;
      break;
    case undefined:
      iconElement = undefined;
      break;
    default:
      iconElement = <img src={icon} width={20} height={20} />;
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
    <Link href={href} target="_blank">
      <Space>
        {iconElement}
        {text}
      </Space>
    </Link>
  );
};

interface Props {
  className?: string;
}

const FooterBar: React.FC<Props> = ({ className }) => {
  const footerContent: Array<FooterLinkProps> = React.useMemo(() => {
    const footerJson = env("NEXT_PUBLIC_FOOTER_CONTENT_JSON");
    if (!footerJson) return [];
    try {
      return JSON.parse(footerJson);
    } catch (error) {
      console.error("Failed to parse NEXT_PUBLIC_FOOTER_CONTENT_JSON:", error);
      return [];
    }
  }, []);

  return (
    <Layout.Footer className={`${className} ${styles.footerBar}`}>
      <Space size="large">
        {footerContent.map((item: FooterLinkProps, index: number) => (
          <FooterLink
            key={index}
            text={item.text}
            href={item.href}
            icon={item.icon}
          />
        ))}
      </Space>
    </Layout.Footer>
  );
};

export default FooterBar;
