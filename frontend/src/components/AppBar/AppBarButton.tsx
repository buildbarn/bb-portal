'use client';

import React from 'react';
import { Button, Tooltip } from 'antd';

type Props = {
  icon: React.ReactNode;
  title: string;
  href?: string;
  onMouseDown?: Function;
};

const AppBarButton: React.FC<Props> = ({ icon, title, href, onMouseDown }) => {
  return (
    <Tooltip key={title} placement="bottom" title={title}>
      <Button type="text" href={href} onMouseDown={onMouseDown ? () => onMouseDown() : undefined}>
        {icon}
      </Button>
    </Tooltip>
  );
};

export default AppBarButton;
