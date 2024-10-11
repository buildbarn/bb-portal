'use client';

import React from 'react';
import Content from '@/components/Content';
import PortalCard from '@/components/PortalCard';
import { Space } from 'antd';
import { ExperimentFilled } from '@ant-design/icons';
import TestGrid from '@/components/TestGrid';

const Page: React.FC = () => {
    return (
        <Content
            content={
                <Space direction="vertical" size="middle" style={{ display: 'flex' }}>
                    <PortalCard
                        icon={<ExperimentFilled />}
                        titleBits={[<span key="title">Tests Overview</span>]}>
                        <TestGrid />
                    </PortalCard>
                </Space >
            }
        />
    );
}

export default Page;
