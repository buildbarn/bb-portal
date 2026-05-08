import {
  ClusterOutlined,
  DashboardOutlined,
  NodeCollapseOutlined,
  PieChartOutlined,
} from "@ant-design/icons";
import { Flex } from "antd";
import type { BazelInvocationActionSummaryFragment } from "@/graphql/__generated__/graphql";
import ActionCacheMissMetrics from "../ActionCacheMissMetrics";
import ActionCacheOverview from "../ActionCacheOverview";
import ActionRunnerMetrics from "../ActionRunnerMetrics";
import ActionTypeMetrics from "../ActionTypeMetrics";
import PortalCard from "../PortalCard";

type Props = {
  actionSummary: BazelInvocationActionSummaryFragment;
};

const ActionStatisticsDisplay: React.FC<Props> = ({ actionSummary }) => {
  return (
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
        titleBits={["Cache Miss Breakdown"]}
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
        <ActionRunnerMetrics runnerMetrics={actionSummary.runnerCount ?? []} />
      </PortalCard>

      <PortalCard
        type="inner"
        icon={<NodeCollapseOutlined />}
        titleBits={["Action Types"]}
      >
        <ActionTypeMetrics actionData={actionSummary?.actionData} />
      </PortalCard>
    </Flex>
  );
};

export default ActionStatisticsDisplay;
