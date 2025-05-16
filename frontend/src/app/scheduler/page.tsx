"use client";

import Content from "@/components/Content";
import PortalCard from "@/components/PortalCard";
import SchedulerGrid from "@/components/SchedulerGrid";
import { FeatureType, isFeatureEnabled } from "@/utils/isFeatureEnabled";
import { CalendarFilled } from "@ant-design/icons";
import { notFound } from "next/navigation";
import type React from "react";

const Page: React.FC = () => {
  if (!isFeatureEnabled(FeatureType.SCHEDULER)) {
    return notFound();
  }

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
