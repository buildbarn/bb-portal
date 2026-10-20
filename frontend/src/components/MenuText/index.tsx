import { CaretDownOutlined } from "@ant-design/icons";
import { Dropdown, type MenuProps, Typography } from "antd";
import type React from "react";

interface Props {
  items: MenuProps["items"] | undefined;
  children: React.ReactNode;
}

const MenuText: React.FC<Props> = ({ items, children }) =>
  items ? (
    <Dropdown trigger={["click", "contextMenu"]} menu={{ items }}>
      <Typography.Text style={{ cursor: "pointer" }}>
        {children} <CaretDownOutlined style={{ opacity: 0.6 }} />
      </Typography.Text>
    </Dropdown>
  ) : (
    <Typography.Text>{children}</Typography.Text>
  );

export default MenuText;
