
import React, { useMemo } from 'react';
import { useLocation } from '@tanstack/react-router';
import { Breadcrumb } from 'antd';
import type { ItemType } from 'antd/es/breadcrumb/Breadcrumb';
import { Link } from '@tanstack/react-router';
import styles from '@/components/Breadcrumbs/index.module.css';
import BuildbarnIcon from '../BuildbarnIcon';

const itemRender = (currentRoute: ItemType) => {
  return <Link to={currentRoute.path}>{currentRoute.title}</Link>
}

export const Breadcrumbs: React.FC = () => {
  const { pathname } = useLocation();

  const breadcrumbItems = useMemo(() => {
    const items: ItemType[] = [{
      path: '/',
      title: <BuildbarnIcon />,
    }];

    let cumulativePath = '';
    pathname.split('/')
      .filter(segment => segment !== '')
      .forEach(segment => {
        cumulativePath += `/${segment}`;
        items.push({
          path: cumulativePath,
          title: decodeURIComponent(segment),
        });
      });
    return items;
  }, [pathname]);

  return (
    <nav aria-label="Breadcrumb" className={styles.breadcrumbs}>
      <Breadcrumb items={breadcrumbItems} itemRender={itemRender} />
    </nav>
  );
};
