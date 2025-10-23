import type { OperationState } from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import themeStyles from "@/theme/theme.module.css";
import { protobufToObjectWithTypeField } from "@/utils/protobufToObject";
import {
  readableDurationFromDates,
  readableDurationFromSeconds,
} from "@/utils/time";
import { ExclamationCircleFilled } from "@ant-design/icons";
import { Descriptions, Space, Tag } from "antd";
import Link from "next/link";
import OperationStatusTag from "../OperationStatusTag";
import { operationsStateToActionPageUrl } from "../OperationsGrid/utils";
import PropertyTagList from "../PropertyTagList";
import {
  historicalExecuteResponseDigestFromUrl,
  historicalExecuteResponseUrlFromOperation,
} from "./utils";

interface Props {
  operation: OperationState;
}

const OperationStateDisplay: React.FC<Props> = ({ operation }) => {
  const invocationMetadata = operation.invocationName?.ids?.map((value) => {
    return JSON.stringify(protobufToObjectWithTypeField(value, false));
  });

  const historical_execute_response_url =
    historicalExecuteResponseUrlFromOperation(operation);
  const historical_execute_response_digest =
    historicalExecuteResponseDigestFromUrl(historical_execute_response_url);

  return (
    <Space direction="vertical" size="middle" style={{ display: "flex" }}>
      <Descriptions column={1} size="small" bordered>
        <Descriptions.Item label="Instance name prefix">
          {
            operation.invocationName?.sizeClassQueueName?.platformQueueName
              ?.instanceNamePrefix
          }
        </Descriptions.Item>
        <Descriptions.Item label="Instance name suffix">
          {operation.instanceNameSuffix}
        </Descriptions.Item>
        <Descriptions.Item label="Platform properties">
          <PropertyTagList
            propertyList={
              operation.invocationName?.sizeClassQueueName?.platformQueueName
                ?.platform?.properties
            }
          />
        </Descriptions.Item>
        <Descriptions.Item label="Size class">
          {operation.invocationName?.sizeClassQueueName?.sizeClass}
        </Descriptions.Item>
        <Descriptions.Item label="Invocation IDs">
          <ul>
            {invocationMetadata?.map((value) => (
              <li key={value}>
                <Link
                  href={{
                    pathname: "/operations",
                    query: {
                      filter_invocation_id: value,
                    },
                  }}
                >
                  {value}
                </Link>
              </li>
            ))}
          </ul>
        </Descriptions.Item>
        <Descriptions.Item label="Action digest">
          {operation.actionDigest && (
            <Link
              href={operationsStateToActionPageUrl(operation) || ""}
            >{`${operation.actionDigest.hash}-${operation.actionDigest.sizeBytes}`}</Link>
          )}
        </Descriptions.Item>
        {historical_execute_response_url &&
          historical_execute_response_digest && (
            <Descriptions.Item label="Historical execute response digest">
              <Link href={historical_execute_response_url}>
                {historical_execute_response_digest}
              </Link>
            </Descriptions.Item>
          )}
        <Descriptions.Item label="Timeout">
          {operation.timeout &&
            readableDurationFromDates(new Date(), operation.timeout, {
              precision: 1,
              smallestUnit: "s",
            })}
        </Descriptions.Item>
        <Descriptions.Item label="Target ID">
          {operation.targetId}
        </Descriptions.Item>
        <Descriptions.Item label="Priority">
          {operation.priority}
        </Descriptions.Item>
        <Descriptions.Item label="Expected duration">
          {operation.expectedDuration &&
            readableDurationFromSeconds(
              Number.parseInt(operation.expectedDuration.seconds),
              { precision: 1, smallestUnit: "s" },
            )}
        </Descriptions.Item>
        <Descriptions.Item label="Age">
          {operation.queuedTimestamp &&
            readableDurationFromDates(operation.queuedTimestamp, new Date(), {
              precision: 1,
              smallestUnit: "s",
            })}
        </Descriptions.Item>
        <Descriptions.Item label="Stage">
          <OperationStatusTag operation={operation} />
          {operation.completed?.status?.message && (
            <Tag
              icon={<ExclamationCircleFilled />}
              color="default"
              className={themeStyles.tag}
            >
              <>Status message: {operation.completed?.status?.message}</>
            </Tag>
          )}
        </Descriptions.Item>
      </Descriptions>
    </Space>
  );
};

export default OperationStateDisplay;
