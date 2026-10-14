import {
  ApartmentOutlined,
  CalendarFilled,
  UnorderedListOutlined,
} from "@ant-design/icons";
import type React from "react";
import InvocationChildrenTable from "@/components/InvocationChildrenTable";
import InvocationInfo from "@/components/InvocationInfo";
import { PortalCard } from "@/components/PortalCard";
import QueuedOperationsTable from "@/components/QueuedOperationsTable";
import type { InvocationSearchParams } from "@/routes/scheduler.invocation";

export const SchedulerInvocationPage: React.FC<InvocationSearchParams> = ({
  invocationName,
  childrenFilter,
  cursor,
}) => {
  return (
    <PortalCard
      icon={<CalendarFilled />}
      titleBits={[
        <span key="title">
          {invocationName.ids.length ? "Invocation" : "Root invocation"}
        </span>,
      ]}
    >
      <InvocationInfo invocationName={invocationName} />
      <PortalCard
        type="inner"
        icon={<ApartmentOutlined />}
        titleBits={[<span key="title">Child invocations</span>]}
      >
        <InvocationChildrenTable
          invocationName={invocationName}
          childrenFilter={childrenFilter}
        />
      </PortalCard>
      <PortalCard
        type="inner"
        icon={<UnorderedListOutlined />}
        titleBits={[<span key="title">Queued operations</span>]}
      >
        <QueuedOperationsTable
          invocationName={invocationName}
          cursor={cursor}
        />
      </PortalCard>
    </PortalCard>
  );
};
