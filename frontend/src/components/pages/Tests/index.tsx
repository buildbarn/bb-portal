
import Content from '@/components/Content';
import PortalCard from '@/components/PortalCard';
import TestGrid from '@/components/TestGrid';
import { ExperimentFilled } from '@ant-design/icons';
import { Alert, Space } from 'antd';
import type React from 'react';

export const TestsPage: React.FC = () => {
    return (
        <Content
            content={
                <Space direction="vertical" size="middle" style={{ display: 'flex' }}>
                    <PortalCard
                        icon={<ExperimentFilled />}
                        extraBits={[<Alert
                            key="search-by-label"
                            showIcon
                            message = "Search by label and/or instance name to further refine your result"
                            type = "info"
                        />]}
                        titleBits={[<span key="title">Tests Overview</span>]}>
                        <TestGrid />
                    </PortalCard>
                </Space >
            }
        />
    );
}
