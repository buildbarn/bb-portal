import { Link } from "@tanstack/react-router";
import { Menu } from "antd";
import type { MenuMode } from "rc-menu/es/interface";
import type React from "react";
import {
  filterPortalMenuItems,
  type PortalMenuItem,
  usePortalMenuSelectedKeys,
} from "@/types/PortalMenuItem";
import { env } from "@/utils/env";
import styles from "./index.module.css";

const MENU_ITEMS: PortalMenuItem[] = [
  {
    key: "/builds",
    label: <Link to="/builds">Builds</Link>,
    requiredFeatures: [env.featureFlags?.bes?.pageBuilds],
  },
  {
    key: "/bazel-invocations",
    label: <Link to="/bazel-invocations">Invocations</Link>,
    requiredFeatures: [env.featureFlags?.bes?.pageInvocations],
  },
  {
    key: "/trends",
    label: <Link to="/trends">Trends</Link>,
    requiredFeatures: [env.featureFlags?.bes?.pageTrends],
  },
  {
    key: "/tests",
    label: <Link to="/tests">Tests</Link>,
    requiredFeatures: [env.featureFlags?.bes?.pageTests],
  },
  {
    key: "/targets",
    label: <Link to="/targets">Targets</Link>,
    requiredFeatures: [env.featureFlags?.bes?.pageTargets],
  },
  {
    key: "/browser",
    label: <Link to="/browser">Browser</Link>,
    requiredFeatures: [env.featureFlags?.browser],
  },
  {
    key: "/scheduler",
    label: <Link to="/scheduler">Scheduler</Link>,
    requiredFeatures: [env.featureFlags?.scheduler],
  },
  {
    key: "/operations",
    label: <Link to="/operations">Operations</Link>,
    requiredFeatures: [env.featureFlags?.scheduler],
  },
];

const APP_BAR_MENU_ITEMS: PortalMenuItem[] = filterPortalMenuItems(MENU_ITEMS);

interface Props {
  mode: MenuMode;
}

export const AppBarMenu: React.FC<Props> = ({ mode }) => {
  const selectedKeys = usePortalMenuSelectedKeys(APP_BAR_MENU_ITEMS);
  return (
    <Menu
      selectedKeys={selectedKeys}
      mode={mode}
      items={APP_BAR_MENU_ITEMS}
      className={styles.menu}
    />
  );
};
