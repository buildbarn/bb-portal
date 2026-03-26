import { Link } from "@tanstack/react-router";
import { Descriptions, Flex, Space, Typography } from "antd";
import type React from "react";
import type {
  Command,
  Digest,
} from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import {
  type BrowserPageParams,
  BrowserPageType,
} from "@/types/BrowserPageType";
import { generateBrowserSplat } from "@/utils/urlGenerator";
import CopyBbClientdCommandButton from "./CopyBbClientdCommandButton";
import DownloadAsShellScriptButton from "./DownloadAsShellScriptButton";

interface Params {
  browserPageParams: BrowserPageParams;
  command: Command;
  commandDigest: Digest | undefined;
  showTitle: boolean;
}

const BrowserCommandDescription: React.FC<Params> = ({
  browserPageParams,
  command,
  commandDigest,
  showTitle,
}) => {
  return (
    <Space direction="vertical" size="middle">
      {showTitle && (
        <Typography.Title level={2}>
          {commandDigest ? (
            <Link
              to="/browser/$"
              params={{
                _splat: generateBrowserSplat(
                  browserPageParams.instanceName,
                  browserPageParams.digestFunction,
                  commandDigest,
                  BrowserPageType.Command,
                ),
              }}
              style={{ textDecoration: "underline" }}
            >
              Command
            </Link>
          ) : (
            "Command"
          )}
        </Typography.Title>
      )}
      <Descriptions
        column={1}
        size="small"
        bordered
        styles={{ label: { width: "25%" }, content: { width: "75%" } }}
      >
        <Descriptions.Item label="Arguments">
          <Flex wrap>
            {command.arguments.map((arg, index) => (
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
        <Descriptions.Item label="Environment variables">
          {command.environmentVariables.map((env) => (
            <pre key={env.name} style={{ textWrap: "wrap" }}>
              <b>{env.name}</b>
              {`=${env.value}`}
            </pre>
          ))}
        </Descriptions.Item>
        {command.workingDirectory !== "" && (
          <Descriptions.Item label="Working directory">
            {command.workingDirectory}
          </Descriptions.Item>
        )}
      </Descriptions>
      {commandDigest && (
        <Space direction="horizontal">
          <CopyBbClientdCommandButton
            browserPageParams={browserPageParams}
            commandDigest={commandDigest}
          />
          <DownloadAsShellScriptButton
            browserPageParams={browserPageParams}
            commandDigest={commandDigest}
          />
        </Space>
      )}
    </Space>
  );
};

export default BrowserCommandDescription;
