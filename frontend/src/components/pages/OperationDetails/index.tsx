import { CodeFilled } from "@ant-design/icons";
import type React from "react";
import OperationDetails from "@/components/OperationDetails";
import PortalCard from "@/components/PortalCard";

interface Params {
  operationID: string;
}

export const OperationDetailsPage: React.FC<Params> = ({ operationID }) => {
  return (
    <PortalCard
      icon={<CodeFilled />}
      titleBits={[<span key="title">{`Operation ${operationID}`}</span>]}
    >
      <OperationDetails operationID={operationID} />
    </PortalCard>
  );
};
