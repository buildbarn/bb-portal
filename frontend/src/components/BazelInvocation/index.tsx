import BuildStepResultTag, {
  BuildStepResultEnum,
} from "@/components/BuildStepResultTag";
import Link from "@/components/Link";
import ArtifactsDataMetrics from "../Artifacts";
import MemoryMetricsDisplay from "../MemoryMetrics";
import SystemMetricsDisplay from "../SystemMetricsDisplay";
import TestMetricsDisplay from "../TestsMetrics";
import CommandLineDisplay from "../CommandLine";
import SourceControlDisplay from "../SourceControlDisplay";
import InvocationOverviewDisplay from "../InvocationOverviewDisplay";
import BuildProblems from "../Problems";
import PortalCard from "@/components/PortalCard";
import PortalDuration from "@/components/PortalDuration";
import {
  BazelInvocationInfoFragment,
  RunnerCount
} from "@/graphql/__generated__/graphql";
import themeStyles from "@/theme/theme.module.css";
import { FeatureType, isFeatureEnabled } from "@/utils/isFeatureEnabled";
import {
  AreaChartOutlined,
  BranchesOutlined,
  BuildOutlined,
  CodeOutlined,
  DeploymentUnitOutlined,
  ExclamationCircleOutlined,
  ExperimentOutlined,
  FieldTimeOutlined,
  FileSearchOutlined,
  InfoCircleOutlined, LineChartOutlined,
  RadiusUprightOutlined
} from "@ant-design/icons";
import { Space, Tabs, Typography } from "antd";
import type { TabsProps } from "antd/lib";
import React, { useMemo, useState } from "react";
import ActionStatisticsDisplay from "../ActionStatisticsDisplay";
import styles from "../AppBar/index.module.css";
import BuildLogsDisplay from "../BuildLogsDisplay";
import ProfileDropdown from "../ProfileDropdown";
import { InvocationTargetsTab } from "../InvocationTargets/InvocationTargetsTab";
import UserStatusIndicator from "../UserStatusIndicator";

const DEFAULT_TAB_KEY = "BazelInvocationTabs-Overview";

