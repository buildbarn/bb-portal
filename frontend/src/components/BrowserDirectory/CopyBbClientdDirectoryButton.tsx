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
  inputRootDigest: Digest;
}

const CopyBbClientdDirectoryButton: React.FC<Params> = ({
  instanceName,
  digestFunction,
  inputRootDigest,
}) => {
  const { copyToClipboard } = useBbPortalMessage();

  const inputRootBbClientdPath = getBBClientdPath(
    instanceName,
    digestFunction,
    inputRootDigest,
    "directory",
  );

  return (
    <Button
      type="primary"
      onClick={() => copyToClipboard(inputRootBbClientdPath)}
    >
      Copy bb_clientd path to clipboard
    </Button>
  );
};

export default CopyBbClientdDirectoryButton;
