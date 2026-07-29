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
  log: string | undefined;
  logSizeBytes?: number;
  loading?: boolean;
  error?: Error | null;
  title: string;
  logDownloadUrl?: string;
  fileName: string;
}

export const LogViewerCard: React.FC<Props> = ({
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

  const renderContent = useMemo(() => {
    if (loading) {
      return (
        <Spin>
          <pre />
        </Spin>
      );
    }
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
      return (
        <PortalAlert
          type="error"
          message="Output is too large to display."
          description={`The size of the output is ${readableFileSize(logSizeBytes)}. Download the output to view it.`}
          showIcon
          className={styles.alert}
        />
      );
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
    return <AnsiScrollingWindow log={log} />;
  }, [loading, error, logSizeBytes, log]);

  return (
    <PortalCard
      type="inner"
      styles={{ body: { padding: "0px" } }}
      reservedTitleWidth={300}
      icon={<FileSearchOutlined />}
      titleBits={[
        <div className={styles.titleWrapper} key="title">
          {title}
        </div>,
      ]}
      extraBits={[
        historicalExecuteResponseUrl && (
          <Tooltip
            key="historical-url"
            title="This URL was extracted from the log, so there are no guarantees that it is correct. It should point to a historical execute response stored in the CAS."
          >
            <Button
              type="primary"
              href={historicalExecuteResponseUrl}
              target="_blank"
              rel="noopener noreferrer"
              icon={<WarningOutlined />}
            >
              View Historical Execute Response
            </Button>
          </Tooltip>
        ),
        logDownloadUrl && (
          <Button
            key="download"
            icon={<DownloadOutlined />}
            type="primary"
            href={logDownloadUrl}
            download={fileName}
          >
            Download Log
          </Button>
        ),
        log && (
          <Button
            key="copy"
            type="primary"
            onClick={() => copyToClipboard(log)}
          >
            Copy to clipboard
          </Button>
        ),
      ]}
    >
      {renderContent}
    </PortalCard>
  );
};
