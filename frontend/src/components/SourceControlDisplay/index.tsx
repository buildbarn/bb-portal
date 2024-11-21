import React from "react";
import { Space, Row, Descriptions } from 'antd';
import { BranchesOutlined } from '@ant-design/icons';
import { SourceControl } from "@/graphql/__generated__/graphql";
import PortalCard from "../PortalCard";
import Link from "next/link";
import { env } from 'next-runtime-env';


const SourceControlDisplay: React.FC<{
    stepLabel: string | undefined | null,
    sourceControlData: SourceControl | undefined | null }> = ({
        sourceControlData: sourceControlData,
        stepLabel: stepLabel
    }) => {
    //build urls
    var ghUrl = env('NEXT_PUBLIC_GITHUB_URL') ?? "https://github.com/"
    if (!ghUrl.endsWith("/")) {
        ghUrl += "/"
    }
    const repoUrl = ghUrl + sourceControlData?.repoURL
    const runURL = repoUrl + "/actions/runs/" + sourceControlData?.runID
    const actorURL = ghUrl + sourceControlData?.actor
    const branchURL = repoUrl + "/tree/" + sourceControlData?.branch
    const commitURL = repoUrl + "/commit/" + sourceControlData?.commitSha
    const prParts = sourceControlData?.refs?.split("/") ?? ""
    const prURL = repoUrl + "/" + prParts[1] + "/" + prParts[2]

    return (
        <Space direction="vertical" size="middle" style={{ display: 'flex' }} >
            <PortalCard type="inner" icon={<BranchesOutlined />} titleBits={["Source Control Information"]}>
                <Row>
                    <Space size="large">
                        <Descriptions bordered column={1}>
                            <Descriptions.Item label="Repository">
                                <Link target="_blank" href={repoUrl}>{sourceControlData?.repoURL}</Link>
                            </Descriptions.Item>
                            <Descriptions.Item label="Branch">
                                <Link target="_blank" href={branchURL}> {sourceControlData?.branch}</Link>
                            </Descriptions.Item>
                            <Descriptions.Item label="Commit SHA">
                                <Link target="_blank" href={commitURL}>  {sourceControlData?.commitSha}</Link></Descriptions.Item>
                            <Descriptions.Item label="Actor"><Link target="_blank" href={actorURL}>  {sourceControlData?.actor}</Link>
                            </Descriptions.Item>
                            <Descriptions.Item label="Run ID">
                                <Link target="_blank" href={runURL}>{sourceControlData?.runID}</Link>
                            </Descriptions.Item>
                            <Descriptions.Item label="Pull Request">
                                <Link target="_blank" href={prURL}>#{prParts[2]}</Link>
                            </Descriptions.Item>
                            <Descriptions.Item label="Step Label">
                                {stepLabel}
                            </Descriptions.Item>
                        </Descriptions>
                        <Descriptions bordered column={1}>
                            <Descriptions.Item label="Workflow">
                                {sourceControlData?.workflow}
                            </Descriptions.Item>
                            <Descriptions.Item label="Job">
                                {sourceControlData?.job}
                            </Descriptions.Item>
                            <Descriptions.Item label="Action">
                                {sourceControlData?.action}
                            </Descriptions.Item>
                            <Descriptions.Item label="Event Name">
                                {sourceControlData?.eventName}
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
    )
}

export default SourceControlDisplay;