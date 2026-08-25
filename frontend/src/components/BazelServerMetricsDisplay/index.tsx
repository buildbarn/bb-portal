import { DatabaseOutlined } from "@ant-design/icons";
import { Space, Statistic } from "antd";
import type React from "react";
import type {
  BazelInvocationMetricsCumulativeMetricsFragment,
  BazelInvocationMetricsPackageMetricsFragment,
} from "@/graphql/__generated__/graphql";
import PortalCard from "../PortalCard";

interface Props {
  packageMetrics?: BazelInvocationMetricsPackageMetricsFragment | null;
  cumulativeMetrics?: BazelInvocationMetricsCumulativeMetricsFragment | null;
  cardStyle?: React.CSSProperties;
}

export const BazelServerMetricsDisplay: React.FC<Props> = ({
  packageMetrics,
  cumulativeMetrics,
  cardStyle,
}) => {
  return (
    <PortalCard
      type="inner"
      icon={<DatabaseOutlined />}
      style={cardStyle}
      titleBits={["Package and Bazel Server Metrics"]}
    >
      <Space size="large">
        <Statistic
          title="Packages Loaded"
          value={packageMetrics?.packagesLoaded ?? 0}
        />
        <Statistic
          title="Server Analyses"
          value={cumulativeMetrics?.numAnalyses ?? 0}
        />
        <Statistic
          title="Server Builds"
          value={cumulativeMetrics?.numBuilds ?? 0}
        />
      </Space>
    </PortalCard>
  );
};
