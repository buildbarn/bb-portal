import { Digest } from '@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution';
import { BrowserPageParams } from '@/types/BrowserPageType';
import { generateCommandShellScriptUrl } from '@/utils/urlGenerator';
import { Button } from 'antd';
import React from 'react';

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
