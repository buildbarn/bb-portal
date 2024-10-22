import React from "react";
import { Space, Table, Row, Col, Statistic, List, Descriptions } from 'antd';
import { BranchesOutlined, CodeOutlined, DeploymentUnitOutlined, SearchOutlined } from '@ant-design/icons';
import type { StatisticProps, TableColumnsType } from "antd/lib";
import CountUp from 'react-countup';
import { BazelCommand, SourceControl, TargetMetrics, TargetPair } from "@/graphql/__generated__/graphql";
import PortalCard from "../PortalCard";
import { SearchFilterIcon, SearchWidget } from '@/components/SearchWidgets';
import NullBooleanTag from "../NullableBooleanTag";
import styles from "../../theme/theme.module.css"
import { millisecondsToTime } from "../Utilities/time";
import Link from "next/link";



const SourceControlDisplay: React.FC<{ sourceControlData: SourceControl | undefined | null }> = ({
    sourceControlData: sourceControlData
}) => {

    const createUnorderedList = (items: string[]): JSX.Element => {
        return (
            <ul>
                {items.map((item, index) => (
                    <li key={index}>{item}</li>
                ))}
            </ul>
        );
    };

    const runURL = sourceControlData?.repoURL + "/actions/runs/" + "10976627244"
    const actorURL = "https://github.com/" + sourceControlData?.actor
    const branchURL = sourceControlData?.repoURL + "/tree/" + sourceControlData?.branch
    const commitURL = sourceControlData?.repoURL + "/commit/" + sourceControlData?.commitSha
    var prref = "refs/pull/12718/merge"
    var prParts = prref.split("/")
    const prURL = sourceControlData?.repoURL + "/" + prParts[1] + "/" + prParts[2]
    return (

        <Space direction="vertical" size="middle" style={{ display: 'flex' }} >
            <PortalCard type="inner" icon={<BranchesOutlined />} titleBits={["Source Control Information"]}>
                <Row>
                    <Space size="large">
                        <Descriptions bordered column={1}>
                            <Descriptions.Item label="Repository URL"><Link target="_blank" href={sourceControlData?.repoURL ?? ""}>{sourceControlData?.repoURL}</Link> </Descriptions.Item>
                            <Descriptions.Item label="Branch Name"><Link target="_blank" href={branchURL}> {sourceControlData?.branch}</Link></Descriptions.Item>
                            <Descriptions.Item label="Commit SHA"><Link target="_blank" href={commitURL}>  {sourceControlData?.commitSha}</Link></Descriptions.Item>
                            <Descriptions.Item label="Actor"><Link target="_blank" href={actorURL}>  {sourceControlData?.actor}</Link></Descriptions.Item>
                            <Descriptions.Item label="Run ID"> <Link target="_blank" href={runURL}>10976627244</Link></Descriptions.Item>
                            <Descriptions.Item label="Pull Request"> <Link target="_blank" href={prURL}>#{prParts[2]}</Link></Descriptions.Item>

                        </Descriptions>
                    </Space>
                </Row>
            </PortalCard>
        </Space>
    )
}

export default SourceControlDisplay;