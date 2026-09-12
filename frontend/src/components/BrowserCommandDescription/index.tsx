import { Link } from "@tanstack/react-router";
import { Descriptions, Grid, Space, Typography } from "antd";
import type React from "react";
import type { ReactNode } from "react";
import type {
  Command,
  Command_EnvironmentVariable,
  Digest,
} from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import type { BrowserPageParams } from "@/types/BrowserPageParams";
import { BrowserPageType } from "@/types/BrowserPageType";
import { generateBrowserSplat } from "@/utils/urlGenerator";
import CodeText from "../CodeText";
import CopyableIcon from "../CopyableIcon";

interface Params {
  browserPageParams: BrowserPageParams;
  command: Omit<Command, "environmentVariables" | "workingDirectory"> & {
    environmentVariables: (Command_EnvironmentVariable & {
      style?: React.CSSProperties;
    })[];
    argumentStyles?: React.CSSProperties[];
    workingDirectory: string | ReactNode;
  };
  commandDigest: Digest;
  showTitle: boolean;
}

const InnerBrowserCommandDescription: React.FC<Params> = ({
  browserPageParams,
  command,
  commandDigest,
  showTitle,
}) => {
  const screens = Grid.useBreakpoint();

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
        layout={screens.md ? "horizontal" : "vertical"}
        column={1}
        size="small"
        bordered
        styles={{ label: { width: "25%" }, content: { width: "75%" } }}
      >
        <Descriptions.Item label="Arguments">
          {command.arguments.map((arg, index) => (
            <CodeText
              // biome-ignore lint/suspicious/noArrayIndexKey: Since there are dupliate args, we need to use index
              key={index}
              style={{
                whiteSpace: "pre-wrap",
                overflowWrap: "anywhere",
                ...(command.argumentStyles?.at(index) || {}),
              }}
            >
              {index === 0 ? <strong>{arg}</strong> : `${arg}`}
            </CodeText>
          ))}{" "}
          <CopyableIcon text={command.arguments.join("")} />
        </Descriptions.Item>
        <Descriptions.Item label="Environment variables">
          {command.environmentVariables.map((env, index) => (
            <>
              <CodeText
                // biome-ignore lint/suspicious/noArrayIndexKey: Since there are dupliate args, we need to use index
                key={`${env.name}-${index}`}
                style={{
                  textWrap: "wrap",
                  ...env.style,
                }}
              >
                <b>{env.name}</b>
                {`=${env.value} `}
              </CodeText>
              <br />
            </>
          ))}
        </Descriptions.Item>
        {command.workingDirectory && (
          <Descriptions.Item label="Working directory">
            {typeof command.workingDirectory === "string" && (
              <p>{command.workingDirectory}</p>
            )}
            {typeof command.workingDirectory !== "string" &&
              command.workingDirectory}
          </Descriptions.Item>
        )}
      </Descriptions>
    </Space>
  );
};

interface BrowserCommandDescriptionParams {
  browserPageParams: BrowserPageParams;
  command: Command;
  commandDigest: Digest;
  showTitle: boolean;
}
const BrowserCommandDescription: React.FC<BrowserCommandDescriptionParams> = ({
  browserPageParams,
  command,
  commandDigest,
  showTitle,
}: BrowserCommandDescriptionParams) => {
  const argumentsWithWhitespace = command.arguments.map((arg) =>
    arg.length > 0 && arg !== " " ? `${arg}\xa0` : arg,
  );
  return (
    <InnerBrowserCommandDescription
      browserPageParams={browserPageParams}
      command={{ ...command, arguments: argumentsWithWhitespace }}
      commandDigest={commandDigest}
      showTitle={showTitle}
    />
  );
};

export { BrowserCommandDescription, InnerBrowserCommandDescription };
