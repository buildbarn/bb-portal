import { CodeFilled } from "@ant-design/icons";
import type React from "react";
import OperationsGrid from "@/components/OperationsGrid";
import PortalCard from "@/components/PortalCard";
import type { OperationsFilterParams } from "@/routes/operations.index";

interface Props {
  filter: OperationsFilterParams;
}

export const OperationsPage: React.FC<Props> = ({ filter }) => {
  return (
    <PortalCard
      icon={<CodeFilled />}
      titleBits={[<span key="title">Operations</span>]}
    >
      <OperationsGrid filter={filter} />
    </PortalCard>
  );
};