const getTabItems = (invocationOverview: BazelInvocationInfoFragment): TabsProps["items"] => {
  const {
    invocationID,
    instanceName,
    state,
    bazelCommand,
    bazelVersion,
    sourceControl,
    metrics,
    testCollection,
    numFetches,
    cpu,
    configurationMnemonic,
    stepLabel,
    hostname,
    isCiWorker,
  } = invocationOverview;

  var runnerMetrics: RunnerCount[] = [];
  metrics?.actionSummary?.runnerCount?.map((item: RunnerCount) =>
    runnerMetrics.push(item)
  );

  const hideActionsTab: boolean = metrics?.actionSummary == undefined || metrics?.actionSummary == null;
  const hideLogsTab: boolean = false;
  const hideArtifactsTab: boolean = metrics?.artifactMetrics == undefined || metrics?.artifactMetrics == null;
  const hideMemoryTab: boolean = metrics?.memoryMetrics == undefined || metrics?.memoryMetrics == null;
  const hideSystemMetricsTab: boolean =
    (metrics?.timingMetrics == undefined || metrics?.timingMetrics == null)
    && (metrics?.networkMetrics == undefined || metrics?.networkMetrics == null);
  const hideTargetsTab: boolean = !isFeatureEnabled(FeatureType.BES_PAGE_TARGETS);
  const hideTestsTab: boolean = (testCollection == undefined || testCollection == null || testCollection.length == 0)
  const hideSourceControlTab: boolean = sourceControl == undefined || sourceControl == null;
  const hideProblemsTab: boolean = state.exitCode?.name == "SUCCESS";

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
            command={[
              bazelCommand.executable,
              bazelCommand.command,
              bazelCommand.residual,
              bazelCommand.explicitCmdLine,
            ].join(" ").trim()}
            cpu={cpu ?? ""}
            invocationId={invocationID}
            instanceName={instanceName.name}
            configuration={configurationMnemonic ?? ""}
            numFetches={numFetches ?? 0}
            startedAt={invocationOverview.startedAt}
            endedAt={invocationOverview.endedAt}
            hostname={hostname ?? ""}
            isCiWorker={isCiWorker ?? false}
            stepLabel={stepLabel ?? ""}
            status={state.exitCode?.name ?? ""}
            bazelVersion={bazelVersion ?? ""}
          />
        </PortalCard>
      </Space>
    ),
  });
  if (!hideActionsTab) items.push({
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
  if (!hideLogsTab) items.push({
    key: "BazelInvocationTabs-Logs",
    label: "Logs",
    icon: <FileSearchOutlined />,
    children: (
      <Space direction="vertical" size="middle" className={themeStyles.space}>
        <BuildLogsDisplay invocationId={invocationID} />
      </Space>
    ),
  });
  if (!hideArtifactsTab) items.push({
    key: "BazelInvocationTabs-Artifacts",
    label: "Artifacts",
    icon: <RadiusUprightOutlined />,
    children: (
      <Space direction="vertical" size="middle" className={themeStyles.space}>
        <ArtifactsDataMetrics artifactMetrics={metrics?.artifactMetrics ?? undefined} />
      </Space>
    ),
  });
  if (!hideMemoryTab) items.push({
    key: "BazelInvocationTabs-Memory",
    label: "Memory",
    icon: <AreaChartOutlined />,
    children: (
      <Space direction="vertical" size="middle" className={themeStyles.space}>
        <MemoryMetricsDisplay memoryMetrics={metrics?.memoryMetrics ?? undefined} />
      </Space>
    ),
  });
  if (!hideSystemMetricsTab) items.push({
    key: "BazelInvocationTabs-SystemMetrics",
    label: "System Metrics",
    icon: <FieldTimeOutlined />,
    children: (
      <Space direction="vertical" size="middle" className={themeStyles.space}>
        <SystemMetricsDisplay
          timingMetrics={metrics?.timingMetrics ?? undefined}
          systemNetworkStats={metrics?.networkMetrics?.systemNetworkStats ?? undefined}
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
  if (!hideTestsTab) items.push({
    key: "BazelInvocationTabs-Tests",
    label: "Tests",
    icon: <ExperimentOutlined />,
    children: (
      <Space direction="vertical" size="middle" className={themeStyles.space}>
        <TestMetricsDisplay
          testMetrics={testCollection ?? undefined}
          invocationId={invocationID}
        />
      </Space>
    ),
  });
  items.push({
    key: "BazelInvocationTabs-CommandLine",
    label: "Command Line",
    icon: <CodeOutlined />,
    children: (
      <Space direction="vertical" size="middle" className={themeStyles.space}>
        <CommandLineDisplay commandLineData={bazelCommand ?? undefined} />
      </Space>
    ),
  });
  if (!hideSourceControlTab) items.push({
    key: "BazelInvocationTabs-SourceControl",
    label: "Source Control",
    icon: <BranchesOutlined />,
    children: (
      <Space direction="vertical" size="middle" className={themeStyles.space}>
        <SourceControlDisplay
          sourceControlData={sourceControl}
          stepLabel={stepLabel}
        />
      </Space>
    ),
  });
  if (!hideProblemsTab) items.push({
    key: "BazelInvocationTabs-Problems",
    label: "Problems",
    icon: <ExclamationCircleOutlined />,
    children: (
      <Space direction="vertical" size="middle" className={themeStyles.space}>
        <BuildProblems
          invocationId={invocationID}
          instanceName={instanceName.name}
        />
      </Space>
    ),
  });

  return items;
}

const getTitleBits = (invocationOverview: BazelInvocationInfoFragment): React.ReactNode[] => {
  const {
    invocationID,
    state,
    authenticatedUser,
    user,
  } = invocationOverview;

  const titleBits: React.ReactNode[] = []
  if (user?.LDAP && user?.LDAP != "") titleBits.push(
    <span key="label">
      User:{" "}
      <Typography.Text type="secondary" className={styles.normalWeight}>
        <UserStatusIndicator authenticatedUser={authenticatedUser} user={user} />
      </Typography.Text>
    </span>
  );
  if (invocationID && invocationID != "") titleBits.push(
    <span key="label">
      Invocation ID:{" "}
      <Typography.Text type="secondary" className={styles.normalWeight} copyable={{ text: invocationID ?? "Copy" }}>
        {invocationID}
      </Typography.Text>{" "}
    </span>
  );
  if (state.exitCode?.name) titleBits.push(
    <BuildStepResultTag
      key="result"
      result={state.exitCode?.name as BuildStepResultEnum}
    />
  );
  return titleBits;
}

const getExtraBits = (invocationOverview: BazelInvocationInfoFragment, isNestedWithinBuildCard?: boolean): React.ReactNode[] => {
  const {
    invocationID,
    instanceName,
    build,
    profile,
  } = invocationOverview;

  const extraBits: React.ReactNode[] = [];
  extraBits.push(
    <PortalDuration
      key="duration"
      from={invocationOverview.startedAt}
      to={invocationOverview.endedAt}
      includeIcon
      includePopover
    />
  )
  if (profile) extraBits.push(
    <ProfileDropdown
      instanceName={instanceName.name}
      profile={profile}
      invocationID={invocationID}
    />,
  );
  if (!isNestedWithinBuildCard && build?.buildUUID) {
    extraBits.unshift(
      <span key="build">
        Build <Link href={`/builds/${build.buildUUID}`}>{build.buildUUID}</Link>
      </span>
    );
  }
  return extraBits;
}


const BazelInvocation: React.FC<{
  invocationOverview: BazelInvocationInfoFragment;
  isNestedWithinBuildCard?: boolean;
  collapsed?: boolean;
}> = ({ invocationOverview, isNestedWithinBuildCard }) => {
  const [activeKey, setActiveKey] = useState(
    localStorage.getItem("bazelInvocationViewActiveTabKey") ??
    "BazelInvocationTabs-Overview"
  );

  const onTabChange = (key: string) => {
    setActiveKey(key);
    localStorage.setItem("bazelInvocationViewActiveTabKey", key);
  };

  const titleBits = useMemo(
    () => getTitleBits(invocationOverview),
    [invocationOverview]
  );

  const extraBits = useMemo(
    () => getExtraBits(invocationOverview, isNestedWithinBuildCard),
    [invocationOverview, isNestedWithinBuildCard]
  );

  const items = useMemo(
    () => getTabItems(invocationOverview),
    [invocationOverview, isNestedWithinBuildCard]
  );

  function checkIfNotHidden(key: string) {
    const hidden = items?.findIndex((x) => x.key == key) == -1
    return hidden ? DEFAULT_TAB_KEY : key;
  }

  return (
    <PortalCard
      type={isNestedWithinBuildCard ? "inner" : undefined}
      icon={<BuildOutlined />}
      titleBits={titleBits}
      extraBits={extraBits}
    >
      <Tabs
        items={items}
        activeKey={checkIfNotHidden(activeKey)}
        onChange={onTabChange}
        defaultActiveKey={DEFAULT_TAB_KEY}
        // TODO(isakstenstrom): Remove this when we can fetch partial logs.
        destroyInactiveTabPane={true}
      />
    </PortalCard>
  );
};

export default BazelInvocation;
