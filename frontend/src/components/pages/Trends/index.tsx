import { ClockCircleFilled, LineChartOutlined } from "@ant-design/icons";
import { useQuery } from "@apollo/client/react";
import { Row, Space, Statistic } from "antd";
import type React from "react";
import { useState } from "react";
import {
  Area,
  AreaChart,
  CartesianGrid,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import Content from "@/components/Content";
import PortalCard from "@/components/PortalCard";
import type {
  BazelInvocationNodeFragment,
  FindBuildTimesQueryVariables,
} from "@/graphql/__generated__/graphql";
import { readableDurationFromMilliseconds } from "@/utils/time";
import FIND_BUILD_DURATIONS from "./index.graphql";

export const TrendsPage: React.FC = () => {
  const [variables, _setVariables] = useState<FindBuildTimesQueryVariables>({
    first: 1000,
  });

  const { loading, data, previousData, error } = useQuery(
    FIND_BUILD_DURATIONS,
    {
      variables,
      pollInterval: 120000,
      fetchPolicy: "cache-and-network",
    },
  );

  const activeData = loading ? previousData : data;
  let _emptyText = "No builds match the specified search criteria";
  let dataSource: BazelInvocationNodeFragment[] = [];

  if (error) {
    _emptyText = error.message;
    dataSource = [];
  } else {
    const buildTimes =
      activeData?.findBazelInvocations.edges?.flatMap((edge) => edge?.node) ??
      [];
    dataSource = buildTimes.filter(
      (x): x is BazelInvocationNodeFragment => !!x,
    );
  }

  interface graphPoint {
    invocationId: string;
    from: string;
    to: string;
    duration: number;
  }

  const dataPoints: graphPoint[] = [];

  dataSource.forEach((x) => {
    var point: graphPoint = {
      invocationId: x.invocationID,
      from: x.startedAt,
      to: x.endedAt,
      duration: new Date(x.endedAt).getTime() - new Date(x.startedAt).getTime(),
    };
    // if there are empty/nil dates they get set to max epoch start time
    // which throws the graph off.
    if (point.duration > 0) {
      dataPoints.push(point);
    }
  });

  var avg: number =
    dataPoints.reduce((sum, item) => sum + item.duration, 0) /
    dataPoints.length;
  var medianVals = dataPoints.map((x) => x.duration).sort((a, b) => a - b);
  var medianMid = Math.floor(dataPoints.length / 2);
  var median: number;

  if (medianVals.length % 2 === 0) {
    median = (medianVals[medianMid - 1] + medianVals[medianMid]) / 2;
  } else {
    median = medianVals[medianMid];
  }

  var max: number = Math.max(...dataPoints.map((x) => x.duration));
  var min: number = Math.min(...dataPoints.map((x) => x.duration));

  return (
    <Content
      content={
        <Space direction="vertical" size="middle" style={{ display: "flex" }}>
          <PortalCard
            icon={<LineChartOutlined />}
            titleBits={[<span key="title">Trends</span>]}
          >
            <PortalCard
              type="inner"
              icon={<ClockCircleFilled />}
              titleBits={[
                <span key="invocation-durations">Invocation Durations</span>,
              ]}
            >
              <Row>
                <Space size="large">
                  <Statistic
                    title="Total"
                    value={data?.findBazelInvocations.totalCount}
                    valueStyle={{ color: "#82ca9d" }}
                  />
                  <Statistic
                    title="Average"
                    value={
                      dataPoints.length === 0
                        ? "-"
                        : readableDurationFromMilliseconds(avg, {
                            smallestUnit: "ms",
                          })
                    }
                    valueStyle={{ color: "#82ca9d" }}
                  />
                  <Statistic
                    title="Median"
                    value={
                      dataPoints.length === 0
                        ? "-"
                        : readableDurationFromMilliseconds(median, {
                            smallestUnit: "ms",
                          })
                    }
                    valueStyle={{ color: "#8884d8" }}
                  />
                  <Statistic
                    title="Max"
                    value={
                      dataPoints.length === 0
                        ? "-"
                        : readableDurationFromMilliseconds(max, {
                            smallestUnit: "ms",
                          })
                    }
                  />
                  <Statistic
                    title="Min"
                    value={
                      dataPoints.length === 0
                        ? "-"
                        : readableDurationFromMilliseconds(min, {
                            smallestUnit: "ms",
                          })
                    }
                  />
                </Space>
              </Row>
              <Row>
                <AreaChart
                  width={1500}
                  height={250}
                  data={dataPoints}
                  margin={{ top: 10, right: 30, left: 0, bottom: 0 }}
                >
                  <defs>
                    <linearGradient id="colorUv" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="5%" stopColor="#8884d8" stopOpacity={0.8} />
                      <stop offset="95%" stopColor="#8884d8" stopOpacity={0} />
                    </linearGradient>
                  </defs>
                  <XAxis dataKey="name" />
                  <YAxis />
                  <CartesianGrid vertical={false} strokeDasharray="3 3" />
                  <Tooltip />
                  <Area
                    type="monotone"
                    dataKey="duration"
                    stroke="#8884d8"
                    fillOpacity={1}
                    fill="url(#colorUv)"
                  />
                </AreaChart>
              </Row>
            </PortalCard>
          </PortalCard>
        </Space>
      }
    />
  );
};
