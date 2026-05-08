import {
  ClusterOutlined,
  DashboardOutlined,
  NodeCollapseOutlined,
  PieChartOutlined,
  ToolOutlined,
} from "@ant-design/icons";
import { Flex } from "antd";
import type { BazelInvocationMetricsActionSummaryFragment } from "@/graphql/__generated__/graphql";
import ActionCacheMissMetrics from "../ActionCacheMissMetrics";
import ActionCacheOverview from "../ActionCacheOverview";
import ActionRunnerMetrics from "../ActionRunnerMetrics";
import ActionTypeMetrics from "../ActionTypeMetrics";
import PortalCard from "../PortalCard";

type Props = {
  actionSummary: BazelInvocationMetricsActionSummaryFragment;
};

const ActionStatisticsDisplay: React.FC<Props> = ({ actionSummary }) => {
  return (
    <PortalCard
      type="inner"
      icon={<ToolOutlined />}
      titleBits={["Action Metrics"]}
    >
      <Flex vertical={false} gap="small" wrap={true}>
        <PortalCard
          type="inner"
          icon={<DashboardOutlined />}
          titleBits={["Action Cache Overview"]}
        >
          <ActionCacheOverview
            acStatistics={actionSummary?.actionCacheStatistics}
          />
        </PortalCard>

        <PortalCard
          type="inner"
          icon={<PieChartOutlined />}
          titleBits={["Action Cache Miss Breakdown"]}
        >
          <ActionCacheMissMetrics
            acStatistics={actionSummary?.actionCacheStatistics}
          />
        </PortalCard>
        <PortalCard
          type="inner"
          icon={<ClusterOutlined />}
          titleBits={["Action Runners Breakdown"]}
        >
          <ActionRunnerMetrics
            runnerMetrics={actionSummary.runnerCount ?? []}
          />
        </PortalCard>

        <PortalCard
          type="inner"
          icon={<NodeCollapseOutlined />}
          titleBits={["Action Types"]}
        >
          <ActionTypeMetrics actionData={actionSummary?.actionData} />
        </PortalCard>
      </Flex>
    </PortalCard>
  );
};

export default ActionStatisticsDisplay;
