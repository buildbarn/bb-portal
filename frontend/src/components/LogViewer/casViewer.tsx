import { useQuery } from "@tanstack/react-query";
import { casByteStreamClient } from "@/grpc/casByteStreamClient";
import { digestFunction_ValueFromJSON } from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import { digestFunctionValueFromString } from "@/utils/digestFunctionUtils";
import { fetchCasObject } from "@/utils/fetchCasObject";
import { readableFileSize } from "@/utils/filesize";
import { generateFileUrl } from "@/utils/urlGenerator";
import { LogViewerCard } from ".";

const SIZE_BYTE_LIMIT = 1000000; // 1MiB

const fetchLog = async (
  instanceName: string,
  digestFunction: string,
  digest: string,
  sizeBytes: number,
): Promise<string | undefined> => {
  if (sizeBytes > SIZE_BYTE_LIMIT) {
    return undefined;
  }

  const data = await fetchCasObject(
    casByteStreamClient,
    instanceName,
    digestFunction_ValueFromJSON(digestFunction.toUpperCase()),
    {
      hash: digest,
      sizeBytes: sizeBytes.toString(),
    },
  );
  return new TextDecoder().decode(data);
};

interface Props {
  instanceName: string;
  digestFunction: string;
  digest: string;
  sizeBytes: number;
  title: string;
  fileName: string;
}

// Takes a digest and displays the content with useful buttons
const CasViewer: React.FC<Props> = ({
  instanceName,
  digestFunction,
  digest,
  sizeBytes,
  title,
  fileName,
}) => {
  const { data, isLoading, error } = useQuery({
    queryKey: ["casLog", digest],
    queryFn: async () => {
      return await fetchLog(instanceName, digestFunction, digest, sizeBytes);
    },
  });

  const validActionOutputLink = digestFunction && digest && sizeBytes;

  return (
    <LogViewerCard
      log={data}
      title={title}
      logDownloadUrl={
        validActionOutputLink
          ? generateFileUrl(
              instanceName,
              digestFunctionValueFromString(digestFunction),
              {
                hash: digest,
                sizeBytes: sizeBytes.toString(),
              },
              fileName,
            )
          : undefined
      }
      fileName={fileName}
      error={
        !data && sizeBytes > SIZE_BYTE_LIMIT
          ? Error("Output is too large to display.", {
              cause: `The size of the output is ${readableFileSize(
                sizeBytes,
              )}. ${validActionOutputLink && "Download the output to view it."}`,
            })
          : error
      }
      loading={isLoading}
    />
  );
};

export { CasViewer };
