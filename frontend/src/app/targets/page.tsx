'use client';

import React from 'react';
import Content from '@/components/Content';
import PortalCard from '@/components/PortalCard';
import { Space } from 'antd';
import { DeploymentUnitOutlined } from '@ant-design/icons';
import TargetGrid from '@/components/Targets/TargetGrid';

const Page: React.FC = () => {
    return (
        <Content
            content={
                <Space direction="vertical" size="middle" style={{ display: 'flex' }}>
                    <PortalCard
                        icon={<DeploymentUnitOutlined />}
                        titleBits={[<span key="title">Targets Overview</span>]}>
                        <TargetGrid />
                    </PortalCard>
                </Space >
            }
        />
    );
}

export default Page;
