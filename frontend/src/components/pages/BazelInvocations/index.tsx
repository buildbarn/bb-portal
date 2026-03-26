
import type React from 'react';
import Content from '@/components/Content';
import PortalCard from '@/components/PortalCard';
import { BuildFilled } from '@ant-design/icons';
import BazelInvocationsTable from "@/components/BazelInvocationsTable";

export const BazelInvocationsPage: React.FC = () => {
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
