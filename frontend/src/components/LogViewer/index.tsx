import PortalAlert from "@/components/PortalAlert";
import type { ApolloError } from "@apollo/client";
import { AnsiUp } from "ansi_up";
import { Card, type CardProps, Spin } from "antd";
import type React from "react";
import type { RefAttributes } from "react";
import { JSX } from "react/jsx-runtime";
import styles from "./index.module.css";
import IntrinsicAttributes = JSX.IntrinsicAttributes;

const ansi = new AnsiUp();

interface Props {
  log?: string | null;
  loading?: boolean;
  error?: ApolloError | undefined;
}

const LogViewer: React.FC<Props> = ({ log, loading, error }) => {
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

  const innerHTML = ansi.ansi_to_html(log);
  return <pre dangerouslySetInnerHTML={{ __html: innerHTML }} />;
};

type LogViewerCardProps = Props &
  IntrinsicAttributes &
  CardProps &
  RefAttributes<HTMLDivElement>;

export const LogViewerCard: React.FC<LogViewerCardProps> = ({
  log,
  bordered,
  ...props
}) => {
  return (
    <Card bordered={false} {...props}>
      <LogViewer log={log} />
    </Card>
  );
};

export default LogViewer;
