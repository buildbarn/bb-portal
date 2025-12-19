import {
  BorderInnerOutlined,
  FieldTimeOutlined,
} from "@ant-design/icons/lib/icons";
import { useQuery } from "@apollo/client";
import { Descriptions, Space, Typography } from "antd";
import React, { useMemo } from "react";
import {
  Area,
  AreaChart,
  CartesianGrid,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import type { GetTestDetailsQuery } from "@/graphql/__generated__/graphql";
import { parseGraphqlEdgeList } from "@/utils/parseGraphqlEdgeList";
import { CursorTable, getNewPaginationVariables } from "../CursorTable";
import type { PaginationVariables } from "../CursorTable/types";
import PortalAlert from "../PortalAlert";
import PortalCard from "../PortalCard";
import { columns, type TestDetailsRowType } from "./columns";
import { GET_TEST_DETAILS } from "./graphql";

interface Props {
  instanceName: string;
  label: string;
  aspect: string;
  targetKind: string;
}

export const TestDetails: React.FC<Props> = ({
  instanceName,
  label,
  aspect,
  targetKind,
}) => {
  const [paginationVariables, setPaginationVariables] =
    React.useState<PaginationVariables>(getNewPaginationVariables());

  const { data, loading, error } = useQuery<GetTestDetailsQuery>(
    GET_TEST_DETAILS,
    {
      variables: {
        where: {
          hasInvocationTargetWith: {
            hasTargetWith: {
              hasInstanceNameWith: {
                name: instanceName,
              },
              label: label,
              aspect: aspect,
              targetKind: targetKind,
            },
          },
        },
        ...paginationVariables,
      },
      fetchPolicy: "network-only",
    },
  );

  const testSummaries: TestDetailsRowType[] = useMemo(() => {
    return parseGraphqlEdgeList(data?.findTestSummaries).map((ts) => {
      console.table(ts);
      var cachedLocally: boolean | null = null;
      var cachedRemotely: boolean | null = null;
      if (ts.testResults !== null && ts.testResults !== undefined) {
        if (ts.testResults.every((tr) => tr.cachedLocally === true)) {
          cachedLocally = true;
        } else if (ts.testResults.some((tr) => tr.cachedLocally === false)) {
          cachedLocally = false;
        }
        if (ts.testResults.every((tr) => tr.cachedRemotely === true)) {
          cachedRemotely = true;
        } else if (ts.testResults.some((tr) => tr.cachedRemotely === false)) {
          cachedRemotely = false;
        }
      }
      return {
        ...ts,
        cachedLocally: cachedLocally,
        cachedRemotely: cachedRemotely,
      };
    });
  }, [data]);

  if (error) {
    return (
      <PortalAlert
        showIcon
        type="error"
        message="Error fetching test details"
        description={
          error?.message ||
          "Unknown error occurred while fetching data from the server."
        }
      />
    );
  }

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
        <CursorTable<TestDetailsRowType>
          rowKey="id"
          loading={loading}
          dataSource={testSummaries}
          columns={columns}
          size="small"
          pagination={{
            position: "bottom",
            justify: "end",
            size: "middle",
          }}
          pageInfo={{
            startCursor: data?.findTestSummaries.pageInfo.startCursor || "",
            endCursor: data?.findTestSummaries.pageInfo.endCursor || "",
            hasNextPage: data?.findTestSummaries.pageInfo.hasNextPage || false,
            hasPreviousPage:
              data?.findTestSummaries.pageInfo.hasPreviousPage || false,
          }}
          paginationVariables={paginationVariables}
          setPaginationVariables={setPaginationVariables}
        />
      </PortalCard>
    </Space>
  );
};
export default TestDetails;
