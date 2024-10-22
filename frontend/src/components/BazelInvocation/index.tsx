import {
  ActionSummary,
  ArtifactMetrics,
  BazelInvocationInfoFragment,
  ProblemInfoFragment,
  RunnerCount,
  TargetMetrics,
  MemoryMetrics,
  TimingMetrics,
  NetworkMetrics,
  TestCollection,
  TargetPair,
  BuildGraphMetrics,
  BazelCommand,
} from "@/graphql/__generated__/graphql";
import styles from "../AppBar/index.module.css"
import React from "react";
import PortalDuration from "@/components/PortalDuration";
import PortalCard from "@/components/PortalCard";
import { Space, Tabs, Typography } from "antd";
import type { TabsProps } from "antd/lib";
import CopyTextButton from "@/components/CopyTextButton";
import PortalAlert from "@/components/PortalAlert";
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
} from "@ant-design/icons";
import themeStyles from '@/theme/theme.module.css';
import { debugMode } from "@/components/Utilities/debugMode";
import DebugInfo from "@/components/DebugInfo";
import BuildStepResultTag, { BuildStepResultEnum } from "@/components/BuildStepResultTag";
import DownloadButton from '@/components/DownloadButton'
import Link from '@/components/Link';
import { LogViewerCard } from "../LogViewer";
import RunnerMetrics from "../RunnerMetrics";
import AcMetrics from "../ActionCacheMetrics";
import TargetMetricsDisplay from "../TargetMetrics";
import ActionDataMetrics from "../ActionDataMetrics";
import ArtifactsDataMetrics from "../Artifacts";
import MemoryMetricsDisplay from "../MemoryMetrics";
import TimingMetricsDisplay from "../TimingMetrics";
import NetworkMetricsDisplay from "../NetworkMetrics";
import TestMetricsDisplay from "../TestsMetrics";
import { env } from 'next-runtime-env';
import CommandLineDisplay from "../CommandLine";
import SourceControlDisplay from "../SourceControlDisplay";


