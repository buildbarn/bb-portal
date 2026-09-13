import { useMatches } from "@tanstack/react-router";
import type { ItemType } from "antd/lib/menu/interface";
import { useMemo } from "react";
import type { Empty } from "@/lib/grpc-client/google/protobuf/empty";

export type PortalMenuItem = NonNullable<ItemType> & {
  key: string;
  hidden?: boolean;
  requiredFeatures?: (Empty | undefined)[];
};

export const filterPortalMenuItems = (
  items: PortalMenuItem[],
): PortalMenuItem[] => {
  return (
    items
      // Removes all hidden items.
      .filter((item) => !item.hidden)
      // Removes all items where not every required feature is enabled.
      .filter((item) => item.requiredFeatures?.every((feat) => !!feat) ?? true)
  );
};

export const usePortalMenuSelectedKeys = (
  menuItems: PortalMenuItem[],
): string[] => {
  const matches = useMatches();
  return useMemo(() => {
    const routeIds: string[] = matches.map((match) => match.routeId);
    return menuItems
      .map((item) => item.key)
      .filter((itemKey) => routeIds.includes(itemKey));
  }, [matches, menuItems]);
};
