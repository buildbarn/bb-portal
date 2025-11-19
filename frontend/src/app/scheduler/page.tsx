"use client";

import { CalendarFilled } from "@ant-design/icons";
import { Row, Space } from "antd";
import type React from "react";
import Content from "@/components/Content";
import PageDisabled from "@/components/PageDisabled";
import PlatformQueuesTable from "@/components/PlatformQueuesTable";
import PortalCard from "@/components/PortalCard";
import { SchedulerStatistics } from "@/components/SchedulerStatistics";
import { FeatureType, isFeatureEnabled } from "@/utils/isFeatureEnabled";

const Page: React.FC = () => {
  if (!isFeatureEnabled(FeatureType.SCHEDULER)) {
    return <PageDisabled />;
  }

  return (
    <Content
      content={
        <PortalCard
          icon={<CalendarFilled />}
          titleBits={[<span key="title">Scheduler</span>]}
        >
          <Space direction="vertical" size="middle" style={{ display: "flex" }}>
            <Row>
              <SchedulerStatistics />
            </Row>
            <Row>
              <PlatformQueuesTable />
            </Row>
          </Space>
        </PortalCard>
      }
    />
  );
};

export default Page;
