import { Descriptions, Space, Tag } from "antd";
import type React from "react";
import type { ExecuteResponse } from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import type { POSIXResourceUsage } from "@/lib/grpc-client/buildbarn/resourceusage/resourceusage";
import type { BrowserPageParams } from "@/types/BrowserPageType";
import { digestFunctionValueToString } from "@/utils/digestFunctionUtils";
import { CasViewer } from "../LogViewer/casViewer";

interface Params {
  browserPageParams: BrowserPageParams;
  executeResponse: ExecuteResponse;
  posixResourceUsage: POSIXResourceUsage | undefined;
}

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
      {executeResponse.result?.stdoutDigest?.hash &&
        executeResponse.result?.stdoutDigest?.sizeBytes && (
          <CasViewer
            instanceName={browserPageParams.instanceName}
            digestFunction={digestFunctionValueToString(
              browserPageParams.digestFunction,
            )}
            digest={executeResponse.result.stdoutDigest.hash}
            sizeBytes={Number.parseInt(
              executeResponse.result.stdoutDigest.sizeBytes,
              10,
            )}
            title="Standard Output"
            fileName="standard_output.txt"
          />
        )}
      {executeResponse.result?.stderrDigest?.hash &&
        executeResponse.result?.stderrDigest?.sizeBytes && (
          <CasViewer
            instanceName={browserPageParams.instanceName}
            digestFunction={digestFunctionValueToString(
              browserPageParams.digestFunction,
            )}
            digest={executeResponse.result?.stderrDigest?.hash}
            sizeBytes={Number.parseInt(
              executeResponse.result?.stderrDigest?.sizeBytes,
              10,
            )}
            title="Standard Error"
            fileName="standard_error.txt"
          />
        )}
    </Space>
  );
};

export default BrowserResultDescription;
