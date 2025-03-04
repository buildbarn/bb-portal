import type {
  Digest,
  DigestFunction_Value,
} from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import { getBBClientdPath } from "@/utils/getBbClientdPath";
import { Button, message } from "antd";
import type React from "react";

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
  const [messageApi, contextHolder] = message.useMessage();

  const inputRootBbClientdPath = getBBClientdPath(
    instanceName,
    digestFunction,
    inputRootDigest,
    "directory",
  );

  return (
    <>
      {contextHolder}
      <Button
        type="primary"
        onClick={() => {
          navigator.clipboard.writeText(inputRootBbClientdPath);
          messageApi.success("Copied command to clipboard", 1.5);
        }}
      >
        Copy bb_clientd path to clipboard
      </Button>
    </>
  );
};

export default CopyBbClientdDirectoryButton;
