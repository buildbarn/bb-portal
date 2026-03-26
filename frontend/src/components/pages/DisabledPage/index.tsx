import type React from 'react';
import { Space, Typography } from 'antd';
import Content from '../../Content';
import { Link } from '@tanstack/react-router';
import themeStyles from '@/theme/theme.module.css';

export const DisabledPage: React.FC = () => {
  return <Content
    content={
      <Space direction='vertical' size="large" className={themeStyles.errorPageContainer}>
        <Typography.Title>This page is disabled.</Typography.Title>
        <Typography.Paragraph>
          This page is currently disabled. If you are the system administrator, you can enable it by changing the frontend configuration.
        </Typography.Paragraph>
        <Link to="/">Go Back Home</Link>
      </Space>
    }
  />
};
