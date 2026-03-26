
import Content from "@/components/Content";
import OperationsGrid from "@/components/OperationsGrid";
import PortalCard from "@/components/PortalCard";
import type { OperationsFilterParams } from "@/routes/operations.index";
import { CodeFilled } from "@ant-design/icons";
import type React from "react";

interface Props {
  filter: OperationsFilterParams;
}

export const OperationsPage: React.FC<Props> = ({ filter }) => {
  return (
    <Content
      content={
        <PortalCard
          icon={<CodeFilled />}
          titleBits={[<span key="title">Operations</span>]}
        >
          <OperationsGrid filter={filter}/>
        </PortalCard>
      }
    />
  );
};
