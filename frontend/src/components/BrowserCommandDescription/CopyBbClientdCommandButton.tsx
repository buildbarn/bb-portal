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
  commandDigest: Digest;
}

const CopyBbClientdCommandButton: React.FC<Params> = ({
  instanceName,
  digestFunction,
  commandDigest,
}) => {
  const { copyToClipboard } = useBbPortalMessage();

  const commandBbClientdPath = getBBClientdPath(
    instanceName,
    digestFunction,
    commandDigest,
    "command",
  );

  return (
    <Button
      type="primary"
      onClick={() => copyToClipboard(commandBbClientdPath)}
    >
      Copy bb_clientd path of shell script to clipboard
    </Button>
  );
};

export default CopyBbClientdCommandButton;
