'use client';

import React from 'react';
import { Layout, Popover, Space } from 'antd';
import { SlackOutlined } from '@ant-design/icons';
import Link from 'next/link';
import styles from '@/components/FooterBar/index.module.css';

interface Props {
  className?: string;
  linkItemClassName?: string;
}

const FooterBar: React.FC<Props> = ({ className, linkItemClassName }) => {
  const linkClassName = linkItemClassName ? linkItemClassName : styles.footerLink;
  return (
    <Layout.Footer className={`${className} ${styles.footerBar}`}>
      <Popover content="#buildbarn @ buildteamworld.slack.com" className={linkClassName}>
        <Link href="https://bit.ly/2SG1amT" target="_blank">
          <Space>
            <SlackOutlined />
            Buildbarn Slack Channel
          </Space>
        </Link>
      </Popover>
    </Layout.Footer>
  );
};

export default FooterBar;
