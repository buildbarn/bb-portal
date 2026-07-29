import { Button } from "antd";
import type React from "react";
import { useBbPortalMessage } from "@/context/MessageContext";
import type { Digest } from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import type { BrowserPageParams } from "@/types/BrowserPageType";
import { getBBClientdPath } from "@/utils/getBbClientdPath";

interface Params {
  browserPageParams: BrowserPageParams;
  commandDigest: Digest;
}

const CopyBbClientdCommandButton: React.FC<Params> = ({
  browserPageParams,
  commandDigest,
}) => {
  const { copyToClipboard } = useBbPortalMessage();

  const commandBbClientdPath = getBBClientdPath(
    browserPageParams.instanceName,
    browserPageParams.digestFunction,
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
