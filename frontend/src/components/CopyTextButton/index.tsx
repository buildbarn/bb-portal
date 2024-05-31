import React, { useEffect, useState } from 'react';
import { Button, Typography } from 'antd';
import { CheckOutlined, CopyOutlined } from '@ant-design/icons';
import styles from './index.module.css';

type ButtonTypes = 'text' | 'link' | 'default' | 'primary' | 'dashed' | undefined;

const CopyTextButton: React.FC<{ buttonType?: ButtonTypes; buttonText: string; copyText: string }> = ({
  buttonType,
  buttonText,
  copyText,
}) => {
  const [copied, setCopied] = useState<boolean>(false);

  useEffect(() => {
    setTimeout(() => setCopied(false), 1000);
  }, [copied, setCopied]);

  const button = (
    <Button icon={copied ? <CheckOutlined /> : <CopyOutlined />} type={buttonType}>
      {buttonText}
    </Button>
  );

  return (
    <Typography.Text
      copyable={{ text: copyText, icon: [button, button], tooltips: false, onCopy: () => setCopied(true) }}
      className={styles.buttonWrapper}
    />
  );
};

export default CopyTextButton;
