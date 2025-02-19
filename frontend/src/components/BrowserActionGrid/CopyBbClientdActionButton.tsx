"use client";

import type { Digest } from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import type { BrowserPageParams } from "@/types/BrowserPageType";
import { getBBClientdPath } from "@/utils/getBbClientdPath";
import { Button, message } from "antd";
import type React from "react";

interface Params {
  browserPageParams: BrowserPageParams;
  actionDigest: Digest;
  commandDigest: Digest;
  inputRootDigest: Digest;
}

const CopyBbClientdActionButton: React.FC<Params> = ({
  browserPageParams,
  actionDigest,
  commandDigest,
  inputRootDigest,
}) => {
  const [messageApi, contextHolder] = message.useMessage();

  const commandBbClientdPath = getBBClientdPath(
    browserPageParams.instanceName,
    browserPageParams.digestFunction,
    commandDigest,
    "command",
  );

  const inputRootBbClientdPath = getBBClientdPath(
    browserPageParams.instanceName,
    browserPageParams.digestFunction,
    inputRootDigest,
    "directory",
  );

  const script = `rsync \\
    --delete \\
    --link-dest ${inputRootBbClientdPath}/ \\
    --progress \\
    --recursive \\
    ${inputRootBbClientdPath}/ \\
    ~/bb_clientd/scratch/${actionDigest.hash}-${actionDigest.sizeBytes} &&
cd ~/bb_clientd/scratch/${actionDigest.hash}-${actionDigest.sizeBytes} &&
${commandBbClientdPath}`;

  return (
    <>
      {contextHolder}
      <Button
        type="primary"
        onClick={() => {
          navigator.clipboard.writeText(script);
          messageApi.success("Copied command to clipboard", 1.5);
        }}
      >
        Copy bb_clientd command for running action locally to clipboard
      </Button>
    </>
  );
};

export default CopyBbClientdActionButton;
