import { CalendarFilled } from "@ant-design/icons";
import type React from "react";
import Content from "@/components/Content";
import PortalCard from "@/components/PortalCard";
import WorkersGrid from "@/components/WorkersGrid";
import type { WorkerSearchParams } from "@/routes/scheduler.worker";

export const SchedulerWorkersPage: React.FC<WorkerSearchParams> = ({
  workerStatusFilter,
  sizeClassQueueName,
  cursor,
}) => {
  return (
    <Content
      content={
        <PortalCard
          icon={<CalendarFilled />}
          titleBits={[<span key="title">Workers</span>]}
        >
          <WorkersGrid
            workerStatusFilter={workerStatusFilter}
            sizeClassQueueName={sizeClassQueueName}
            cursor={cursor}
          />
        </PortalCard>
      }
    />
  );
};
