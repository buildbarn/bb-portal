import {
  ClusterOutlined,
  DashboardOutlined,
  NodeCollapseOutlined,
  NumberOutlined,
  PieChartOutlined,
  ToolOutlined,
} from "@ant-design/icons";
import { Flex, Space, Statistic } from "antd";
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
          icon={<NumberOutlined />}
          titleBits={["Action Overview"]}
        >
          <Space size="large">
            <Statistic
              title="Created"
              value={actionSummary.actionsCreated ?? 0}
            />
            <Statistic
              title="Created Without Aspects"
              value={actionSummary.actionsCreatedNotIncludingAspects ?? 0}
            />
            <Statistic
              title="Executed"
              value={actionSummary.actionsExecuted ?? 0}
            />
          </Space>
        </PortalCard>

        <PortalCard
          type="inner"
          icon={<DashboardOutlined />}
          titleBits={["Action Cache Overview"]}
        >
          <ActionCacheOverview
            acStatistics={actionSummary.actionCacheStatistics}
          />
        </PortalCard>

        <PortalCard
          type="inner"
          icon={<PieChartOutlined />}
          titleBits={["Action Cache Miss Breakdown"]}
        >
          <ActionCacheMissMetrics
            acStatistics={actionSummary.actionCacheStatistics}
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
          titleBits={["Action Types Created"]}
        >
          <ActionTypeMetrics
            actionData={actionSummary.actionData}
            countField="actionsCreated"
          />
        </PortalCard>

        <PortalCard
          type="inner"
          icon={<NodeCollapseOutlined />}
          titleBits={["Action Types Executed"]}
        >
          <ActionTypeMetrics
            actionData={actionSummary.actionData}
            countField="actionsExecuted"
          />
        </PortalCard>
      </Flex>
    </PortalCard>
  );
};

export default ActionStatisticsDisplay;
