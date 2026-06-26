import {
  DownloadOutlined,
  FileSearchOutlined,
  WarningOutlined,
} from "@ant-design/icons";
import { Button, Spin, Tooltip } from "antd";
import type React from "react";
import { useMemo } from "react";
import PortalAlert from "@/components/PortalAlert";
import { useBbPortalMessage } from "@/context/MessageContext";
import { readableFileSize } from "@/utils/filesize";
import PortalCard from "../PortalCard";
import { AnsiScrollingWindow } from "./ansiScrollWindow";
import styles from "./index.module.css";

export const SIZE_BYTE_LIMIT = 1_000_000; // 1MB

const HISTORICAL_EXECUTE_RESPONSE_REGEX =
  /https?:\/\/[-a-zA-Z0-9.]{1,256}(:[0-9]+)?[-a-zA-Z0-9()@:%_+.~#?&/=]*\/blobs\/[a-zA-Z0-9]{0,20}\/historical_execute_response\/[0-9a-f]{64}-[0-9]*\//;

interface Props {
  log?: string | undefined;
  logSizeBytes?: number | undefined;
  loading?: boolean;
  error?: Error | null;
  title: string;
  logDownloadUrl: string | undefined;
  fileName?: string;
}

const LogViewerCard: React.FC<Props> = ({
  log,
  logSizeBytes,
  loading,
  error,
  title,
  logDownloadUrl,
  fileName,
}) => {
  const { copyToClipboard } = useBbPortalMessage();

  const historicalExecuteResponseUrl = useMemo(() => {
    return log?.match(HISTORICAL_EXECUTE_RESPONSE_REGEX)?.[0];
  }, [log]);

  if (loading === true)
    return (
      <Spin>
        <pre />
      </Spin>
    );
  if (error) {
    return (
      <PortalAlert
        type="error"
        message={error.message}
        description={error.cause?.toString()}
        showIcon
        className={styles.alert}
      />
    );
  }
  if (logSizeBytes !== undefined && logSizeBytes > SIZE_BYTE_LIMIT) {
    <PortalAlert
      type="error"
      message={"Output is too large to display."}
      description={`The size of the output is ${readableFileSize(
        logSizeBytes,
      )}. Download the output to view it.`}
      showIcon
      className={styles.alert}
    />;
  }
  if (!log) {
    return (
      <PortalAlert
        message="There is no log information to display"
        type="warning"
        showIcon
        className={styles.alert}
      />
    );
  }
  return (
    <PortalCard
      type="inner"
      styles={{
        body: {
          padding: "0px",
        },
      }}
      reservedTitleWidth={300}
      className={!error ? styles.compactCard : undefined}
      icon={<FileSearchOutlined />}
      titleBits={[
        <div className={styles.titleWrapper} key={"title"}>
          {title}
        </div>,
      ]}
      extraBits={[
        historicalExecuteResponseUrl && (
          <Tooltip title="This URL was extracted from the log, so there are no guarantees that it is correct. It should point to a historical execute response stored in the CAS.">
            <Button
              type="primary"
              href={historicalExecuteResponseUrl}
              target="_blank"
              rel="noopener noreferrer"
            >
              View Historical Execute Response
              <WarningOutlined />
            </Button>
          </Tooltip>
        ),
        logDownloadUrl && (
          <Button icon={<DownloadOutlined />} type="primary">
            <a
              href={logDownloadUrl}
              download={fileName || "log.txt"}
              target="_self"
            >
              Download Log
            </a>
          </Button>
        ),
        log && (
          <Button type="primary" onClick={() => copyToClipboard(log)}>
            Copy to clipboard
          </Button>
        ),
      ]}
    >
      <AnsiScrollingWindow log={log} />
    </PortalCard>
  );
};

export { LogViewerCard };
