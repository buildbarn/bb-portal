import { useGrpcClients } from "@/context/GrpcClientsContext";
import { FileSystemAccessProfileReference } from "@/lib/grpc-client/buildbarn/query/query";
import type { BrowserPageParams } from "@/types/BrowserPageType";
import { useQuery } from "@tanstack/react-query";
import { Spin, Typography } from "antd";
import { useSearchParams } from "next/navigation";
import BrowserDirectory from "../BrowserDirectory";
import PortalAlert from "../PortalAlert";

interface Params {
  browserPageParams: BrowserPageParams;
}

const BrowserDirectoryPage: React.FC<Params> = ({ browserPageParams }) => {
  const { fileSystemAccessCacheClient } = useGrpcClients();
  const searchParams = useSearchParams();
  const params = searchParams.get("fileSystemAccessProfile");
  let fileSystemAccessProfileReference:
    | FileSystemAccessProfileReference
    | undefined = undefined;

  if (params) {
    try {
      fileSystemAccessProfileReference =
        FileSystemAccessProfileReference.fromJSON(
          JSON.parse(decodeURIComponent(params)),
        );
    } catch (error) {
      console.error("Could not parse query parameters");
    }
  }

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
        className="error"
        message={
          <>
            <Typography.Text>
              There was a problem communicating with the backend server:
            </Typography.Text>
            <pre>{String(error)}</pre>
          </>
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
