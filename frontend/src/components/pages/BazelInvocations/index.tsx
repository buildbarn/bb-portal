import { BuildFilled } from "@ant-design/icons";
import type React from "react";
import BazelInvocationsTable from "@/components/BazelInvocationsTable";
import PortalCard from "@/components/PortalCard";

type Props = React.ComponentProps<typeof BazelInvocationsTable>;

export const BazelInvocationsPage: React.FC<Props> = (props) => {
  return (
    <PortalCard
      icon={<BuildFilled />}
      titleBits={[<span key="title">Bazel Invocations</span>]}
    >
      <BazelInvocationsTable {...props} />
    </PortalCard>
  );
};
