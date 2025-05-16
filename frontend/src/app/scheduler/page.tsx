"use client";

import Content from "@/components/Content";
import PageDisabled from "@/components/PageDisabled";
import PortalCard from "@/components/PortalCard";
import SchedulerGrid from "@/components/SchedulerGrid";
import { FeatureType, isFeatureEnabled } from "@/utils/isFeatureEnabled";
import { CalendarFilled } from "@ant-design/icons";
import type React from "react";

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
          <SchedulerGrid />
        </PortalCard>
      }
    />
  );
};

export default Page;
