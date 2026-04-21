import {
  AreaChartOutlined,
  BranchesOutlined,
  BuildOutlined,
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
import { Link } from "@tanstack/react-router";
import { Space, Tabs, Typography } from "antd";
import type { TabsProps } from "antd/lib";
import type React from "react";
import { useMemo, useState } from "react";
import PortalCard from "@/components/PortalCard";
import PortalDuration from "@/components/PortalDuration";
import type {
  BazelInvocationInfoFragment,
  RunnerCount,
} from "@/graphql/__generated__/graphql";
import themeStyles from "@/theme/theme.module.css";
import { commandLineDataToString } from "@/utils/commandLineDataToString";
import { env } from "@/utils/env";
import { parseGraphqlEdgeList } from "@/utils/parseGraphqlEdgeList";
import ActionStatisticsDisplay from "../ActionStatisticsDisplay";
import { ActionsTab } from "../ActionsTab";
import styles from "../AppBar/index.module.css";
import ArtifactsDataMetrics from "../Artifacts";
import BuildLogsDisplay from "../BuildLogsDisplay";
import CommandLineDisplay from "../CommandLine";
import InvocationOverviewDisplay from "../InvocationOverviewDisplay";
import { InvocationResultTag } from "../InvocationResultTag";
import { InvocationTagTab } from "../InvocationTagTab";
import { InvocationTargetsTab } from "../InvocationTargets/InvocationTargetsTab";
import MemoryMetricsDisplay from "../MemoryMetrics";
import ProfileDropdown from "../ProfileDropdown";
import SourceControlDisplay from "../SourceControlDisplay";
import SystemMetricsDisplay from "../SystemMetricsDisplay";
import { TestTab } from "../TestTab";
import UserStatusIndicator from "../UserStatusIndicator";

const DEFAULT_TAB_KEY = "BazelInvocationTabs-Overview";

const getTabItems = (
  invocationOverview: BazelInvocationInfoFragment,
): TabsProps["items"] => {
  const {
    actions,
    invocationID,
    instanceName,
    canonicalCommandLine,
    originalCommandLine,
    optionsParsed,
    bazelVersion,
    sourceControl,
    metrics,
    numFetches,
    configurations,
    hostname,
    tags,
  } = invocationOverview;

  var runnerMetrics: RunnerCount[] = [];
  metrics?.actionSummary?.runnerCount?.map((item: RunnerCount) =>
    runnerMetrics.push(item),
  );

  const tagList = parseGraphqlEdgeList(tags);

  const hideActionStatisticsTab: boolean =
    metrics?.actionSummary === undefined || metrics?.actionSummary == null;
  const hideLogsTab: boolean = false;
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
  const hideTagsTab: boolean = tagList.length === 0;

  const command = commandLineDataToString(originalCommandLine);

  const items: TabsProps["items"] = [];
  items.push({
    key: "BazelInvocationTabs-Overview",
    label: "Overview",
    icon: <InfoCircleOutlined />,
    children: (
      <Space direction="vertical" size="middle" className={themeStyles.space}>
        <PortalCard
          type="inner"
          icon={<FileSearchOutlined />}
          titleBits={["Invocation Overview"]}
        >
          <InvocationOverviewDisplay
            command={command}
            invocationId={invocationID}
            instanceName={instanceName.name}
            configurations={configurations || undefined}
            numFetches={numFetches ?? 0}
            startedAt={invocationOverview.startedAt}
            endedAt={invocationOverview.endedAt}
            hostname={hostname ?? ""}
            exitCodeName={invocationOverview.exitCodeName || undefined}
            connectionLastOpenAt={
              invocationOverview.connectionMetadata?.connectionLastOpenAt ||
              undefined
            }
            timeSinceLastConnectionMillis={
              invocationOverview.connectionMetadata
                ?.timeSinceLastConnectionMillis || undefined
            }
            bazelVersion={bazelVersion ?? ""}
          />
        </PortalCard>
      </Space>
    ),
  });
  if (!hideActionStatisticsTab)
    items.push({
      key: "BazelInvocationTabs-ActionStatistics",
      label: "Action Statistics",
      icon: <LineChartOutlined />,
      children: (
        <Space direction="vertical" size="middle" className={themeStyles.space}>
          <ActionStatisticsDisplay
            actionData={metrics?.actionSummary ?? undefined}
            runnerMetrics={runnerMetrics}
          />
        </Space>
      ),
    });
  if (!hideLogsTab)
    items.push({
      key: "BazelInvocationTabs-Logs",
      label: "Logs",
      icon: <FileSearchOutlined />,
      children: (
        <Space direction="vertical" size="middle" className={themeStyles.space}>
          <BuildLogsDisplay rawCommand={command} invocationId={invocationID} />
        </Space>
      ),
    });
  if (!hideArtifactsTab)
    items.push({
      key: "BazelInvocationTabs-Artifacts",
      label: "Artifacts",
      icon: <RadiusUprightOutlined />,
      children: (
        <Space direction="vertical" size="middle" className={themeStyles.space}>
          <ArtifactsDataMetrics
            artifactMetrics={metrics?.artifactMetrics ?? undefined}
          />
        </Space>
      ),
    });
  if (!hideMemoryTab)
    items.push({
      key: "BazelInvocationTabs-Memory",
      label: "Memory",
      icon: <AreaChartOutlined />,
      children: (
        <Space direction="vertical" size="middle" className={themeStyles.space}>
          <MemoryMetricsDisplay
            memoryMetrics={metrics?.memoryMetrics ?? undefined}
          />
        </Space>
      ),
    });
  if (!hideSystemMetricsTab)
    items.push({
      key: "BazelInvocationTabs-SystemMetrics",
      label: "System Metrics",
      icon: <FieldTimeOutlined />,
      children: (
        <Space direction="vertical" size="middle" className={themeStyles.space}>
          <SystemMetricsDisplay
            timingMetrics={metrics?.timingMetrics ?? undefined}
            systemNetworkStats={
              metrics?.networkMetrics?.systemNetworkStats ?? undefined
            }
          />
        </Space>
      ),
    });
  if (!hideTargetsTab)
    items.push({
      key: "BazelInvocationTabs-Targets",
      label: "Targets",
      icon: <DeploymentUnitOutlined />,
      children: (
        <Space direction="vertical" size="middle" className={themeStyles.space}>
          <InvocationTargetsTab
            invocationId={invocationID}
            targetMetrics={metrics?.targetMetrics ?? undefined}
          />
        </Space>
      ),
    });
  if (!hideTestsTab)
    items.push({
      key: "BazelInvocationTabs-Tests",
      label: "Tests",
      icon: <ExperimentOutlined />,
      children: (
        <Space direction="vertical" size="middle" className={themeStyles.space}>
          <TestTab invocationId={invocationID} />
        </Space>
      ),
    });
  items.push({
    key: "BazelInvocationTabs-CommandLine",
    label: "Command Line",
    icon: <CodeOutlined />,
    children: (
      <Space direction="vertical" size="middle" className={themeStyles.space}>
        <CommandLineDisplay
          parsedOptions={optionsParsed}
          rawCommand={command}
          canonicalCommandLine={canonicalCommandLine}
        />
      </Space>
    ),
  });
  if (!hideSourceControlTab)
    items.push({
      key: "BazelInvocationTabs-SourceControl",
      label: "Source Control",
      icon: <BranchesOutlined />,
      children: (
        <Space direction="vertical" size="middle" className={themeStyles.space}>
          <SourceControlDisplay sourceControlData={sourceControl} />
        </Space>
      ),
    });
  if (!hideTagsTab)
    items.push({
      key: "BazelInvocationTabs-Tags",
      label: "Tags",
      icon: <TagsOutlined />,
      children: (
        <Space direction="vertical" size="middle" className={themeStyles.space}>
          <InvocationTagTab tags={tagList} />
        </Space>
      ),
    });
  if (!hideFailedActionsTab)
    items.push({
      key: "BazelInvocationTabs-Actions",
      label: "Failed Actions",
      icon: <DatabaseOutlined />,
      children: (
        <Space direction="vertical" size="middle" className={themeStyles.space}>
          <ActionsTab instanceName={instanceName.name} actions={actions} />
        </Space>
      ),
    });

  return items;
};

const getTitleBits = (
  invocationOverview: BazelInvocationInfoFragment,
): React.ReactNode[] => {
  const { invocationID, authenticatedUser, username } = invocationOverview;

  const titleBits: React.ReactNode[] = [];
  if (username && username !== "")
    titleBits.push(
      <span key="label">
        User:{" "}
        <Typography.Text type="secondary" className={styles.normalWeight}>
          <UserStatusIndicator
            authenticatedUser={authenticatedUser}
            username={username}
          />
        </Typography.Text>
      </span>,
    );
  if (invocationID && invocationID !== "")
    titleBits.push(
      <span key="label">
        Invocation ID:{" "}
        <Typography.Text
          type="secondary"
          className={styles.normalWeight}
          copyable={{ text: invocationID ?? "Copy" }}
        >
          {invocationID}
        </Typography.Text>{" "}
      </span>,
    );
  titleBits.push(
    <InvocationResultTag
      key="result"
      exitCodeName={invocationOverview.exitCodeName || undefined}
      timeSinceLastConnectionMillis={
        invocationOverview.connectionMetadata?.timeSinceLastConnectionMillis ||
        undefined
      }
    />,
  );
  return titleBits;
};

const getExtraBits = (
  invocationOverview: BazelInvocationInfoFragment,
): React.ReactNode[] => {
  const { invocationID, instanceName, build, profile } = invocationOverview;

  const extraBits: React.ReactNode[] = [];
  extraBits.push(
    <PortalDuration
      key="duration"
      from={invocationOverview.startedAt || undefined}
      to={
        invocationOverview.endedAt
          ? invocationOverview.endedAt
          : invocationOverview.connectionMetadata?.connectionLastOpenAt
      }
      includeIcon
      includePopover
      formatConfig={{ smallestUnit: "s" }}
    />,
  );
  if (profile)
    extraBits.push(
      <ProfileDropdown
        instanceName={instanceName.name}
        profile={profile}
        invocationID={invocationID}
      />,
    );
  if (build?.buildUUID) {
    extraBits.unshift(
      <span key="build">
        Build{" "}
        <Link to={`/builds/$buildUUID`} params={{ buildUUID: build.buildUUID }}>
          {build.buildUUID}
        </Link>
      </span>,
    );
  }
  return extraBits;
};

interface Props {
  invocationOverview: BazelInvocationInfoFragment;
}

const BazelInvocation: React.FC<Props> = ({ invocationOverview }) => {
  const [activeKey, setActiveKey] = useState(
    localStorage.getItem("bazelInvocationViewActiveTabKey") ??
      "BazelInvocationTabs-Overview",
  );

  const onTabChange = (key: string) => {
    setActiveKey(key);
    localStorage.setItem("bazelInvocationViewActiveTabKey", key);
  };

  const titleBits = useMemo(
    () => getTitleBits(invocationOverview),
    [invocationOverview],
  );

  const extraBits = useMemo(
    () => getExtraBits(invocationOverview),
    [invocationOverview],
  );

  const items = useMemo(
    () => getTabItems(invocationOverview),
    [invocationOverview],
  );

  function checkIfNotHidden(key: string) {
    const hidden = items?.findIndex((x) => x.key === key) === -1;
    return hidden ? DEFAULT_TAB_KEY : key;
  }

  return (
    <PortalCard
      icon={<BuildOutlined />}
      titleBits={titleBits}
      extraBits={extraBits}
    >
      <Tabs
        items={items}
        activeKey={checkIfNotHidden(activeKey)}
        onChange={onTabChange}
        defaultActiveKey={DEFAULT_TAB_KEY}
      />
    </PortalCard>
  );
};

export default BazelInvocation;
