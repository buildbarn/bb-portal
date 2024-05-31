import {
  BazelInvocationInfoFragment,
  ProblemInfoFragment,
} from "@/graphql/__generated__/graphql";
import React from "react";
import PortalDuration from "@/components/PortalDuration";
import PortalCard from "@/components/PortalCard";
import {Space} from "antd";
import CopyTextButton from "@/components/CopyTextButton";
import PortalAlert from "@/components/PortalAlert";
import {BuildOutlined} from "@ant-design/icons";
import themeStyles from '@/theme/theme.module.css';
import {debugMode} from "@/components/Utilities/debugMode";
import DebugInfo from "@/components/DebugInfo";
import BuildStepResultTag, {BuildStepResultEnum} from "@/components/BuildStepResultTag";
import Link from '@/components/Link';

const BazelInvocation: React.FC<{
  invocationOverview: BazelInvocationInfoFragment;
  problems?: ProblemInfoFragment[] | null | undefined;
  children?: React.ReactNode;
  isNestedWithinBuildCard?: boolean;
}> = ({ invocationOverview, problems, children, isNestedWithinBuildCard }) => {
  const { invocationID, build, state, stepLabel, bazelCommand, relatedFiles } = invocationOverview;
  let { exitCode } = state;

  exitCode = exitCode ?? null;

  const titleBits: React.ReactNode[] = [<span key="label">Invocation: {invocationID}</span>];
  if (exitCode?.name) {
    titleBits.push(<BuildStepResultTag key="result" result={exitCode?.name as BuildStepResultEnum} />);
  }

  const extraBits: React.ReactNode[] = [
    <PortalDuration key="duration" from={invocationOverview.startedAt} to={invocationOverview.endedAt} includeIcon includePopover />,
  ];
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
      <Space direction="vertical" size="middle" className={themeStyles.space}>
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
      </Space>
    </PortalCard>
  );
};

export default BazelInvocation;
