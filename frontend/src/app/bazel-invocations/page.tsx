'use client';

import React from 'react';
import Content from '@/components/Content';
import PortalCard from '@/components/PortalCard';
import { BuildFilled } from '@ant-design/icons';
import BazelInvocationsTable from "@/components/BazelInvocationsTable";
import { isFeatureEnabled, FeatureType } from '@/utils/isFeatureEnabled';
import PageDisabled from '@/components/PageDisabled';

const Page: React.FC = () => {
  if (!isFeatureEnabled(FeatureType.BES) || !isFeatureEnabled(FeatureType.BES_PAGE_INVOCATIONS)) {
    return <PageDisabled />;
  }
  return (
    <Content
      content={
        <PortalCard
          icon={<BuildFilled />}
          titleBits={[<span key="title">Bazel Invocations</span>]}
        >
          <BazelInvocationsTable />
        </PortalCard>
      }
    />
  );
}

export default Page;
