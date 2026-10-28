import { CalendarFilled } from "@ant-design/icons";
import type React from "react";
import DrainsTable from "@/components/DrainsTable";
import { PortalCard } from "@/components/PortalCard";
import WorkersInfo from "@/components/WorkersInfo";
import type { DrainSearchParams } from "@/routes/scheduler.drain";

export const SchedulerDrainsPage: React.FC<DrainSearchParams> = ({
  sizeClassQueueName,
}) => {
  return (
    <PortalCard
      icon={<CalendarFilled />}
      titleBits={[<span key="title">Drains</span>]}
    >
      <WorkersInfo sizeClassQueueName={sizeClassQueueName} />
      <DrainsTable sizeClassQueueName={sizeClassQueueName} />
    </PortalCard>
  );
};
