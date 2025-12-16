import {
  ExclamationCircleOutlined,
  FileSearchOutlined,
} from "@ant-design/icons";
import { useQuery } from "@tanstack/react-query";
import { Tooltip } from "antd";
import { env } from "next-runtime-env";
import DownloadButton from "@/components/DownloadButton";
import PortalCard from "@/components/PortalCard";
import LogViewer from "../LogViewer";

interface Props {
  invocationId: string;
}

const fetchLog = async (id: string, start = 0, end = -1): Promise<string> => {
  const params = new URLSearchParams({
    start_line: start.toString(),
    end_line: end.toString(),
  });
  const uri = `${env("NEXT_PUBLIC_BES_BACKEND_URL")}/api/v1/invocations/${id}/log?${params}`;
  const response = await fetch(uri);
  if (!response.ok) throw new Error("Failed to fetch logs");
  return response.text();
};

const BuildLogsDisplay: React.FC<Props> = ({ invocationId }) => {
  // TODO: Only fetch the currently viewed parts of the log.
  const { data, error, isLoading } = useQuery({
    queryKey: ["getLogs", invocationId],
    queryFn: () => fetchLog(invocationId),
  });

  const logDownloadUrl = `${env("NEXT_PUBLIC_BES_BACKEND_URL")}/api/v1/invocations/${invocationId}/log`;

  return (
    <PortalCard
      type="inner"
      icon={<FileSearchOutlined />}
      titleBits={["Raw Build Logs"]}
      extraBits={[
        <Tooltip
          key="tooltip"
          title="Bazel emits logs in ANSI format a screen at a time.  They are presented here concatenated for your convenience."
        >
          <ExclamationCircleOutlined />
        </Tooltip>,
        logDownloadUrl && (
          <DownloadButton
            key="downloadButton"
            enabled={true}
            buttonLabel="Download Log"
            fileName="log.txt"
            url={logDownloadUrl}
          />
        ),
      ]}
    >
      <LogViewer loading={isLoading} error={error} log={data} />
    </PortalCard>
  );
};

export default BuildLogsDisplay;
