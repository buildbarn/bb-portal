import { Button } from "antd";
import type React from "react";
import type {
  Digest,
  DigestFunction_Value,
} from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import { generateCommandShellScriptUrl } from "@/utils/urlGenerator";

interface Params {
  instanceName: string;
  digestFunction: DigestFunction_Value;
  commandDigest: Digest;
}

const DownloadAsShellScriptButton: React.FC<Params> = ({
  instanceName,
  digestFunction,
  commandDigest,
}) => {
  return (
    <Button
      type="primary"
      href={generateCommandShellScriptUrl(
        instanceName,
        digestFunction,
        commandDigest,
      )}
    >
      Download as shell script
    </Button>
  );
};

export default DownloadAsShellScriptButton;
