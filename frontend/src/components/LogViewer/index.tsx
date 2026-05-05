import { DownloadOutlined, FileSearchOutlined } from "@ant-design/icons";
import { Button, Spin } from "antd";
import type React from "react";
import PortalAlert from "@/components/PortalAlert";
import { useBbPortalMessage } from "@/context/MessageContext";
import PortalCard from "../PortalCard";
import { AnsiScrollingWindow } from "./ansiScrollWindow";
import styles from "./index.module.css";

interface Props {
  log?: string | undefined;
  loading?: boolean;
  error?: Error | null;
  title: string;
  logDownloadUrl: string | undefined;
  fileName?: string;
}

const LogViewerCard: React.FC<Props> = ({
  log,
  loading,
  error,
  title,
  logDownloadUrl,
  fileName,
}) => {
  const { copyToClipboard } = useBbPortalMessage();

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
      className={!error ? styles.compactCard : undefined}
      icon={<FileSearchOutlined />}
      titleBits={[title]}
      extraBits={[
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
