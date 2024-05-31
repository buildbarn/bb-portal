'use client';

import React, { useContext, useEffect, useRef, useState } from 'react';
import { useQuery } from '@apollo/client';
import { FindBuildByUuidQuery } from '@/graphql/__generated__/graphql';
import Content from '@/components/Content';
import { UpdateSidebarMenuExpandedWidthFunction } from '@/components/Utilities/navigation';
import { getBreadcrumbSegmentTitles, getBuildMenuItems } from '@/app/builds/[buildUUID]/[[...slugs]]/layout.helpers';
import { SetExtraAppBarMenuItemsContext } from '@/components/AppBar';
import {FIND_BUILD_BY_UUID_QUERY} from "@/app/builds/[buildUUID]/[[...slugs]]/index.graphql";

export default function Page({
  children,
  params,
}: {
  children: React.ReactNode;
  params: { buildUUID: string; slugs?: string[] };
}) {
  const [step] = params.slugs ?? [];

  const { loading, error, data } = useQuery<FindBuildByUuidQuery>(FIND_BUILD_BY_UUID_QUERY, {
    variables: { uuid: params.buildUUID },
    fetchPolicy: 'no-cache',
    pollInterval: 60000,
  });

  const pathBase = `/builds/${params.buildUUID}`;

  const activeMenuItemRef = useRef<HTMLDivElement>(null);
  const [sidebarMenuExpandedWidth, setSidebarMenuExpandedWidth] = useState<number>(0);
  const updateSidebarMenuExpandedWidth: UpdateSidebarMenuExpandedWidthFunction = (
    updatedSidebarMenuExpandedWidth: number,
  ) => {
    if (updatedSidebarMenuExpandedWidth && updatedSidebarMenuExpandedWidth > sidebarMenuExpandedWidth) {
      setSidebarMenuExpandedWidth(updatedSidebarMenuExpandedWidth);
    }
  };
  const sidebarMenuItems = getBuildMenuItems(
    '/builds',
    0,
    loading,
    error,
    data,
    activeMenuItemRef,
    updateSidebarMenuExpandedWidth,
  );

  const setExtraAppBarMenuItems = useContext(SetExtraAppBarMenuItemsContext);
  if (setExtraAppBarMenuItems) {
    const sidebarMenuItemsWithoutChildren = sidebarMenuItems.map(item => {
      if (item && 'children' in item) {
        return { ...item, children: undefined };
      }
      return item;
    });
    setExtraAppBarMenuItems(sidebarMenuItemsWithoutChildren);
  }
  useEffect(() => {
    return () => {
      if (setExtraAppBarMenuItems) {
        setExtraAppBarMenuItems([]);
      }
    };
  }, [setExtraAppBarMenuItems]);

  const breadcrumbSegmentTitles = getBreadcrumbSegmentTitles(data, step);

  useEffect(() => {
    if (activeMenuItemRef.current)
      activeMenuItemRef.current.scrollIntoView({
        behavior: 'smooth',
        block: 'center',
      });
  });

  return (
    <Content
      breadcrumbSegmentTitles={breadcrumbSegmentTitles}
      sidebarMenuItems={sidebarMenuItems}
      sidebarMenuDefaultOpenKeys={[pathBase]}
      sidebarMenuExpandedWidth={sidebarMenuExpandedWidth}
      content={children}
    />
  );
}
