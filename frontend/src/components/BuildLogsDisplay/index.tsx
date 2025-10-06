import DownloadButton from "@/components/DownloadButton";
import { GET_BUILD_LOGS } from "@/components/LogViewer/graphql";
import PortalCard from "@/components/PortalCard";
import type { GetBuildLogsQuery } from "@/graphql/__generated__/graphql";
import {
  ExclamationCircleOutlined,
  FileSearchOutlined,
} from "@ant-design/icons";
import { useQuery } from "@apollo/client";
import ansiRegex from "ansi-regex";
import { Tooltip } from "antd";
import { useMemo } from "react";
import LogViewer from "../LogViewer";

interface Props {
  invocationId: string;
}

const BuildLogsDisplay: React.FC<Props> = ({ invocationId }) => {
  const { data, error, loading } = useQuery<GetBuildLogsQuery>(GET_BUILD_LOGS, {
    variables: { invocationId: invocationId },
    fetchPolicy: "cache-and-network",
    notifyOnNetworkStatusChange: true,
  });

  const ansiEscapeRegex = ansiRegex();
  const logDownloadUrl = useMemo(
    () =>
      data?.bazelInvocation.buildLogs
        ? `data:text/plain;charset=utf-8,${encodeURIComponent(data?.bazelInvocation.buildLogs.replace(ansiEscapeRegex, ""))}`
        : undefined,
    [data?.bazelInvocation.buildLogs, ansiEscapeRegex],
  );

  return (
    <PortalCard
      type="inner"
      icon={<FileSearchOutlined />}
      titleBits={["Raw Build Logs"]}
      extraBits={[
        <Tooltip title="Bazel emits logs in ANSI format a screen at a time.  They are presented here concatenated for your convenience.">
          <ExclamationCircleOutlined />
        </Tooltip>,
        logDownloadUrl && (
          <DownloadButton
            enabled={true}
            buttonLabel="Download Log"
            fileName="log.txt"
            url={logDownloadUrl}
          />
        ),
      ]}
    >
      <LogViewer
        loading={loading}
        error={error}
        log={data?.bazelInvocation.buildLogs}
      />
    </PortalCard>
  );
};

export default BuildLogsDisplay;
