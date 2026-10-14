import { CalendarFilled } from "@ant-design/icons";
import type React from "react";
import { OperatorModeSwitch } from "@/components/OperatorModeSwitch";
import { PortalCard } from "@/components/PortalCard";
import WorkersGrid from "@/components/WorkersGrid";
import type { WorkerSearchParams } from "@/routes/scheduler.worker";

export const SchedulerWorkersPage: React.FC<WorkerSearchParams> = ({
  workerStatusFilter,
  sizeClassQueueName,
  invocationIds,
  cursor,
}) => {
  return (
    <PortalCard
      icon={<CalendarFilled />}
      titleBits={[<span key="title">Workers</span>]}
      extraBits={[<OperatorModeSwitch key="operatorMode" />]}
    >
      <WorkersGrid
        workerStatusFilter={workerStatusFilter}
        sizeClassQueueName={sizeClassQueueName}
        invocationIds={invocationIds}
        cursor={cursor}
      />
    </PortalCard>
  );
};
