import React from "react";
import { Table, Row, Col, Space } from 'antd';
import type { TableColumnsType } from "antd/lib";
import { ArtifactMetrics } from "@/graphql/__generated__/graphql";
import PortalCard from "../PortalCard";
import { RadiusUprightOutlined } from "@ant-design/icons";
import styles from "../../theme/theme.module.css"
import { readableFileSize } from "@/utils/filesize";


const artifacts_columns: TableColumnsType<ArtifactMetricsTableData> = [
    {
        title: "Type",
        dataIndex: "name"
    },
    {
        title: "Size",
        dataIndex: "sizeInBytes",
        align: "right",
        render: (_, record) => <span className={styles.numberFormat} >{readableFileSize(record.sizeInBytes)}</span>,
        sorter: (a, b) => (a.sizeInBytes ?? 0) - (b.sizeInBytes ?? 0),
    },
    {
        title: "Count",
        dataIndex: "count",
        align: "right",
        render: (_, record) => <span className={styles.numberFormat} >{record.count}</span>,
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
            sizeInBytes: artifactMetrics?.sourceArtifactsRead?.sizeInBytes ?? 0,
            count: artifactMetrics?.sourceArtifactsRead?.count ?? 0
        },
        {
            name: "Output Artifacts From Action Cache",
            sizeInBytes: artifactMetrics?.outputArtifactsFromActionCache?.sizeInBytes ?? 0,
            count: artifactMetrics?.outputArtifactsFromActionCache?.count ?? 0
        },
        {
            name: "Output Artifacts Seen",
            sizeInBytes: artifactMetrics?.outputArtifactsSeen?.sizeInBytes ?? 0,
            count: artifactMetrics?.outputArtifactsSeen?.count ?? 0
        },
        {
            name: "Top Level Artifacts",
            sizeInBytes: artifactMetrics?.topLevelArtifacts?.sizeInBytes ?? 0,
            count: artifactMetrics?.topLevelArtifacts?.count ?? 0
        },
    )


    const actionsTitle: React.ReactNode[] = [<span key="label">Artifacts</span>];

    return (
        <Space direction="vertical" size="middle" style={{ display: 'flex' }} >
            <PortalCard type="inner" icon={<RadiusUprightOutlined />} titleBits={actionsTitle}>
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