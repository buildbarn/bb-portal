import { Button } from "antd";
import type React from "react";
import type { Digest } from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import type { BrowserPageParams } from "@/types/BrowserPageType";
import { generateCommandShellScriptUrl } from "@/utils/urlGenerator";

interface Params {
  browserPageParams: BrowserPageParams;
  commandDigest: Digest;
}

const DownloadAsShellScriptButton: React.FC<Params> = ({
  browserPageParams,
  commandDigest,
}) => {
  return (
    <Button
      type="primary"
      href={generateCommandShellScriptUrl(
        browserPageParams.instanceName,
        browserPageParams.digestFunction,
        commandDigest,
      )}
    >
      Download as shell script
    </Button>
  );
};

export default DownloadAsShellScriptButton;
