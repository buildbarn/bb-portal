'use client';

import React from 'react';
import { Spin } from 'antd';
import styles from '@/components/FullScreenLoader/index.module.css';

interface Props {
  tip?: React.ReactNode;
}

const FullScreenLoader: React.FC<Props> = ({ tip }) => {
  return (
    <div className={styles.container}>
      <div className={styles.centered}>
        <Spin size="large" tip={<div className={styles.tip}>{tip}</div>}>
          {tip ? ' ' : null}
        </Spin>
      </div>
    </div>
  );
};

export default FullScreenLoader;
