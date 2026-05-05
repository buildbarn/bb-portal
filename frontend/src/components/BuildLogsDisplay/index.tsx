import { useQuery } from "@tanstack/react-query";
import { LogViewerCard } from "../LogViewer";

interface Props {
  invocationId: string;
  rawCommand: string | null;
}

const fetchLog = async (id: string, start = 0, end = -1): Promise<string> => {
  const params = new URLSearchParams({
    start_line: start.toString(),
    end_line: end.toString(),
  });
  const uri = `/api/v1/invocations/${id}/log?${params}`;
  const response = await fetch(uri);
  if (!response.ok) throw new Error("Failed to fetch logs");
  return response.text();
};

const BuildLogsDisplay: React.FC<Props> = ({ invocationId, rawCommand }) => {
  // TODO: Only fetch the currently viewed parts of the log.
  const { data, error, isLoading } = useQuery({
    queryKey: ["getLogs", invocationId],
    queryFn: () => fetchLog(invocationId),
  });

  const logDownloadUrl = `/api/v1/invocations/${invocationId}/log`;

  return (
    <LogViewerCard
      loading={isLoading}
      error={error}
      log={data}
      logDownloadUrl={logDownloadUrl}
      title={`Raw Build Logs for ${rawCommand}`}
      fileName="log.txt"
    />
  );
};

export default BuildLogsDisplay;
