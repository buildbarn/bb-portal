import { useQuery } from "@tanstack/react-query";
import { casByteStreamClient } from "@/grpc/casByteStreamClient";
import { digestFunction_ValueFromJSON } from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import { digestFunctionValueFromString } from "@/utils/digestFunctionUtils";
import { fetchCasObject } from "@/utils/fetchCasObject";
import { generateFileUrl } from "@/utils/urlGenerator";
import { LogViewerCard, SIZE_BYTE_LIMIT } from ".";

export const fetchLog = async (
  instanceName: string,
  digestFunction: string,
  hash: string,
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
      hash,
      sizeBytes: sizeBytes.toString(),
    },
  );
  return new TextDecoder().decode(data);
};

interface Props {
  instanceName: string;
  digestFunction: string;
  hash: string;
  sizeBytes: number;
  title: string;
  fileName: string;
}

// Takes a digest and displays the content with useful buttons
const CasViewer: React.FC<Props> = ({
  instanceName,
  digestFunction,
  hash,
  sizeBytes,
  title,
  fileName,
}) => {
  const { data, isLoading, error } = useQuery({
    queryKey: ["casLog", hash],
    queryFn: async () => {
      return await fetchLog(instanceName, digestFunction, hash, sizeBytes);
    },
  });

  return (
    <LogViewerCard
      log={data}
      logSizeBytes={sizeBytes}
      title={title}
      logDownloadUrl={generateFileUrl(
        instanceName,
        digestFunctionValueFromString(digestFunction),
        {
          hash,
          sizeBytes: sizeBytes.toString(),
        },
        fileName,
      )}
      fileName={fileName}
      error={error}
      loading={isLoading}
    />
  );
};

export { CasViewer };
