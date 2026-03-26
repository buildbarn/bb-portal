import { ExclamationCircleFilled } from "@ant-design/icons";
import { Link } from "@tanstack/react-router";
import { Descriptions, Space, Tag } from "antd";
import { CodeLink } from "@/components/CodeLink";
import type { OperationState } from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import type { OperationsFilterParams } from "@/routes/operations.index";
import themeStyles from "@/theme/theme.module.css";
import { BrowserPageType } from "@/types/BrowserPageType";
import { protobufToObjectWithTypeField } from "@/utils/protobufToObject";
import {
  readableDurationFromDates,
  readableDurationFromSeconds,
} from "@/utils/time";
import { generateBrowserSplat } from "@/utils/urlGenerator";
import OperationStatusTag from "../OperationStatusTag";
import { operationsStateToBrowserSplat } from "../OperationsGrid/utils";
import PropertyTagList from "../PropertyTagList";
import { historicalExecuteResponseDigestFromOperation } from "./utils";

interface Props {
  operation: OperationState;
}

const OperationStateDisplay: React.FC<Props> = ({ operation }) => {
  const invocationMetadata = operation.invocationName?.ids?.map((value) => {
    return JSON.stringify(protobufToObjectWithTypeField(value, false));
  });

  const historicalExecuteResponseBrowserPageParams =
    historicalExecuteResponseDigestFromOperation(operation);

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
            {invocationMetadata?.map((value) => {
              let metadataObject: OperationsFilterParams;
              try {
                metadataObject = JSON.parse(value);
              } catch {
                console.error(
                  "Failed to deserialize invocation metadata object",
                );
              }
              return (
                <li key={value}>
                  <Link
                    to="/operations"
                    search={{
                      filter: metadataObject,
                    }}
                  >
                    {value}
                  </Link>
                </li>
              );
            })}
          </ul>
        </Descriptions.Item>
        <Descriptions.Item label="Action digest">
          {operation.actionDigest && (
            <CodeLink
              text={`${operation.actionDigest.hash}-${operation.actionDigest.sizeBytes}`}
              link={{
                to: "/browser/$",
                params: {
                  _splat: operationsStateToBrowserSplat(operation),
                },
              }}
            />
          )}
        </Descriptions.Item>
        {historicalExecuteResponseBrowserPageParams && (
          <Descriptions.Item label="Historical execute response digest">
            <CodeLink
              text={historicalExecuteResponseBrowserPageParams.digest.hash}
              link={{
                to: "/browser/$",
                params: {
                  _splat: generateBrowserSplat(
                    historicalExecuteResponseBrowserPageParams.instanceName,
                    historicalExecuteResponseBrowserPageParams.digestFunction,
                    historicalExecuteResponseBrowserPageParams.digest,
                    BrowserPageType.HistoricalExecuteResponse,
                  ),
                },
              }}
            />
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
              Number.parseInt(operation.expectedDuration.seconds, 10),
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
              Status message: {operation.completed?.status?.message}
            </Tag>
          )}
        </Descriptions.Item>
      </Descriptions>
    </Space>
  );
};

export default OperationStateDisplay;
