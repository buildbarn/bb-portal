import type { Digest } from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import type { BrowserPageParams } from "@/types/BrowserPageType";
import { getBBClientdPath } from "@/utils/getBbClientdPath";
import { Button, message } from "antd";
import type React from "react";

interface Params {
  browserPageParams: BrowserPageParams;
  commandDigest: Digest;
}

const CopyBbClientdCommandButton: React.FC<Params> = ({
  browserPageParams,
  commandDigest,
}) => {
  const [messageApi, contextHolder] = message.useMessage();

  const commandBbClientdPath = getBBClientdPath(
    browserPageParams.instanceName,
    browserPageParams.digestFunction,
    commandDigest,
    "command",
  );

  return (
    <>
      {contextHolder}
      <Button
        type="primary"
        onClick={() => {
          navigator.clipboard.writeText(commandBbClientdPath);
          messageApi.success("Copied path to clipboard", 1.5);
        }}
      >
        Copy bb_clientd path of shell script to clipboard
      </Button>
    </>
  );
};

export default CopyBbClientdCommandButton;
