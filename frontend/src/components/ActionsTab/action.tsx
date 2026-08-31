import { Descriptions, Flex, Space, Typography } from "antd";
import type { BazelInvocationActionExecutionFragment } from "@/graphql/__generated__/graphql";
import {
  type GraphqlFile,
  generateActionUrlFromGraphqlDigest,
  generateFileUrlFromGraphqlFile,
} from "@/utils/urlGenerator";
import { getActionCacheStatus } from "./cache";
import { getActionExecutionKind } from "./execution";

interface Props {
  action: BazelInvocationActionExecutionFragment;
}

interface OutputLinkProps {
  file?: GraphqlFile | null;
  children: React.ReactNode;
}

const OutputLink: React.FC<OutputLinkProps> = ({ file, children }) => {
  return file ? (
    <Typography.Link href={generateFileUrlFromGraphqlFile(file)}>
      {children}
    </Typography.Link>
  ) : (
    children
  );
};

export const ActionDetails: React.FC<Props> = ({ action }) => {
  const actionHref = generateActionUrlFromGraphqlDigest(action.actionDigest);
  const cacheStatus = getActionCacheStatus(action.cacheHit, action.runner);

  return (
    <Space direction="vertical" size="middle" style={{ width: "100%" }}>
      <Descriptions
        bordered
        column={1}
        size="small"
        styles={{ label: { width: "20%" }, content: { width: "90%" } }}
      >
        {action.type && (
          <Descriptions.Item label="Action mnemonic">
            {action.type}
          </Descriptions.Item>
        )}
        <Descriptions.Item label="Execution">
          {getActionExecutionKind(action.runner)}
          {action.runner && ` (${action.runner})`}
        </Descriptions.Item>
        <Descriptions.Item label="Cache result">
          {cacheStatus.label}
        </Descriptions.Item>
        {actionHref && (
          <Descriptions.Item label="Action">
            <Typography.Link href={actionHref}>
              View action in bb-browser
            </Typography.Link>
          </Descriptions.Item>
        )}
        {action.success !== null && action.success !== undefined && (
          <Descriptions.Item label="Success">
            {action.success ? "Yes" : "No"}
          </Descriptions.Item>
        )}
        {action.exitCode !== null && action.exitCode !== undefined && (
          <Descriptions.Item label="Exit code">
            {action.exitCode}
          </Descriptions.Item>
        )}
        {action.failureCode && (
          <Descriptions.Item label="Failure code">
            {action.failureCode}
          </Descriptions.Item>
        )}
        {action.failureMessage && (
          <Descriptions.Item label="Failure message">
            {action.failureMessage}
          </Descriptions.Item>
        )}
        {action.primaryOutput && (
          <Descriptions.Item label="Primary output">
            <OutputLink file={action.primaryOutputFile}>
              {action.primaryOutput}
            </OutputLink>
          </Descriptions.Item>
        )}
        {action.stdout && (
          <Descriptions.Item label="Standard output">
            <OutputLink file={action.stdout}>
              Download standard output
            </OutputLink>
          </Descriptions.Item>
        )}
        {action.stderr && (
          <Descriptions.Item label="Standard error">
            <OutputLink file={action.stderr}>
              Download standard error
            </OutputLink>
          </Descriptions.Item>
        )}
        {action.commandLine && (
          <Descriptions.Item label="Command line">
            <Flex wrap>
              {action.commandLine.map((arg, index) => (
                <pre
                  // biome-ignore lint/suspicious/noArrayIndexKey: duplicate arguments require the index
                  key={`${arg}-${index}`}
                  style={{ textWrap: "wrap", paddingRight: "0.7em" }}
                >
                  {index === 0 ? <strong>{arg}</strong> : arg}
                </pre>
              ))}
            </Flex>
          </Descriptions.Item>
        )}
        {action.configuration?.cpu && (
          <Descriptions.Item label="Configuration CPU">
            {action.configuration.cpu}
          </Descriptions.Item>
        )}
        {action.configuration?.platformName && (
          <Descriptions.Item label="Configuration platform">
            {action.configuration.platformName}
          </Descriptions.Item>
        )}
        {action.configuration?.mnemonic && (
          <Descriptions.Item label="Configuration mnemonic">
            {action.configuration.mnemonic}
          </Descriptions.Item>
        )}
        {action.configuration?.makeVariables &&
          Object.keys(action.configuration.makeVariables).length > 0 && (
            <Descriptions.Item label="Configuration make variables">
              <Space direction="vertical" size="small">
                {Object.entries(action.configuration.makeVariables).map(
                  ([key, value]) => (
                    <span key={key}>
                      <strong>{key}=</strong>
                      {`${value}`}
                    </span>
                  ),
                )}
              </Space>
            </Descriptions.Item>
          )}
      </Descriptions>
    </Space>
  );
};
