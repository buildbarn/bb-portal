import type React from "react";
import PortalAlert from "@/components/PortalAlert";

interface Props {
  type: string;
}

export const InvocationDataNotFoundAlert: React.FC<Props> = ({ type }) => {
  return (
    <PortalAlert
      showIcon
      title={`This invocation has not reported any ${type} yet.`}
    />
  );
};
