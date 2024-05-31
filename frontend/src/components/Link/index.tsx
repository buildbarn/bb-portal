import React from 'react';
import { LinkOutlined } from '@ant-design/icons';
import NextLink from 'next/link';
import { Space, Tooltip } from 'antd';
import styles from './index.module.css';

export interface RequiredLinkProps {
  href: string;
  tooltipTitle?: string;
  hideFlair?: boolean;
}

export type LinkProps = RequiredLinkProps & React.AnchorHTMLAttributes<HTMLAnchorElement>;

const EXTERNAL_LINK_PREFIXES = ['https://', 'http://', 'mailto:'];

const Link: React.FC<LinkProps> = ({ href, children, tooltipTitle, hideFlair, ...props }) => {
  const external = EXTERNAL_LINK_PREFIXES.some(prefix => href.startsWith(prefix));
  const flair = props.download || hideFlair ? null : <LinkOutlined className={styles.flair} />;
  if (external) {
    return (
      <Tooltip title={tooltipTitle} mouseEnterDelay={0.5}>
        <a href={href} target="_blank" rel="noopener noreferrer" {...props}>
          <Space size={4} className={styles.content}>
            {!!children && <span className={styles.content}>{children}</span>}
            <span className={styles.content}>{flair}</span>
          </Space>
        </a>
      </Tooltip>
    );
  } else {
    return (
      <Tooltip title={tooltipTitle} mouseEnterDelay={0.5}>
        <NextLink href={href} {...props}>
          {children}
        </NextLink>
      </Tooltip>
    );
  }
};
export default Link;
