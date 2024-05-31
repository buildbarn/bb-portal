'use client';

import React from 'react';
import { BulbOutlined, GithubOutlined } from '@ant-design/icons';
import AppBarButton from '@/components/AppBar/AppBarButton';
import styles from '@/components/AppBar/index.module.css';

type Props = {
  toggleTheme: () => void;
  prefersDark: boolean;
};

const AppBarButtons: React.FC<Props> = ({ toggleTheme, prefersDark }) => {
  return (
    <div className={styles.buttonContainer}>
      <AppBarButton icon={<GithubOutlined />} title="Github" href="https://github.com/buildbarn/bb-portal" />
      <AppBarButton
        icon={<BulbOutlined />}
        title={`${prefersDark ? 'Light' : 'Dark'} Mode`}
        onMouseDown={toggleTheme}
      />
    </div>
  );
};

export default AppBarButtons;
