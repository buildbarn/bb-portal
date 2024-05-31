import React, { ForwardedRef, forwardRef, useLayoutEffect, useRef } from 'react';
import Link from 'next/link';
import { Tag } from 'antd';
import { usePathname } from 'next/navigation';
import styles from '@/components/MenuItemLabel/index.module.css';
import { UpdateSidebarMenuExpandedWidthFunction } from '@/components/Utilities/navigation';
import { SIDEBAR_MENU_INLINE_INDENT } from '@/components/SiderBar';

export interface MenuItemTag {
  label: string;
  color: string;
}

export interface Props {
  depth: number;
  href: string;
  title: string;
  hasIcon?: boolean;
  hasExpander?: boolean;
  tag?: MenuItemTag;
  updateMenuItemWidth?: UpdateSidebarMenuExpandedWidthFunction;
}

export const MenuItemLabel = forwardRef((props: Props, ref: ForwardedRef<HTMLDivElement>) => {
  const pathname = usePathname();
  const menuItemLabelRef = useRef<HTMLAnchorElement>(null);
  useLayoutEffect(() => {
    const handleResize = () => {
      if (props?.updateMenuItemWidth && menuItemLabelRef.current?.clientWidth) {
        const itemMarginInline = 8;
        const iconWidth = props.hasIcon ? 16 : 0;
        const iconMarginInlineEnd = props.hasIcon ? 10 : 0;
        const marginBetweenLabelAndExpandWidth = 16;
        const submenuExpanderWidth = props.hasExpander ? 34 : 0;
        const fullMenuItemWidth =
          itemMarginInline +
          iconWidth +
          iconMarginInlineEnd +
          (props.depth + 1) * SIDEBAR_MENU_INLINE_INDENT +
          menuItemLabelRef.current.clientWidth +
          marginBetweenLabelAndExpandWidth +
          submenuExpanderWidth +
          itemMarginInline;
        props.updateMenuItemWidth(fullMenuItemWidth);
      }
    };
    handleResize();
    window.addEventListener('resize', handleResize);
    return () => {
      window.removeEventListener('resize', handleResize);
    };
  }, [props]);
  let label = (
    <Link ref={menuItemLabelRef} href={props.href} className={styles.menuItemLabel}>
      {props.title}
    </Link>
  );
  if (props.tag) {
    label = (
      <span>
        {label}
        <span className={styles.menuItemTag}>
          <Tag bordered={false} color={props.tag.color}>
            {props.tag.label}
          </Tag>
        </span>
      </span>
    );
  }
  if (pathname === props.href) {
    label = <div ref={ref}>{label}</div>;
  }
  return label;
});
MenuItemLabel.displayName = 'MenuItemLabel';

export default MenuItemLabel;
