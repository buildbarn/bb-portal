import { RocketFilled } from "@ant-design/icons";
import type React from "react";
import BuildsTable from "@/components/BuildsTable";
import PortalCard from "@/components/PortalCard";

type Props = React.ComponentProps<typeof BuildsTable>;

export const BuildsPage: React.FC<Props> = (props) => {
  return (
    <PortalCard
      icon={<RocketFilled rotate={20} />}
      titleBits={[<span key="title">Bazel Builds</span>]}
    >
      <BuildsTable {...props} />
    </PortalCard>
  );
};
