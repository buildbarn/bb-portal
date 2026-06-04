import { BorderInnerOutlined, DeploymentUnitOutlined } from "@ant-design/icons";
import { Descriptions, Space, Statistic, Typography } from "antd";
import type React from "react";
import { PageCursorTable } from "@/components/PageCursorTable";
import type { GetPaginationUpdateLinkType } from "@/components/PageCursorTable/types";
import { tableFiltersToGraphqlWhere } from "@/components/PageCursorTable/utils";
import type {
  InvocationTargetDetailsFragment,
  InvocationTargetWhereInput,
  PageInfo,
  TargetDetailsFragment,
} from "@/graphql/__generated__/graphql";
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
          <Statistic
            title="Total Runs"
            value={invocationTargets.length || "-"}
          />
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
