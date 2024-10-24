import React from "react";
import { Space, Row, Descriptions } from 'antd';
import { BranchesOutlined } from '@ant-design/icons';
import { SourceControl } from "@/graphql/__generated__/graphql";
import PortalCard from "../PortalCard";
import Link from "next/link";
import { env } from 'next-runtime-env';


const SourceControlDisplay: React.FC<{ sourceControlData: SourceControl | undefined | null }> = ({
    sourceControlData: sourceControlData
}) => {
    //build urls
    var ghUrl = env('NEXT_PUBLIC_GITHUB_URL') ?? "https://github.com/"
    if (!ghUrl.endsWith("/")) {
        ghUrl += "/"
    }
    const runURL = ghUrl + sourceControlData?.repoURL + "/actions/runs/" + sourceControlData?.runID
    const actorURL = ghUrl + sourceControlData?.actor
    const branchURL = ghUrl + sourceControlData?.repoURL + "/tree/" + sourceControlData?.branch
    const commitURL = ghUrl + sourceControlData?.repoURL + "/commit/" + sourceControlData?.commitSha
    const prParts = sourceControlData?.refs?.split("/") ?? ""
    const prURL = ghUrl + sourceControlData?.repoURL + "/" + prParts[1] + "/" + prParts[2]

    return (
        <Space direction="vertical" size="middle" style={{ display: 'flex' }} >
            <PortalCard type="inner" icon={<BranchesOutlined />} titleBits={["Source Control Information"]}>
                <Row>
                    <Space size="large">
                        <Descriptions bordered column={1}>
                            <Descriptions.Item label="Repository URL">
                                <Link target="_blank" href={sourceControlData?.repoURL ?? ""}>{sourceControlData?.repoURL}</Link>
                            </Descriptions.Item>
                            <Descriptions.Item label="Branch Name">
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
                        </Descriptions>
                    </Space>
                </Row>
            </PortalCard>
        </Space>
    )
}

export default SourceControlDisplay;