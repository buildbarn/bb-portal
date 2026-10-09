import { useQuery, useQueryClient } from "@tanstack/react-query";
import { LogViewerCard } from "../LogViewer";
import { fetchTailedLog } from "./tailLog";

const LIVE_REFETCH_INTERVAL_MS = 5000;

interface Props {
  invocationId: string;
  rawCommand: string | null;
  isLive: boolean;
}

const fetchLog = async (
  id: string,
  start: number,
  end = -1,
): Promise<string> => {
  const params = new URLSearchParams({
    start_line: start.toString(),
    end_line: end.toString(),
  });
  const uri = `/api/v1/invocations/${id}/log?${params}`;
  const response = await fetch(uri);
  if (!response.ok) throw new Error("Failed to fetch logs");
  return response.text();
};

const BuildLogsDisplay: React.FC<Props> = ({
  invocationId,
  rawCommand,
  isLive,
}) => {
  const queryClient = useQueryClient();
  const queryKey = ["getLogs", invocationId];

  // TODO: Only fetch the currently viewed parts of the log.
  const { data, error, isLoading } = useQuery({
    queryKey,
    queryFn: () =>
      fetchTailedLog(queryClient.getQueryData<string>(queryKey), (startLine) =>
        fetchLog(invocationId, startLine),
      ),
    refetchInterval: isLive ? LIVE_REFETCH_INTERVAL_MS : false,
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
      isLive={isLive}
    />
  );
};

export default BuildLogsDisplay;
