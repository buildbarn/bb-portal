"use client";

import Content from "@/components/Content";
import OperationsGrid from "@/components/OperationsGrid";
import PortalCard from "@/components/PortalCard";
import useScreenSize from "@/utils/screen";
import { CodeFilled } from "@ant-design/icons";
import type React from "react";

const Page: React.FC = () => {
  const screenSize = useScreenSize();
  return (
    <Content
      content={
        <PortalCard
          icon={<CodeFilled />}
          titleBits={[<span key="title">Operations</span>]}
        >
          <OperationsGrid />
        </PortalCard>
      }
    />
  );
};

export default Page;
