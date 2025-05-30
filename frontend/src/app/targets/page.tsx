'use client';

import React from 'react';
import Content from '@/components/Content';
import PortalCard from '@/components/PortalCard';
import { Alert, Space } from 'antd';
import { DeploymentUnitOutlined } from '@ant-design/icons';
import TargetGrid from '@/components/Targets/TargetGrid';
import { isFeatureEnabled, FeatureType } from '@/utils/isFeatureEnabled';
import { notFound } from 'next/navigation';

const Page: React.FC = () => {
    if (!isFeatureEnabled(FeatureType.BES)) {
        return notFound();
    }

    return (
        <Content
            content={
                <Space direction="vertical" size="middle" style={{ display: 'flex' }}>
                    <PortalCard
                        icon={<DeploymentUnitOutlined />}
                        extraBits={[<Alert
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

export default Page;
