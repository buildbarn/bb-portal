import React from "react";
import { Table, Row, Col, Space } from 'antd';
import type { TableColumnsType } from "antd/lib";
import { ArtifactMetrics } from "@/graphql/__generated__/graphql";
import PortalCard from "../PortalCard";
import {
    RadiusUprightOutlined

} from "@ant-design/icons";


const artifacts_columns: TableColumnsType<ArtifactMetricsTableData> = [
    {
        title: "Type",
        dataIndex: "name"
    },
    {
        title: "Size (bytes)",
        dataIndex: "sizeInBytes",
        sorter: (a, b) => (a.sizeInBytes ?? 0) - (b.sizeInBytes ?? 0),
    },
    {
        title: "Count",
        dataIndex: "count",
        sorter: (a, b) => (a.count ?? 0) - (b.count ?? 0),
    },
]


interface ArtifactMetricsTableData {
    name: string;
    sizeInBytes: number;
    count: number;
}

const ArtifactsDataMetrics: React.FC<{ artifactMetrics: ArtifactMetrics | undefined; }> = ({ artifactMetrics }) => {

    const artifacts_data: ArtifactMetricsTableData[] = [];
    artifacts_data.push(
        {
            name: "Source Artifacts Read",
            sizeInBytes: artifactMetrics?.sourceArtifactsRead?.at(0)?.sizeInBytes ?? 0,
            count: artifactMetrics?.sourceArtifactsRead?.at(0)?.count ?? 0
        },
        {
            name: "Output Artifacts From Action Cache",
            sizeInBytes: artifactMetrics?.outputArtifactsFromActionCache?.at(0)?.sizeInBytes ?? 0,
            count: artifactMetrics?.outputArtifactsFromActionCache?.at(0)?.count ?? 0
        },
        {
            name: "Output Artifacts Seen",
            sizeInBytes: artifactMetrics?.outputArtifactsSeen?.at(0)?.sizeInBytes ?? 0,
            count: artifactMetrics?.outputArtifactsSeen?.at(0)?.count ?? 0
        },
        {
            name: "Top Level Artifacts",
            sizeInBytes: artifactMetrics?.topLevelArtifacts?.at(0)?.sizeInBytes ?? 0,
            count: artifactMetrics?.topLevelArtifacts?.at(0)?.count ?? 0
        },
    )


    const actionsTitle: React.ReactNode[] = [<span key="label">Artifacts</span>];

    return (
        <Space direction="vertical" size="middle" style={{ display: 'flex' }} >
            <PortalCard icon={<RadiusUprightOutlined />} titleBits={actionsTitle}>
                <Row justify="space-around" align="middle">
                    <Col span="18">
                        <Table
                            columns={artifacts_columns}
                            dataSource={artifacts_data}
                            showSorterTooltip={{ target: 'sorter-icon' }}
                            pagination={false}
                        />
                    </Col>
                    <Col span="6" />
                </Row>
            </PortalCard>
        </Space>
    )
}

export default ArtifactsDataMetrics;