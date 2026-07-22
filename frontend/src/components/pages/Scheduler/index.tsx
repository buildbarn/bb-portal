import { CalendarFilled } from "@ant-design/icons";
import { Row } from "antd";
import type React from "react";
import PlatformQueuesTable from "@/components/PlatformQueuesTable";
import { PortalCard } from "@/components/PortalCard";
import { SchedulerStatistics } from "@/components/SchedulerStatistics";

export const SchedulerPage: React.FC = () => {
  return (
    <PortalCard
      icon={<CalendarFilled />}
      titleBits={[<span key="title">Scheduler</span>]}
    >
      <Row>
        <SchedulerStatistics />
      </Row>
      <Row>
        <PlatformQueuesTable />
      </Row>
    </PortalCard>
  );
};
