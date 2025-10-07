import type { PerSizeClassStats } from "@/lib/grpc-client/buildbarn/iscc/iscc";
import { readableDurationFromProtobufDuration } from "@/utils/time";
import { Space, Typography } from "antd";
import SizeClassOutcomeTag from "../SizeClassOutcomeTag";

interface Props {
  sizeClassStats: PerSizeClassStats;
}

const SizeClassOutcome: React.FC<Props> = ({ sizeClassStats }) => {
  return (
    <Space direction="vertical" size="small">
      <Space direction="horizontal" size="small" wrap>
        {sizeClassStats.previousExecutions.map((val, index) => {
          if (val.succeeded) {
            return (
              // biome-ignore lint/suspicious/noArrayIndexKey: We have nothing better to use
              <SizeClassOutcomeTag color="success" key={index}>
                Succeeded: {readableDurationFromProtobufDuration(val.succeeded)}
              </SizeClassOutcomeTag>
            );
          }
          if (val.timedOut) {
            return (
              // biome-ignore lint/suspicious/noArrayIndexKey: We have nothing better to use
              <SizeClassOutcomeTag color="warning" key={index}>
                Timed out: {readableDurationFromProtobufDuration(val.timedOut)}
              </SizeClassOutcomeTag>
            );
          }
          if (val.failed) {
            return (
              // biome-ignore lint/suspicious/noArrayIndexKey: We have nothing better to use
              <SizeClassOutcomeTag color="error" key={index}>
                Failed
              </SizeClassOutcomeTag>
            );
          }
        })}
      </Space>
      <Typography.Title level={5}>
        {`Initial PageRank probability: ${sizeClassStats.initialPageRankProbability}`}
      </Typography.Title>
    </Space>
  );
};

export default SizeClassOutcome;
