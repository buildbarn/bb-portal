import PortalAlert from "@/components/PortalAlert";
import { AnsiUp } from "ansi_up";
import { Card, type CardProps } from "antd";
import type React from "react";
import type { RefAttributes } from "react";
import { JSX } from "react/jsx-runtime";
import styles from "./index.module.css";
import IntrinsicAttributes = JSX.IntrinsicAttributes;

const ansi = new AnsiUp();

interface Props {
  log?: string | null;
}

const LogViewer: React.FC<Props> = ({ log }) => {
  if (!log) {
    return (
      <PortalAlert
        message="There is no log information to display"
        type="warning"
        showIcon
        className={styles.alert}
      />
    );
  } else {
    const innerHTML = ansi.ansi_to_html(log);
    return <pre dangerouslySetInnerHTML={{ __html: innerHTML }} />;
  }
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
