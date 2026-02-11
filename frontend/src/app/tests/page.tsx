'use client';

import React from 'react';
import Content from '@/components/Content';
import PortalCard from '@/components/PortalCard';
import { Alert, Space } from 'antd';
import { ExperimentFilled } from '@ant-design/icons';
import TestGrid from '@/components/TestGrid';
import { isFeatureEnabled, FeatureType } from '@/utils/isFeatureEnabled';
import PageDisabled from '@/components/PageDisabled';

const Page: React.FC = () => {
    if (!isFeatureEnabled(FeatureType.BES) || !isFeatureEnabled(FeatureType.BES_PAGE_TESTS)) {
        return <PageDisabled />;
    }

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

export default Page;
