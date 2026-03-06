import React from 'react';
import WorkersGrid from '@/components/WorkersGrid';
import type { WorkerSearchParams } from '@/routes/scheduler.worker';
import Content from '@/components/Content';
import PortalCard from '@/components/PortalCard';
import { CalendarFilled } from '@ant-design/icons';

export const SchedulerWorkersPage: React.FC<WorkerSearchParams> = ({
  workerStatusFilter,
  sizeClassQueueName,
  cursor
}) => {
  return <Content
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
};
