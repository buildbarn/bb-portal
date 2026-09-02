import { useQuery } from "@tanstack/react-query";
import { Space, Spin } from "antd";
import { fileSystemAccessCacheClient } from "@/grpc/fileSystemAccessCacheClient";
import type { BrowserPageParams } from "@/types/BrowserPageParams";
import type { FileSystemAccessProfileReference } from "@/types/FileSystemAccessProfileReference";
import { BrowserDirectory } from "../BrowserDirectory";
import CopyBbClientdDirectoryButton from "../BrowserDirectory/CopyBbClientdDirectoryButton";
import DownloadAsTarballButton from "../BrowserDirectory/DownloadAsTarballButton";
import { DirectoryPrefetchDescription } from "../BrowserDirectory/directoryPrefetchDescription";
import PortalAlert from "../PortalAlert";

interface Params {
  browserPageParams: BrowserPageParams;
  fileSystemAccessProfileReference:
    | FileSystemAccessProfileReference
    | undefined;
}

const BrowserDirectoryPage: React.FC<Params> = ({
  browserPageParams,
  fileSystemAccessProfileReference,
}) => {
  const { data, isError, error, isLoading } = useQuery({
    queryKey: [
      "fileSystemAccessProfile",
      browserPageParams,
      fileSystemAccessProfileReference,
    ],
    queryFn: fileSystemAccessCacheClient.getFileSystemAccessProfile.bind(
      {},
      {
        instanceName: browserPageParams.instanceName,
        digestFunction: browserPageParams.digestFunction,
        reducedActionDigest: fileSystemAccessProfileReference?.digest,
      },
    ),
    enabled: fileSystemAccessProfileReference !== undefined,
  });

  if (isLoading) {
    return <Spin />;
  }

  if (isError) {
    return (
      <PortalAlert
        showIcon
        type="error"
        message="Error fetching directory"
        description={
          error.message ||
          "Unknown error occurred while fetching data from the server."
        }
      />
    );
  }

  return (
    <Space direction="vertical">
      <BrowserDirectory
        baseData={{
          instanceName: browserPageParams.instanceName,
          digestFunction: browserPageParams.digestFunction,
          digest: browserPageParams.digest,
          fileSystemAccessProfile: data,
          fileSystemAccessProfileReference: fileSystemAccessProfileReference,
        }}
        openDirsString=""
        useBloomFilter={true}
      />
      <Space direction="vertical" size="small">
        <DirectoryPrefetchDescription prefetchDataExists={!!data} />
        <Space direction="horizontal">
          <CopyBbClientdDirectoryButton
            instanceName={browserPageParams.instanceName}
            digestFunction={browserPageParams.digestFunction}
            inputRootDigest={browserPageParams.digest}
          />
          <DownloadAsTarballButton
            instanceName={browserPageParams.instanceName}
            digestFunction={browserPageParams.digestFunction}
            directoryDigest={browserPageParams.digest}
          />
        </Space>
      </Space>
    </Space>
  );
};

export default BrowserDirectoryPage;
