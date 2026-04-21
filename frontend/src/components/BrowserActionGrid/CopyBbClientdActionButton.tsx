import { Button } from "antd";
import type React from "react";
import { useBbPortalMessage } from "@/context/MessageContext";
import type {
  Digest,
  DigestFunction_Value,
} from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import { getBBClientdPath } from "@/utils/getBbClientdPath";

interface Params {
  instanceName: string;
  digestFunction: DigestFunction_Value;
  actionDigest: Digest;
  commandDigest: Digest;
  inputRootDigest: Digest;
}

const CopyBbClientdActionButton: React.FC<Params> = ({
  instanceName,
  digestFunction,
  actionDigest,
  commandDigest,
  inputRootDigest,
}) => {
  const { copyToClipboard } = useBbPortalMessage();

  const commandBbClientdPath = getBBClientdPath(
    instanceName,
    digestFunction,
    commandDigest,
    "command",
  );

  const inputRootBbClientdPath = getBBClientdPath(
    instanceName,
    digestFunction,
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
    <Button type="primary" onClick={() => copyToClipboard(script)}>
      Copy bb_clientd command for running action locally to clipboard
    </Button>
  );
};

export default CopyBbClientdActionButton;
