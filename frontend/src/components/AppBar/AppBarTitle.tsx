'use client';

import React from 'react';
import { Typography } from 'antd';
import Link from 'next/link';
import styles from '@/components/AppBar/index.module.css';

const AppBarTitle = () => {
  return (
    <div className={styles.title}>
      <Link href="/">
        <Typography.Title level={3}>Buildbarn Portal</Typography.Title>
      </Link>
    </div>
  );
};

export default AppBarTitle;
