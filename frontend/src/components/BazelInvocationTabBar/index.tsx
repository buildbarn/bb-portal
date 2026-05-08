import {
  AreaChartOutlined,
  BranchesOutlined,
  CodeOutlined,
  DatabaseOutlined,
  DeploymentUnitOutlined,
  ExperimentOutlined,
  FieldTimeOutlined,
  FileSearchOutlined,
  InfoCircleOutlined,
  LineChartOutlined,
  RadiusUprightOutlined,
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
  const { invocationID, actions, sourceControl, metrics, tags } = invocation;

  const hideActionStatisticsTab: boolean =
    metrics?.actionSummary === undefined || metrics?.actionSummary == null;
  const hideArtifactsTab: boolean =
    metrics?.artifactMetrics === undefined || metrics?.artifactMetrics == null;
  const hideMemoryTab: boolean =
    metrics?.memoryMetrics === undefined || metrics?.memoryMetrics == null;
  const hideSystemMetricsTab: boolean =
    (metrics?.timingMetrics === undefined || metrics?.timingMetrics == null) &&
    (metrics?.networkMetrics === undefined || metrics?.networkMetrics == null);
  const hideFailedActionsTab: boolean =
    actions === undefined || actions == null || actions.length === 0;
  const hideTargetsTab: boolean = !env.featureFlags?.bes?.pageTargets;
  const hideTestsTab: boolean = !env.featureFlags?.bes?.pageTests;
  const hideSourceControlTab: boolean =
    sourceControl === undefined ||
    sourceControl == null ||
    sourceControl.length === 0;
  const hideTagsTab: boolean = tags.totalCount === 0;

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
  if (!hideActionStatisticsTab)
    items.push({
      key: "action-statistics",
      icon: <LineChartOutlined />,
      label: (
        <Link
          to="/bazel-invocations/$invocationID/action-statistics"
          params={{ invocationID }}
        >
          Action Statistics
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
  if (!hideArtifactsTab)
    items.push({
      key: "artifacts",
      icon: <RadiusUprightOutlined />,
      label: (
        <Link
          to="/bazel-invocations/$invocationID/artifacts"
          params={{ invocationID }}
        >
          Artifacts
        </Link>
      ),
    });
  if (!hideMemoryTab)
    items.push({
      key: "memory",
      icon: <AreaChartOutlined />,
      label: (
        <Link
          to="/bazel-invocations/$invocationID/memory"
          params={{ invocationID }}
        >
          Memory
        </Link>
      ),
    });
  if (!hideSystemMetricsTab)
    items.push({
      key: "system-metrics",
      icon: <FieldTimeOutlined />,
      label: (
        <Link
          to="/bazel-invocations/$invocationID/system-metrics"
          params={{ invocationID }}
        >
          System Metrics
        </Link>
      ),
    });
  if (!hideTargetsTab)
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
  if (!hideTestsTab)
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
  if (!hideSourceControlTab)
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
  if (!hideTagsTab)
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
  if (!hideFailedActionsTab)
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
