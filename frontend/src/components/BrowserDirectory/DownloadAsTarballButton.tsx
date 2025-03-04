import {
  Digest,
  DigestFunction_Value,
} from '@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution';
import { generateDirectoryTarballUrl } from '@/utils/urlGenerator';
import { Button } from 'antd';
import React from 'react';

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
