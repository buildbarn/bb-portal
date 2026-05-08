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
import { Link, useLocation } from "@tanstack/react-router";
import { Menu } from "antd";
import type { ItemType } from "antd/lib/menu/interface";
import type React from "react";
import { useMemo } from "react";
import type { BazelInvocationCommonFragment } from "@/graphql/__generated__/graphql";
import { env } from "@/utils/env";

const getMenuItems = (
  invocation: BazelInvocationCommonFragment,
): ItemType[] => {
  const { invocationID } = invocation;

  const showActionsTab = !!invocation.actions?.length;
  const showMetricsTab = !!invocation.metrics;
  const showSourceControlTab = !!invocation.sourceControl?.length;
  const showTagsTab = !!invocation.tags?.totalCount;
  const showTargetsTab = !!env.featureFlags?.bes?.pageTargets;
  const showTestsTab = !!env.featureFlags?.bes?.pageTests;

  const items: ItemType[] = [];
  items.push({
    key: "overview",
    icon: <InfoCircleOutlined />,
    label: (
      <Link to="/bazel-invocations/$invocationID" params={{ invocationID }}>
        Overview
      </Link>
    ),
  });
  items.push({
    key: "log",
    icon: <FileSearchOutlined />,
    label: (
      <Link to="/bazel-invocations/$invocationID/log" params={{ invocationID }}>
        Log
      </Link>
    ),
  });
  if (showMetricsTab)
    items.push({
      key: "metrics",
      icon: <LineChartOutlined />,
      label: (
        <Link
          to="/bazel-invocations/$invocationID/metrics"
          params={{ invocationID }}
        >
          Metrics
        </Link>
      ),
    });
  if (showTargetsTab)
    items.push({
      key: "targets",
      icon: <DeploymentUnitOutlined />,
      label: (
        <Link
          to="/bazel-invocations/$invocationID/targets"
          params={{ invocationID }}
        >
          Targets
        </Link>
      ),
    });
  if (showTestsTab)
    items.push({
      key: "tests",
      icon: <ExperimentOutlined />,
      label: (
        <Link
          to="/bazel-invocations/$invocationID/tests"
          params={{ invocationID }}
        >
          Tests
        </Link>
      ),
    });
  items.push({
    key: "command-line",
    icon: <CodeOutlined />,
    label: (
      <Link
        to="/bazel-invocations/$invocationID/command-line"
        params={{ invocationID }}
      >
        Command Line
      </Link>
    ),
  });
  if (showSourceControlTab)
    items.push({
      key: "source-control",
      icon: <BranchesOutlined />,
      label: (
        <Link
          to="/bazel-invocations/$invocationID/source-control"
          params={{ invocationID }}
        >
          Source Control
        </Link>
      ),
    });
  if (showTagsTab)
    items.push({
      key: "tags",
      icon: <TagsOutlined />,
      label: (
        <Link
          to="/bazel-invocations/$invocationID/tags"
          params={{ invocationID }}
        >
          Tags
        </Link>
      ),
    });
  if (showActionsTab)
    items.push({
      key: "actions",
      icon: <DatabaseOutlined />,
      label: (
        <Link
          to="/bazel-invocations/$invocationID/actions"
          params={{ invocationID }}
        >
          Failed Actions
        </Link>
      ),
    });
  return items;
};

interface Props {
  invocation: BazelInvocationCommonFragment;
}

export const BazelInvocationTabBar: React.FC<Props> = ({ invocation }) => {
  const { pathname } = useLocation();

  const menuItems = useMemo(() => getMenuItems(invocation), [invocation]);

  const selectedKey = useMemo(() => {
    const lastPathSegment = pathname.split("/").pop();
    const activeItem = menuItems.find((item) => item?.key === lastPathSegment);
    return activeItem?.key?.toString() ?? "overview";
  }, [pathname, menuItems]);

  return (
    <Menu
      mode="horizontal"
      style={{ background: "inherit" }}
      selectedKeys={[selectedKey]}
      items={menuItems}
    />
  );
};
