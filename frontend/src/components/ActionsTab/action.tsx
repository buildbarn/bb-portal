import { WarningOutlined } from "@ant-design/icons";
import { useQuery } from "@tanstack/react-query";
import { Button, Descriptions, Flex, Space, Tooltip } from "antd";
import type { BazelInvocationActionsFragment } from "@/graphql/__generated__/graphql";
import { casByteStreamClient } from "@/grpc/casByteStreamClient";
import { digestFunction_ValueFromJSON } from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import type { ByteStreamClient } from "@/lib/grpc-client/google/bytestream/bytestream";
import { digestFunctionValueFromString } from "@/utils/digestFunctionUtils";
import { fetchCasObject } from "@/utils/fetchCasObject";
import { readableFileSize } from "@/utils/filesize";
import { generateFileUrl } from "@/utils/urlGenerator";
import { LogViewerCard } from "../LogViewer";

const SIZE_BYTE_LIMIT = 1000000; // 1MiB

const fetchLog = async (
  casByteStreamClient: ByteStreamClient,
  instanceName: string,
  digestFunction: string | undefined | null,
  digest: string | undefined | null,
  sizeBytes: number | undefined | null,
): Promise<string | undefined> => {
  if (!digest || !sizeBytes || !digestFunction || sizeBytes > SIZE_BYTE_LIMIT) {
    return undefined;
  }

  const data = await fetchCasObject(
    casByteStreamClient,
    instanceName,
    digestFunction_ValueFromJSON(digestFunction.toUpperCase()),
    {
      hash: digest,
      sizeBytes: sizeBytes.toString(),
    },
  );
  return new TextDecoder().decode(data);
};

interface Props {
  instanceName: string;
  action: BazelInvocationActionsFragment;
}

export const ActionDetails: React.FC<Props> = ({ instanceName, action }) => {
  const { data } = useQuery({
    queryKey: ["actionLogs", action.id],
    queryFn: async () => {
      const stdoutPromise = fetchLog(
        casByteStreamClient,
        instanceName,
        action.stdoutHashFunction,
        action.stdoutHash,
        action.stdoutSizeBytes,
      );
      const stderrPromise = fetchLog(
        casByteStreamClient,
        instanceName,
        action.stderrHashFunction,
        action.stderrHash,
        action.stderrSizeBytes,
      );
      const [stdout, stderr] = await Promise.all([
        stdoutPromise,
        stderrPromise,
      ]);

      // Regex to match historical_execute_response URLs
      const re =
        /https?:\/\/[-a-zA-Z0-9.]{1,256}(:[0-9]+)?[-a-zA-Z0-9()@:%_+.~#?&/=]*\/blobs\/[a-zA-Z0-9]{0,20}\/historical_execute_response\/[0-9a-f]{64}-[0-9]*\//;
      const historicalUrl = stdout?.match(re)?.[0] || stderr?.match(re)?.[0];

      return { stdout, stderr, historicalUrl };
    },
  });

  const validActionOutputLink =
    action.stdoutHash && action.stdoutHashFunction && action.stdoutSizeBytes;
  const validErrorOutputLink =
    action.stderrHash && action.stderrHashFunction && action.stderrSizeBytes;

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
      <LogViewerCard
        log={data?.stdout}
        title="Standard output"
        logDownloadUrl={
          action.stdoutHashFunction &&
          action.stdoutHash &&
          action.stdoutSizeBytes
            ? generateFileUrl(
                instanceName,
                digestFunctionValueFromString(action.stdoutHashFunction),
                {
                  hash: action.stdoutHash,
                  sizeBytes: action.stdoutSizeBytes.toString(),
                },
                "standard_output",
              )
            : undefined
        }
        fileName="standard_output.txt"
        error={
          !data?.stdout &&
          action.stdoutSizeBytes &&
          action.stdoutSizeBytes > SIZE_BYTE_LIMIT
            ? Error("Standard output is too large to display.", {
                cause: `The standard output is ${readableFileSize(
                  action.stdoutSizeBytes,
                )}. ${validActionOutputLink && "Please download the output to view it."}`,
              })
            : undefined
        }
      />
      <LogViewerCard
        log={data?.stderr}
        title="Standard error"
        logDownloadUrl={
          action.stderrHashFunction &&
          action.stderrHash &&
          action.stderrSizeBytes
            ? generateFileUrl(
                instanceName,
                digestFunctionValueFromString(action.stderrHashFunction),
                {
                  hash: action.stderrHash,
                  sizeBytes: action.stderrSizeBytes.toString(),
                },
                "error_output",
              )
            : undefined
        }
        fileName="error_output.txt"
        error={
          !data?.stderr &&
          action.stderrSizeBytes &&
          action.stderrSizeBytes > SIZE_BYTE_LIMIT
            ? Error("Standard error is too large to display.", {
                cause: `The standard error output is ${readableFileSize(
                  action.stderrSizeBytes,
                )}. ${validErrorOutputLink && "Please download the output to view it."}`,
              })
            : undefined
        }
      />
      {data?.historicalUrl && (
        <Space
          size="small"
          direction="vertical"
          style={{ width: "100%" }}
          align="end"
        >
          <Tooltip title="This URL was extracted from the action's stdout or stderr, so there are no guarantees that it is correct. It points to a historical execute response stored in the CAS.">
            <Button
              type="primary"
              href={data.historicalUrl}
              target="_blank"
              rel="noopener noreferrer"
            >
              View Historical Execute Response
              <WarningOutlined />
            </Button>
          </Tooltip>
        </Space>
      )}
    </Space>
  );
};
