import React from 'react';
import { Space, Typography } from 'antd';
import Content from '../Content';
import Link from 'next/link';
import themeStyles from '@/theme/theme.module.css';

const PageDisabled: React.FC = () => {
  return (
    <>
      <title>Buildbarn Portal - Page disabled</title>
      <Content
        content={
          <Space direction='vertical' size='large' className={themeStyles.errorPageContainer}>
            <Typography.Title>This page is disabled.</Typography.Title>
            <div>
              <Typography.Paragraph>
                This page is currently disabled. If you are the system administrator, you can enable it by changing the environment variables of the portal frontend process.
              </Typography.Paragraph>
              <Typography.Paragraph>
                Available environment variables can be found <Link href='https://github.com/buildbarn/bb-portal/blob/main/frontend/.env'>here</Link>.
              </Typography.Paragraph>
            </div>
            <Link href="/">Go Back Home</Link>
          </Space>
        }
      />
    </>
  );
};

export default PageDisabled;
