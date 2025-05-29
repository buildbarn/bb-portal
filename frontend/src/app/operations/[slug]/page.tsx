"use client";

import Content from "@/components/Content";
import OperationDetails from "@/components/OperationDetails";
import PortalCard from "@/components/PortalCard";
import { CodeFilled } from "@ant-design/icons";
import { Space } from "antd";
import type React from "react";

interface PageParams {
  params: {
    slug: string;
  };
}

const Page: React.FC<PageParams> = ({ params }) => {
  const operationID = decodeURIComponent(params.slug);

  return (
    <Content
      content={
        <Space direction="vertical" size="middle" style={{ display: "flex" }}>
          <PortalCard
            icon={<CodeFilled />}
            titleBits={[<span key="title">{`Operation ${operationID}`}</span>]}
          >
            <OperationDetails operationID={operationID} />
          </PortalCard>
        </Space>
      }
    />
  );
};

export default Page;
