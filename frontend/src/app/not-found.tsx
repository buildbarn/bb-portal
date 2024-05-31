'use client';

import React from 'react';
import Link from 'next/link';
import { Space, Typography } from 'antd';
import Content from '@/components/Content';
import themeStyles from '@/theme/theme.module.css';

const NotFound: React.FC = () => {
  return (
    <>
      <title>Buildbarn Portal - Page Not Found</title>
      <Content
        content={
          <div className={themeStyles.errorPageContainer}>
            <Typography.Title>The page you’re looking for can’t be found.</Typography.Title>
            <Space className={themeStyles.errorPageSubtitle}>
              <Link href="/">Go Back Home</Link>
            </Space>
          </div>
        }
      />
    </>
  );
};

export default NotFound;
