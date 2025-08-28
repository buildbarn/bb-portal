import {
  type SourceControl,
  SourceControlProvider,
} from "@/graphql/__generated__/graphql";
import { BranchesOutlined } from "@ant-design/icons";
import { Descriptions, Row, Space } from "antd";
import Link from "next/link";
import type React from "react";
import PortalCard from "../PortalCard";

const getRepoUrl = (
  sc: SourceControl | undefined | null,
): string | undefined => {
  if (
    sc?.instanceURL === null ||
    sc?.instanceURL === undefined ||
    sc?.instanceURL === "" ||
    sc?.repo === null ||
    sc?.repo === undefined ||
    sc?.repo === ""
  ) {
    return undefined;
  }
  return `${sc.instanceURL}/${sc.repo}`;
};

const getRefLabelAndUrl = (
  sc: SourceControl | undefined | null,
  repoUrl: string | undefined,
): [string | undefined, string | undefined] => {
  if (
    sc?.refs === null ||
    sc?.refs === undefined ||
    sc?.refs === "" ||
    repoUrl === undefined
  ) {
    return [undefined, undefined];
  }
  switch (sc?.provider) {
    case SourceControlProvider.Github:
      if (sc.refs.startsWith("refs/heads/")) {
        return ["Branch", `${repoUrl}/tree/${sc.refs.substring("refs/heads/".length)}`];
      }
      if (sc.refs.startsWith("refs/tags/")) {
        return ["Tag",`${repoUrl}/tree/${sc.refs.substring("refs/tags/".length)}`];
      }
      if (sc.refs.startsWith("refs/pull/")) {
        const prNumber = sc.refs
          .substring("refs/pull/".length)
          .split("/")[0];
        return ["Pull request",`${repoUrl}/pull/${prNumber}`];
      }
    case SourceControlProvider.Gitlab:
      return ["Branch",`${repoUrl}/-/tree/${sc.refs}`];
    default:
      return [undefined, undefined];
  }
};

const getCommitUrl = (
  sc: SourceControl | undefined | null,
  repoUrl: string | undefined,
): string | undefined => {
  if (
    sc?.commitSha === null ||
    sc?.commitSha === undefined ||
    sc?.commitSha === "" ||
    repoUrl === undefined
  ) {
    return undefined;
  }
  switch (sc?.provider) {
    case SourceControlProvider.Github:
      return `${repoUrl}/commit/${sc.commitSha}`;
    case SourceControlProvider.Gitlab:
      return `${repoUrl}/-/commit/${sc.commitSha}`;
    default:
      return undefined;
  }
};

const getActorUrl = (
  sc: SourceControl | undefined | null,
): string | undefined => {
  if (
    sc?.actor === null ||
    sc?.actor === undefined ||
    sc?.actor === "" ||
    sc?.instanceURL === null ||
    sc?.instanceURL === undefined ||
    sc?.instanceURL === ""
  ) {
    return undefined;
  }
  return `${sc.instanceURL}/${sc.actor}`;
};

const getRunUrl = (
  sc: SourceControl | undefined | null,
  repoUrl: string | undefined,
): string | undefined => {
  if (
    sc?.runID === null ||
    sc?.runID === undefined ||
    sc?.runID === "" ||
    repoUrl === undefined
  ) {
    return undefined;
  }
  switch (sc?.provider) {
    case SourceControlProvider.Github:
      return `${repoUrl}/actions/runs/${sc.runID}`;
    case SourceControlProvider.Gitlab:
      return `${repoUrl}/-/jobs/${sc.runID}`;
    default:
      return undefined;
  }
};

type RepoLinkProps = {
  text: string;
  url?: string;
};

const RepoLink: React.FC<RepoLinkProps> = ({ text, url }) => {
  if (url) {
    return (
      <Link target="_blank" href={url}>
        {text}
      </Link>
    );
  }
  return <>{text}</>;
};

const SourceControlDisplay: React.FC<{
  stepLabel: string | undefined | null;
  sourceControlData: SourceControl | undefined | null;
}> = ({ sourceControlData, stepLabel }) => {
  const repoUrl = getRepoUrl(sourceControlData);
  const [refLabel, refUrl] = getRefLabelAndUrl(sourceControlData, repoUrl);
  const commitUrl = getCommitUrl(sourceControlData, repoUrl);
  const actorUrl = getActorUrl(sourceControlData);
  const runUrl = getRunUrl(sourceControlData, repoUrl);

  let workflowLabel = sourceControlData?.workflow || "";
  const runNumber = sourceControlData?.runNumber || "";
  if (workflowLabel != "" && runNumber != "") {
    workflowLabel = `${workflowLabel} #${runNumber}`;
  }

  return (
    <Space direction="vertical" size="middle" style={{ display: "flex" }}>
      <PortalCard
        type="inner"
        icon={<BranchesOutlined />}
        titleBits={["Source Control Information"]}
      >
        <Row>
          <Space size="large">
            <Descriptions bordered column={1}>
              <Descriptions.Item label="Repository">
                <RepoLink url={repoUrl} text={sourceControlData?.repo || ""} />
              </Descriptions.Item>
              <Descriptions.Item label={refLabel || "Ref"}>
                <RepoLink url={refUrl} text={sourceControlData?.refs || ""} />
              </Descriptions.Item>
              <Descriptions.Item label="Commit SHA">
                <RepoLink url={commitUrl} text={sourceControlData?.commitSha || ""} />
              </Descriptions.Item>
              <Descriptions.Item label="Actor">
                <RepoLink url={actorUrl} text={sourceControlData?.actor || ""} />
              </Descriptions.Item>
              <Descriptions.Item label="Event Name">
                {sourceControlData?.eventName}
              </Descriptions.Item>
              <Descriptions.Item label="Workflow">
                <RepoLink url={runUrl} text={workflowLabel} />
              </Descriptions.Item>
            </Descriptions>
            <Descriptions bordered column={1}>
              <Descriptions.Item label="Run ID">
                <RepoLink url={runUrl} text={sourceControlData?.runID || ""} />
              </Descriptions.Item>
              <Descriptions.Item label="Job">
                {sourceControlData?.job}
              </Descriptions.Item>
              <Descriptions.Item label="Action">
                {sourceControlData?.action}
              </Descriptions.Item>
              <Descriptions.Item label="Runner Name">
                {sourceControlData?.runnerName}
              </Descriptions.Item>
              <Descriptions.Item label="Runner Architecture">
                {sourceControlData?.runnerArch}
              </Descriptions.Item>
              <Descriptions.Item label="Runner Operating System">
                {sourceControlData?.runnerOs}
              </Descriptions.Item>
            </Descriptions>
          </Space>
        </Row>
      </PortalCard>
    </Space>
  );
};

export default SourceControlDisplay;
