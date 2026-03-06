
import type React from 'react';
import { Link } from '@tanstack/react-router';
import { Space, Typography } from 'antd';
import Content from '@/components/Content';
import themeStyles from '@/theme/theme.module.css';

interface Props {
  error: Error;
}

export const ErrorPage: React.FC<Props> = ({ error }) => {
  return <Content
    content={
      <Space direction='vertical' size="large" className={themeStyles.errorPageContainer}>
        <Typography.Title>Something went wrong!</Typography.Title>
        {error && 
          <Typography.Text>Error: {error.message}</Typography.Text>
        }
        <Link to="/">Go Back Home</Link>
      </Space>
    }
  />
};
