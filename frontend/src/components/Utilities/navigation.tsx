'use client';

import React, { ReactNode, Key } from 'react';
import { ItemType } from 'antd/es/menu/interface';
import { MenuItemLabel, MenuItemTag } from '@/components/MenuItemLabel';
import { FeatureType, isFeatureEnabled } from '@/utils/isFeatureEnabled';

export type UpdateSidebarMenuExpandedWidthFunction = (updatedSidebarMenuExpandedWidth: number) => void;

interface ItemProps {
  depth: number;
  href: string;
  title: string;
  children?: ItemType[];
  icon?: ReactNode;
  tag?: MenuItemTag;
  danger?: boolean;
  disabled?: boolean;
  activeMenuItemRef?: React.RefObject<HTMLDivElement>;
  updateMenuItemWidth?: UpdateSidebarMenuExpandedWidthFunction;
  requiredFeatures?: FeatureType[];
}

export const getItem = ({
  depth,
  href,
  title,
  children,
  icon,
  tag,
  danger,
  disabled,
  activeMenuItemRef,
  updateMenuItemWidth,
  requiredFeatures,
}: ItemProps): ItemType | undefined => {
  for (const feature of requiredFeatures || []) {
    if (!isFeatureEnabled(feature)) {
      return undefined;
    }
  }

  return {
    key: href,
    icon: icon,
    label: (
      <MenuItemLabel
        ref={activeMenuItemRef}
        depth={depth}
        href={href}
        title={title}
        hasIcon={!!icon}
        hasExpander={!!children?.length}
        tag={tag}
        updateMenuItemWidth={updateMenuItemWidth}
      />
    ),
    danger,
    disabled,
    children,
  };
};

// Recurse to build a flattened array of all children of a given menu tree
const getFlattenedMenuItems = (array: ItemType[]): ItemType[] => {
  let flattenedMenuItems: ItemType[] = [];
  array.forEach(item => {
    flattenedMenuItems.push(item);
    if (item && 'children' in item && item.children) {
      flattenedMenuItems = flattenedMenuItems.concat(getFlattenedMenuItems(item.children));
    }
  });
  return flattenedMenuItems;
};

// Return an array of all keys of a given menu tree
export const getFlattedMenuKeys = (items: ItemType[]): Key[] => {
  return getFlattenedMenuItems(items)
    .map((item: ItemType) => item?.key)
    .filter((key): key is Key => !!key);
};

// Return an array of all items marked with danger from a given menu tree
const getFlattenedDangerMenuItems = (items: ItemType[]): ItemType[] => {
  return getFlattenedMenuItems(items).filter(item => item && 'danger' in item && item.danger);
};

// Return an array of all keys of items marked with danger from a given mnu tree
export const getFlattedDangerMenuKeys = (items: ItemType[]): Key[] => {
  return getFlattenedDangerMenuItems(items)
    .map((item: ItemType) => item?.key)
    .filter((key): key is Key => !!key);
};

// Given a menu item key, find the closest corresponding menu item from a hierarchical array of menu items
// An exact match is not guaranteed so perform a recursive search to find the nearest menu item
// For example, if the provided menu items have keys ["/foo", "/foo/bar/baz"]:
//    - Searching for "/foo" is to return the menu item with key "/foo"
//    - Searching for "/foo/bar" is to return the menu item with key "/foo"
//    - Searching for "/foo/bar/baz" is to return the menu item with key "/foo/bar/baz"
//    - Searching for "/foo/bar/baz/qux" is to return the menu item with key "/foo/bar/baz"
// Return null if there is no possible match, e.g. searching the above for "/banana"
const getClosestItem = (key: Key, items: ItemType[]): ItemType | null => {
  return (
    items?.reduce((previousItem?: ItemType, item?: ItemType) => {
      if (item && 'children' in item && item.children) {
        const closestChildItem = getClosestItem(key, item.children);
        if (closestChildItem) {
          return closestChildItem;
        }
      }
      if (previousItem) {
        return previousItem;
      }
      if (item?.key?.toString() === key.toString().slice(0, item?.key?.toString().length)) {
        return item;
      }
    }, null) || null
  );
};

// Given a menu item key, find the closest corresponding menu item
// An exact match is not guaranteed so perform a recursive search to find the nearest menu item
// For example, if the provided menu items have keys ["/foo", "/foo/bar/baz"]:
//    - Searching for "/foo" is to return "/foo"
//    - Searching for "/foo/bar" is to return "/foo"
//    - Searching for "/foo/bar/baz" is to return "/foo/bar/baz"
//    - Searching for "/foo/bar/baz/qux" is to return "/foo/bar/baz"
// Use getClosestItem() to perform the search and return null if there is no corresponding item
export const getClosestKey = (key: string, items: ItemType[]): Key | null => {
  const closestItem = getClosestItem(key, items);
  return closestItem?.key ? closestItem.key : null;
};

// Given a menu item key, return a list of menu item keys in the subtree leading up to the menu item
export const getFlattenedParentKeys = (key: Key, items: ItemType[]): Key[] => {
  return getFlattedMenuKeys(items).filter(potentialKey => {
    return (
      potentialKey.toString() === key.toString().slice(0, potentialKey.toString().length) &&
      key.toString().length > potentialKey.toString().length
    );
  });
};

// Given a menu item key, return a list of menu item keys of any child menu items
export const getFirstChildrenKeys = (key: Key, items: ItemType[]): Key[] => {
  const closestItem = getClosestItem(key, items);
  if (closestItem && 'children' in closestItem && closestItem.children) {
    return closestItem.children
      .map(child => child?.key)
      .reduce<Key[]>((accumulator, key) => {
        if (key) {
          accumulator.push(key);
        }
        return accumulator;
      }, []);
  }
  return [];
};
