import { Button, Tooltip } from "antd";
import type React from "react";

type Props = {
  icon: React.ReactNode;
  title: string;
  href?: string;
  onMouseDown?: () => void;
};

const AppBarButton: React.FC<Props> = ({ icon, title, href, onMouseDown }) => {
  return (
    <Tooltip key={title} placement="bottom" title={title}>
      <Button
        type="text"
        size="large"
        icon={icon}
        href={href}
        onMouseDown={onMouseDown ? () => onMouseDown() : undefined}
      />
    </Tooltip>
  );
};

export default AppBarButton;
