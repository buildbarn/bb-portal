import { Descriptions, Flex, Space } from "antd";
import { getFragmentData } from "@/graphql/__generated__";
import type { BazelInvocationActionsFragment } from "@/graphql/__generated__/graphql";
import { FILE_DETAILS_FRAGMENT } from "@/types/GraphqlFileFragment";
import { CasGqlFileViewer } from "../LogViewer/casGqlFileViewer";

interface Props {
  action: BazelInvocationActionsFragment;
}

export const ActionDetails: React.FC<Props> = ({ action }) => {
  const stdoutFile = getFragmentData(FILE_DETAILS_FRAGMENT, action.stdout);
  const stderrFile = getFragmentData(FILE_DETAILS_FRAGMENT, action.stderr);

  return (
    <Space direction="vertical" size="middle" style={{ width: "100%" }}>
      <Descriptions
        bordered
        column={1}
        size="small"
        styles={{ label: { width: "20%" }, content: { width: "90%" } }}
      >
        {action.type && (
          <Descriptions.Item label="Type">{action.type}</Descriptions.Item>
        )}
        {action.success !== null && action.success !== undefined && (
          <Descriptions.Item label="Success">
            {action.success ? "Yes" : "No"}
          </Descriptions.Item>
        )}
        {action.exitCode !== null && action.exitCode !== undefined && (
          <Descriptions.Item label="Exit Code">
            {action.exitCode}
          </Descriptions.Item>
        )}
        {action.failureCode && (
          <Descriptions.Item label="Failure Code">
            {action.failureCode}
          </Descriptions.Item>
        )}
        {action.failureMessage && (
          <Descriptions.Item label="Failure Message">
            {action.failureMessage}
          </Descriptions.Item>
        )}
        {action.commandLine && (
          <Descriptions.Item label="Command Line">
            <Flex wrap>
              {action.commandLine.map((arg, index) => (
                <pre
                  // biome-ignore lint/suspicious/noArrayIndexKey: Since there are dupliate args, we need to use index
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
          <Descriptions.Item label="Configuration Platform Name">
            {action.configuration.platformName}
          </Descriptions.Item>
        )}
        {action.configuration?.mnemonic && (
          <Descriptions.Item label="Configuration Mnemonic">
            {action.configuration.mnemonic}
          </Descriptions.Item>
        )}
        {action.configuration?.makeVariables &&
          Object.keys(action.configuration.makeVariables).length > 0 && (
            <Descriptions.Item label="Configuration Make Variables">
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
      {stdoutFile && (
        <CasGqlFileViewer
          file={stdoutFile}
          title="Standard output"
          fileName="standard_output.txt"
        />
      )}
      {stderrFile && (
        <CasGqlFileViewer
          file={stderrFile}
          title="Standard error"
          fileName="standard_error.txt"
        />
      )}
    </Space>
  );
};
