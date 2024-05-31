'use client';

import React from 'react';
import { ItemType } from 'antd/lib/menu/hooks/useItems';
import { ExclamationCircleOutlined, LoadingOutlined, QuestionCircleOutlined } from '@ant-design/icons';
import { ApolloError } from '@apollo/client';
import { getItem, UpdateSidebarMenuExpandedWidthFunction } from '@/components/Utilities/navigation';
import {
  BazelInvocationInfoFragment,
  FindBuildByUuidQuery, FullBazelInvocationDetailsFragment
} from '@/graphql/__generated__/graphql';
import { getFragmentData } from '@/graphql/__generated__';
import BuildStepStatusIcon from "@/components/BuildStepStatusIcon";
import {
  BAZEL_INVOCATION_FRAGMENT,
  FULL_BAZEL_INVOCATION_DETAILS
} from "@/app/bazel-invocations/[invocationID]/index.graphql";
import {BuildStepResultEnum} from "@/components/BuildStepResultTag";

const getBuildStepsMenuItems = (
  pathBase: string,
  menuItemDepth: number,
  invocations: readonly BazelInvocationInfoFragment[],
  activeMenuItemRef?: React.RefObject<HTMLDivElement>,
  updateMenuItemWidth?: UpdateSidebarMenuExpandedWidthFunction,
): ItemType[] => {
  return [...invocations]
    .sort((a: BazelInvocationInfoFragment, b: BazelInvocationInfoFragment) => a.startedAt - b.startedAt)
    .map(invocation => {
      return getItem({
        depth: menuItemDepth,
        href: `${pathBase}/${encodeURIComponent(invocation.invocationID)}`,
        title: invocation.invocationID,
        icon: <BuildStepStatusIcon status={invocation.state.exitCode?.name as BuildStepResultEnum} />,
        danger: invocation.state.exitCode?.name !== BuildStepResultEnum.SUCCESS,
        activeMenuItemRef,
        updateMenuItemWidth,
      });
    });
};

export const getBuildMenuItems = (
  pathBase: string,
  menuItemDepth: number,
  loading: boolean,
  error?: ApolloError,
  buildQueryResults?: FindBuildByUuidQuery,
  activeMenuItemRef?: React.RefObject<HTMLDivElement>,
  updateMenuItemWidth?: UpdateSidebarMenuExpandedWidthFunction,
): ItemType[] => {
  const build = buildQueryResults?.getBuild

  if (loading) {
    return [getItem({ depth: menuItemDepth, href: pathBase, title: 'Loading...', icon: <LoadingOutlined /> })];
  } else if (error) {
    console.error(error);
    return [
      getItem({
        depth: menuItemDepth,
        href: pathBase,
        title: 'Error',
        icon: <ExclamationCircleOutlined />,
        danger: true,
      }),
    ];
  } else if (!build) {
    return [
      getItem({ depth: 0, href: pathBase, title: 'Awaiting Build Information', icon: <QuestionCircleOutlined /> }),
    ];
  } else {
    const invocations = getFragmentData(FULL_BAZEL_INVOCATION_DETAILS, build.invocations);
    if (!invocations) {
      return [
        getItem({
          depth: menuItemDepth,
          href: pathBase,
          title: 'Awaiting Build Information',
          icon: <QuestionCircleOutlined />,
        }),
      ];
    }

    const invocationOverviews = invocations?.map(i => getFragmentData(BAZEL_INVOCATION_FRAGMENT, i))
    return getBuildStepsMenuItems("/bazel-invocations/", menuItemDepth, invocationOverviews, activeMenuItemRef, updateMenuItemWidth);
  }
};

export const getBreadcrumbSegmentTitles = (buildQueryResults: FindBuildByUuidQuery | undefined, step?: string) => {
  const breadcrumbSegmentTitles: string[] = ['BUILDS'];
  if (buildQueryResults) {
    breadcrumbSegmentTitles.push(buildQueryResults.getBuild?.buildUUID);
  }
  return breadcrumbSegmentTitles;
};
