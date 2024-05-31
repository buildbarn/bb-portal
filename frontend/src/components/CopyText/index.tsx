import React, { MouseEventHandler, useState } from 'react';
import { Space, Typography } from 'antd';
import { CheckOutlined, CopyOutlined } from '@ant-design/icons';
import styles from './index.module.css';

interface Props {
  text: string;
  copyText: string;
  onClick?: MouseEventHandler;
}

const CopyText: React.FC<Props> = ({ text, copyText, onClick }) => {
  const [copied, setCopied] = useState<boolean>(false);
  return (
    <Typography.Text
      className={styles.copyText}
      onClick={e => {
        navigator.clipboard.writeText(copyText).then(null, () => console.log('Failed to write to clipboard'));
        setCopied(true);
        setTimeout(() => {
          setCopied(false);
        }, 1500);
        onClick && onClick(e);
      }}
    >
      <Space size={4}>
        {text}
        {copied ? <CheckOutlined /> : <CopyOutlined />}
      </Space>
    </Typography.Text>
  );
};

export default CopyText;