const BazelInvocation: React.FC<{
  invocationOverview: BazelInvocationInfoFragment;
  problems?: ProblemInfoFragment[] | null | undefined;
  children?: React.ReactNode;
  isNestedWithinBuildCard?: boolean;
}> = ({ invocationOverview, problems, children, isNestedWithinBuildCard }) => {
  const {
    invocationID,
    build,
    state,
    stepLabel,
    bazelCommand,
    profile,
    relatedFiles,
    sourceControl,
    user,
    metrics,
    testCollection,
    targets

  } = invocationOverview;

  var buildLogs = "tmp"
  //data for runner metrics
  var runnerMetrics: RunnerCount[] = [];
  metrics?.actionSummary?.runnerCount?.map((item: RunnerCount) => runnerMetrics.push(item));

  //data for ac metrics
  var acMetrics: ActionSummary | undefined = metrics?.actionSummary ?? undefined;

  //artifact metrics
  var artifactMetrics: ArtifactMetrics | undefined = metrics?.artifactMetrics ?? undefined;

  //data for target metrics
  var targetMetrics: TargetMetrics | undefined | null = metrics?.targetMetrics ?? undefined

  //memory metrics
  var memoryMetrics: MemoryMetrics | undefined = metrics?.memoryMetrics ?? undefined

  //build graph metrics
  var buildGraphMetrics: BuildGraphMetrics | undefined = metrics?.buildGraphMetrics ?? undefined

  //timing metrics
  var timingMetrics: TimingMetrics | undefined = metrics?.timingMetrics ?? undefined

  //netowrk metrics
  var networkMetrics: NetworkMetrics | undefined = metrics?.networkMetrics ?? undefined
  const bytesRecv = networkMetrics?.systemNetworkStats?.bytesRecv ?? 0
  const bytesSent = networkMetrics?.systemNetworkStats?.bytesSent ?? 0
  const hideNetworkMetricsTab: boolean = bytesRecv == 0 && bytesSent == 0

  //test data
  var testCollections: TestCollection[] | undefined | null = testCollection
  var targetData: TargetPair[] | undefined | null = targets
  var targetTimes: Map<string, number> = new Map<string, number>();

  targetData?.map(x => { targetTimes.set(x.label ?? "", x.durationInMs ?? 0) })

  //show/hide tabs
  const testCount = testCollection?.length ?? 0
  const hideTestsTab: boolean = testCount == 0
  const hideTargetsTab: boolean = (targetData?.length ?? 0) == 0 ? true : false
  const hideSourceControlTab: boolean = sourceControl == undefined || sourceControl == null ? true : false

  let { exitCode } = state;
  exitCode = exitCode ?? null;
  const titleBits: React.ReactNode[] = [<span key="label">User: {user?.LDAP}</span>];
  titleBits.push(<span key="label">Invocation: {invocationID}</span>)
  titleBits.push(<span className={styles.copyIcon}>
    <Typography.Text copyable={{ text: invocationID ?? "Copy" }}></Typography.Text>
  </span>)
  if (exitCode?.name) {
    titleBits.push(<BuildStepResultTag key="result" result={exitCode?.name as BuildStepResultEnum} />);
  }

  //logs
  const logs: string = buildLogs ?? "no build log data found..."

  var items: TabsProps['items'] = [
    {
      key: 'BazelInvocationTabs-Problems',
      label: 'Problems',
      icon: <ExclamationCircleOutlined />,
      children: <Space direction="vertical" size="middle" className={themeStyles.space}>
        {debugMode() && <DebugInfo invocationId={invocationID} exitCode={exitCode} />}
        {exitCode === null || exitCode.code !== 0 ? (
          children
        ) : (

          <PortalAlert
            message="There is no debug information to display because there are no reported failures with the build step"
            type="success"
            showIcon
          />
        )}

      </Space>,
    },
    {
      key: 'BazelInvocationTabs-Logs',
      label: 'Logs',
      icon: <FileSearchOutlined />,
      children: <Space direction="vertical" size="middle" className={themeStyles.space}>
        <PortalCard type="inner" icon={<FileSearchOutlined />} titleBits={["Build Logs"]} extraBits={["test"]}>
          <LogViewerCard log={logs} copyable={true} />
        </PortalCard>
      </Space>,
    },
    {
      key: 'BazelInvocationTabs-Runners',
      label: 'Runners',
      icon: <ClusterOutlined />,
      children: <Space direction="vertical" size="middle" className={themeStyles.space}>
        <RunnerMetrics runnerMetrics={runnerMetrics} />
      </Space>,
    },
    {
      key: 'BazelInvocationTabs-ActionCache',
      label: 'Action Cache',
      icon: <HddOutlined />,
      children: <Space direction="vertical" size="middle" className={themeStyles.space}>
        <AcMetrics acMetrics={acMetrics} />
      </Space>,
    },
    {
      key: 'BazelInvocationTabs-ActionsData',
      label: 'Actions Data',
      icon: <NodeCollapseOutlined />,
      children: <Space direction="vertical" size="middle" className={themeStyles.space}>
        <ActionDataMetrics acMetrics={acMetrics} />
      </Space>,
    },

    {
      key: 'BazelInvocationTabs-Artifacts',
      label: 'Artifacts',
      icon: <RadiusUprightOutlined />,
      children: <Space direction="vertical" size="middle" className={themeStyles.space}>
        <ArtifactsDataMetrics artifactMetrics={artifactMetrics} />
      </Space>,
    },
    {
      key: 'BazelInvocationTabs-Memory',
      label: 'Memory',
      icon: <AreaChartOutlined />,
      children: <Space direction="vertical" size="middle" className={themeStyles.space}>
        <MemoryMetricsDisplay memoryMetrics={memoryMetrics} />
      </Space>,
    },
    {
      key: 'BazelInvocationTabs-Timing',
      label: 'Timing',
      icon: <FieldTimeOutlined />,
      children: <Space direction="vertical" size="middle" className={themeStyles.space}>
        <TimingMetricsDisplay timingMetrics={timingMetrics} buildGraphMetrics={buildGraphMetrics} />
      </Space>,
    },

    {
      key: 'BazelInvocationTabs-Targets',
      label: 'Targets',
      icon: <DeploymentUnitOutlined />,
      children: <Space direction="vertical" size="middle" className={themeStyles.space}>
        <TargetMetricsDisplay targetMetrics={targetMetrics} targetData={targetData} />
      </Space>,
    },
    {
      key: 'BazelInvocationTabs-Tests',
      label: 'Tests',
      icon: <ExperimentOutlined />,
      children: <Space direction="vertical" size="middle" className={themeStyles.space}>
        <TestMetricsDisplay testMetrics={testCollections} targetTimes={targetTimes} />
      </Space>,
    },
    {
      key: 'BazelInvocationTabs-Network',
      label: 'Network',
      icon: <WifiOutlined />,
      children: <Space direction="vertical" size="middle" className={themeStyles.space}>
        <NetworkMetricsDisplay networkMetrics={networkMetrics} />
      </Space>,
    },
    {
      key: 'BazelInvocationTabs-CommandLine',
      label: 'Command Line',
      icon: <CodeOutlined />,
      children: <Space direction="vertical" size="middle" className={themeStyles.space}>
        <CommandLineDisplay commandLineData={bazelCommand ?? undefined} />
      </Space>,
    },
    {
      key: 'BazelInvocationTabs-SourceControl',
      label: 'Source Control',
      icon: <BranchesOutlined />,
      children: <Space direction="vertical" size="middle" className={themeStyles.space}>
        <SourceControlDisplay sourceControlData={sourceControl} />
      </Space>,
    },
  ];

  const hideLogs = true //hide the logs tab for now
  if (hideLogs == true) {
    var idx = items.findIndex((x, _) => x.key == "BazelInvocationTabs-Logs")
    if (idx > -1) {
      items.splice(idx, 1);
    }
  }

  if (hideTestsTab == true) {
    var idx = items.findIndex((x, _) => x.key == "BazelInvocationTabs-Tests")
    if (idx > -1) {
      items.splice(idx, 1);
    }
  }

  if (hideNetworkMetricsTab == true) {
    var idx = items.findIndex((x, _) => x.key == "BazelInvocationTabs-Network")
    if (idx > -1) {
      items.splice(idx, 1);
    }
  }

  const extraBits: React.ReactNode[] = [
    <PortalDuration key="duration" from={invocationOverview.startedAt} to={invocationOverview.endedAt} includeIcon includePopover />,
  ];

  if (env('NEXT_PUBLIC_BROWSER_URL') && profile) {
    var url = new URL(`blobs/sha256/file/${profile.digest}-${profile.sizeInBytes}/${profile.name}`, env('NEXT_PUBLIC_BROWSER_URL'))
    extraBits.push(
      <DownloadButton url={url.toString()} fileName="profile" buttonLabel="Profile" enabled={true} />
    );
  }

  if (problems?.length) {
    extraBits.push(
      <CopyTextButton buttonText="Copy Problems" copyText={problems.map(problem => problem.label).join(' ')} />
    );
  }

  if (!isNestedWithinBuildCard && build?.buildUUID) {
    extraBits.unshift(<span key="build">Build <Link href={`/builds/${build.buildUUID}`}>{build.buildUUID}</Link></span>);
  }

  return (
    <PortalCard type={isNestedWithinBuildCard ? "inner" : undefined} icon={<BuildOutlined />} titleBits={titleBits} extraBits={extraBits}>
      <Tabs defaultActiveKey="1" items={items} />
    </PortalCard >
  );
};

export default BazelInvocation;