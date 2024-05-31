'use client';

import React from 'react';
import Content from '@/components/Content';
import useScreenSize from '@/utils/screen';
import PortalCard from '@/components/PortalCard';
import { RocketFilled } from '@ant-design/icons';
import BuildsTable from "@/components/BuildsTable";

const Page: React.FC = () => {
  const screenSize = useScreenSize();
  return (
    <Content
      content={
        <PortalCard
          icon={<RocketFilled rotate={20}/>}
          titleBits={[<span key="title">Bazel Builds</span>]}
        >
          <BuildsTable height={screenSize.height - 370} />
        </PortalCard>
      }
    />
  );
}

export default Page;
