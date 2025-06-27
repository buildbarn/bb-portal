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
} from "@/graphql/__generated__/graphql";
import styles from "../AppBar/index.module.css";
import React, { useState } from "react";
import PortalDuration from "@/components/PortalDuration";
import PortalCard from "@/components/PortalCard";
import { Space, Tabs, Tooltip, Typography } from "antd";
import type { TabsProps } from "antd/lib";
import {
  BuildOutlined,
  FileSearchOutlined,
  ClusterOutlined,
  ExclamationCircleOutlined,
  NodeCollapseOutlined,
  DeploymentUnitOutlined,
  ExperimentOutlined,
  RadiusUprightOutlined,
  AreaChartOutlined,
  FieldTimeOutlined,
  WifiOutlined,
  HddOutlined,
  CodeOutlined,
  BranchesOutlined,
  InfoCircleOutlined
} from "@ant-design/icons";
import themeStyles from "@/theme/theme.module.css";
import BuildStepResultTag, {
  BuildStepResultEnum,
} from "@/components/BuildStepResultTag";
import DownloadButton from "@/components/DownloadButton";
import Link from "@/components/Link";
import LogViewer from "../LogViewer";
import RunnerMetrics from "../RunnerMetrics";
import AcMetrics from "../ActionCacheMetrics";
import TargetMetricsDisplay from "../TargetMetrics";
import ActionDataMetrics from "../ActionDataMetrics";
import ArtifactsDataMetrics from "../Artifacts";
import MemoryMetricsDisplay from "../MemoryMetrics";
import TimingMetricsDisplay from "../TimingMetrics";
import NetworkMetricsDisplay from "../NetworkMetrics";
import TestMetricsDisplay from "../TestsMetrics";
import CommandLineDisplay from "../CommandLine";
import SourceControlDisplay from "../SourceControlDisplay";
import InvocationOverviewDisplay from "../InvocationOverviewDisplay";
import BuildProblems from "../Problems";
import { generateFileUrl } from "@/utils/urlGenerator";
import { DigestFunction_Value } from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";

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
  var networkMetrics: NetworkMetrics | undefined =
    metrics?.networkMetrics ?? undefined;
  const bytesRecv = networkMetrics?.systemNetworkStats?.bytesRecv ?? 0;
  const bytesSent = networkMetrics?.systemNetworkStats?.bytesSent ?? 0;

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
    const url = generateFileUrl(
      instanceName ?? undefined, 
      DigestFunction_Value.SHA256, 
      {
        hash: profile.digest,
        sizeBytes: profile.sizeInBytes.toString()
      },
      profile.name
    )
    extraBits.push(
      <DownloadButton
        url={url}
        fileName="profile"
        buttonLabel="Profile"
        enabled={true}
      />
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
  const hideNetworkTab: boolean = bytesRecv == 0 && bytesSent == 0;
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
  const hideActionsDataTab: boolean = acMetrics?.actionsExecuted == 0;
  const hideActionCacheTab: boolean =
    acMetrics?.actionCacheStatistics?.hits == 0 &&
    acMetrics?.actionCacheStatistics?.misses == 0;
  const hideRunnersTab: boolean = runnerMetrics.length == 0;
  const hideTimingTab: boolean =
    timingMetrics?.wallTimeInMs == 0 &&
    timingMetrics.executionPhaseTimeInMs == 0 &&
    timingMetrics.analysisPhaseTimeInMs == 0 &&
    timingMetrics.cpuTimeInMs == 0 &&
    timingMetrics.actionsExecutionStartInMs == 0 &&
    buildGraphMetrics?.actionCount == 0 &&
    buildGraphMetrics.actionLookupValueCount == 0 &&
    buildGraphMetrics.actionCountNotIncludingAspects == 0 &&
    buildGraphMetrics.inputFileConfiguredTargetCount == 0 &&
    buildGraphMetrics.outputArtifactCount == 0 &&
    buildGraphMetrics.postInvocationSkyframeNodeCount == 0 &&
    buildGraphMetrics.outputFileConfiguredTargetCount == 0;

  interface TabShowHideDisplay {
    hide: boolean;
    key: string;
  }

  const showHideTabs: TabShowHideDisplay[] = [
    { key: "BazelInvocationTabs-Tests", hide: hideTestsTab },
    { key: "BazelInvocationTabs-Targets", hide: hideTargetsTab },
    { key: "BazelInvocationTabs-Network", hide: hideNetworkTab },
    { key: "BazelInvocationTabs-SourceControl", hide: hideSourceControlTab },
    { key: "BazelInvocationTabs-Logs", hide: hideLogsTab },
    { key: "BazelInvocationTabs-Memory", hide: hideMemoryTab },
    { key: "BazelInvocationTabs-Problems", hide: hideProblemsTab },
    { key: "BazelInvocationTabs-Artifacts", hide: hideArtifactsTab },
    { key: "BazelInvocationTabs-ActionsData", hide: hideActionsDataTab },
    { key: "BazelInvocationTabs-ActionCache", hide: hideActionCacheTab },
    { key: "BazelInvocationTabs-Timing", hide: hideTimingTab },
    { key: "BazelInvocationTabs-Runners", hide: hideRunnersTab },
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
            ]}
          >
            <LogViewer log={buildLogs} />
          </PortalCard>
        </Space>
      ),
    },
    {
      key: "BazelInvocationTabs-Runners",
      label: "Runners",
      icon: <ClusterOutlined />,
      children: (
        <Space direction="vertical" size="middle" className={themeStyles.space}>
          <RunnerMetrics runnerMetrics={runnerMetrics} />
        </Space>
      ),
    },
    {
      key: "BazelInvocationTabs-ActionCache",
      label: "Action Cache",
      icon: <HddOutlined />,
      children: (
        <Space direction="vertical" size="middle" className={themeStyles.space}>
          <AcMetrics acMetrics={acMetrics} />
        </Space>
      ),
    },
    {
      key: "BazelInvocationTabs-ActionsData",
      label: "Actions Data",
      icon: <NodeCollapseOutlined />,
      children: (
        <Space direction="vertical" size="middle" className={themeStyles.space}>
          <ActionDataMetrics acMetrics={acMetrics} />
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
      key: "BazelInvocationTabs-Timing",
      label: "Timing",
      icon: <FieldTimeOutlined />,
      children: (
        <Space direction="vertical" size="middle" className={themeStyles.space}>
          <TimingMetricsDisplay
            timingMetrics={timingMetrics}
            buildGraphMetrics={buildGraphMetrics}
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
      key: "BazelInvocationTabs-Network",
      label: "Network",
      icon: <WifiOutlined />,
      children: (
        <Space direction="vertical" size="middle" className={themeStyles.space}>
          <NetworkMetricsDisplay networkMetrics={networkMetrics} />
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
