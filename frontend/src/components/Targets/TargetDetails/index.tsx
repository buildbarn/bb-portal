import {
  BorderInnerOutlined,
  FieldTimeOutlined,
} from "@ant-design/icons/lib/icons";
import { useQuery } from "@apollo/client";
import { Descriptions, Row, Space, Statistic, Typography } from "antd";
import type { FilterValue } from "antd/es/table/interface";
import React from "react";
import {
  Area,
  AreaChart,
  CartesianGrid,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import {
  CursorTable,
  getNewPaginationVariables,
} from "@/components/CursorTable";
import type { PaginationVariables } from "@/components/CursorTable/types";
import { TargetDurationWarning } from "@/components/TargetDurationWarning";
import type {
  InvocationTargetAbortReason,
  InvocationTargetWhereInput,
} from "@/graphql/__generated__/graphql";
import { parseGraphqlEdgeList } from "@/utils/parseGraphqlEdgeList";
import { readableDurationFromMilliseconds } from "@/utils/time";
import PortalAlert from "../../PortalAlert";
import PortalCard from "../../PortalCard";
import { columns, type InvocationTargetRowType } from "./Columns";
import { GET_TARGET_DETAILS } from "./graphql";

interface Props {
  instanceName: string;
  label: string;
  aspect: string;
  targetKind: string;
}

export const TargetDetails: React.FC<Props> = ({
  instanceName,
  label,
  aspect,
  targetKind,
}) => {
  const [paginationVariables, setPaginationVariables] =
    React.useState<PaginationVariables>(getNewPaginationVariables());
  const [filterVariables, setFilterVariables] =
    React.useState<InvocationTargetWhereInput>({});

  const { data, loading, error } = useQuery(GET_TARGET_DETAILS, {
    variables: {
      instanceName: instanceName,
      label: label,
      aspect: aspect,
      targetKind: targetKind,
      where: filterVariables,
      ...paginationVariables,
    },
    fetchPolicy: "network-only",
  });

  const onFilterChange = (filters: Record<string, FilterValue | null>) => {
    const newFilters: InvocationTargetWhereInput[] = [];
    Object.entries(filters).forEach(([key, value]) => {
      if (value && value.length > 0) {
        switch (key) {
          case "success":
            if (value.length === 1) {
              newFilters.push({ success: value[0] as boolean });
            }
            break;
          case "abort-reason":
            newFilters.push({
              abortReasonIn: value as InvocationTargetAbortReason[],
            });
            break;
        }
      }
    });
    setFilterVariables({ and: newFilters });
  };

  if (error) {
    return (
      <PortalAlert
        showIcon
        type="error"
        message="Error fetching target details"
        description={
          error?.message ||
          "Unknown error occurred while fetching data from the server."
        }
      />
    );
  }

  const invocationTargetsCount = data?.getTarget?.invocationTargets.totalCount;
  const invocationTargetsTotalDurationMillis =
    data?.getTarget?.invocationTargetsTotalDurationMillis;
  const averageDuration =
    invocationTargetsCount && invocationTargetsTotalDurationMillis
      ? invocationTargetsTotalDurationMillis / invocationTargetsCount
      : undefined;

  const invocationTargets = parseGraphqlEdgeList<InvocationTargetRowType>(
    data?.getTarget?.invocationTargets,
  );

  return (
    <Space direction="vertical" size="middle" style={{ display: "flex" }}>
      <Descriptions column={1}>
        <Descriptions.Item label="Instance Name">
          {instanceName || "-"}
        </Descriptions.Item>
        <Descriptions.Item label="Target Kind">
          {targetKind || "-"}
        </Descriptions.Item>
        <Descriptions.Item label="Target Label">
          <Typography.Text copyable>{label || "-"}</Typography.Text>
        </Descriptions.Item>
        <Descriptions.Item label="Target Aspect">
          {aspect || "-"}
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
          <Statistic title="Total Runs" value={invocationTargetsCount || "-"} />
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
        <CursorTable<InvocationTargetRowType>
          rowKey="id"
          loading={loading}
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
                <pre>{record.failureMessage}</pre>
              </Space>
            ),
          }}
          onChange={(_pagination, filters, _sorter, _extra) =>
            onFilterChange(filters)
          }
          pagination={{
            position: "bottom",
            justify: "end",
            size: "middle",
          }}
          pageInfo={{
            startCursor:
              data?.getTarget?.invocationTargets.pageInfo.startCursor || "",
            endCursor:
              data?.getTarget?.invocationTargets.pageInfo.endCursor || "",
            hasNextPage:
              data?.getTarget?.invocationTargets.pageInfo.hasNextPage || false,
            hasPreviousPage:
              data?.getTarget?.invocationTargets.pageInfo.hasPreviousPage ||
              false,
          }}
          paginationVariables={paginationVariables}
          setPaginationVariables={setPaginationVariables}
        />
      </PortalCard>
    </Space>
  );
};
