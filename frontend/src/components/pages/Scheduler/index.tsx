import { CalendarFilled } from "@ant-design/icons";
import { Row, Space, Typography } from "antd";
import type React from "react";
import PlatformQueuesTable from "@/components/PlatformQueuesTable";
import PortalCard from "@/components/PortalCard";
import { SchedulerStatistics } from "@/components/SchedulerStatistics";
import WorkerUtilizationChart from "@/components/WorkerUtilizationChart";

export const SchedulerPage: React.FC = () => {
  return (
    <PortalCard
      icon={<CalendarFilled />}
      titleBits={[<span key="title">Scheduler</span>]}
    >
      <Space direction="vertical" size="middle" style={{ display: "flex" }}>
        <Row>
          <SchedulerStatistics />
        </Row>
        <Row style={{ width: "100%" }}>
          <Space direction="vertical" style={{ width: "100%" }}>
            <Typography.Text strong>Worker Utilization</Typography.Text>
            <WorkerUtilizationChart />
          </Space>
        </Row>
        <Row>
          <PlatformQueuesTable />
        </Row>
      </Space>
    </PortalCard>
  );
};
