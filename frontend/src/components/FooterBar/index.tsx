'use client';

import React from 'react';
import { Layout, Popover, Space } from 'antd';
import { SlackOutlined } from '@ant-design/icons';
import Link from 'next/link';
import styles from '@/components/FooterBar/index.module.css';
import { env } from 'next-runtime-env';

interface Props {
  className?: string;
  linkItemClassName?: string;
}


const FooterBar: React.FC<Props> = ({ className, linkItemClassName }) => {
  const linkClassName = linkItemClassName ? linkItemClassName : styles.footerLink;
  return (
    <Layout.Footer className={`${className} ${styles.footerBar}`}>
      <Space size="large">
        <Popover content={"#" + env('NEXT_PUBLIC_COMPANY_SLACK_CHANNEL_NAME')} className={linkClassName}>
          <Link href={env('NEXT_PUBLIC_COMPANY_SLACK_CHANNEL_URL') ?? ""} target="_blank" hidden={env('NEXT_PUBLIC_COMPANY_SLACK_CHANNEL_NAME') == undefined}>
            <Space>
              <SlackOutlined />
              {env('NEXT_PUBLIC_COMPANY_SLACK_CHANNEL_NAME')}
            </Space>
          </Link>
        </Popover>
        <Popover content="#buildteam" className={linkClassName}>
          <Link href="https://buildteamworld.slack.com/archives/CD6HZC750" target="_blank">
            <Space>
              <SlackOutlined />
              Buildteam
            </Space>
          </Link>
        </Popover>
      </Space>
    </Layout.Footer>
  );
};

export default FooterBar;
