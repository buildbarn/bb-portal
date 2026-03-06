import { useQuery } from "@tanstack/react-query";
import { Spin } from "antd";
import { useGrpcClients } from "@/context/GrpcClientsContext";
import { FileSystemAccessProfileReference } from "@/lib/grpc-client/buildbarn/query/query";
import type { BrowserPageParams } from "@/types/BrowserPageType";
import BrowserDirectory from "../BrowserDirectory";
import PortalAlert from "../PortalAlert";

interface Params {
  browserPageParams: BrowserPageParams;
  fileSystemAccessProfileReference: FileSystemAccessProfileReference | undefined
}

const BrowserDirectoryPage: React.FC<Params> = ({ browserPageParams, fileSystemAccessProfileReference }) => {
  const { fileSystemAccessCacheClient } = useGrpcClients();

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
    <BrowserDirectory
      instanceName={browserPageParams.instanceName}
      digestFunction={browserPageParams.digestFunction}
      inputRootDigest={browserPageParams.digest}
      fileSystemAccessProfile={data}
      fileSystemAccessProfileReference={fileSystemAccessProfileReference}
    />
  );
};

export default BrowserDirectoryPage;
