import { AnsiUp } from "ansi_up";
import { Spin } from "antd";
import React from "react";
import { WindowVirtualizer } from "virtua";
import PortalAlert from "@/components/PortalAlert";
import styles from "./index.module.css";
import { v4 as uuidv4 } from 'uuid';

const ansi = new AnsiUp();

interface Props {
  log?: string | null;
  loading?: boolean;
  error?: Error | null;
}

const LogViewer: React.FC<Props> = ({ log, loading, error }) => {
  const lines = React.useMemo(() => {
    if (!log) return [];
    return ansi.ansi_to_html(log).split("\n").map(line => ({
      line,
      key: uuidv4()
    }));
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
    <pre className={styles.logContainer}>
      <WindowVirtualizer>
        {lines.map((line) => (
          // TODO: Remove the danger
          // biome-ignore lint/security/noDangerouslySetInnerHtml: Should be reworked
          <span key={line.key} dangerouslySetInnerHTML={{ __html: line.line }} />
        ))}
      </WindowVirtualizer>
    </pre>
  );
};

export default LogViewer;
