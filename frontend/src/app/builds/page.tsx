'use client';

import type React from 'react';
import Content from '@/components/Content';
import PortalCard from '@/components/PortalCard';
import { RocketFilled } from '@ant-design/icons';
import BuildsTable from "@/components/BuildsTable";
import { isFeatureEnabled, FeatureType } from '@/utils/isFeatureEnabled';
import PageDisabled from '@/components/PageDisabled';

const Page: React.FC = () => {
  if (!isFeatureEnabled(FeatureType.BES) || !isFeatureEnabled(FeatureType.BES_PAGE_BUILDS)) {
    return <PageDisabled />;
  }

  return (
    <Content
      content={
        <PortalCard
          icon={<RocketFilled rotate={20}/>}
          titleBits={[<span key="title">Bazel Builds</span>]}
        >
          <BuildsTable />
        </PortalCard>
      }
    />
  );
}

export default Page;
