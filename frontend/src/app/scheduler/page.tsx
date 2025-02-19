"use client";

import Content from "@/components/Content";
import PortalCard from "@/components/PortalCard";
import SchedulerGrid from "@/components/SchedulerGrid";
import { CalendarFilled } from "@ant-design/icons";
import type React from "react";

const Page: React.FC = () => {
  return (
    <Content
      content={
        <PortalCard
          icon={<CalendarFilled />}
          titleBits={[<span key="title">Scheduler</span>]}
        >
          <SchedulerGrid />
        </PortalCard>
      }
    />
  );
};

export default Page;
