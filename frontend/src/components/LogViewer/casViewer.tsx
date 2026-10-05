import { useQuery } from "@tanstack/react-query";
import { digestFunctionValueFromString } from "@/utils/digestFunctionUtils";
import { generateFileUrl } from "@/utils/urlGenerator";
import { LogViewerCard, SIZE_BYTE_LIMIT } from ".";

export const fetchLog = async (
  instanceName: string,
  digestFunction: string,
  hash: string,
  sizeBytes: number,
  fileName: string,
): Promise<string | undefined> => {
  if (sizeBytes > SIZE_BYTE_LIMIT) {
    return undefined;
  }

  const res = await fetch(
    generateFileUrl(
      instanceName,
      digestFunctionValueFromString(digestFunction),
      {
        hash,
        sizeBytes: sizeBytes.toString(),
      },
      fileName,
    ),
  );
  if (!res.ok) {
    throw new Error(`Failed to fetch log: ${res.statusText}`);
  }
  return await res.text();
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
    enabled: sizeBytes <= SIZE_BYTE_LIMIT,
    queryFn: async () => {
      return await fetchLog(
        instanceName,
        digestFunction,
        hash,
        sizeBytes,
        fileName,
      );
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
