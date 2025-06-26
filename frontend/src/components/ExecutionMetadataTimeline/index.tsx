"use client";

import type { ExecutedActionMetadata } from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import { readableDurationFromDates } from "@/utils/time";
import { ClockCircleOutlined } from "@ant-design/icons";
import { Flex, Space, Typography } from "antd";
import type React from "react";

interface Params {
  executionMetadata: ExecutedActionMetadata;
}

const formatTimelineElement = (
  timestamp: Date,
  previous: Date | undefined = undefined,
) => {
  if (timestamp.getTime() === previous?.getTime()) {
    return null;
  }

  return (
    <Space direction="horizontal">
      <Typography.Text strong>{timestamp.toISOString()}</Typography.Text>
      {previous && (
        <>
          <ClockCircleOutlined />
          <Typography.Text type="secondary">
            {" "}
            (+
            {readableDurationFromDates(previous, timestamp, {
              smallestUnit: "ms",
            })}
            )
          </Typography.Text>
        </>
      )}
    </Space>
  );
};

const ExecutionMetadataTimeline: React.FC<Params> = ({
  executionMetadata: em,
}) => {
  return (
    <Flex vertical>
      {em.queuedTimestamp && formatTimelineElement(em.queuedTimestamp)}
      <Typography.Text>Action added to the queue.</Typography.Text>
      {em.workerStartTimestamp &&
        formatTimelineElement(em.workerStartTimestamp, em.queuedTimestamp)}
      <Typography.Text>Worker received the action.</Typography.Text>
      {em.inputFetchStartTimestamp &&
        formatTimelineElement(
          em.inputFetchStartTimestamp,
          em.workerStartTimestamp,
        )}
      <Typography.Text>Worker started fetching action inputs.</Typography.Text>
      {em.inputFetchCompletedTimestamp &&
        formatTimelineElement(
          em.inputFetchCompletedTimestamp,
          em.inputFetchStartTimestamp,
        )}
      <Typography.Text>Worker finished fetching action inputs.</Typography.Text>
      {em.executionStartTimestamp &&
        formatTimelineElement(
          em.executionStartTimestamp,
          em.inputFetchCompletedTimestamp,
        )}
      <Typography.Text>
        Worker started executing the action command.
      </Typography.Text>
      {em.executionCompletedTimestamp &&
        formatTimelineElement(
          em.executionCompletedTimestamp,
          em.executionStartTimestamp,
        )}
      <Typography.Text>
        Worker completed executing the action command.
      </Typography.Text>
      {em.outputUploadStartTimestamp &&
        formatTimelineElement(
          em.outputUploadStartTimestamp,
          em.executionCompletedTimestamp,
        )}
      <Typography.Text>
        Worker started uploading action outputs.
      </Typography.Text>
      {em.outputUploadCompletedTimestamp &&
        formatTimelineElement(
          em.outputUploadCompletedTimestamp,
          em.outputUploadStartTimestamp,
        )}
      <Typography.Text>
        Worker completed uploading action outputs.
      </Typography.Text>
      {em.workerCompletedTimestamp &&
        formatTimelineElement(
          em.workerCompletedTimestamp,
          em.outputUploadCompletedTimestamp,
        )}
      <Typography.Text>
        Worker completed the action, including all stages.
      </Typography.Text>
    </Flex>
  );
};

export default ExecutionMetadataTimeline;
