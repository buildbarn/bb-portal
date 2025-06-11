'use client';

import React from 'react';
import Content from '@/components/Content';
import PortalCard from '@/components/PortalCard';
import { Space } from 'antd';
import { DeploymentUnitOutlined } from '@ant-design/icons';
import TargetDetails from '@/components/Targets/TargetDetails';
import { isFeatureEnabled, FeatureType } from '@/utils/isFeatureEnabled';
import PageDisabled from '@/components/PageDisabled';

interface PageParams {
    params: {
        slug: string
    }
}

const Page: React.FC<PageParams> = ({ params }) => {
    if (!isFeatureEnabled(FeatureType.BES) || !isFeatureEnabled(FeatureType.BES_PAGE_TARGETS)) {
        return <PageDisabled />;
    }

    const label = decodeURIComponent(atob(decodeURIComponent(params.slug)))
    return (
        <Content
            content={
                <Space direction="vertical" size="middle" style={{ display: 'flex' }}>
                    <PortalCard
                        icon={<DeploymentUnitOutlined />}
                        titleBits={[<span key="title">Target Details</span>]}
                    >
                        <TargetDetails label={label} />
                    </PortalCard>
                </Space >
            }
        />
    );
}

export default Page;
