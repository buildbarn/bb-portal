import { Button } from "antd";
import type React from "react";
import type {
  Digest,
  DigestFunction_Value,
} from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import { generateDirectoryTarballUrl } from "@/utils/urlGenerator";

interface Params {
  instanceName: string;
  digestFunction: DigestFunction_Value;
  directoryDigest: Digest;
}

const DownloadAsTarballButton: React.FC<Params> = ({
  instanceName,
  digestFunction,
  directoryDigest,
}) => {
  return (
    <Button
      type="primary"
      href={generateDirectoryTarballUrl(
        instanceName,
        digestFunction,
        directoryDigest,
      )}
    >
      Download as tarball
    </Button>
  );
};

export default DownloadAsTarballButton;
