import { Link, useLocation } from "@tanstack/react-router";
import { Breadcrumb } from "antd";
import type { ItemType } from "antd/es/breadcrumb/Breadcrumb";
import { type CSSProperties, type FC, useMemo } from "react";
import BuildbarnIcon from "../BuildbarnIcon";

interface Props {
  style?: CSSProperties;
}

const itemRender = (currentRoute: ItemType) => {
  return <Link to={currentRoute.path}>{currentRoute.title}</Link>;
};

export const Breadcrumbs: FC<Props> = ({ style }) => {
  const { pathname } = useLocation();

  const breadcrumbItems = useMemo(() => {
    const items: ItemType[] = [
      {
        path: "/",
        title: <BuildbarnIcon />,
      },
    ];

    let cumulativePath = "";
    pathname
      .split("/")
      .filter((segment) => segment !== "")
      .forEach((segment) => {
        cumulativePath += `/${segment}`;
        items.push({
          path: cumulativePath,
          title: decodeURIComponent(segment),
        });
      });
    return items;
  }, [pathname]);

  return (
    <nav aria-label="Breadcrumb" style={style}>
      <Breadcrumb items={breadcrumbItems} itemRender={itemRender} />
    </nav>
  );
};
