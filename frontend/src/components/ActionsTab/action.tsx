import { WarningOutlined } from "@ant-design/icons";
import { useQuery } from "@tanstack/react-query";
import { Button, Descriptions, Flex, Space, Tooltip, Typography } from "antd";
import { useGrpcClients } from "@/context/GrpcClientsContext";
import type { BazelInvocationInfoFragment } from "@/graphql/__generated__/graphql";
import { digestFunction_ValueFromJSON } from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import type { ByteStreamClient } from "@/lib/grpc-client/google/bytestream/bytestream";
import { fetchCasObject } from "@/utils/fetchCasObject";
import LogViewer from "../LogViewer";

export type ActionDetailsData = NonNullable<
  BazelInvocationInfoFragment["actions"]
>[number];

const fetchLog = async (
  casByteStreamClient: ByteStreamClient,
  instanceName: string,
  digestFunction: string | undefined | null,
  digest: string | undefined | null,
  sizeBytes: number | undefined | null,
): Promise<string | undefined> => {
  if (!digest || !sizeBytes || !digestFunction) {
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
  action: ActionDetailsData;
}

export const ActionDetails: React.FC<Props> = ({ instanceName, action }) => {
  const { casByteStreamClient } = useGrpcClients();

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
      {data?.stdout && (
        <Space size="small" direction="vertical" style={{ width: "100%" }}>
          <Typography.Title level={4}>Standard output:</Typography.Title>
          <LogViewer log={data.stdout} />
        </Space>
      )}
      {data?.stderr && (
        <Space size="small" direction="vertical" style={{ width: "100%" }}>
          <Typography.Title level={4}>Standard error:</Typography.Title>
          <LogViewer log={data.stderr} />
        </Space>
      )}
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
