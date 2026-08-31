import { Descriptions, Space, Tag } from "antd";
import type React from "react";
import type {
  Digest,
  ExecuteResponse,
} from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import type { POSIXResourceUsage } from "@/lib/grpc-client/buildbarn/resourceusage/resourceusage";
import type { BrowserPageParams } from "@/types/BrowserPageType";
import { digestFunctionValueToString } from "@/utils/digestFunctionUtils";
import { LogViewerCard } from "../LogViewer";
import { CasViewer } from "../LogViewer/casViewer";

interface Params {
  browserPageParams: BrowserPageParams;
  executeResponse: ExecuteResponse;
  posixResourceUsage: POSIXResourceUsage | undefined;
}

interface ConsoleOutputParams {
  browserPageParams: BrowserPageParams;
  digest: Digest | undefined;
  rawOutput: Uint8Array;
  title: string;
  streamName: string;
  fileName: string;
}

const ConsoleOutput: React.FC<ConsoleOutputParams> = ({
  browserPageParams,
  digest,
  rawOutput,
  title,
  streamName,
  fileName,
}) => {
  const sizeBytes = Number.parseInt(digest?.sizeBytes ?? "0", 10);

  if (digest?.hash && Number.isFinite(sizeBytes) && sizeBytes > 0) {
    return (
      <CasViewer
        instanceName={browserPageParams.instanceName}
        digestFunction={digestFunctionValueToString(
          browserPageParams.digestFunction,
        )}
        hash={digest.hash}
        sizeBytes={sizeBytes}
        title={title}
        fileName={fileName}
      />
    );
  }

  if (rawOutput.length > 0) {
    return (
      <LogViewerCard
        log={new TextDecoder().decode(rawOutput)}
        logSizeBytes={rawOutput.length}
        title={title}
        fileName={fileName}
      />
    );
  }

  return (
    <LogViewerCard
      log={undefined}
      title={title}
      fileName={fileName}
      emptyMessage={
        digest?.hash
          ? `The action produced no ${streamName}. No log file is available for download.`
          : `No ${streamName} log was uploaded.`
      }
    />
  );
};

const BrowserResultDescription: React.FC<Params> = ({
  browserPageParams,
  executeResponse,
  posixResourceUsage,
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

  return (
    <Space direction="vertical" size="small" style={{ width: "100%" }}>
      <Descriptions
        column={1}
        size="small"
        bordered
        styles={{ label: { width: "25%" }, content: { width: "75%" } }}
      >
        {renderResult()}
      </Descriptions>
      {executeResponse.result && (
        <ConsoleOutput
          browserPageParams={browserPageParams}
          digest={executeResponse.result.stdoutDigest}
          rawOutput={executeResponse.result.stdoutRaw}
          title="Standard Output"
          streamName="standard output"
          fileName="standard_output.txt"
        />
      )}
      {executeResponse.result && (
        <ConsoleOutput
          browserPageParams={browserPageParams}
          digest={executeResponse.result.stderrDigest}
          rawOutput={executeResponse.result.stderrRaw}
          title="Standard Error"
          streamName="standard error"
          fileName="standard_error.txt"
        />
      )}
    </Space>
  );
};

export default BrowserResultDescription;
