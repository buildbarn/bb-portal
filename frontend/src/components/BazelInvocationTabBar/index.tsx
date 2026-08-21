import {
  BranchesOutlined,
  CodeOutlined,
  DatabaseOutlined,
  DeploymentUnitOutlined,
  ExperimentOutlined,
  FileSearchOutlined,
  InfoCircleOutlined,
  LineChartOutlined,
  TagsOutlined,
} from "@ant-design/icons";
import { Link } from "@tanstack/react-router";
import { Menu } from "antd";
import type React from "react";
import { useMemo } from "react";
import type { BazelInvocationCommonFragment } from "@/graphql/__generated__/graphql";
import {
  filterPortalMenuItems,
  type PortalMenuItem,
  usePortalMenuSelectedKeys,
} from "@/types/PortalMenuItem";
import { env } from "@/utils/env";

const getMenuItems = (
  invocation: BazelInvocationCommonFragment,
): PortalMenuItem[] => {
  const { invocationID } = invocation;

  const hideActionsTab = !invocation.actions?.length;
  const hideMetricsTab = !invocation.metrics;
  const hideSourceControlTab = !invocation.sourceControl?.length;

  return filterPortalMenuItems([
    {
      key: "/bazel-invocations/$invocationID/",
      icon: <InfoCircleOutlined />,
      label: (
        <Link to="/bazel-invocations/$invocationID" params={{ invocationID }}>
          Overview
        </Link>
      ),
    },
    {
      key: "/bazel-invocations/$invocationID/log",
      icon: <FileSearchOutlined />,
      label: (
        <Link
          to="/bazel-invocations/$invocationID/log"
          params={{ invocationID }}
        >
          Log
        </Link>
      ),
    },
    {
      key: "/bazel-invocations/$invocationID/metrics",
      icon: <LineChartOutlined />,
      label: (
        <Link
          to="/bazel-invocations/$invocationID/metrics"
          params={{ invocationID }}
        >
          Metrics
        </Link>
      ),
      hidden: hideMetricsTab,
    },
    {
      key: "/bazel-invocations/$invocationID/targets",
      icon: <DeploymentUnitOutlined />,
      label: (
        <Link
          to="/bazel-invocations/$invocationID/targets"
          params={{ invocationID }}
        >
          Targets
        </Link>
      ),
      requiredFeatures: [env.featureFlags?.bes?.pageTargets],
    },
    {
      key: "/bazel-invocations/$invocationID/tests",
      icon: <ExperimentOutlined />,
      label: (
        <Link
          to="/bazel-invocations/$invocationID/tests"
          params={{ invocationID }}
        >
          Tests
        </Link>
      ),
      requiredFeatures: [env.featureFlags?.bes?.pageTests],
    },
    {
      key: "/bazel-invocations/$invocationID/command-line",
      icon: <CodeOutlined />,
      label: (
        <Link
          to="/bazel-invocations/$invocationID/command-line"
          params={{ invocationID }}
        >
          Command Line
        </Link>
      ),
    },
    {
      key: "/bazel-invocations/$invocationID/source-control",
      icon: <BranchesOutlined />,
      label: (
        <Link
          to="/bazel-invocations/$invocationID/source-control"
          params={{ invocationID }}
        >
          Source Control
        </Link>
      ),
      hidden: hideSourceControlTab,
    },
    // Previously we counted the number of tags to determine if we should show
    // this link, but due to a bug in (probably) ApolloClient, this was removed.
    // The bug caused a cache warning when fetching the count, which in turn
    // caused a rerender of the build details page when you hovered over a link
    // to a invocation.
    {
      key: "/bazel-invocations/$invocationID/tags",
      icon: <TagsOutlined />,
      label: (
        <Link
          to="/bazel-invocations/$invocationID/tags"
          params={{ invocationID }}
        >
          Tags
        </Link>
      ),
    },
    {
      key: "/bazel-invocations/$invocationID/actions",
      icon: <DatabaseOutlined />,
      label: (
        <Link
          to="/bazel-invocations/$invocationID/actions"
          params={{ invocationID }}
        >
          Failed Actions
        </Link>
      ),
      hidden: hideActionsTab,
    },
  ]);
};

interface Props {
  invocation: BazelInvocationCommonFragment;
}

export const BazelInvocationTabBar: React.FC<Props> = ({ invocation }) => {
  const menuItems = useMemo(() => getMenuItems(invocation), [invocation]);
  const selectedKeys = usePortalMenuSelectedKeys(menuItems);
  return (
    <Menu
      mode="horizontal"
      style={{ background: "inherit" }}
      selectedKeys={selectedKeys}
      items={menuItems}
    />
  );
};
