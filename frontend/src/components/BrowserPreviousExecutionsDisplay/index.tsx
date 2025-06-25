import type { Digest } from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import type { PreviousExecutionStats } from "@/lib/grpc-client/buildbarn/iscc/iscc";
import type { BrowserPageParams } from "@/types/BrowserPageType";
import { digestFunctionValueToString } from "@/utils/digestFunctionUtils";
import { Descriptions, Space, Typography } from "antd";
import Link from "next/link";
import PreviousExecutionsPlot from "../PreviousExecuteStatsPlot";
import SizeClassOutcome from "../SizeClassOutcome";

interface Props {
  browserParams: BrowserPageParams;
  reducedActionDigest: Digest;
  previousExecutionStats: PreviousExecutionStats;
  showTitle: boolean;
}

const BrowserPreviousExecutionsDisplay: React.FC<Props> = ({
  browserParams,
  previousExecutionStats,
  showTitle,
  reducedActionDigest,
}) => (
  <Space direction="vertical" size="middle" style={{ width: "100%" }}>
    {showTitle && (
      <Typography.Title level={2}>
        <Link
          href={`/browser/${
            browserParams.instanceName
          }/blobs/${digestFunctionValueToString(
            browserParams.digestFunction,
          )}/previous_execution_stats/${
            reducedActionDigest.hash
          }-${reducedActionDigest.sizeBytes}`}
          style={{ textDecoration: "underline" }}
        >
          Previous execution stats
        </Link>
      </Typography.Title>
    )}

    <Descriptions
      column={1}
      size="small"
      bordered
      styles={{ label: { width: "25%" }, content: { width: "75%" } }}
    >
      {previousExecutionStats.lastSeenFailure && (
        <Descriptions.Item label="Last seen failure">
          {previousExecutionStats.lastSeenFailure.toISOString()}
        </Descriptions.Item>
      )}
      {Object.entries(previousExecutionStats.sizeClasses).map((value) => (
        <Descriptions.Item
          key={value[0]}
          label={`Outcomes on size class ${value[0]}`}
        >
          <SizeClassOutcome sizeClassStats={value[1]} />
        </Descriptions.Item>
      ))}
    </Descriptions>
    <PreviousExecutionsPlot prevStats={previousExecutionStats} />
  </Space>
);

export default BrowserPreviousExecutionsDisplay;
