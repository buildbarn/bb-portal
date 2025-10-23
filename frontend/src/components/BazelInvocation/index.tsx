import {
  ActionSummary,
  ArtifactMetrics,
  BazelInvocationInfoFragment,
  RunnerCount,
  TargetMetrics,
  MemoryMetrics,
  TimingMetrics,
  NetworkMetrics,
  TestCollection,
  TargetPair,
  BuildGraphMetrics,
  SystemNetworkStats,
} from "@/graphql/__generated__/graphql";
import styles from "../AppBar/index.module.css";
import React, { useMemo, useState } from "react";
import PortalDuration from "@/components/PortalDuration";
import PortalCard from "@/components/PortalCard";
import { Space, Tabs, Tooltip, Typography } from "antd";
import type { TabsProps } from "antd/lib";
import {
  BuildOutlined,
  FileSearchOutlined, ExclamationCircleOutlined, DeploymentUnitOutlined,
  ExperimentOutlined,
  RadiusUprightOutlined,
  AreaChartOutlined,
  FieldTimeOutlined,
  WifiOutlined, CodeOutlined,
  BranchesOutlined,
  InfoCircleOutlined, LineChartOutlined
} from "@ant-design/icons";
import themeStyles from "@/theme/theme.module.css";
import BuildStepResultTag, {
  BuildStepResultEnum,
} from "@/components/BuildStepResultTag";
import DownloadButton from "@/components/DownloadButton";
import Link from "@/components/Link";
import LogViewer from "../LogViewer";
import TargetMetricsDisplay from "../TargetMetrics";
import ArtifactsDataMetrics from "../Artifacts";
import MemoryMetricsDisplay from "../MemoryMetrics";
import SystemMetricsDisplay from "../SystemMetricsDisplay";
import TestMetricsDisplay from "../TestsMetrics";
import CommandLineDisplay from "../CommandLine";
import SourceControlDisplay from "../SourceControlDisplay";
import InvocationOverviewDisplay from "../InvocationOverviewDisplay";
import BuildProblems from "../Problems";
import ActionStatisticsDisplay from "../ActionStatisticsDisplay";
import ProfileDropdown from "../ProfileDropdown";
import ansiRegex from 'ansi-regex';

const ansiEscapeRegex = ansiRegex();

