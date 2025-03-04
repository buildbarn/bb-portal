import type { ExecuteResponse } from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import type { POSIXResourceUsage } from "@/lib/grpc-client/buildbarn/resourceusage/resourceusage";
import type { BrowserPageParams } from "@/types/BrowserPageType";
import { formatFileSizeFromString } from "@/utils/formatValues";
import { generateFileUrl } from "@/utils/urlGenerator";
import { Descriptions, Space, Tag, Typography } from "antd";
import Paragraph from "antd/es/typography/Paragraph";
import Link from "next/link";
import type React from "react";
import type { ActionConsoleOutput } from "../BrowserActionGrid/types";

interface Params {
  browserPageParams: BrowserPageParams;
  executeResponse: ExecuteResponse;
  posixResourceUsage: POSIXResourceUsage | undefined;
  consoleOutputs: ActionConsoleOutput[];
}

const BrowserResultDescription: React.FC<Params> = ({
  browserPageParams,
  executeResponse,
  posixResourceUsage,
  consoleOutputs,
}) => {
  const renderResult = () => {
    if (executeResponse.status !== undefined) {
      return (
        <Descriptions.Item label="Status Code">
          {`Code ${executeResponse.status.code}: ${executeResponse.status.message}`}
        </Descriptions.Item>
      );
    }

    if (
      posixResourceUsage?.terminationSignal !== undefined &&
      posixResourceUsage?.terminationSignal !== ""
    ) {
      return (
        <Descriptions.Item label="Termination signal">
          <Tag color="red">{`SIG${posixResourceUsage.terminationSignal}`}</Tag>
        </Descriptions.Item>
      );
    }

    return (
      <Descriptions.Item label="Exit code">
        <Space>
          {executeResponse.result?.exitCode}
          <Tag color={executeResponse.result?.exitCode === 0 ? "green" : "red"}>
            {executeResponse.result?.exitCode === 0 ? "Success" : "Failure"}
          </Tag>
        </Space>
      </Descriptions.Item>
    );
  };

  const renderConsoleOutput = (consoleOutput: ActionConsoleOutput) => {
    const logLinkHref = consoleOutput.digest
      ? generateFileUrl(
          browserPageParams.instanceName,
          browserPageParams.digestFunction,
          consoleOutput.digest,
          "log.txt",
        )
      : undefined;

    const label = () => {
      if (logLinkHref) {
        return <Link href={logLinkHref}>{consoleOutput.name}</Link>;
      }
      return consoleOutput.name;
    };

    const content = () => {
      if (consoleOutput.notFound) {
        return "The log file for this action could not be found.";
      }
      if (consoleOutput.tooLarge) {
        if (consoleOutput.digest && logLinkHref) {
          return (
            <Typography.Text>
              The <Link href={logLinkHref}>log file</Link> for this action is to
              large to display (
              {formatFileSizeFromString(consoleOutput.digest.sizeBytes)}).
            </Typography.Text>
          );
        }
        return "The log file for this action is to large to display.";
      }
      return (
        <Paragraph style={{ margin: 0 }}>
          <pre style={{ margin: 0 }}>{consoleOutput.content}</pre>
        </Paragraph>
      );
    };

    return (
      <Descriptions.Item label={label()} key={consoleOutput.name}>
        {content()}
      </Descriptions.Item>
    );
  };

  return (
    <Descriptions
      column={1}
      size="small"
      bordered
      styles={{ label: { width: "25%" }, content: { width: "75%" } }}
    >
      {renderResult()}
      {consoleOutputs.map(renderConsoleOutput)}
    </Descriptions>
  );
};

export default BrowserResultDescription;
