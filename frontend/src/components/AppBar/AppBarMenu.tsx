
import { getClosestKey } from "@/components/Utilities/navigation";
import { Menu } from "antd";
import type { ItemType } from "antd/lib/menu/interface";
import type { MenuMode } from "rc-menu/es/interface";
import type React from "react";
import styles from "./index.module.css";
import { useLocation } from "@tanstack/react-router";

type Props = {
  mode: MenuMode;
  items: ItemType[];
};

const AppBarMenu: React.FC<Props> = ({ mode, items }) => {
  const { pathname } = useLocation();
  const closestKeyToPathname = getClosestKey(pathname, items);
  const currentKeys = closestKeyToPathname
    ? [closestKeyToPathname.toString()]
    : [];
  return (
    <Menu
      selectedKeys={currentKeys}
      mode={mode}
      items={items}
      className={styles.menu}
    />
  );
};

export default AppBarMenu;
