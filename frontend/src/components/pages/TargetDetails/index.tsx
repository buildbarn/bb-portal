import {
  BorderInnerOutlined,
  DeploymentUnitOutlined,
  FieldTimeOutlined,
} from "@ant-design/icons";
import { Descriptions, Row, Space, Statistic, Typography } from "antd";
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
import { tableFiltersToGraphqlWhere } from "@/components/PageCursorTable/utils";
import { TargetDurationWarning } from "@/components/TargetDurationWarning";
import type {
  InvocationTargetDetailsFragment,
  InvocationTargetWhereInput,
  PageInfo,
  TargetDetailsFragment,
} from "@/graphql/__generated__/graphql";
import { readableDurationFromMilliseconds } from "@/utils/time";
import PortalCard from "../../PortalCard";
import { columns } from "./Columns";

interface Props {
  targetData: TargetDetailsFragment;
  onFilterChange: (where: InvocationTargetWhereInput[]) => void;
  getPaginationUpdateLink: GetPaginationUpdateLinkType;
  pageSize: number;
  pageInfo: PageInfo;
  invocationTargets: InvocationTargetDetailsFragment[];
}

export const TargetDetailsPage: React.FC<Props> = ({
  targetData,
  onFilterChange,
  getPaginationUpdateLink,
  pageSize,
  pageInfo,
  invocationTargets,
}) => {
  const averageDuration =
    targetData.invocationTargetsTotalDurationMillis / invocationTargets.length;

  return (
    <Space direction="vertical" size="middle" style={{ display: "flex" }}>
      <PortalCard
        icon={<DeploymentUnitOutlined />}
        titleBits={[<span key="title">Target Details</span>]}
      >
        <Space direction="vertical" size="middle" style={{ display: "flex" }}>
          <Descriptions column={1}>
            <Descriptions.Item label="Instance Name">
              {targetData.instanceName.name || "-"}
            </Descriptions.Item>
            <Descriptions.Item label="Target Kind">
              {targetData.targetKind || "-"}
            </Descriptions.Item>
            <Descriptions.Item label="Target Label">
              <Typography.Text copyable>
                {targetData.label || "-"}
              </Typography.Text>
            </Descriptions.Item>
            <Descriptions.Item label="Target Aspect">
              {targetData.aspect || "-"}
            </Descriptions.Item>
          </Descriptions>
          <Row>
            <Space size="large">
              <Statistic
                title={<TargetDurationWarning text="Average Duration" />}
                value={
                  (averageDuration &&
                    readableDurationFromMilliseconds(averageDuration, {
                      smallestUnit: "ms",
                    })) ||
                  "-"
                }
              />
              <Statistic
                title="Total Runs"
                value={invocationTargets.length || "-"}
              />
            </Space>
          </Row>
          <PortalCard
            icon={<FieldTimeOutlined />}
            titleBits={[
              <TargetDurationWarning
                key="title"
                text="Target Duration Over Time"
              />,
            ]}
          >
            <AreaChart
              width={800}
              height={250}
              data={invocationTargets}
              margin={{ top: 10, right: 30, left: 0, bottom: 0 }}
            >
              <defs>
                <linearGradient x1="0" y1="0" x2="0" y2="1">
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
                dataKey="durationInMs"
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
            <PageCursorTable<InvocationTargetDetailsFragment>
              rowKey="id"
              dataSource={invocationTargets}
              columns={columns}
              size="small"
              expandable={{
                rowExpandable: (record) =>
                  record.failureMessage !== null &&
                  record.failureMessage !== undefined &&
                  record.failureMessage !== "",
                expandedRowRender: (record) => (
                  <Space direction="vertical" style={{ paddingLeft: "48px" }}>
                    <strong>Failure Message:</strong>
                    <pre
                      style={{
                        whiteSpace: "pre-wrap",
                        overflowWrap: "break-word",
                      }}
                    >
                      {record.failureMessage}
                    </pre>
                  </Space>
                ),
              }}
              onChange={(_pagination, filters, _sorter, _extra) =>
                onFilterChange(tableFiltersToGraphqlWhere(columns, filters))
              }
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
