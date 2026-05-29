import { RocketFilled } from "@ant-design/icons";
import type React from "react";
import BuildsTable from "@/components/BuildsTable";
import PortalCard from "@/components/PortalCard";

export const BuildsPage: React.FC = () => {
  return (
    <PortalCard
      icon={<RocketFilled rotate={20} />}
      titleBits={[<span key="title">Bazel Builds</span>]}
    >
      <BuildsTable />
    </PortalCard>
  );
};