const BazelInvocation: React.FC<{
  invocationOverview: BazelInvocationInfoFragment;
  isNestedWithinBuildCard?: boolean;
  collapsed?: boolean;
}> = ({ invocationOverview, isNestedWithinBuildCard }) => {
  const {
    invocationID,
    instanceName,
    build,
    state,
    bazelCommand,
    profile,
    sourceControl,
    user,
    metrics,
    testCollection,
    targets,
    numFetches,
    cpu,
    configurationMnemonic,
    stepLabel,
    hostname,
    isCiWorker,
    buildLogs,

    //relatedFiles,
  } = invocationOverview;

  const logDownloadUrl = useMemo(
    () => buildLogs ? `data:text/plain;charset=utf-8,${encodeURIComponent(buildLogs.replace(ansiEscapeRegex, ""))}` : undefined,
    [buildLogs]
  );

  //data for runner metrics
  var runnerMetrics: RunnerCount[] = [];
  metrics?.actionSummary?.runnerCount?.map((item: RunnerCount) =>
    runnerMetrics.push(item)
  );

  //data for ac metrics
  var acMetrics: ActionSummary | undefined =
    metrics?.actionSummary ?? undefined;

  //artifact metrics
  var artifactMetrics: ArtifactMetrics | undefined =
    metrics?.artifactMetrics ?? undefined;

  //memory metrics
  var memoryMetrics: MemoryMetrics | undefined =
    metrics?.memoryMetrics ?? undefined;

  //build graph metrics
  var buildGraphMetrics: BuildGraphMetrics | undefined =
    metrics?.buildGraphMetrics ?? undefined;

  //timing metrics
  var timingMetrics: TimingMetrics | undefined =
    metrics?.timingMetrics ?? undefined;

  //netowrk metrics
  var systemNetworkStats: SystemNetworkStats | undefined =
    metrics?.networkMetrics?.systemNetworkStats ?? undefined;

  //test data
  var testCollections: TestCollection[] | undefined | null = testCollection;

  //data for target metrics
  var targetMetrics: TargetMetrics | undefined | null =
    metrics?.targetMetrics ?? undefined;
  var targetData: TargetPair[] | undefined | null = targets;
  var targetTimes: Map<string, number> = new Map<string, number>();
  targetData?.map((x) => {
    targetTimes.set(x.label ?? "", x.durationInMs ?? 0);
  });

  //build the title
  let { exitCode } = state;
  exitCode = exitCode ?? null;
  const titleBits: React.ReactNode[] = [
    <span key="label">
      User:{" "}
      <Typography.Text type="secondary" className={styles.normalWeight}>
        {user?.LDAP}
      </Typography.Text>
    </span>,
  ];
  titleBits.push(
    <span key="label">
      Invocation ID:{" "}
      <Typography.Text type="secondary" className={styles.normalWeight}>
        {invocationID}
      </Typography.Text>{" "}
    </span>
  );
  titleBits.push(
    <span className={styles.copyIcon}>
      <Typography.Text
        copyable={{ text: invocationID ?? "Copy" }}
      ></Typography.Text>
    </span>
  );
  if (exitCode?.name) {
    titleBits.push(
      <BuildStepResultTag
        key="result"
        result={exitCode?.name as BuildStepResultEnum}
      />
    );
  }

  const extraBits: React.ReactNode[] = [
    <PortalDuration
      key="duration"
      from={invocationOverview.startedAt}
      to={invocationOverview.endedAt}
      includeIcon
      includePopover
    />,
  ];

  if (profile) {
    extraBits.push(
      <ProfileDropdown
        instanceName={instanceName ?? undefined}
        profile={profile}
        invocationID={invocationID}
      />,
    );
  }

  if (!isNestedWithinBuildCard && build?.buildUUID) {
    extraBits.unshift(
      <span key="build">
        Build <Link href={`/builds/${build.buildUUID}`}>{build.buildUUID}</Link>
      </span>
    );
  }

  const hideTestsTab: boolean = (testCollection?.length ?? 0) == 0;
  const hideTargetsTab: boolean = (targetData?.length ?? 0) == 0 ? true : false;
  const hideSourceControlTab: boolean =
    sourceControl?.runID == undefined ||
      sourceControl.runID == null ||
      sourceControl.runID == ""
      ? true
      : false;
  const hideLogsTab: boolean = false;
  const hideMemoryTab: boolean =
    (memoryMetrics?.peakPostGcHeapSize ?? 0) == 0 &&
    (memoryMetrics?.peakPostGcHeapSize ?? 0) == 0 &&
    (memoryMetrics?.usedHeapSizePostBuild ?? 0) == 0;
  const hideProblemsTab: boolean = exitCode?.name == "SUCCESS";
  const hideArtifactsTab: boolean =
    (artifactMetrics?.outputArtifactsSeen?.count ?? 0) == 0 &&
    (artifactMetrics?.sourceArtifactsRead?.count ?? 0) == 0 &&
    (artifactMetrics?.outputArtifactsFromActionCache?.count ?? 0) == 0 &&
    (artifactMetrics?.topLevelArtifacts?.count ?? 0) == 0;
  const hideActionsTab: boolean =
    (acMetrics?.actionsExecuted == 0) &&
    (acMetrics?.actionCacheStatistics?.hits == 0) &&
    (acMetrics?.actionCacheStatistics?.misses == 0);
  const hideSystemMetricsTab: boolean =
    (timingMetrics == undefined &&
      systemNetworkStats == undefined) || (
      timingMetrics?.wallTimeInMs == 0 &&
      timingMetrics?.executionPhaseTimeInMs == 0 &&
      timingMetrics?.analysisPhaseTimeInMs == 0 &&
      timingMetrics?.cpuTimeInMs == 0 &&
      timingMetrics?.actionsExecutionStartInMs == 0 &&
      buildGraphMetrics?.actionCount == 0 &&
      buildGraphMetrics.actionLookupValueCount == 0 &&
      buildGraphMetrics.actionCountNotIncludingAspects == 0 &&
      buildGraphMetrics.inputFileConfiguredTargetCount == 0 &&
      buildGraphMetrics.outputArtifactCount == 0 &&
      buildGraphMetrics.postInvocationSkyframeNodeCount == 0 &&
      buildGraphMetrics.outputFileConfiguredTargetCount == 0 &&
      systemNetworkStats?.bytesRecv == 0 &&
      systemNetworkStats?.bytesSent == 0 &&
      systemNetworkStats?.packetsRecv == 0 &&
      systemNetworkStats?.packetsSent == 0 &&
      systemNetworkStats?.peakBytesRecvPerSec == 0 &&
      systemNetworkStats?.peakBytesSentPerSec == 0 &&
      systemNetworkStats?.peakPacketsRecvPerSec == 0 &&
      systemNetworkStats?.peakPacketsSentPerSec == 0
    );

  interface TabShowHideDisplay {
    hide: boolean;
    key: string;
  }

  const showHideTabs: TabShowHideDisplay[] = [
    { key: "BazelInvocationTabs-Tests", hide: hideTestsTab },
    { key: "BazelInvocationTabs-Targets", hide: hideTargetsTab },
    { key: "BazelInvocationTabs-SourceControl", hide: hideSourceControlTab },
    { key: "BazelInvocationTabs-Logs", hide: hideLogsTab },
    { key: "BazelInvocationTabs-Memory", hide: hideMemoryTab },
    { key: "BazelInvocationTabs-Problems", hide: hideProblemsTab },
    { key: "BazelInvocationTabs-Artifacts", hide: hideArtifactsTab },
    { key: "BazelInvocationTabs-ActionStatistics", hide: hideActionsTab },
    { key: "BazelInvocationTabs-SystemMetrics", hide: hideSystemMetricsTab },
  ];

  const [activeKey, setActiveKey] = useState(
    localStorage.getItem("bazelInvocationViewActiveTabKey") ??
      "BazelInvocationTabs-Overview"
  );
  function checkIfNotHidden(key: string) {
    var hidden: boolean =
      showHideTabs.filter((x) => x.key == key).at(0)?.hide ?? false;
    return hidden ? "BazelInvocationTabs-Overview" : key;
  }
  const onTabChange = (key: string) => {
    setActiveKey(key);
    localStorage.setItem("bazelInvocationViewActiveTabKey", key);
  };

  //tabs
  var items: TabsProps["items"] = [
    {
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
              ].join(" ")}
              targets={targetTimes.size}
              cpu={cpu ?? ""}
              user={user?.LDAP ?? ""}
              invocationId={invocationID}
              instanceName={instanceName ?? undefined}
              configuration={configurationMnemonic ?? ""}
              numFetches={numFetches ?? 0}
              startedAt={invocationOverview.startedAt}
              endedAt={invocationOverview.endedAt}
              hostname={hostname ?? ""}
              isCiWorker={isCiWorker ?? false}
              stepLabel={stepLabel ?? ""}
              status={state.exitCode?.name ?? ""}
            />
          </PortalCard>
        </Space>
      ),
    },
    {
      key: "BazelInvocationTabs-ActionStatistics",
      label: "Action Statistics",
      icon: <LineChartOutlined />,
      children: (
        <Space direction="vertical" size="middle" className={themeStyles.space}>
          <ActionStatisticsDisplay
            actionData={acMetrics}
            runnerMetrics={runnerMetrics}
          />
        </Space>
      ),
    },
    {
      key: "BazelInvocationTabs-Logs",
      label: "Logs",
      icon: <FileSearchOutlined />,
      children: (
        <Space direction="vertical" size="middle" className={themeStyles.space}>
          <PortalCard
            type="inner"
            icon={<FileSearchOutlined />}
            titleBits={["Raw Build Logs"]}
            extraBits={[
              <Tooltip title="Bazel emits logs in ANSI format a screen at a time.  They are presented here concatenated for your convenience.">
                <ExclamationCircleOutlined />
              </Tooltip>,
              logDownloadUrl && (
                <DownloadButton
                  enabled={true}
                  buttonLabel="Download Log"
                  fileName="log.txt"
                  url={logDownloadUrl}
                />
              ),
            ]}
          >
            <LogViewer log={buildLogs} />
          </PortalCard>
        </Space>
      ),
    },
    {
      key: "BazelInvocationTabs-Artifacts",
      label: "Artifacts",
      icon: <RadiusUprightOutlined />,
      children: (
        <Space direction="vertical" size="middle" className={themeStyles.space}>
          <ArtifactsDataMetrics artifactMetrics={artifactMetrics} />
        </Space>
      ),
    },
    {
      key: "BazelInvocationTabs-Memory",
      label: "Memory",
      icon: <AreaChartOutlined />,
      children: (
        <Space direction="vertical" size="middle" className={themeStyles.space}>
          <MemoryMetricsDisplay memoryMetrics={memoryMetrics} />
        </Space>
      ),
    },
    {
      key: "BazelInvocationTabs-SystemMetrics",
      label: "System Metrics",
      icon: <FieldTimeOutlined />,
      children: (
        <Space direction="vertical" size="middle" className={themeStyles.space}>
          <SystemMetricsDisplay
            timingMetrics={timingMetrics}
            buildGraphMetrics={buildGraphMetrics}
            systemNetworkStats={systemNetworkStats}
          />
        </Space>
      ),
    },
    {
      key: "BazelInvocationTabs-Targets",
      label: "Targets",
      icon: <DeploymentUnitOutlined />,
      children: (
        <Space direction="vertical" size="middle" className={themeStyles.space}>
          <TargetMetricsDisplay
            targetMetrics={targetMetrics}
            targetData={targetData}
          />
        </Space>
      ),
    },
    {
      key: "BazelInvocationTabs-Tests",
      label: "Tests",
      icon: <ExperimentOutlined />,
      children: (
        <Space direction="vertical" size="middle" className={themeStyles.space}>
          <TestMetricsDisplay
            testMetrics={testCollections}
            targetTimes={targetTimes}
          />
        </Space>
      ),
    },
    {
      key: "BazelInvocationTabs-CommandLine",
      label: "Command Line",
      icon: <CodeOutlined />,
      children: (
        <Space direction="vertical" size="middle" className={themeStyles.space}>
          <CommandLineDisplay commandLineData={bazelCommand ?? undefined} />
        </Space>
      ),
    },
    {
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
    },
    {
      key: "BazelInvocationTabs-Problems",
      label: "Problems",
      icon: <ExclamationCircleOutlined />,
      children: (
        <Space direction="vertical" size="middle" className={themeStyles.space}>
          <BuildProblems
            invocationId={invocationID}
            instanceName={instanceName ?? undefined}
            onTabChange={onTabChange}
          />
        </Space>
      ),
    },
  ];

  //show/hide tabs

  for (var i in showHideTabs) {
    var tab = showHideTabs[i];
    if (tab.hide == true) {
      var idx = items.findIndex((x, _) => x.key == tab.key);
      if (idx > -1) {
        items.splice(idx, 1);
      }
    }
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
        defaultActiveKey="BazelInvocationTabs-Overview"
      />
    </PortalCard>
  );
};

export default BazelInvocation;
