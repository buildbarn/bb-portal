import {
  ClusterOutlined,
  DashboardOutlined,
  NodeCollapseOutlined,
  NumberOutlined,
  PieChartOutlined,
  ToolOutlined,
} from "@ant-design/icons";
import { Flex, Space, Statistic, Typography } from "antd";
import type { BazelInvocationMetricsActionSummaryFragment } from "@/graphql/__generated__/graphql";
import ActionCacheMissMetrics from "../ActionCacheMissMetrics";
import ActionCacheOverview from "../ActionCacheOverview";
import ActionRunnerMetrics from "../ActionRunnerMetrics";
import ActionTypeMetrics from "../ActionTypeMetrics";
import { PortalCard } from "../PortalCard";
import { PortalTooltip } from "../PortalTooltip";

type Props = {
  actionSummary: BazelInvocationMetricsActionSummaryFragment;
};

const ACTION_OVERVIEW_CARD_STYLE: React.CSSProperties = {
  flex: "1 1 20rem",
};

const ACTION_CACHE_OVERVIEW_CARD_STYLE: React.CSSProperties = {
  flex: "2 1 40rem",
};

const HALF_ROW_CARD_STYLE: React.CSSProperties = {
  flex: "1 1 30rem",
};

const ActionStatisticsDisplay: React.FC<Props> = ({ actionSummary }) => {
  return (
    <PortalCard
      type="inner"
      icon={<ToolOutlined />}
      titleBits={["Action Metrics"]}
    >
      <Flex vertical={true} gap="small">
        <Flex vertical={false} gap="small" wrap={true}>
          <PortalCard
            type="inner"
            icon={<NumberOutlined />}
            style={ACTION_OVERVIEW_CARD_STYLE}
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
            style={ACTION_CACHE_OVERVIEW_CARD_STYLE}
            titleBits={[
              <Space key="title" orientation="horizontal">
                <Typography.Text>Local Action Cache Overview</Typography.Text>
                <PortalTooltip text="This shows the local Bazel action cache, and not remote caching. Remote execution cache hits can be viewed in the Action Runners Breakdown" />
              </Space>,
            ]}
          >
            <ActionCacheOverview
              acStatistics={actionSummary.actionCacheStatistics}
            />
          </PortalCard>
        </Flex>

        <Flex vertical={false} gap="small" wrap={true}>
          <PortalCard
            type="inner"
            icon={<PieChartOutlined />}
            style={HALF_ROW_CARD_STYLE}
            titleBits={["Action Cache Miss Breakdown"]}
          >
            <ActionCacheMissMetrics
              acStatistics={actionSummary.actionCacheStatistics}
            />
          </PortalCard>
          <PortalCard
            type="inner"
            icon={<ClusterOutlined />}
            style={HALF_ROW_CARD_STYLE}
            titleBits={["Action Runners Breakdown"]}
          >
            <ActionRunnerMetrics
              runnerMetrics={actionSummary.runnerCount ?? []}
            />
          </PortalCard>
        </Flex>

        <Flex vertical={false} gap="small" wrap={true}>
          <PortalCard
            type="inner"
            icon={<NodeCollapseOutlined />}
            style={HALF_ROW_CARD_STYLE}
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
            style={HALF_ROW_CARD_STYLE}
            titleBits={["Action Types Executed"]}
          >
            <ActionTypeMetrics
              actionData={actionSummary.actionData}
              countField="actionsExecuted"
            />
          </PortalCard>
        </Flex>
      </Flex>
    </PortalCard>
  );
};

export default ActionStatisticsDisplay;
