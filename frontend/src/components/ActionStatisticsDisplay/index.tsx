import type {
  ActionSummary,
  RunnerCount,
} from "@/graphql/__generated__/graphql";
import {
  ClusterOutlined,
  DashboardOutlined,
  NodeCollapseOutlined,
  PieChartOutlined,
} from "@ant-design/icons";
import { Flex } from "antd";
import ActionCacheMissMetrics from "../ActionCacheMissMetrics";
import ActionCacheOverview from "../ActionCacheOverview";
import ActionRunnerMetrics from "../ActionRunnerMetrics";
import ActionTypeMetrics from "../ActionTypeMetrics";
import PortalCard from "../PortalCard";

type Props = {
  runnerMetrics: RunnerCount[];
  actionData?: ActionSummary;
};

const ActionStatisticsDisplay: React.FC<Props> = ({
  runnerMetrics,
  actionData,
}) => {
  return (
    <Flex vertical={false} gap="small" wrap={true}>
      <PortalCard
        type="inner"
        icon={<DashboardOutlined />}
        titleBits={["Action Cache Overview"]}
      >
        <ActionCacheOverview acStatistics={actionData?.actionCacheStatistics} />
      </PortalCard>

      <PortalCard
        type="inner"
        icon={<PieChartOutlined />}
        titleBits={["Cache Miss Breakdown"]}
      >
        <ActionCacheMissMetrics
          acStatistics={actionData?.actionCacheStatistics}
        />
      </PortalCard>
      <PortalCard
        type="inner"
        icon={<ClusterOutlined />}
        titleBits={["Action Runners Breakdown"]}
      >
        <ActionRunnerMetrics runnerMetrics={runnerMetrics} />
      </PortalCard>

      <PortalCard
        type="inner"
        icon={<NodeCollapseOutlined />}
        titleBits={["Action Types"]}
      >
        <ActionTypeMetrics actionData={actionData?.actionData} />
      </PortalCard>
    </Flex>
  );
};

export default ActionStatisticsDisplay;
