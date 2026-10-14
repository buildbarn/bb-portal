import { CalendarFilled } from "@ant-design/icons";
import { Row, Typography } from "antd";
import type React from "react";
import PlatformQueuesTable from "@/components/PlatformQueuesTable";
import { PortalCard } from "@/components/PortalCard";

interface Props {
  title: string;
}

export const SchedulerQueueChooserPage: React.FC<Props> = ({ title }) => {
  return (
    <PortalCard
      icon={<CalendarFilled />}
      titleBits={[<span key="title">{title}</span>]}
    >
      <Row>
        <Typography.Text>Choose a size class queue</Typography.Text>
      </Row>
      <Row>
        <PlatformQueuesTable />
      </Row>
    </PortalCard>
  );
};
