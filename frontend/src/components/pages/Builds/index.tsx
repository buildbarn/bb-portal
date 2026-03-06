
import BuildsTable from "@/components/BuildsTable";
import Content from '@/components/Content';
import PortalCard from '@/components/PortalCard';
import { RocketFilled } from '@ant-design/icons';
import type React from 'react';

export const BuildsPage: React.FC = () => {
  return (
    <Content
      content={
        <PortalCard
          icon={<RocketFilled rotate={20} />}
          titleBits={[<span key="title">Bazel Builds</span>]}
        >
          <BuildsTable />
        </PortalCard>
      }
    />
  );
}

