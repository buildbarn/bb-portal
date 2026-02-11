'use client';

import React from 'react';
import Content from '@/components/Content';
import PortalCard from '@/components/PortalCard';
import { Alert, Space } from 'antd';
import { DeploymentUnitOutlined } from '@ant-design/icons';
import TargetGrid from '@/components/Targets/TargetGrid';
import { isFeatureEnabled, FeatureType } from '@/utils/isFeatureEnabled';
import PageDisabled from '@/components/PageDisabled';

const Page: React.FC = () => {
    if (!isFeatureEnabled(FeatureType.BES) || !isFeatureEnabled(FeatureType.BES_PAGE_TARGETS)) {
        return <PageDisabled />;
    }

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

export default Page;
