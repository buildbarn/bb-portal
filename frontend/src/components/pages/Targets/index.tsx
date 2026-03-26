
import Content from '@/components/Content';
import PortalCard from '@/components/PortalCard';
import TargetGrid from '@/components/Targets/TargetGrid';
import { DeploymentUnitOutlined } from '@ant-design/icons';
import { Alert, Space } from 'antd';
import type React from 'react';

export const TargetsPage: React.FC = () => {
    return (
        <Content
            content={
                <Space direction="vertical" size="middle" style={{ display: 'flex' }}>
                    <PortalCard
                        icon={<DeploymentUnitOutlined />}
                        extraBits={[<Alert
                            key="search-by-label"
                            showIcon
                            message = "Search by label to further refine your result"
                            type = "info"
                        />]}
                        titleBits={[<span key="title">Targets Overview  </span>]}>
                        <TargetGrid />
                    </PortalCard>
                </Space >
            }
        />
    );
}
