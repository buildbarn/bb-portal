import type { ApolloError } from "@apollo/client";
import { AnsiUp } from "ansi_up";
import { Card, type CardProps, Spin } from "antd";
import type { RefAttributes } from "react";
import React from "react";
import { JSX } from "react/jsx-runtime";
import { WindowVirtualizer } from "virtua";
import PortalAlert from "@/components/PortalAlert";
import styles from "./index.module.css";

import IntrinsicAttributes = JSX.IntrinsicAttributes;

const ansi = new AnsiUp();

interface Props {
  log?: string | null;
  loading?: boolean;
  error?: ApolloError | Error | null;
}

const LogViewer: React.FC<Props> = ({ log, loading, error }) => {
  const lines = React.useMemo(() => {
    if (!log) return [];
    return ansi.ansi_to_html(log).split("\n");
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
        {lines.map((line, index) => (
          <span key={index} dangerouslySetInnerHTML={{ __html: line }} />
        ))}
      </WindowVirtualizer>
    </pre>
  );
};

export default LogViewer;
