'use client';

import React from 'react';
import Content from '@/components/Content';
import useScreenSize from '@/utils/screen';
import PortalCard from '@/components/PortalCard';
import { BuildFilled } from '@ant-design/icons';
import BazelInvocationsTable from "@/components/BazelInvocationsTable";

const Page: React.FC = () => {
  const screenSize = useScreenSize();
  return (
    <Content
      content={
        <PortalCard
          icon={<BuildFilled />}
          titleBits={[<span key="title">Bazel Invocations</span>]}
        >
          <BazelInvocationsTable height={screenSize.height - 370} />
        </PortalCard>
      }
    />
  );
}

export default Page;
