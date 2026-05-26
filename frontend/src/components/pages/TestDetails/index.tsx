import {
  BorderInnerOutlined,
  ExperimentFilled,
  FieldTimeOutlined,
} from "@ant-design/icons";
import { Descriptions, Space, Typography } from "antd";
import type React from "react";
import {
  Area,
  AreaChart,
  CartesianGrid,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import { PageCursorTable } from "@/components/PageCursorTable";
import type { GetPaginationUpdateLinkType } from "@/components/PageCursorTable/types";
import type {
  PageInfo,
  TestSummaryTargetDetailsFragment,
} from "@/graphql/__generated__/graphql";
import PortalCard from "../../PortalCard";
import { columns, type TestDetailsRowType } from "./columns";

interface Props {
  target: TestSummaryTargetDetailsFragment;
  testSummaries: TestDetailsRowType[];
  getPaginationUpdateLink: GetPaginationUpdateLinkType;
  pageSize: number;
  pageInfo: PageInfo;
}

export const TestDetailsPage: React.FC<Props> = ({
  target,
  testSummaries,
  getPaginationUpdateLink,
  pageSize,
  pageInfo,
}) => {
  return (
    <Space direction="vertical" size="middle" style={{ display: "flex" }}>
      <PortalCard
        icon={<ExperimentFilled />}
        titleBits={[<span key="title">Test Details</span>]}
      >
        <Space direction="vertical" size="middle" style={{ display: "flex" }}>
          <Descriptions column={1}>
            <Descriptions.Item label="Instance Name">
              {target.instanceName.name || "-"}
            </Descriptions.Item>
            <Descriptions.Item label="Target Kind">
              {target.targetKind || "-"}
            </Descriptions.Item>
            <Descriptions.Item label="Target Label">
              <Typography.Text copyable>{target.label || "-"}</Typography.Text>
            </Descriptions.Item>
            <Descriptions.Item label="Target Aspect">
              {target.aspect || "-"}
            </Descriptions.Item>
          </Descriptions>
          <PortalCard
            icon={<FieldTimeOutlined />}
            titleBits={["Test Duration Over Time"]}
          >
            <AreaChart
              width={900}
              height={250}
              data={testSummaries}
              margin={{ top: 10, right: 30, left: 0, bottom: 0 }}
            >
              <defs>
                <linearGradient id="colorUv" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#8884d8" stopOpacity={0.8} />
                  <stop offset="95%" stopColor="#8884d8" stopOpacity={0} />
                </linearGradient>
              </defs>
              <XAxis />
              <YAxis />
              <CartesianGrid strokeDasharray="3 3" vertical={false} />
              <Tooltip />
              <Area
                type="monotone"
                dataKey="totalRunDurationInMs"
                stroke="#8884d8"
                fillOpacity={1}
                fill="url(#colorUv)"
              />
            </AreaChart>
          </PortalCard>
          <PortalCard
            icon={<BorderInnerOutlined />}
            titleBits={["Per Invocation Details"]}
          >
            <PageCursorTable<TestDetailsRowType>
              rowKey="id"
              dataSource={testSummaries}
              columns={columns}
              size="small"
              pageInfo={pageInfo}
              getPaginationUpdateLink={getPaginationUpdateLink}
              pageSize={pageSize}
            />
          </PortalCard>
        </Space>
      </PortalCard>
    </Space>
  );
};
