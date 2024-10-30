'use client';

import React from 'react';
import { Menu } from 'antd';
import { MenuMode } from 'rc-menu/es/interface';
import { usePathname } from 'next/navigation';
//import { ItemType } from 'antd/lib/menu/hooks/useItems';
import styles from './index.module.css';
import { getClosestKey } from '@/components/Utilities/navigation';
import { ItemType } from 'antd/lib/menu/interface';

type Props = {
  mode: MenuMode;
  items: ItemType[];
  className?: string;
};

const AppBarMenu: React.FC<Props> = ({ mode, items, className }) => {
  const closestKeyToPathname = getClosestKey(usePathname(), items);
  const currentKeys = closestKeyToPathname ? [closestKeyToPathname.toString()] : [];
  const classNames = [styles.menu];
  if (className) {
    classNames.push(className);
  }
  return <Menu selectedKeys={currentKeys} mode={mode} style={{ minWidth: 0, flex: "auto" }} items={items} className={classNames.join(' ')} />;
};

export default AppBarMenu;
