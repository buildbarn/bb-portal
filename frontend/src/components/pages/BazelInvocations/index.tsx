import { BuildFilled } from "@ant-design/icons";
import type React from "react";
import BazelInvocationsTable from "@/components/BazelInvocationsTable";
import Content from "@/components/Content";
import PortalCard from "@/components/PortalCard";

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
};
