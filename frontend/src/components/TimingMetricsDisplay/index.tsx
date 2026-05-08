import { FieldTimeOutlined } from "@ant-design/icons";
import { Row, Space, Statistic } from "antd";
import type React from "react";
import type { BazelInvocationMetricsTimingMetricsFragment } from "@/graphql/__generated__/graphql";
import { readableDurationFromMilliseconds } from "@/utils/time";
import PortalCard from "../PortalCard";

interface Props {
  timingMetrics: BazelInvocationMetricsTimingMetricsFragment | undefined | null;
}

export const TimingMetricsDisplay: React.FC<Props> = ({ timingMetrics }) => {
  return (
    <PortalCard
      type="inner"
      titleBits={["Timing Metrics"]}
      icon={<FieldTimeOutlined />}
    >
      <Row>
        <Space size={"large"}>
          <Statistic
            title="Wall Time"
            value={readableDurationFromMilliseconds(
              timingMetrics?.wallTimeInMs ?? 0,
              { smallestUnit: "ms" },
            )}
          />
          <Statistic
            title="Analysis"
            value={readableDurationFromMilliseconds(
              timingMetrics?.analysisPhaseTimeInMs ?? 0,
              { smallestUnit: "ms" },
            )}
          />
          <Statistic
            title="CPU Time"
            value={readableDurationFromMilliseconds(
              timingMetrics?.cpuTimeInMs ?? 0,
              { smallestUnit: "ms" },
            )}
          />
          <Statistic
            title="Execution"
            value={readableDurationFromMilliseconds(
              timingMetrics?.executionPhaseTimeInMs ?? 0,
              { smallestUnit: "ms" },
            )}
          />
          <Statistic
            title="Actions Execution Start"
            value={readableDurationFromMilliseconds(
              timingMetrics?.actionsExecutionStartInMs ?? 0,
              { smallestUnit: "ms" },
            )}
          />
        </Space>
      </Row>
    </PortalCard>
  );
};
